package retrieval

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"athenamind/internal/index"
	"athenamind/internal/types"
)

func TestIsSemanticConfident(t *testing.T) {
	if IsSemanticConfident(0.90, 0.82) {
		t.Fatal("expected low margin to fail")
	}
	if !IsSemanticConfident(0.90, 0.60) {
		t.Fatal("expected clear margin to pass")
	}
}

func TestRetrieveUsesEmbeddingSimilarityWhenAvailable(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "net",
		Title:  "Network Runbook",
		Type:   "prompt",
		Domain: "ops",
		Body:   "network incident runbook with rollback",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed net entry: %v", err)
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "payroll",
		Title:  "Payroll Policy",
		Type:   "prompt",
		Domain: "ops",
		Body:   "payroll tax policy and handbook",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed payroll entry: %v", err)
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode embedding req: %v", err)
		}
		prompt, _ := req["prompt"].(string)
		vec := []float64{0.0, 1.0}
		if strings.Contains(strings.ToLower(prompt), "network") {
			vec = []float64{1.0, 0.0}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": vec})
	}))
	defer server.Close()

	if warning, err := IndexEntryEmbedding(root, "net", server.URL, "sess-1"); err != nil || warning != "" {
		t.Fatalf("index embedding net failed warning=%q err=%v", warning, err)
	}
	if warning, err := IndexEntryEmbedding(root, "payroll", server.URL, "sess-1"); err != nil || warning != "" {
		t.Fatalf("index embedding payroll failed warning=%q err=%v", warning, err)
	}

	got, warning, err := RetrieveWithEmbeddingEndpoint(root, "network outage playbook", "ops", server.URL)
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if warning != "" {
		t.Fatalf("unexpected warning: %s", warning)
	}
	if got.SelectedID != "net" {
		t.Fatalf("expected net selected, got %+v", got)
	}
	if got.SelectionMode != "embedding_semantic" {
		t.Fatalf("expected embedding_semantic mode, got %s", got.SelectionMode)
	}
}

func TestRetrieveFallsBackWhenEmbeddingEndpointUnavailable(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry",
		Title:  "Entry",
		Type:   "prompt",
		Domain: "ops",
		Body:   "fallback body",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed entry failed: %v", err)
	}

	got, warning, err := RetrieveWithEmbeddingEndpoint(root, "fallback", "ops", "http://127.0.0.1:1")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if strings.TrimSpace(warning) == "" {
		t.Fatal("expected warning when embedding endpoint is unavailable")
	}
	if got.SelectedID == "" {
		t.Fatalf("expected a deterministic fallback selection, got %+v", got)
	}
}

func TestEmbeddingScoringImprovesOverTokenBaseline(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "alpha-doc",
		Title:  "Alpha Policy",
		Type:   "prompt",
		Domain: "ops",
		Body:   "contains uncommon token zzzzalpha",
		Stage:  "pm",
	}, policy)
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "beta-doc",
		Title:  "Beta Policy",
		Type:   "prompt",
		Domain: "ops",
		Body:   "contains uncommon token zzzzbeta",
		Stage:  "pm",
	}, policy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		prompt, _ := req["prompt"].(string)
		p := strings.ToLower(prompt)
		vec := []float64{1.0, 0.0}
		if strings.Contains(p, "zzzzbeta") || strings.Contains(p, "incident") || strings.Contains(p, "restore") || strings.Contains(p, "downtime") {
			vec = []float64{0.0, 1.0}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": vec})
	}))
	defer server.Close()

	if _, err := IndexEntryEmbedding(root, "alpha-doc", server.URL, "sess-2"); err != nil {
		t.Fatalf("index alpha embedding: %v", err)
	}
	if _, err := IndexEntryEmbedding(root, "beta-doc", server.URL, "sess-2"); err != nil {
		t.Fatalf("index beta embedding: %v", err)
	}

	queries := []types.EvaluationQuery{
		{Query: "incident response guide", Domain: "ops", ExpectedID: "beta-doc"},
		{Query: "restore procedure", Domain: "ops", ExpectedID: "beta-doc"},
		{Query: "downtime escalation", Domain: "ops", ExpectedID: "beta-doc"},
	}

	tokenHits := 0
	for _, q := range queries {
		got, _, err := RetrieveWithEmbeddingEndpoint(root, q.Query, q.Domain, "http://127.0.0.1:1")
		if err != nil {
			t.Fatalf("token baseline retrieve failed: %v", err)
		}
		if got.SelectedID == q.ExpectedID {
			tokenHits++
		}
	}

	embedHits := 0
	for _, q := range queries {
		got, _, err := RetrieveWithEmbeddingEndpoint(root, q.Query, q.Domain, server.URL)
		if err != nil {
			t.Fatalf("embedding retrieve failed: %v", err)
		}
		if got.SelectedID == q.ExpectedID {
			embedHits++
		}
	}

	if embedHits <= tokenHits {
		t.Fatalf("expected embedding mode to improve over token baseline, token=%d embed=%d", tokenHits, embedHits)
	}
}

func TestRetrieveSkipsIncompatibleEmbeddingDimensions(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry",
		Title:  "Ops Guide",
		Type:   "prompt",
		Domain: "ops",
		Body:   "fallback body content",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed entry failed: %v", err)
	}

	indexServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": []float64{1.0, 0.0}})
	}))
	defer indexServer.Close()
	if _, err := IndexEntryEmbedding(root, "entry", indexServer.URL, "sess-compat"); err != nil {
		t.Fatalf("index embedding failed: %v", err)
	}

	retrieveServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": []float64{1.0, 0.0, 0.0}})
	}))
	defer retrieveServer.Close()

	got, warning, err := RetrieveWithEmbeddingEndpoint(root, "ops guide", "ops", retrieveServer.URL)
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if got.SelectionMode == "embedding_semantic" {
		t.Fatalf("expected embedding scoring to be skipped on dimension mismatch, got %+v", got)
	}
	if strings.TrimSpace(warning) == "" {
		t.Fatal("expected warning when embeddings cannot be applied")
	}
}

func TestRetrievePrefersSessionFreshnessWhenScoresTie(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "a-entry",
		Title:  "Entry A",
		Type:   "prompt",
		Domain: "ops",
		Body:   "same semantic body",
		Stage:  "pm",
	}, policy)
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "b-entry",
		Title:  "Entry B",
		Type:   "prompt",
		Domain: "ops",
		Body:   "same semantic body",
		Stage:  "pm",
	}, policy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": []float64{1.0, 0.0}})
	}))
	defer server.Close()

	if _, err := IndexEntryEmbedding(root, "a-entry", server.URL, "sess-a"); err != nil {
		t.Fatalf("index a-entry: %v", err)
	}
	if _, err := IndexEntryEmbedding(root, "b-entry", server.URL, "sess-b"); err != nil {
		t.Fatalf("index b-entry: %v", err)
	}

	got, _, err := RetrieveWithEmbeddingEndpointAndSession(root, "same semantic body", "ops", server.URL, "sess-b")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if got.SelectedID != "b-entry" {
		t.Fatalf("expected session-fresh entry to win tie, got %+v", got)
	}
}
