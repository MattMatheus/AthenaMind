package retrieval

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
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

func TestIsEmbeddingSemanticConfident(t *testing.T) {
	if IsEmbeddingSemanticConfident(0.19, 0.01) {
		t.Fatal("expected low embedding confidence to fail")
	}
	if IsEmbeddingSemanticConfident(0.25, 0.24) {
		t.Fatal("expected low embedding margin to fail")
	}
	if !IsEmbeddingSemanticConfident(0.25, 0.20) {
		t.Fatal("expected embedding confidence gate to pass")
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

func TestFallbackPrefersHigherScoreThenPath(t *testing.T) {
	cs := []candidate{
		{
			Entry: types.IndexEntry{ID: "b", Path: "z/path-b.md"},
			Score: 0.61,
		},
		{
			Entry: types.IndexEntry{ID: "a", Path: "a/path-a.md"},
			Score: 0.61,
		},
		{
			Entry: types.IndexEntry{ID: "c", Path: "b/path-c.md"},
			Score: 0.42,
		},
	}
	got := chooseDeterministicFallback(cs)
	if got.Entry.ID != "a" {
		t.Fatalf("expected path tie-break for top score candidates, got %+v", got.Entry)
	}
}

func TestGenerateEmbeddingsChunksLongInputs(t *testing.T) {
	var calls int
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		prompt, _ := req["prompt"].(string)
		if len(prompt) > 3000 {
			w.WriteHeader(http.StatusInternalServerError)
			_, _ = w.Write([]byte(`{"error":"input length exceeds the context length"}`))
			return
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": []float64{1.0, 3.0}})
	}))
	defer server.Close()

	longText := strings.Repeat("very long markdown paragraph ", 600)
	vecs, err := GenerateEmbeddings(server.URL, []string{longText})
	if err != nil {
		t.Fatalf("GenerateEmbeddings failed: %v", err)
	}
	if len(vecs) != 1 || len(vecs[0]) != 2 {
		t.Fatalf("unexpected embedding shape: %+v", vecs)
	}
	if calls < 2 {
		t.Fatalf("expected long input to be chunked across multiple requests, calls=%d", calls)
	}
}

func TestGenerateEmbeddingsUsesExplicitEndpointEvenWhenAzureEnvSet(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "https://azure.example.invalid")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT_NAME", "test-deployment")
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": []float64{2.0, 4.0}})
	}))
	defer server.Close()

	vecs, err := GenerateEmbeddings(server.URL, []string{"short text"})
	if err != nil {
		t.Fatalf("GenerateEmbeddings failed: %v", err)
	}
	if len(vecs) != 1 || len(vecs[0]) != 2 {
		t.Fatalf("unexpected embedding shape: %+v", vecs)
	}
}

