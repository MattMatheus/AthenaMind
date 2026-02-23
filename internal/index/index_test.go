package index

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"athenamind/internal/types"
)

func TestLoadIndexInitializesMissingFile(t *testing.T) {
	root := t.TempDir()
	idx, err := LoadIndex(root)
	if err != nil {
		t.Fatalf("LoadIndex failed: %v", err)
	}
	if idx.SchemaVersion == "" || len(idx.Entries) != 0 {
		t.Fatalf("unexpected initial index: %+v", idx)
	}
	if _, err := os.Stat(filepath.Join(root, "index.db")); err != nil {
		t.Fatalf("expected sqlite index.db to be initialized: %v", err)
	}
}

func TestLoadIndexMigratesLegacyYAMLToSQLite(t *testing.T) {
	root := t.TempDir()
	now := "2026-02-23T00:00:00Z"

	entry := types.IndexEntry{
		ID:           "entry-1",
		Type:         "prompt",
		Domain:       "ops",
		Path:         "prompts/ops/entry-1.md",
		MetadataPath: "metadata/entry-1.yaml",
		Status:       "approved",
		UpdatedAt:    now,
		Title:        "Entry 1",
	}
	if err := os.MkdirAll(filepath.Join(root, "prompts", "ops"), 0o755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "metadata"), 0o755); err != nil {
		t.Fatalf("mkdir metadata: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, filepath.FromSlash(entry.Path)), []byte("# Entry 1\n\nbody\n"), 0o644); err != nil {
		t.Fatalf("write entry body: %v", err)
	}
	meta := types.MetadataFile{
		SchemaVersion: types.DefaultSchema,
		ID:            entry.ID,
		Title:         entry.Title,
		Status:        "approved",
		UpdatedAt:     now,
		Review: types.ReviewMeta{
			ReviewedBy: "qa",
			ReviewedAt: now,
			Decision:   "approved",
		},
	}
	if err := WriteJSONAsYAML(filepath.Join(root, filepath.FromSlash(entry.MetadataPath)), meta); err != nil {
		t.Fatalf("write metadata: %v", err)
	}
	legacy := types.IndexFile{
		SchemaVersion: types.DefaultSchema,
		UpdatedAt:     now,
		Entries:       []types.IndexEntry{entry},
	}
	if err := WriteJSONAsYAML(filepath.Join(root, "index.yaml"), legacy); err != nil {
		t.Fatalf("write legacy index: %v", err)
	}

	idx, err := LoadIndex(root)
	if err != nil {
		t.Fatalf("LoadIndex failed: %v", err)
	}
	if len(idx.Entries) != 1 || idx.Entries[0].ID != entry.ID {
		t.Fatalf("expected migrated entry, got %+v", idx.Entries)
	}
	if _, err := os.Stat(filepath.Join(root, "index.db")); err != nil {
		t.Fatalf("expected index.db after migration: %v", err)
	}
}

func TestUpsertEntryEnablesSQLiteWAL(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	err := UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-1",
		Title:  "Entry 1",
		Type:   "prompt",
		Domain: "ops",
		Body:   "body",
		Stage:  "pm",
	}, policy)
	if err != nil {
		t.Fatalf("UpsertEntry failed: %v", err)
	}

	cmd := exec.Command("sqlite3", filepath.Join(root, "index.db"), "PRAGMA journal_mode;")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("sqlite pragma failed: %v %s", err, string(out))
	}
	if got := strings.TrimSpace(string(out)); got != "wal" {
		t.Fatalf("expected WAL journal mode, got %q", got)
	}
}
