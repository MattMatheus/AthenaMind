package index

import "testing"

func TestLoadIndexInitializesMissingFile(t *testing.T) {
	root := t.TempDir()
	idx, err := LoadIndex(root)
	if err != nil {
		t.Fatalf("LoadIndex failed: %v", err)
	}
	if idx.SchemaVersion == "" || len(idx.Entries) != 0 {
		t.Fatalf("unexpected initial index: %+v", idx)
	}
}
