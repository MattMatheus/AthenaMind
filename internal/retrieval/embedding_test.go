package retrieval

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestChunkTextBySemanticBoundariesPrefersHeadings(t *testing.T) {
	text := strings.Join([]string{
		"# Alpha",
		"first section content",
		"",
		"## Beta",
		"second section content",
		"",
		"### Gamma",
		"third section content",
	}, "\n")

	chunks := chunkTextBySemanticBoundaries(text, 60)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	if !strings.Contains(chunks[0], "# Alpha") {
		t.Fatalf("expected first chunk to include first heading, got %q", chunks[0])
	}
	foundBeta := false
	for _, chunk := range chunks {
		if strings.Contains(chunk, "## Beta") {
			foundBeta = true
			break
		}
	}
	if !foundBeta {
		t.Fatalf("expected at least one chunk to include second heading: %+v", chunks)
	}
}

func TestChunkTextBySemanticBoundariesRespectsByteCap(t *testing.T) {
	text := strings.Join([]string{
		"# Header",
		strings.Repeat("a", 250),
		"",
		"func Build() {",
		strings.Repeat("b", 250),
		"}",
		"",
		"tail",
	}, "\n")

	maxBytes := 128
	chunks := chunkTextBySemanticBoundaries(text, maxBytes)
	if len(chunks) < 2 {
		t.Fatalf("expected multiple chunks, got %d", len(chunks))
	}
	for i, chunk := range chunks {
		if len([]byte(chunk)) > maxBytes {
			t.Fatalf("chunk %d exceeds cap: %d > %d", i, len([]byte(chunk)), maxBytes)
		}
	}
}

func TestGenerateEmbeddingsChunksAndAveragesOllama(t *testing.T) {
	t.Setenv("AZURE_OPENAI_ENDPOINT", "")
	t.Setenv("AZURE_OPENAI_API_KEY", "")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT_NAME", "")
	t.Setenv("AZURE_OPENAI_API_VERSION", "")

	prompts := make([]string, 0, 8)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		prompt, _ := req["prompt"].(string)
		prompts = append(prompts, prompt)
		_ = json.NewEncoder(w).Encode(map[string]any{
			"embedding": []float64{float64(len(prompt)), 1},
		})
	}))
	defer server.Close()

	text := strings.Join([]string{
		"# One",
		strings.Repeat("alpha ", 350),
		"",
		"## Two",
		strings.Repeat("beta ", 350),
		"",
		"### Three",
		strings.Repeat("gamma ", 350),
	}, "\n")

	vecs, err := GenerateEmbeddings(server.URL, []string{text})
	if err != nil {
		t.Fatalf("GenerateEmbeddings failed: %v", err)
	}
	if len(prompts) < 2 {
		t.Fatalf("expected chunked embedding calls, got %d", len(prompts))
	}
	if len(vecs) != 1 || len(vecs[0]) != 2 {
		t.Fatalf("unexpected vector shape: %+v", vecs)
	}

	var sum float64
	for _, p := range prompts {
		sum += float64(len(p))
	}
	expected := sum / float64(len(prompts))
	if math.Abs(vecs[0][0]-expected) > 1e-9 {
		t.Fatalf("unexpected averaged embedding value: got=%f want=%f", vecs[0][0], expected)
	}
}

func TestGenerateEmbeddingsChunksAzureInputsAndAverages(t *testing.T) {
	maxSeen := 0
	inputsSeen := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var req map[string]any
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Fatalf("decode request: %v", err)
		}
		inputs, _ := req["input"].([]any)
		inputsSeen = len(inputs)
		data := make([]map[string]any, 0, len(inputs))
		for i, in := range inputs {
			s, _ := in.(string)
			if len([]byte(s)) > maxSeen {
				maxSeen = len([]byte(s))
			}
			data = append(data, map[string]any{
				"index":     i,
				"embedding": []float64{float64(len(s)), 1},
			})
		}
		_ = json.NewEncoder(w).Encode(map[string]any{"data": data})
	}))
	defer server.Close()

	t.Setenv("AZURE_OPENAI_ENDPOINT", server.URL)
	t.Setenv("AZURE_OPENAI_API_KEY", "test-key")
	t.Setenv("AZURE_OPENAI_DEPLOYMENT_NAME", "embed-test")
	t.Setenv("AZURE_OPENAI_API_VERSION", "2024-05-01-preview")

	text := strings.Join([]string{
		"# One",
		strings.Repeat("alpha ", 350),
		"",
		"## Two",
		strings.Repeat("beta ", 350),
		"",
		"### Three",
		strings.Repeat("gamma ", 350),
	}, "\n")

	vecs, err := GenerateEmbeddings("http://unused-for-azure", []string{text})
	if err != nil {
		t.Fatalf("GenerateEmbeddings failed: %v", err)
	}
	if inputsSeen < 2 {
		t.Fatalf("expected Azure request to include chunked inputs, got %d", inputsSeen)
	}
	if maxSeen > embeddingChunkMaxBytes {
		t.Fatalf("azure chunk exceeded cap: %d > %d", maxSeen, embeddingChunkMaxBytes)
	}
	if len(vecs) != 1 || len(vecs[0]) != 2 {
		t.Fatalf("unexpected vector shape: %+v", vecs)
	}
}
