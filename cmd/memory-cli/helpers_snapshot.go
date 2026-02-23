package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func createSnapshot(root, createdBy, reason string) (snapshotManifest, error) {
	idx, err := loadIndex(root)
	if err != nil {
		return snapshotManifest{}, err
	}

	now := time.Now().UTC()
	snapshotID := fmt.Sprintf("snapshot-%d", now.UnixNano())
	baseDir := filepath.Join(root, "snapshots", snapshotID)
	payloadDir := filepath.Join(baseDir, "payload")

	refs := collectSnapshotRefs(idx)
	if len(refs) == 0 {
		return snapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: no payload references found for snapshot")
	}

	checksums := make([]snapshotChecksum, 0, len(refs))
	for _, rel := range refs {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(payloadDir, filepath.FromSlash(rel))
		if err := copyFile(src, dst); err != nil {
			return snapshotManifest{}, err
		}
		sum, err := fileSHA256(dst)
		if err != nil {
			return snapshotManifest{}, err
		}
		checksums = append(checksums, snapshotChecksum{Path: rel, SHA256: sum})
	}

	manifest := snapshotManifest{
		SnapshotID:    snapshotID,
		CreatedAt:     now.Format(time.RFC3339),
		CreatedBy:     strings.TrimSpace(createdBy),
		SchemaVersion: defaultSchema,
		IndexVersion:  idx.SchemaVersion,
		Scope:         "full",
		Reason:        strings.TrimSpace(reason),
		Checksums:     checksums,
		PayloadRefs:   refs,
	}
	if err := validateSnapshotManifest(manifest); err != nil {
		return snapshotManifest{}, err
	}
	if err := writeJSONAsYAML(filepath.Join(baseDir, "manifest.json"), manifest); err != nil {
		return snapshotManifest{}, err
	}
	return manifest, nil
}

func listSnapshots(root string) ([]snapshotListRow, error) {
	snapRoot := filepath.Join(root, "snapshots")
	if _, err := os.Stat(snapRoot); errors.Is(err, os.ErrNotExist) {
		return []snapshotListRow{}, nil
	}

	entries, err := os.ReadDir(snapRoot)
	if err != nil {
		return nil, err
	}
	rows := make([]snapshotListRow, 0, len(entries))
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		manifest, err := loadSnapshotManifest(root, ent.Name())
		if err != nil {
			return nil, err
		}
		rows = append(rows, snapshotListRow{
			SnapshotID:    manifest.SnapshotID,
			CreatedAt:     manifest.CreatedAt,
			CreatedBy:     manifest.CreatedBy,
			SchemaVersion: manifest.SchemaVersion,
			IndexVersion:  manifest.IndexVersion,
			Scope:         manifest.Scope,
			Reason:        manifest.Reason,
		})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].CreatedAt > rows[j].CreatedAt })
	return rows, nil
}

func restoreSnapshot(root, snapshotID string) error {
	manifest, err := loadSnapshotManifest(root, snapshotID)
	if err != nil {
		return err
	}
	if err := checkSnapshotCompatibility(root, manifest); err != nil {
		return err
	}
	if err := verifySnapshotChecksums(root, manifest); err != nil {
		return err
	}

	// restore as forward revision by replacing active payload paths while preserving audits/snapshots history
	for _, rel := range manifest.PayloadRefs {
		src := filepath.Join(root, "snapshots", snapshotID, "payload", filepath.FromSlash(rel))
		dst := filepath.Join(root, filepath.FromSlash(rel))
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func loadSnapshotManifest(root, snapshotID string) (snapshotManifest, error) {
	path := filepath.Join(root, "snapshots", snapshotID, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return snapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: snapshot manifest not found")
		}
		return snapshotManifest{}, err
	}
	var m snapshotManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return snapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: cannot parse manifest")
	}
	if err := validateSnapshotManifest(m); err != nil {
		return snapshotManifest{}, err
	}
	return m, nil
}

func validateSnapshotManifest(m snapshotManifest) error {
	if strings.TrimSpace(m.SnapshotID) == "" || strings.TrimSpace(m.CreatedAt) == "" || strings.TrimSpace(m.CreatedBy) == "" ||
		strings.TrimSpace(m.SchemaVersion) == "" || strings.TrimSpace(m.IndexVersion) == "" || strings.TrimSpace(m.Scope) == "" || strings.TrimSpace(m.Reason) == "" {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: missing required manifest fields")
	}
	if m.Scope != "full" {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: unsupported snapshot scope")
	}
	if _, err := time.Parse(time.RFC3339, m.CreatedAt); err != nil {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: created_at must be RFC3339")
	}
	if len(m.PayloadRefs) == 0 || len(m.Checksums) == 0 {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: payload references and checksums are required")
	}
	return nil
}

func checkSnapshotCompatibility(root string, m snapshotManifest) error {
	if err := validateSchemaVersion(m.SchemaVersion); err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	if err := validateSchemaVersion(m.IndexVersion); err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	cur, err := loadIndex(root)
	if err != nil {
		return err
	}
	snapMajor, _, err := parseMajorMinor(m.IndexVersion)
	if err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	curMajor, _, err := parseMajorMinor(cur.SchemaVersion)
	if err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	if snapMajor != curMajor {
		return errors.New("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: snapshot index major version does not match current index major")
	}
	return nil
}

func verifySnapshotChecksums(root string, m snapshotManifest) error {
	expected := map[string]string{}
	for _, c := range m.Checksums {
		expected[c.Path] = c.SHA256
	}
	for _, rel := range m.PayloadRefs {
		sum, err := fileSHA256(filepath.Join(root, "snapshots", m.SnapshotID, "payload", filepath.FromSlash(rel)))
		if err != nil {
			return fmt.Errorf("ERR_SNAPSHOT_INTEGRITY_CHECK_FAILED: %w", err)
		}
		if expected[rel] == "" || expected[rel] != sum {
			return errors.New("ERR_SNAPSHOT_INTEGRITY_CHECK_FAILED: checksum mismatch")
		}
	}
	return nil
}

func collectSnapshotRefs(idx indexFile) []string {
	refs := []string{"index.yaml"}
	for _, e := range idx.Entries {
		refs = append(refs, e.Path, e.MetadataPath)
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		r = filepath.ToSlash(strings.TrimSpace(r))
		if r == "" {
			continue
		}
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

func writeSnapshotAudit(root string, ev snapshotAuditEvent) error {
	name := strings.ReplaceAll(ev.EventName, ".", "-")
	path := filepath.Join(root, "audits", fmt.Sprintf("%s-%d.json", name, time.Now().UTC().UnixNano()))
	return writeJSONAsYAML(path, ev)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
}