func TestLoadEmbeddingsCachedDoesNotCollideAcrossCandidateSets(t *testing.T) {
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
		Domain: "ops-a",
		Body:   "entry a",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed entry-a failed: %v", err)
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-b",
		Title:  "Entry B",
		Type:   "prompt",
		Domain: "ops-b",
		Body:   "entry b",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed entry-b failed: %v", err)
	}
	if err := index.UpsertEmbedding(root, "entry-a", []float64{1, 0}); err != nil {
		t.Fatalf("seed embedding entry-a failed: %v", err)
	}
	if err := index.UpsertEmbedding(root, "entry-b", []float64{0, 1}); err != nil {
		t.Fatalf("seed embedding entry-b failed: %v", err)
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		t.Fatalf("load index failed: %v", err)
	}

	var entryA, entryB types.IndexEntry
	for _, e := range idx.Entries {
		switch e.ID {
		case "entry-a":
			entryA = e
		case "entry-b":
			entryB = e
		}
	}
	clearRetrievalCaches()
	first, err := loadEmbeddingsCached(root, idx, []candidate{{Entry: entryA}})
	if err != nil {
		t.Fatalf("first loadEmbeddingsCached failed: %v", err)
	}
	if _, ok := first["entry-a"]; !ok {
		t.Fatalf("expected entry-a in first embedding set, got %+v", first)
	}

	second, err := loadEmbeddingsCached(root, idx, []candidate{{Entry: entryB}})
	if err != nil {
		t.Fatalf("second loadEmbeddingsCached failed: %v", err)
	}
	if _, ok := second["entry-b"]; !ok {
		t.Fatalf("expected entry-b in second embedding set, got %+v", second)
	}
	if _, wrong := second["entry-a"]; wrong {
		t.Fatalf("unexpected entry-a leak in second embedding set: %+v", second)
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

func TestRetrieveHybridReturnsCandidates(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "alpha",
		Title:  "Alpha Runbook",
		Type:   "prompt",
		Domain: "ops",
		Body:   "network runbook alpha",
		Stage:  "pm",
	}, policy)
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "beta",
		Title:  "Beta Guide",
		Type:   "prompt",
		Domain: "ops",
		Body:   "incident restore beta",
		Stage:  "pm",
	}, policy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		_ = json.NewDecoder(r.Body).Decode(&req)
		prompt, _ := req["prompt"].(string)
		p := strings.ToLower(prompt)
		vec := []float64{1.0, 0.0}
		if strings.Contains(p, "beta") || strings.Contains(p, "incident") {
			vec = []float64{0.0, 1.0}
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": vec})
	}))
	defer server.Close()

	if _, err := IndexEntryEmbedding(root, "alpha", server.URL, "sess-hybrid"); err != nil {
		t.Fatalf("index alpha: %v", err)
	}
	if _, err := IndexEntryEmbedding(root, "beta", server.URL, "sess-hybrid"); err != nil {
		t.Fatalf("index beta: %v", err)
	}

	got, warning, err := RetrieveWithOptionsAndEndpointAndSession(
		root,
		"incident restore",
		"ops",
		server.URL,
		"sess-hybrid",
		RetrieveOptions{Mode: "hybrid", TopK: 2},
	)
	if err != nil {
		t.Fatalf("hybrid retrieve failed: %v", err)
	}
	if warning != "" {
		t.Fatalf("unexpected warning: %s", warning)
	}
	if got.SelectionMode != "hybrid_rrf" {
		t.Fatalf("expected hybrid_rrf, got %s", got.SelectionMode)
	}
	if len(got.Candidates) != 2 {
		t.Fatalf("expected 2 candidates, got %d", len(got.Candidates))
	}
	if got.Candidates[0].ID == "" || got.Candidates[0].FusedScore <= 0 {
		t.Fatalf("expected fused candidate scores, got %+v", got.Candidates[0])
	}
}

func TestRetrieveSkipsCorruptEntryAndReturnsWarning(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "healthy",
		Title:  "Healthy Entry",
		Type:   "prompt",
		Domain: "ops",
		Body:   "healthy retrieval target",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed healthy entry failed: %v", err)
	}
	if err := index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "broken",
		Title:  "Broken Entry",
		Type:   "prompt",
		Domain: "ops",
		Body:   "broken retrieval target",
		Stage:  "pm",
	}, policy); err != nil {
		t.Fatalf("seed broken entry failed: %v", err)
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		t.Fatalf("load index failed: %v", err)
	}
	brokenMeta := ""
	for _, entry := range idx.Entries {
		if entry.ID == "broken" {
			brokenMeta = filepath.Join(root, filepath.FromSlash(entry.MetadataPath))
			break
		}
	}
	if brokenMeta == "" {
		t.Fatal("broken metadata path not found")
	}
	if err := os.WriteFile(brokenMeta, []byte("{"), 0o644); err != nil {
		t.Fatalf("corrupt metadata: %v", err)
	}

	got, warning, err := RetrieveWithEmbeddingEndpoint(root, "healthy retrieval", "ops", "http://127.0.0.1:1")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if got.SelectedID != "healthy" {
		t.Fatalf("expected healthy entry selected, got %+v", got)
	}
	if !strings.Contains(warning, "skipped 1 invalid entries") {
		t.Fatalf("expected skipped-candidate warning, got %q", warning)
	}
}

