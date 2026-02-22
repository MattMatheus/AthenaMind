package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestWriteAndRetrieveSemantic(t *testing.T) {
	root := t.TempDir()

	err := runWrite([]string{
		"--root", root,
		"--id", "getting-started",
		"--title", "Getting Started Prompt",
		"--type", "prompt",
		"--domain", "onboarding",
		"--body", "Use this prompt to onboard a new engineer quickly",
		"--stage", "planning",
		"--reviewer", "maya",
		"--approved=true",
	})
	if err != nil {
		t.Fatalf("runWrite failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "index.yaml")); err != nil {
		t.Fatalf("index.yaml missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "prompts", "onboarding", "getting-started.md")); err != nil {
		t.Fatalf("content markdown missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "metadata", "getting-started.yaml")); err != nil {
		t.Fatalf("metadata missing: %v", err)
	}

	// Retrieval should succeed and return deterministic JSON with required fields.
	// We validate helper behavior directly for determinism instead of parsing stdout.
	idx, err := loadIndex(root)
	if err != nil {
		t.Fatalf("loadIndex failed: %v", err)
	}
	candidates, err := loadCandidates(root, idx.Entries, "")
	if err != nil {
		t.Fatalf("loadCandidates failed: %v", err)
	}
	if len(candidates) != 1 {
		t.Fatalf("expected exactly one candidate, got %d", len(candidates))
	}
	score := semanticScore("onboard engineer prompt", candidates[0])
	if score <= 0 {
		t.Fatalf("expected positive semantic score, got %f", score)
	}
}

func TestWriteFailsDuringAutonomousRun(t *testing.T) {
	t.Setenv("AUTONOMOUS_RUN", "true")
	root := t.TempDir()
	err := runWrite([]string{
		"--root", root,
		"--id", "blocked",
		"--title", "Blocked",
		"--type", "prompt",
		"--domain", "security",
		"--body", "blocked",
		"--stage", "planning",
		"--reviewer", "maya",
		"--approved=true",
	})
	if err == nil {
		t.Fatal("expected write to fail during autonomous run")
	}
}

func TestLoadIndexRejectsUnsupportedMajor(t *testing.T) {
	root := t.TempDir()
	idx := indexFile{
		SchemaVersion: "2.0",
		UpdatedAt:     "2026-02-22T00:00:00Z",
		Entries:       []indexEntry{},
	}
	data, _ := json.Marshal(idx)
	if err := os.WriteFile(filepath.Join(root, "index.yaml"), data, 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	_, err := loadIndex(root)
	if err == nil {
		t.Fatal("expected unsupported major schema to fail")
	}
}

func TestSemanticConfidenceGate(t *testing.T) {
	if isSemanticConfident(0.90, 0.82) {
		t.Fatal("expected low margin to fail confidence gate")
	}
	if !isSemanticConfident(0.90, 0.60) {
		t.Fatal("expected clear top score to pass confidence gate")
	}
}
