package retrieval

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const DefaultEmbeddingEndpoint = "http://localhost:11434"
const embeddingChunkMaxBytes = 2800

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
	if strings.TrimSpace(os.Getenv("AZURE_OPENAI_ENDPOINT")) != "" {
		deployment := strings.TrimSpace(os.Getenv("AZURE_OPENAI_DEPLOYMENT_NAME"))
		if deployment == "" {
			deployment = "text-embedding-3-small"
		}
		return EmbeddingProfile{
			Provider: "azure_openai",
			ModelID:  deployment,
		}
	}
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

	results := make([][]float64, 0, len(texts))
	for _, text := range texts {
		chunks := chunkTextBySemanticBoundaries(text, embeddingChunkMaxBytes)
		// Try Azure first if configured.
		if azureEndpoint := os.Getenv("AZURE_OPENAI_ENDPOINT"); azureEndpoint != "" {
			chunkVectors, err := generateAzureEmbeddings(azureEndpoint, chunks)
			if err == nil {
				avg, avgErr := averageEmbeddings(chunkVectors)
				if avgErr != nil {
					return nil, avgErr
				}
				results = append(results, avg)
				continue
			}
			// Fall through to Ollama if Azure fails.
		}

		chunkVectors, err := generateOllamaEmbeddings(endpoint, chunks)
		if err != nil {
			return nil, err
		}
		avg, avgErr := averageEmbeddings(chunkVectors)
		if avgErr != nil {
			return nil, avgErr
		}
		results = append(results, avg)
	}
	return results, nil
}

func generateOllamaEmbeddings(endpoint string, texts []string) ([][]float64, error) {
	results := make([][]float64, len(texts))
	client := &http.Client{Timeout: 10 * time.Second}
	url := strings.TrimRight(endpoint, "/") + "/api/embeddings"

	for i, text := range texts {
		body, err := json.Marshal(ollamaEmbeddingRequest{
			Model:  "nomic-embed-text",
			Prompt: text,
		})
		if err != nil {
			return nil, err
		}
		req, err := http.NewRequest(http.MethodPost, url, bytes.NewReader(body))
		if err != nil {
			return nil, err
		}
		req.Header.Set("Content-Type", "application/json")

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
		results[i] = parsed.Embedding
	}
	return results, nil
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

func averageEmbeddings(vectors [][]float64) ([]float64, error) {
	if len(vectors) == 0 {
		return nil, errors.New("no embedding returned")
	}
	dim := len(vectors[0])
	if dim == 0 {
		return nil, errors.New("no embedding returned")
	}
	avg := make([]float64, dim)
	for _, vec := range vectors {
		if len(vec) != dim {
			return nil, errors.New("inconsistent embedding dimensions")
		}
		for i := 0; i < dim; i++ {
			avg[i] += vec[i]
		}
	}
	n := float64(len(vectors))
	for i := 0; i < dim; i++ {
		avg[i] /= n
	}
	return avg, nil
}

func chunkTextBySemanticBoundaries(text string, maxChunkBytes int) []string {
	if maxChunkBytes <= 0 {
		maxChunkBytes = embeddingChunkMaxBytes
	}
	normalized := strings.ReplaceAll(strings.ReplaceAll(text, "\r\n", "\n"), "\r", "\n")
	if strings.TrimSpace(normalized) == "" {
		return []string{text}
	}

	units := make([]string, 0, 32)
	lines := strings.Split(normalized, "\n")
	current := make([]string, 0, 16)
	flushCurrent := func() {
		if len(current) == 0 {
			return
		}
		unit := strings.Trim(strings.Join(current, "\n"), "\n")
		current = current[:0]
		if strings.TrimSpace(unit) != "" {
			units = append(units, unit)
		}
	}

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if isSemanticBoundaryLine(trimmed) {
			flushCurrent()
			current = append(current, line)
			continue
		}
		if trimmed == "" {
			flushCurrent()
			continue
		}
		current = append(current, line)
	}
	flushCurrent()
	if len(units) == 0 {
		units = append(units, normalized)
	}

	chunks := make([]string, 0, len(units))
	currentChunk := ""
	for _, unit := range units {
		if len([]byte(unit)) > maxChunkBytes {
			if strings.TrimSpace(currentChunk) != "" {
				chunks = append(chunks, currentChunk)
				currentChunk = ""
			}
			chunks = append(chunks, splitTextByBytes(unit, maxChunkBytes)...)
			continue
		}

		candidate := unit
		if currentChunk != "" {
			candidate = currentChunk + "\n\n" + unit
		}
		if len([]byte(candidate)) <= maxChunkBytes {
			currentChunk = candidate
			continue
		}
		if strings.TrimSpace(currentChunk) != "" {
			chunks = append(chunks, currentChunk)
		}
		currentChunk = unit
	}
	if strings.TrimSpace(currentChunk) != "" {
		chunks = append(chunks, currentChunk)
	}
	if len(chunks) == 0 {
		return []string{text}
	}
	return chunks
}

func isSemanticBoundaryLine(trimmed string) bool {
	if trimmed == "" {
		return true
	}
	if strings.HasPrefix(trimmed, "# ") ||
		strings.HasPrefix(trimmed, "## ") ||
		strings.HasPrefix(trimmed, "### ") ||
		strings.HasPrefix(trimmed, "#### ") ||
		strings.HasPrefix(trimmed, "##### ") ||
		strings.HasPrefix(trimmed, "###### ") {
		return true
	}
	if strings.HasPrefix(trimmed, "func ") ||
		strings.HasPrefix(trimmed, "type ") ||
		strings.HasPrefix(trimmed, "class ") ||
		strings.HasPrefix(trimmed, "def ") ||
		strings.HasPrefix(trimmed, "function ") {
		return true
	}
	return false
}

func splitTextByBytes(text string, maxChunkBytes int) []string {
	if maxChunkBytes <= 0 || len([]byte(text)) <= maxChunkBytes {
		return []string{text}
	}
	chunks := make([]string, 0, (len(text)/maxChunkBytes)+1)
	rest := text
	for len([]byte(rest)) > maxChunkBytes {
		cut := maxPrefixByBytes(rest, maxChunkBytes)
		if cut <= 0 {
			break
		}
		chunk := strings.Trim(rest[:cut], "\n")
		if strings.TrimSpace(chunk) != "" {
			chunks = append(chunks, chunk)
		}
		rest = rest[cut:]
	}
	rest = strings.Trim(rest, "\n")
	if strings.TrimSpace(rest) != "" {
		chunks = append(chunks, rest)
	}
	if len(chunks) == 0 {
		return []string{text}
	}
	return chunks
}

func maxPrefixByBytes(text string, maxBytes int) int {
	if maxBytes <= 0 {
		return 0
	}
	if len([]byte(text)) <= maxBytes {
		return len(text)
	}
	total := 0
	last := 0
	for i, r := range text {
		size := utf8.RuneLen(r)
		if total+size > maxBytes {
			break
		}
		total += size
		last = i + size
	}
	if last == 0 && len(text) > 0 {
		_, size := utf8.DecodeRuneInString(text)
		if size > 0 {
			return size
		}
	}
	return last
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
