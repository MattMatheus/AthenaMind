package retrieval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

const DefaultEmbeddingEndpoint = "http://localhost:11434"

var (
	embedFailureCache = struct {
		mu    sync.Mutex
		until map[string]time.Time
	}{
		until: map[string]time.Time{},
	}

	azureTokenCache = struct {
		mu     sync.Mutex
		token  string
		expiry time.Time
	}{}
)

type ollamaEmbeddingRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
}

type ollamaEmbeddingResponse struct {
	Embedding []float64 `json:"embedding"`
}

type azureEmbeddingRequest struct {
	Input []string `json:"input"`
}

type azureEmbeddingResponse struct {
	Data []struct {
		Embedding []float64 `json:"embedding"`
		Index     int       `json:"index"`
	} `json:"data"`
}

type EmbeddingProfile struct {
	Provider string
	ModelID  string
}

func ActiveEmbeddingProfile(endpoint string) EmbeddingProfile {
	_ = endpoint
	return EmbeddingProfile{
		Provider: "ollama",
		ModelID:  "nomic-embed-text",
	}
}

func GenerateEmbedding(endpoint, text string) ([]float64, error) {
	vecs, err := GenerateEmbeddings(endpoint, []string{text})
	if err != nil {
		return nil, err
	}
	if len(vecs) == 0 {
		return nil, errors.New("no embedding returned")
	}
	return vecs[0], nil
}

func GenerateEmbeddings(endpoint string, texts []string) ([][]float64, error) {
	if len(texts) == 0 {
		return nil, nil
	}

	endpoint = strings.TrimSpace(endpoint)
	if endpoint == "" {
		endpoint = DefaultEmbeddingEndpoint
	}

	if isEndpointTemporarilyUnavailable(endpoint) {
		return nil, fmt.Errorf("embedding endpoint temporarily unavailable: %s", endpoint)
	}

	results := make([][]float64, len(texts))
	client := &http.Client{Timeout: 10 * time.Second}

	for i, text := range texts {
		embedding, err := generateOllamaEmbeddingWithChunkFallback(client, endpoint, text)
		if err != nil {
			return nil, err
		}
		results[i] = embedding
	}
	return results, nil
}

func generateOllamaEmbeddingWithChunkFallback(client *http.Client, endpoint, text string) ([]float64, error) {
	embedding, statusCode, responseBody, err := requestOllamaEmbedding(client, endpoint, text)
	if err == nil {
		return embedding, nil
	}

	if shouldChunkForContextLimit(statusCode, responseBody) {
		chunks := chunkTextByBytes(text, 2800)
		if len(chunks) <= 1 {
			return nil, err
		}
		parts := make([][]float64, 0, len(chunks))
		for _, chunk := range chunks {
			part, _, _, chunkErr := requestOllamaEmbedding(client, endpoint, chunk)
			if chunkErr != nil {
				markEndpointUnavailable(endpoint)
				return nil, chunkErr
			}
			parts = append(parts, part)
		}
		return averageEmbeddings(parts), nil
	}

	if statusCode > 0 || err != nil {
		markEndpointUnavailable(endpoint)
	}
	return nil, err
}

func requestOllamaEmbedding(client *http.Client, endpoint, text string) ([]float64, int, string, error) {
	body, err := json.Marshal(ollamaEmbeddingRequest{
		Model:  "nomic-embed-text",
		Prompt: text,
	})
	if err != nil {
		return nil, 0, "", err
	}
	url := strings.TrimRight(endpoint, "/") + "/api/embeddings"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, 0, "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, "", err
	}
	data, readErr := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if readErr != nil {
		return nil, resp.StatusCode, "", readErr
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, resp.StatusCode, strings.TrimSpace(string(data)), fmt.Errorf("embedding endpoint returned status %d", resp.StatusCode)
	}

	var parsed ollamaEmbeddingResponse
	if err := json.Unmarshal(data, &parsed); err != nil {
		return nil, resp.StatusCode, "", err
	}
	return parsed.Embedding, resp.StatusCode, "", nil
}

func shouldChunkForContextLimit(statusCode int, responseBody string) bool {
	if statusCode < http.StatusBadRequest {
		return false
	}
	msg := strings.ToLower(strings.TrimSpace(responseBody))
	return strings.Contains(msg, "context length") || strings.Contains(msg, "input length exceeds")
}

