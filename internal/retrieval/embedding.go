package retrieval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"
	"time"
)

const DefaultEmbeddingEndpoint = "http://localhost:11434"

var embedFailureCache = struct {
	mu    sync.Mutex
	until map[string]time.Time
}{
	until: map[string]time.Time{},
}

type ollamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

func GenerateEmbedding(endpoint, text string) ([]float64, error) {
	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = DefaultEmbeddingEndpoint
	}
	if strings.TrimSpace(text) == "" {
		return nil, errors.New("empty embedding input")
	}

	if isEndpointTemporarilyUnavailable(endpoint) {
		return nil, fmt.Errorf("embedding endpoint temporarily unavailable: %s", endpoint)
	}

	body, err := json.Marshal(ollamaEmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	})
	if err != nil {
		return nil, err
	}
	url := strings.TrimRight(endpoint, "/") + "/api/embeddings"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{Timeout: 350 * time.Millisecond}
	resp, err := client.Do(req)
	if err != nil {
		markEndpointUnavailable(endpoint)
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		markEndpointUnavailable(endpoint)
		return nil, fmt.Errorf("embedding endpoint returned status %d", resp.StatusCode)
	}

	var parsed ollamaEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		markEndpointUnavailable(endpoint)
		return nil, err
	}
	if len(parsed.Embedding) == 0 {
		markEndpointUnavailable(endpoint)
		return nil, errors.New("embedding response did not include vector")
	}
	return parsed.Embedding, nil
}

func markEndpointUnavailable(endpoint string) {
	embedFailureCache.mu.Lock()
	defer embedFailureCache.mu.Unlock()
	embedFailureCache.until[endpoint] = time.Now().Add(30 * time.Second)
}

func isEndpointTemporarilyUnavailable(endpoint string) bool {
	embedFailureCache.mu.Lock()
	defer embedFailureCache.mu.Unlock()
	until, ok := embedFailureCache.until[endpoint]
	if !ok {
		return false
	}
	if time.Now().Before(until) {
		return true
	}
	delete(embedFailureCache.until, endpoint)
	return false
}

func cosineSimilarity(a, b []float64) float64 {
	if len(a) == 0 || len(b) == 0 || len(a) != len(b) {
		return 0
	}
	var dot float64
	var normA float64
	var normB float64
	for i := range a {
		dot += a[i] * b[i]
		normA += a[i] * a[i]
		normB += b[i] * b[i]
	}
	if normA == 0 || normB == 0 {
		return 0
	}
	return dot / (sqrt(normA) * sqrt(normB))
}

func sqrt(v float64) float64 {
	// Newton iteration keeps dependencies minimal.
	if v <= 0 {
		return 0
	}
	x := v
	for i := 0; i < 8; i++ {
		x = 0.5 * (x + v/x)
	}
	return x
}
