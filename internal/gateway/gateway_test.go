package gateway

import (
	"testing"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func TestAPIRetrieveWithFallbackWithoutGateway(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{Decision: "approved", Reviewer: "maya", Notes: "ok", Reason: "seed", Risk: "low"}
	err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID: "entry-1", Title: "Entry 1", Type: "prompt", Domain: "ops", Body: "retrieve me", Stage: "planning",
	}, policy)
	if err != nil {
		t.Fatalf("seed entry failed: %v", err)
	}
	resp, err := APIRetrieveWithFallback(root, "", types.APIRetrieveRequest{Query: "entry 1", Domain: "ops", SessionID: "s1"}, "trace-1", nil)
	if err != nil {
		t.Fatalf("APIRetrieveWithFallback failed: %v", err)
	}
	if resp.SelectedID == "" || resp.SelectionMode == "" || resp.SourcePath == "" {
		t.Fatalf("expected retrieval payload fields, got %+v", resp)
	}
}
