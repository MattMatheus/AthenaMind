package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

func retrieve(root, query, domain string) (retrieveResult, error) {
	startedAt := time.Now()
	if strings.TrimSpace(query) == "" {
		return retrieveResult{}, errors.New("--query is required")
	}

	idx, err := loadIndex(root)
	if err != nil {
		return retrieveResult{}, err
	}
	if len(idx.Entries) == 0 {
		return retrieveResult{}, errors.New("memory index has no entries")
	}

	candidates, err := loadCandidates(root, idx.Entries, domain)
	if err != nil {
		return retrieveResult{}, err
	}
	if len(candidates) == 0 {
		return retrieveResult{}, errors.New("no candidates found for query/domain")
	}

	q := strings.ToLower(strings.TrimSpace(query))
	for i := range candidates {
		candidates[i].Score = semanticScore(q, candidates[i])
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

	if isLatencyDegraded(time.Since(startedAt).Milliseconds()) {
		sort.SliceStable(candidates, func(i, j int) bool {
			return candidates[i].Entry.Path < candidates[j].Entry.Path
		})
		chosen := candidates[0]
		return retrieveResult{
			SelectedID:    chosen.Entry.ID,
			SelectionMode: "fallback_path_priority",
			SourcePath:    chosen.Entry.Path,
			Confidence:    chosen.Score,
			Reason:        "latency degradation policy forced deterministic fallback",
		}, nil
	}

	if isSemanticConfident(top.Score, second) {
		return retrieveResult{
			SelectedID:    top.Entry.ID,
			SelectionMode: "semantic",
			SourcePath:    top.Entry.Path,
			Confidence:    top.Score,
			Reason:        "semantic confidence gate passed",
		}, nil
	}

	for _, c := range candidates {
		if strings.EqualFold(c.Entry.ID, q) {
			return retrieveResult{
				SelectedID:    c.Entry.ID,
				SelectionMode: "fallback_exact_key",
				SourcePath:    c.Entry.Path,
				Confidence:    c.Score,
				Reason:        "semantic confidence gate failed; exact-key fallback matched",
			}, nil
		}
	}

	sort.SliceStable(candidates, func(i, j int) bool {
		return candidates[i].Entry.Path < candidates[j].Entry.Path
	})
	chosen := candidates[0]
	return retrieveResult{
		SelectedID:    chosen.Entry.ID,
		SelectionMode: "fallback_path_priority",
		SourcePath:    chosen.Entry.Path,
		Confidence:    chosen.Score,
		Reason:        "semantic confidence gate failed; deterministic path-priority fallback used",
	}, nil
}

func readGatewayHandler(root string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req apiRetrieveRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "ERR_API_INPUT_INVALID: invalid JSON payload", http.StatusBadRequest)
			return
		}
		if strings.TrimSpace(req.Query) == "" || strings.TrimSpace(req.SessionID) == "" {
			http.Error(w, "ERR_API_INPUT_INVALID: query and session_id are required", http.StatusBadRequest)
			return
		}

		result, err := retrieve(root, req.Query, req.Domain)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		traceID := fmt.Sprintf("%s:%d", strings.TrimSpace(req.SessionID), time.Now().UTC().UnixNano())
		resp := apiRetrieveResponse{
			SelectedID:    result.SelectedID,
			SelectionMode: result.SelectionMode,
			SourcePath:    result.SourcePath,
			Confidence:    result.Confidence,
			Reason:        result.Reason,
			TraceID:       traceID,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	})
}

