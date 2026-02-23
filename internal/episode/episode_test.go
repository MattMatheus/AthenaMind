package episode

import (
	"testing"

	"athenamind/internal/types"
)

func TestWriteAndListEpisode(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{Decision: "approved", Reviewer: "maya", Notes: "ok", Reason: "seed", Risk: "low"}
	_, err := Write(root, types.WriteEpisodeInput{
		Repo:         "AthenaMind",
		SessionID:    "sess-1",
		CycleID:      "cycle-1",
		StoryID:      "story-1",
		Outcome:      "success",
		Summary:      "completed",
		FilesChanged: "a.go,b.go",
		Decisions:    "kept scope",
		Stage:        "pm",
	}, policy)
	if err != nil {
		t.Fatalf("Write failed: %v", err)
	}
	rows, err := List(root, "athenamind")
	if err != nil {
		t.Fatalf("List failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 episode, got %d", len(rows))
	}
}
