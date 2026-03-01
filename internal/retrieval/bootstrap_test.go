package retrieval

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"athenamind/internal/types"
)

func TestLoadLatestEpisodeFallsBackToRepoLatest(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "episodes", "athenamind")
	if err := os.MkdirAll(repoDir, 0o755); err != nil {
		t.Fatalf("mkdir repo dir: %v", err)
	}

	want := types.EpisodeContext{
		Repo:      "athenamind",
		Scenario:  "episode",
		CycleID:   "cycle-repo",
		StoryID:   "story-repo",
		Outcome:   "success",
		Summary:   "repo-level latest",
		Timestamp: "2026-03-01T00:00:00Z",
	}
	data, _ := json.Marshal(want)
	if err := os.WriteFile(filepath.Join(repoDir, "latest.json"), data, 0o644); err != nil {
		t.Fatalf("write repo latest: %v", err)
	}

	got := loadLatestEpisode(root, "AthenaMind", "engineering")
	if got == nil {
		t.Fatal("expected latest episode from repo fallback path")
	}
	if got.CycleID != "cycle-repo" || got.StoryID != "story-repo" {
		t.Fatalf("unexpected fallback episode payload: %+v", got)
	}
}

func TestLoadLatestEpisodePrefersScenarioLatest(t *testing.T) {
	root := t.TempDir()
	repoDir := filepath.Join(root, "episodes", "athenamind")
	scenarioDir := filepath.Join(repoDir, "engineering")
	if err := os.MkdirAll(scenarioDir, 0o755); err != nil {
		t.Fatalf("mkdir scenario dir: %v", err)
	}

	repoLatest, _ := json.Marshal(types.EpisodeContext{CycleID: "cycle-repo", StoryID: "story-repo"})
	if err := os.WriteFile(filepath.Join(repoDir, "latest.json"), repoLatest, 0o644); err != nil {
		t.Fatalf("write repo latest: %v", err)
	}
	scenarioLatest, _ := json.Marshal(types.EpisodeContext{CycleID: "cycle-scenario", StoryID: "story-scenario"})
	if err := os.WriteFile(filepath.Join(scenarioDir, "latest.json"), scenarioLatest, 0o644); err != nil {
		t.Fatalf("write scenario latest: %v", err)
	}

	got := loadLatestEpisode(root, "AthenaMind", "engineering")
	if got == nil {
		t.Fatal("expected latest episode from scenario path")
	}
	if got.CycleID != "cycle-scenario" || got.StoryID != "story-scenario" {
		t.Fatalf("expected scenario latest to win, got %+v", got)
	}
}
