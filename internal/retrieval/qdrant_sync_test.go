package retrieval

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func TestSyncQdrantCollection(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-a",
		Title:  "Entry A",
		Type:   "prompt",
		Domain: "ops",
		Body:   "a body",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed entry-a: %v", err)
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-b",
		Title:  "Entry B",
		Type:   "prompt",
		Domain: "ops",
		Body:   "b body",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed entry-b: %v", err)
	}

	if err := index.UpsertEmbedding(root, "entry-a", []float64{1.0, 0.0}); err != nil {
		t.Fatalf("seed embedding entry-a: %v", err)
	}

	upserts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Fatalf("unexpected method: %s", r.Method)
		}
		var payload struct {
			Points []map[string]any `json:"points"`
		}
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			t.Fatalf("decode upsert payload: %v", err)
		}
		upserts += len(payload.Points)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))
	defer server.Close()

	report, err := SyncQdrantCollection(root, server.URL, "athena_memories", 16)
	if err != nil {
		t.Fatalf("SyncQdrantCollection failed: %v", err)
	}
	if report.Synced != 1 {
		t.Fatalf("expected 1 synced embedding, got %+v", report)
	}
	if upserts != 1 {
		t.Fatalf("expected one upsert point, got %d", upserts)
	}
}