func apiRetrieveWithFallback(root, gatewayURL string, req apiRetrieveRequest, traceID string, client *http.Client) (apiRetrieveResponse, error) {
	if strings.TrimSpace(gatewayURL) != "" {
		resp, err := gatewayRetrieve(gatewayURL, req, client)
		if err == nil {
			local, localErr := retrieve(root, req.Query, req.Domain)
			if localErr != nil {
				return apiRetrieveResponse{}, localErr
			}
			if resp.SelectedID != local.SelectedID || resp.SelectionMode != local.SelectionMode || resp.SourcePath != local.SourcePath {
				return apiRetrieveResponse{}, errors.New("ERR_API_CLI_PARITY_MISMATCH: gateway response diverged from CLI retrieve contract")
			}
			if strings.TrimSpace(resp.TraceID) == "" {
				resp.TraceID = traceID
			}
			resp.GatewayEndpoint = gatewayURL
			return resp, nil
		}
		local, localErr := retrieve(root, req.Query, req.Domain)
		if localErr != nil {
			return apiRetrieveResponse{}, localErr
		}
		return apiRetrieveResponse{
			SelectedID:      local.SelectedID,
			SelectionMode:   local.SelectionMode,
			SourcePath:      local.SourcePath,
			Confidence:      local.Confidence,
			Reason:          local.Reason,
			TraceID:         traceID,
			FallbackUsed:    true,
			FallbackCode:    "ERR_API_WRAPPER_UNAVAILABLE",
			FallbackReason:  err.Error(),
			GatewayEndpoint: gatewayURL,
		}, nil
	}

	local, err := retrieve(root, req.Query, req.Domain)
	if err != nil {
		return apiRetrieveResponse{}, err
	}
	return apiRetrieveResponse{
		SelectedID:    local.SelectedID,
		SelectionMode: local.SelectionMode,
		SourcePath:    local.SourcePath,
		Confidence:    local.Confidence,
		Reason:        local.Reason,
		TraceID:       traceID,
	}, nil
}

func gatewayRetrieve(gatewayURL string, req apiRetrieveRequest, client *http.Client) (apiRetrieveResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return apiRetrieveResponse{}, err
	}
	endpoint := strings.TrimRight(gatewayURL, "/") + "/memory/retrieve"
	httpReq, err := http.NewRequest(http.MethodPost, endpoint, bytes.NewReader(body))
	if err != nil {
		return apiRetrieveResponse{}, err
	}
	httpReq.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(httpReq)
	if err != nil {
		return apiRetrieveResponse{}, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		data, _ := io.ReadAll(resp.Body)
		return apiRetrieveResponse{}, fmt.Errorf("gateway status %d: %s", resp.StatusCode, strings.TrimSpace(string(data)))
	}

	var out apiRetrieveResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return apiRetrieveResponse{}, err
	}
	if out.SelectedID == "" || out.SelectionMode == "" || out.SourcePath == "" {
		return apiRetrieveResponse{}, errors.New("ERR_API_INPUT_INVALID: gateway response missing required retrieval fields")
	}
	return out, nil
}

func loadCandidates(root string, entries []indexEntry, domain string) ([]candidate, error) {
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

func loadMetadata(root string, entry indexEntry) (metadataFile, error) {
	path := filepath.Join(root, filepath.FromSlash(entry.MetadataPath))
	data, err := os.ReadFile(path)
	if err != nil {
		return metadataFile{}, err
	}
	var meta metadataFile
	if err := json.Unmarshal(data, &meta); err != nil {
		return metadataFile{}, fmt.Errorf("ERR_SCHEMA_VALIDATION: cannot parse metadata %s: %w", path, err)
	}
	if strings.TrimSpace(meta.SchemaVersion) == "" {
		return metadataFile{}, errors.New("ERR_SCHEMA_VERSION_INVALID: metadata schema_version is required")
	}
	if err := validateSchemaVersion(meta.SchemaVersion); err != nil {
		return metadataFile{}, err
	}
	if err := validateMetadata(meta, entry, path); err != nil {
		return metadataFile{}, err
	}
	return meta, nil
}

func validateMetadata(meta metadataFile, entry indexEntry, path string) error {
	if strings.TrimSpace(meta.ID) == "" || strings.TrimSpace(meta.Title) == "" || strings.TrimSpace(meta.Status) == "" || strings.TrimSpace(meta.UpdatedAt) == "" {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata %s missing required fields", path)
	}
	if meta.ID != entry.ID {
		return fmt.Errorf("ERR_SCHEMA_VALIDATION: metadata id %s does not match entry id %s", meta.ID, entry.ID)
	}
	if !isValidStatus(meta.Status) {
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

func isSemanticConfident(top, second float64) bool {
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
