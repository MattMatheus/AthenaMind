package index

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"athenamind/internal/types"
)

const (
	SupportedMajor = 1
	SupportedMinor = 0
)

func LoadIndex(root string) (types.IndexFile, error) {
	path := filepath.Join(root, "index.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(root, 0o755); err != nil {
				return types.IndexFile{}, err
			}
			return types.IndexFile{SchemaVersion: types.DefaultSchema, UpdatedAt: time.Now().UTC().Format(time.RFC3339), Entries: []types.IndexEntry{}}, nil
		}
		return types.IndexFile{}, err
	}

	var idx types.IndexFile
	if err := json.Unmarshal(data, &idx); err != nil {
		return types.IndexFile{}, fmt.Errorf("ERR_SCHEMA_VERSION_INVALID: cannot parse %s: %w", path, err)
	}
	if strings.TrimSpace(idx.SchemaVersion) == "" {
		return types.IndexFile{}, errors.New("ERR_SCHEMA_VERSION_INVALID: index schema_version is required")
	}
	if err := ValidateSchemaVersion(idx.SchemaVersion); err != nil {
		return types.IndexFile{}, err
	}
	if err := ValidateIndex(idx, root); err != nil {
		return types.IndexFile{}, err
	}
	return idx, nil
}

func ValidateIndex(idx types.IndexFile, root string) error {
	if strings.TrimSpace(idx.UpdatedAt) == "" {
		return errors.New("ERR_SCHEMA_VALIDATION: index updated_at is required")
	}
	if _, err := time.Parse(time.RFC3339, idx.UpdatedAt); err != nil {
		return errors.New("ERR_SCHEMA_VALIDATION: index updated_at must be RFC3339")
	}
	seen := map[string]struct{}{}
	for _, e := range idx.Entries {
		if e.ID == "" || e.Type == "" || e.Domain == "" || e.Path == "" || e.MetadataPath == "" || e.Status == "" || e.UpdatedAt == "" {
			return errors.New("ERR_SCHEMA_VALIDATION: index entry missing required fields")
		}
		if _, ok := seen[e.ID]; ok {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: duplicate index entry id %s", e.ID)
		}
		seen[e.ID] = struct{}{}
		if e.Type != "prompt" && e.Type != "instruction" && e.Type != "episode" {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s has invalid type", e.ID)
		}
		if !IsValidStatus(e.Status) {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s has invalid status", e.ID)
		}
		if _, err := time.Parse(time.RFC3339, e.UpdatedAt); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s has invalid updated_at", e.ID)
		}
		if e.Type == "prompt" && !strings.HasPrefix(e.Path, "prompts/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s path must be under prompts/", e.ID)
		}
		if e.Type == "instruction" && !strings.HasPrefix(e.Path, "instructions/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s path must be under instructions/", e.ID)
		}
		if e.Type == "episode" && !strings.HasPrefix(e.Path, "episodes/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s path must be under episodes/", e.ID)
		}
		if !strings.HasPrefix(e.MetadataPath, "metadata/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s metadata path must be under metadata/", e.ID)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(e.Path))); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s content path does not exist", e.ID)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(e.MetadataPath))); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s metadata path does not exist", e.ID)
		}
	}
	return nil
}

func ValidateSchemaVersion(version string) error {
	major, minor, err := ParseMajorMinor(version)
	if err != nil {
		return fmt.Errorf("ERR_SCHEMA_VERSION_INVALID: %w", err)
	}
	if major > SupportedMajor {
		return fmt.Errorf("ERR_SCHEMA_MAJOR_UNSUPPORTED: schema version %s is newer than supported major %d", version, SupportedMajor)
	}
	if major == SupportedMajor && minor > SupportedMinor {
		fmt.Fprintf(os.Stderr, "WARN_SCHEMA_MINOR_NEWER_COMPAT: operating in compatibility mode for schema version %s\n", version)
	}
	return nil
}

func ParseMajorMinor(v string) (int, int, error) {
	trimmed := strings.TrimSpace(v)
	parts := strings.Split(trimmed, ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("version must be MAJOR.MINOR")
	}
	var major, minor int
	_, err := fmt.Sscanf(trimmed, "%d.%d", &major, &minor)
	if err != nil {
		return 0, 0, fmt.Errorf("version must contain numeric MAJOR.MINOR")
	}
	return major, minor, nil
}

