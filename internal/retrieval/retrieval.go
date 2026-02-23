package retrieval

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"athenamind/internal/governance"
	"athenamind/internal/index"
	"athenamind/internal/types"
)

type candidate struct {
	Entry  types.IndexEntry
	Meta   types.MetadataFile
	Body   string
	Score  float64
	Reason string
}

func Retrieve(root, query, domain string) (types.RetrieveResult, error) {
	result, _, err := RetrieveWithEmbeddingEndpoint(root, query, domain, DefaultEmbeddingEndpoint)
	return result, err
}

func RetrieveWithEmbeddingEndpoint(root, query, domain, embeddingEndpoint string) (types.RetrieveResult, string, error) {
	startedAt := time.Now()
	if strings.TrimSpace(query) == "" {
		return types.RetrieveResult{}, "", errors.New("--query is required")
	}

	idx, err := index.LoadIndex(root)
	if err != nil {
		return types.RetrieveResult{}, "", err
	}
	if len(idx.Entries) == 0 {
		return types.RetrieveResult{}, "", errors.New("memory index has no entries")
	}

	candidates, err := loadCandidates(root, idx.Entries, domain)
	if err != nil {
		return types.RetrieveResult{}, "", err
	}
	if len(candidates) == 0 {
		return types.RetrieveResult{}, "", errors.New("no candidates found for query/domain")
	}

	q := strings.ToLower(strings.TrimSpace(query))
	warning := ""
	embeddingScoresApplied := false
	queryEmbedding, embedErr := GenerateEmbedding(embeddingEndpoint, q)
	if embedErr != nil {
		warning = fmt.Sprintf("embedding unavailable; using token-overlap scoring: %v", embedErr)
	}
	ids := make([]string, 0, len(candidates))
	for _, c := range candidates {
		ids = append(ids, c.Entry.ID)
	}
	embeddings, embLoadErr := index.GetEmbeddings(root, ids)
	if embLoadErr != nil {
		warning = fmt.Sprintf("embedding store unavailable; using token-overlap scoring: %v", embLoadErr)
	}

	for i := range candidates {
		candidates[i].Score = semanticScore(q, candidates[i])
		if len(queryEmbedding) > 0 {
			if vec, ok := embeddings[candidates[i].Entry.ID]; ok && len(vec) > 0 {
				candidates[i].Score = cosineSimilarity(queryEmbedding, vec)
				embeddingScoresApplied = true
			}
		}
	}
	if len(queryEmbedding) > 0 && !embeddingScoresApplied && warning == "" {
		warning = "embedding unavailable for candidate entries; using token-overlap scoring"
	}
	sort.SliceStable(candidates, func(i, j int) bool {
		if candidates[i].Score == candidates[j].Score {
			return candidates[i].Entry.ID < candidates[j].Entry.ID
		}
		return candidates[i].Score > candidates[j].Score
	})

	top := candidates[0]
	second := 0.0
	if len(candidates) > 1 {
		second = candidates[1].Score
	}

	if governance.IsLatencyDegraded(time.Since(startedAt).Milliseconds()) {
		sort.SliceStable(candidates, func(i, j int) bool {
			return candidates[i].Entry.Path < candidates[j].Entry.Path
		})
		chosen := candidates[0]
		return types.RetrieveResult{
			SelectedID:    chosen.Entry.ID,
			SelectionMode: "fallback_path_priority",
			SourcePath:    chosen.Entry.Path,
			Confidence:    chosen.Score,
			Reason:        "latency degradation policy forced deterministic fallback",
		}, warning, nil
	}

	if IsSemanticConfident(top.Score, second) {
		mode := "semantic"
		if embeddingScoresApplied {
			mode = "embedding_semantic"
		}
		return types.RetrieveResult{
			SelectedID:    top.Entry.ID,
			SelectionMode: mode,
			SourcePath:    top.Entry.Path,
			Confidence:    top.Score,
			Reason:        "semantic confidence gate passed",
		}, warning, nil
	}

	for _, c := range candidates {
		if strings.EqualFold(c.Entry.ID, q) {
			return types.RetrieveResult{
				SelectedID:    c.Entry.ID,
				SelectionMode: "fallback_exact_key",
				SourcePath:    c.Entry.Path,
				Confidence:    c.Score,
				Reason:        "semantic confidence gate failed; exact-key fallback matched",
			}, warning, nil
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Entry.Path < candidates[j].Entry.Path
	})
	chosen := candidates[0]
	return types.RetrieveResult{
		SelectedID:    chosen.Entry.ID,
		SelectionMode: "fallback_path_priority",
		SourcePath:    chosen.Entry.Path,
		Confidence:    chosen.Score,
		Reason:        "semantic confidence gate failed; deterministic path-priority fallback used",
	}, warning, nil
}