func chunkTextByBytes(text string, maxChunkBytes int) []string {
	if maxChunkBytes <= 0 || len(text) <= maxChunkBytes {
		return []string{text}
	}
	parts := strings.Fields(text)
	if len(parts) == 0 {
		return []string{text}
	}
	chunks := make([]string, 0, len(parts))
	var builder strings.Builder
	for _, part := range parts {
		next := part
		if builder.Len() > 0 {
			next = " " + part
		}
		if builder.Len()+len(next) > maxChunkBytes && builder.Len() > 0 {
			chunks = append(chunks, builder.String())
			builder.Reset()
			builder.WriteString(part)
			continue
		}
		builder.WriteString(next)
	}
	if builder.Len() > 0 {
		chunks = append(chunks, builder.String())
	}
	if len(chunks) == 0 {
		return []string{text}
	}
	return chunks
}

func averageEmbeddings(vectors [][]float64) []float64 {
	if len(vectors) == 0 {
		return nil
	}
	dim := len(vectors[0])
	if dim == 0 {
		return nil
	}
	out := make([]float64, dim)
	count := 0
	for _, vector := range vectors {
		if len(vector) != dim {
			continue
		}
		for i := 0; i < dim; i++ {
			out[i] += vector[i]
		}
		count++
	}
	if count == 0 {
		return nil
	}
	denom := float64(count)
	for i := range out {
		out[i] /= denom
	}
	return out
}

func generateAzureEmbeddings(endpoint string, texts []string) ([][]float64, error) {
	deployment := os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME")
	if deployment == "" {
		deployment = "text-embedding-3-small"
	}
	apiVersion := os.Getenv("AZURE_OPENAI_API_VERSION")
	if apiVersion == "" {
		apiVersion = "2024-05-01-preview"
	}

	url := fmt.Sprintf("%s/openai/deployments/%s/embeddings?api-version=%s",
		strings.TrimRight(endpoint, "/"), deployment, apiVersion)

	body, err := json.Marshal(azureEmbeddingRequest{Input: texts})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	if apiKey := os.Getenv("AZURE_OPENAI_API_KEY"); apiKey != "" {
		req.Header.Set("api-key", apiKey)
	} else {
		token, err := getAzureToken()
		if err != nil {
			return nil, fmt.Errorf("azure auth failed: %w", err)
		}
		req.Header.Set("Authorization", "Bearer "+token)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("azure embedding returned status %d", resp.StatusCode)
	}

	var parsed azureEmbeddingResponse
	if err := json.NewDecoder(resp.Body).Decode(&parsed); err != nil {
		return nil, err
	}

	if len(parsed.Data) == 0 {
		return nil, errors.New("azure response missing data")
	}

	// Azure might return out of order, though usually it's in order
	results := make([][]float64, len(texts))
	for _, item := range parsed.Data {
		if item.Index < len(results) {
			results[item.Index] = item.Embedding
		}
	}
	return results, nil
}

func getAzureToken() (string, error) {
	azureTokenCache.mu.Lock()
	if time.Now().Before(azureTokenCache.expiry) {
		token := azureTokenCache.token
		azureTokenCache.mu.Unlock()
		return token, nil
	}
	azureTokenCache.mu.Unlock()

	tenantID := os.Getenv("AZURE_TENANT_ID")
	clientID := os.Getenv("AZURE_CLIENT_ID")
	clientSecret := os.Getenv("AZURE_CLIENT_SECRET")
	if tenantID == "" || clientID == "" || clientSecret == "" {
		return "", errors.New("missing AZURE_TENANT_ID, AZURE_CLIENT_ID, or AZURE_CLIENT_SECRET")
	}

	url := fmt.Sprintf("https://login.microsoftonline.com/%s/oauth2/v2.0/token", tenantID)
	data := fmt.Sprintf("grant_type=client_credentials&client_id=%s&client_secret=%s&scope=https://cognitiveservices.azure.com/.default",
		clientID, clientSecret)

	req, err := http.NewRequest(http.MethodPost, url, strings.NewReader(data))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var res struct {
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&res); err != nil {
		return "", err
	}

	azureTokenCache.mu.Lock()
	azureTokenCache.token = res.AccessToken
	azureTokenCache.expiry = time.Now().Add(time.Duration(res.ExpiresIn-60) * time.Second)
	azureTokenCache.mu.Unlock()

	return res.AccessToken, nil
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