func IsValidStatus(s string) bool {
	return s == "draft" || s == "approved"
}

func WriteJSONAsYAML(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func UpsertEntry(root string, in types.UpsertEntryInput, policy types.WritePolicyDecision) error {
	if in.ID == "" || in.Title == "" || in.Type == "" || in.Domain == "" {
		return errors.New("--id --title --type --domain are required")
	}
	if in.Type != "prompt" && in.Type != "instruction" && in.Type != "episode" {
		return errors.New("--type must be prompt, instruction, or episode")
	}

	entryBody := strings.TrimSpace(in.Body)
	if in.BodyFile != "" {
		data, err := os.ReadFile(in.BodyFile)
		if err != nil {
			return err
		}
		entryBody = strings.TrimSpace(string(data))
	}
	if entryBody == "" {
		return errors.New("entry body is required via --body or --body-file")
	}

	idx, err := LoadIndex(root)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	dirName := "prompts"
	if in.Type == "instruction" {
		dirName = "instructions"
	}
	if in.Type == "episode" {
		dirName = "episodes"
	}

	relContentPath := filepath.ToSlash(filepath.Join(dirName, in.Domain, in.ID+".md"))
	relMetaPath := filepath.ToSlash(filepath.Join("metadata", in.ID+".yaml"))
	contentPath := filepath.Join(root, filepath.FromSlash(relContentPath))
	metaPath := filepath.Join(root, filepath.FromSlash(relMetaPath))
	auditPath := filepath.Join(root, "audits", fmt.Sprintf("%s-%d.json", in.ID, time.Now().UTC().UnixNano()))

	audit := types.MutationAuditRecord{
		SchemaVersion: types.DefaultSchema,
		ID:            in.ID,
		Stage:         in.Stage,
		Decision:      policy.Decision,
		ReviewedBy:    policy.Reviewer,
		ReviewedAt:    now,
		DecisionNotes: policy.Notes,
		Reason:        policy.Reason,
		Risk:          policy.Risk,
		ReworkNotes:   policy.ReworkNotes,
		ReReviewedBy:  policy.ReReviewedBy,
		ChangedFiles: []string{
			relContentPath,
			relMetaPath,
			"index.yaml",
		},
		Applied: policy.Decision == "approved",
	}
	if err := WriteJSONAsYAML(auditPath, audit); err != nil {
		return err
	}
	if policy.Decision == "rejected" {
		return errors.New("ERR_MUTATION_REJECTED_PENDING_REWORK: write blocked until rework and re-review are completed")
	}

	if err := os.MkdirAll(filepath.Dir(contentPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(metaPath), 0o755); err != nil {
		return err
	}

	markdown := fmt.Sprintf("# %s\n\n%s\n", in.Title, entryBody)
	if err := os.WriteFile(contentPath, []byte(markdown), 0o644); err != nil {
		return err
	}

	meta := types.MetadataFile{
		SchemaVersion: types.DefaultSchema,
		ID:            in.ID,
		Title:         in.Title,
		Status:        "approved",
		UpdatedAt:     now,
		Review: types.ReviewMeta{
			ReviewedBy:   policy.Reviewer,
			ReviewedAt:   now,
			Decision:     "approved",
			DecisionNote: policy.Notes,
		},
	}
	if err := WriteJSONAsYAML(metaPath, meta); err != nil {
		return err
	}

	updated := false
	for i := range idx.Entries {
		if idx.Entries[i].ID == in.ID {
			idx.Entries[i] = types.IndexEntry{
				ID:           in.ID,
				Type:         in.Type,
				Domain:       in.Domain,
				Path:         relContentPath,
				MetadataPath: relMetaPath,
				Status:       "approved",
				UpdatedAt:    now,
				Title:        in.Title,
			}
			updated = true
			break
		}
	}
	if !updated {
		idx.Entries = append(idx.Entries, types.IndexEntry{
			ID:           in.ID,
			Type:         in.Type,
			Domain:       in.Domain,
			Path:         relContentPath,
			MetadataPath: relMetaPath,
			Status:       "approved",
			UpdatedAt:    now,
			Title:        in.Title,
		})
	}
	sort.Slice(idx.Entries, func(i, j int) bool { return idx.Entries[i].ID < idx.Entries[j].ID })
	idx.UpdatedAt = now
	if err := WriteJSONAsYAML(filepath.Join(root, "index.yaml"), idx); err != nil {
		return err
	}

	return nil
}