func loadCandidates(root string, entries []types.IndexEntry, domain string) ([]candidate, error) {
	candidates := make([]candidate, 0, len(entries))
	for _, e := range entries {
		if domain != "" && e.Domain != domain {
			continue
		}
		if e.Status != "approved" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(root, filepath.FromSlash(e.Path)))
		if err != nil {
			return nil, err
		}
		meta, err := loadMetadata(root, e)
		if err != nil {
			return nil, err
		}
		candidates = append(candidates, candidate{Entry: e, Meta: meta, Body: string(data)})
	}
	return candidates, nil
}

func loadMetadata(root string, entry types.IndexEntry) (types.MetadataFile, error) {
	path := filepath.Join(root, filepath.FromSlash(entry.MetadataPath))
	data, err := os.ReadFile(path)
	if err != nil {
		return types.MetadataFile{}, err
	}
	var meta types.MetadataFile
	if err := json.Unmarshal(data, &meta); err != nil {
		return types.MetadataFile{}, fmt.Errorf("ERR_SCHEMA_VALIDATION: cannot parse metadata %s: %w", path, err)
	}
	if strings.TrimSpace(meta.SchemaVersion) == "" {
		return types.MetadataFile{}, errors.New("ERR_SCHEMA_VERSION_INVALID: metadata schema_version is required")
	}
	if err := index.ValidateSchemaVersion(meta.SchemaVersion); err != nil {
		return types.MetadataFile{}, err
	}
	if err := validateMetadata(meta, entry, path); err != nil {
		return types.MetadataFile{}, err
	}
	return meta, nil
}

func validateMetadata(meta types.MetadataFile, entry types.IndexEntry, path string) error {
	if strings.TrimSpace(meta.ID) == "" || strings.TrimSpace(meta.Title) == "" || strings.TrimSpace(meta.Status) == "" || strings.TrimSpace(meta.UpdatedAt) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s missing required fields", path)
	}
	if meta.ID != entry.ID {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata id %s does not match entry id %s", meta.ID, entry.ID)
	}
	if !index.IsValidStatus(meta.Status) {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s has invalid status", path)
	}
	if _, err := time.Parse(time.RFC3339, meta.UpdatedAt); err != nil {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s has invalid updated_at", path)
	}
	if strings.TrimSpace(meta.Review.ReviewedBy) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s missing review.reviewed_by", path)
	}
	if meta.Review.Decision != "approved" && meta.Review.Decision != "rejected" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s review.decision must be approved or rejected", path)
	}
	if meta.Review.Decision == "approved" && strings.TrimSpace(meta.Review.ReviewedAt) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s approved records must include review.reviewed_at", path)
	}
	if strings.TrimSpace(meta.Review.ReviewedAt) != "" {
		if _, err := time.Parse(time.RFC3339, meta.Review.ReviewedAt); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s has invalid review.reviewed_at", path)
		}
	}
	return nil
}

func semanticScore(query string, c candidate) float64 {
	qTokens := tokenSet(query)
	if len(qTokens) == 0 {
		return 0
	}

	hay := strings.ToLower(strings.Join([]string{c.Entry.Title, c.Entry.ID, c.Entry.Domain, c.Meta.Title, c.Body}, " "))
	hits := 0
	for tok := range qTokens {
		if strings.Contains(hay, tok) {
			hits++
		}
	}
	return float64(hits) / float64(len(qTokens))
}

func tokenSet(s string) map[string]struct{} {
	clean := strings.NewReplacer(".", " ", ",", " ", ":", " ", ";", " ", "/", " ", "-", " ", "_", " ").Replace(strings.ToLower(s))
	parts := strings.Fields(clean)
	out := make(map[string]struct{}, len(parts))
	for _, p := range parts {
		if len(p) > 1 {
			out[p] = struct{}{}
		}
	}
	return out
}

func IsSemanticConfident(top, second float64) bool {
	const minConfidence = 0.34
	const minMargin = 0.15
	if top < minConfidence {
		return false
	}
	if top-second < minMargin {
		return false
	}
	return true
}