func TestRetrieveHybridFallsBackWithoutSignal(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-a",
		Title:  "Entry A",
		Type:   "prompt",
		Domain: "ops",
		Body:   "alpha alpha alpha",
		Stage:  "pm",
	}, policy)
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-b",
		Title:  "Entry B",
		Type:   "prompt",
		Domain: "ops",
		Body:   "beta beta beta",
		Stage:  "pm",
	}, policy)

	got, _, err := RetrieveWithOptionsAndEndpointAndSession(
		root,
		"qqqqq",
		"ops",
		"http://127.0.0.1:1",
		"sess-hybrid-fallback",
		RetrieveOptions{Mode: "hybrid", TopK: 3},
	)
	if err != nil {
		t.Fatalf("hybrid fallback retrieve failed: %v", err)
	}
	if got.SelectionMode != "fallback_path_priority" {
		t.Fatalf("expected deterministic fallback in no-signal hybrid, got %s", got.SelectionMode)
	}
	if len(got.Candidates) == 0 {
		t.Fatal("expected candidate traces in fallback output")
	}
}

func TestRetrieveQdrantBackendUnavailableFallsBackGracefully(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "net",
		Title:  "Network Runbook",
		Type:   "prompt",
		Domain: "ops",
		Body:   "network rollback playbook",
		Stage:  "pm",
	}, policy)

	embedServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{"embedding": []float64{1.0, 0.0}})
	}))
	defer embedServer.Close()

	if _, err := IndexEntryEmbedding(root, "net", embedServer.URL, "sess-qdrant"); err != nil {
		t.Fatalf("index embedding failed: %v", err)
	}

	t.Setenv("ATHENA_QDRANT_URL", "http://127.0.0.1:1")
	got, warning, err := RetrieveWithOptionsAndEndpointAndSession(
		root,
		"network rollback",
		"ops",
		embedServer.URL,
		"sess-qdrant",
		RetrieveOptions{
			Mode:    "hybrid",
			TopK:    3,
			Backend: "qdrant",
		},
	)
	if err != nil {
		t.Fatalf("retrieve should succeed with qdrant fallback, got err: %v", err)
	}
	if strings.TrimSpace(warning) == "" || !strings.Contains(warning, "qdrant backend unavailable") {
		t.Fatalf("expected qdrant fallback warning, got: %q", warning)
	}
	if got.SelectedID == "" {
		t.Fatalf("expected fallback/local selection, got %+v", got)
	}
}

func TestRetrieveNeo4jBackendBoostsCandidates(t *testing.T) {
	root := t.TempDir()
	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: "qa",
		Notes:    "ok",
		Reason:   "test",
		Risk:     "low",
	}
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-a",
		Title:  "Entry A",
		Type:   "prompt",
		Domain: "ops",
		Body:   "alpha body",
		Stage:  "pm",
	}, policy)
	_ = index.UpsertEntry(root, types.UpsertEntryInput{
		ID:     "entry-b",
		Title:  "Entry B",
		Type:   "prompt",
		Domain: "ops",
		Body:   "beta body",
		Stage:  "pm",
	}, policy)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(map[string]any{
			"results": []any{
				map[string]any{
					"data": []any{
						map[string]any{"row": []any{"entry-a", 9}},
						map[string]any{"row": []any{"entry-b", 1}},
					},
				},
			},
			"errors": []any{},
		})
	}))
	defer server.Close()

	t.Setenv("ATHENA_NEO4J_HTTP_URL", server.URL)
	t.Setenv("ATHENA_NEO4J_USER", "neo4j")
	t.Setenv("ATHENA_NEO4J_PASSWORD", "devpassword")
	t.Setenv("ATHENA_NEO4J_DATABASE", "neo4j")

	got, warning, err := RetrieveWithOptionsAndEndpointAndSession(
		root,
		"unrelated-token",
		"ops",
		"http://127.0.0.1:1",
		"sess-neo",
		RetrieveOptions{
			Mode:    "hybrid",
			TopK:    2,
			Backend: "neo4j",
		},
	)
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if strings.TrimSpace(warning) == "" {
		t.Fatal("expected embedding warning due unavailable endpoint")
	}
	if got.SelectionMode != "hybrid_rrf" {
		t.Fatalf("expected hybrid mode selection, got %s", got.SelectionMode)
	}
	if len(got.Candidates) == 0 || got.Candidates[0].BackendScore <= 0 {
		t.Fatalf("expected backend score in candidates, got %+v", got.Candidates)
	}
	if got.SelectedID != "entry-a" {
		t.Fatalf("expected neo4j boosted candidate entry-a, got %s", got.SelectedID)
	}
}
