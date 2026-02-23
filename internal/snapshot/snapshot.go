package snapshot

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

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func CreateSnapshot(root, createdBy, reason string) (types.SnapshotManifest, error) {
	idx, err := index.LoadIndex(root)
	if err != nil {
		return types.SnapshotManifest{}, err
	}

	now := time.Now().UTC()
	snapshotID := fmt.Sprintf("snapshot-%d", now.UnixNano())
	baseDir := filepath.Join(root, "snapshots", snapshotID)
	payloadDir := filepath.Join(baseDir, "payload")

	refs := collectSnapshotRefs(idx)
	if len(refs) == 0 {
		return types.SnapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: no payload references found for snapshot")
	}

	checksums := make([]types.SnapshotChecksum, 0, len(refs))
	for _, rel := range refs {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(payloadDir, filepath.FromSlash(rel))
		if err := copyFile(src, dst); err != nil {
			return types.SnapshotManifest{}, err
		}
		sum, err := fileSHA256(dst)
		if err != nil {
			return types.SnapshotManifest{}, err
		}
		checksums = append(checksums, types.SnapshotChecksum{Path: rel, SHA256: sum})
	}

	manifest := types.SnapshotManifest{
		SnapshotID:    snapshotID,
		CreatedAt:     now.Format(time.RFC3339),
		CreatedBy:     strings.TrimSpace(createdBy),
		SchemaVersion: types.DefaultSchema,
		IndexVersion:  idx.SchemaVersion,
		Scope:         "full",
		Reason:        strings.TrimSpace(reason),
		Checksums:     checksums,
		PayloadRefs:   refs,
	}
	if err := ValidateSnapshotManifest(manifest); err != nil {
		return types.SnapshotManifest{}, err
	}
	if err := index.WriteJSONAsYAML(filepath.Join(baseDir, "manifest.json"), manifest); err != nil {
		return types.SnapshotManifest{}, err
	}
	return manifest, nil
}

func ListSnapshots(root string) ([]types.SnapshotListRow, error) {
	snapRoot := filepath.Join(root, "snapshots")
	if _, err := os.Stat(snapRoot); errors.Is(err, os.ErrNotExist) {
		return []types.SnapshotListRow{}, nil
	}

	entries, err := os.ReadDir(snapRoot)
	if err != nil {
		return nil, err
	}
	rows := make([]types.SnapshotListRow, 0, len(entries))
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		manifest, err := LoadSnapshotManifest(root, ent.Name())
		if err != nil {
			return nil, err
		}
		rows = append(rows, types.SnapshotListRow{
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

func RestoreSnapshot(root, snapshotID string) error {
	manifest, err := LoadSnapshotManifest(root, snapshotID)
	if err != nil {
		return err
	}
	if err := checkSnapshotCompatibility(root, manifest); err != nil {
		return err
	}
	if err := verifySnapshotChecksums(root, manifest); err != nil {
		return err
	}

	for _, rel := range manifest.PayloadRefs {
		src := filepath.Join(root, "snapshots", snapshotID, "payload", filepath.FromSlash(rel))
		dst := filepath.Join(root, filepath.FromSlash(rel))
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func LoadSnapshotManifest(root, snapshotID string) (types.SnapshotManifest, error) {
	path := filepath.Join(root, "snapshots", snapshotID, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return types.SnapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: snapshot manifest not found")
		}
		return types.SnapshotManifest{}, err
	}
	var m types.SnapshotManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return types.SnapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: cannot parse manifest")
	}
	if err := ValidateSnapshotManifest(m); err != nil {
		return types.SnapshotManifest{}, err
	}
	return m, nil
}

func ValidateSnapshotManifest(m types.SnapshotManifest) error {
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

func checkSnapshotCompatibility(root string, m types.SnapshotManifest) error {
	if err := index.ValidateSchemaVersion(m.SchemaVersion); err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	if err := index.ValidateSchemaVersion(m.IndexVersion); err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	cur, err := index.LoadIndex(root)
	if err != nil {
		return err
	}
	snapMajor, _, err := index.ParseMajorMinor(m.IndexVersion)
	if err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	curMajor, _, err := index.ParseMajorMinor(cur.SchemaVersion)
	if err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	if snapMajor != curMajor {
		return errors.New("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: snapshot index major version does not match current index major")
	}
	return nil
}

func verifySnapshotChecksums(root string, m types.SnapshotManifest) error {
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

func collectSnapshotRefs(idx types.IndexFile) []string {
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

func WriteSnapshotAudit(root string, ev types.SnapshotAuditEvent) error {
	name := strings.ReplaceAll(ev.EventName, ".", "-")
	path := filepath.Join(root, "audits", fmt.Sprintf("%s-%d.json", name, time.Now().UTC().UnixNano()))
	return index.WriteJSONAsYAML(path, ev)
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
