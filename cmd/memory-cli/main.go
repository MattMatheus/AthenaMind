package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

const (
	supportedMajor = 1
	supportedMinor = 0
	defaultSchema  = "1.0"

	defaultEvaluationQuerySetPath = "cmd/memory-cli/testdata/eval-query-set-v1.json"
	defaultEvaluationCorpusID     = "memory-corpus-v1"
	defaultEvaluationQuerySetID   = "query-set-v1"
	defaultEvaluationConfigID     = "config-v1-confidence-0.34-margin-0.15"
)

type indexFile struct {
	SchemaVersion string       `json:"schema_version"`
	UpdatedAt     string       `json:"updated_at"`
	Entries       []indexEntry `json:"entries"`
}

type indexEntry struct {
	ID           string `json:"id"`
	Type         string `json:"type"`
	Domain       string `json:"domain"`
	Path         string `json:"path"`
	MetadataPath string `json:"metadata_path"`
	Status       string `json:"status"`
	UpdatedAt    string `json:"updated_at"`
	Title        string `json:"title"`
}

type metadataFile struct {
	SchemaVersion string     `json:"schema_version"`
	ID            string     `json:"id"`
	Title         string     `json:"title"`
	Status        string     `json:"status"`
	UpdatedAt     string     `json:"updated_at"`
	Review        reviewMeta `json:"review"`
}

type reviewMeta struct {
	ReviewedBy   string `json:"reviewed_by"`
	ReviewedAt   string `json:"reviewed_at"`
	Decision     string `json:"decision"`
	DecisionNote string `json:"decision_notes"`
}

type retrieveResult struct {
	SelectedID    string  `json:"selected_id"`
	SelectionMode string  `json:"selection_mode"`
	SourcePath    string  `json:"source_path"`
	Confidence    float64 `json:"confidence"`
	Reason        string  `json:"reason"`
}

type apiRetrieveRequest struct {
	Query     string `json:"query"`
	Domain    string `json:"domain,omitempty"`
	SessionID string `json:"session_id"`
}

type apiRetrieveResponse struct {
	SelectedID      string  `json:"selected_id"`
	SelectionMode   string  `json:"selection_mode"`
	SourcePath      string  `json:"source_path"`
	Confidence      float64 `json:"confidence"`
	Reason          string  `json:"reason"`
	TraceID         string  `json:"trace_id"`
	FallbackUsed    bool    `json:"fallback_used"`
	FallbackCode    string  `json:"fallback_code,omitempty"`
	FallbackReason  string  `json:"fallback_reason,omitempty"`
	GatewayEndpoint string  `json:"gateway_endpoint,omitempty"`
}

type candidate struct {
	Entry  indexEntry
	Meta   metadataFile
	Body   string
	Score  float64
	Reason string
}

type mutationAuditRecord struct {
	SchemaVersion string   `json:"schema_version"`
	ID            string   `json:"id"`
	Stage         string   `json:"stage"`
	Decision      string   `json:"decision"`
	ReviewedBy    string   `json:"reviewed_by"`
	ReviewedAt    string   `json:"reviewed_at"`
	DecisionNotes string   `json:"decision_notes"`
	Reason        string   `json:"reason"`
	Risk          string   `json:"risk"`
	ReworkNotes   string   `json:"rework_notes,omitempty"`
	ReReviewedBy  string   `json:"re_reviewed_by,omitempty"`
	ChangedFiles  []string `json:"changed_files"`
	Applied       bool     `json:"applied"`
}

type evaluationQuery struct {
	Query      string `json:"query"`
	Domain     string `json:"domain,omitempty"`
	ExpectedID string `json:"expected_id"`
}

type queryMiss struct {
	Query      string `json:"query"`
	ExpectedID string `json:"expected_id"`
	ActualID   string `json:"actual_id"`
	Mode       string `json:"mode"`
}

type deterministicReplay struct {
	Query      string `json:"query"`
	Mode       string `json:"mode"`
	SelectedID string `json:"selected_id"`
	SourcePath string `json:"source_path"`
	StableRuns int    `json:"stable_runs"`
}

type evaluationMetric struct {
	Numerator   int     `json:"numerator"`
	Denominator int     `json:"denominator"`
	Rate        float64 `json:"rate"`
}

type evaluationReport struct {
	CorpusID                string                `json:"corpus_id"`
	QuerySetID              string                `json:"query_set_id"`
	ConfigID                string                `json:"config_id"`
	Status                  string                `json:"status"`
	Top1UsefulRate          evaluationMetric      `json:"top1_useful_rate"`
	FallbackDeterminism     evaluationMetric      `json:"fallback_determinism"`
	SelectionModeReporting  evaluationMetric      `json:"selection_mode_reporting"`
	SourceTraceCompleteness evaluationMetric      `json:"source_trace_completeness"`
	FailingQueries          []queryMiss           `json:"failing_queries"`
	DeterministicReplay     []deterministicReplay `json:"deterministic_replay_proof"`
}

type writePolicyInput struct {
	Stage        string
	Reviewer     string
	ApprovedFlag bool
	Decision     string
	Notes        string
	Reason       string
	Risk         string
	ReworkNotes  string
	ReReviewedBy string
}

type writePolicyDecision struct {
	Decision     string
	Reviewer     string
	Notes        string
	Reason       string
	Risk         string
	ReworkNotes  string
	ReReviewedBy string
}

func main() {
	if len(os.Args) < 2 {
		exitErr(errors.New("usage: memory-cli <write|retrieve|evaluate> [flags]"))
	}

	var err error
	switch os.Args[1] {
	case "write":
		err = runWrite(os.Args[2:])
	case "retrieve":
		err = runRetrieve(os.Args[2:])
	case "serve-read-gateway":
		err = runServeReadGateway(os.Args[2:])
	case "api-retrieve":
		err = runAPIRetrieve(os.Args[2:])
	case "evaluate":
		err = runEvaluate(os.Args[2:])
	default:
		err = fmt.Errorf("unknown command: %s", os.Args[1])
	}

	if err != nil {
		exitErr(err)
	}
}

func runWrite(args []string) error {
	fs := flag.NewFlagSet("write", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	id := fs.String("id", "", "entry id")
	title := fs.String("title", "", "entry title")
	typeValue := fs.String("type", "", "entry type: prompt|instruction")
	domain := fs.String("domain", "", "entry domain")
	body := fs.String("body", "", "entry body")
	bodyFile := fs.String("body-file", "", "path to markdown body")
	stage := fs.String("stage", "", "workflow stage: planning|architect|pm")
	reviewer := fs.String("reviewer", "", "reviewer identity")
	approved := fs.Bool("approved", false, "legacy flag: equivalent to --decision=approved")
	decision := fs.String("decision", "", "review decision: approved|rejected")
	notes := fs.String("notes", "", "decision notes")
	reason := fs.String("reason", "", "reason for change")
	risk := fs.String("risk", "", "risk and mitigation note")
	reworkNotes := fs.String("rework-notes", "", "required when --decision=rejected")
	reReviewedBy := fs.String("re-reviewed-by", "", "required when --decision=rejected")

	if err := fs.Parse(args); err != nil {
		return err
	}

	policy, err := enforceWritePolicy(writePolicyInput{
		Stage:        *stage,
		Reviewer:     *reviewer,
		ApprovedFlag: *approved,
		Decision:     *decision,
		Notes:        *notes,
		Reason:       *reason,
		Risk:         *risk,
		ReworkNotes:  *reworkNotes,
		ReReviewedBy: *reReviewedBy,
	})
	if err != nil {
		return err
	}

	if *id == "" || *title == "" || *typeValue == "" || *domain == "" {
		return errors.New("--id --title --type --domain are required")
	}
	if *typeValue != "prompt" && *typeValue != "instruction" {
		return errors.New("--type must be prompt or instruction")
	}

	entryBody := strings.TrimSpace(*body)
	if *bodyFile != "" {
		data, err := os.ReadFile(*bodyFile)
		if err != nil {
			return err
		}
		entryBody = strings.TrimSpace(string(data))
	}
	if entryBody == "" {
		return errors.New("entry body is required via --body or --body-file")
	}

	idx, err := loadIndex(*root)
	if err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	dirName := "prompts"
	if *typeValue == "instruction" {
		dirName = "instructions"
	}

	relContentPath := filepath.ToSlash(filepath.Join(dirName, *domain, *id+".md"))
	relMetaPath := filepath.ToSlash(filepath.Join("metadata", *id+".yaml"))
	contentPath := filepath.Join(*root, filepath.FromSlash(relContentPath))
	metaPath := filepath.Join(*root, filepath.FromSlash(relMetaPath))
	auditPath := filepath.Join(*root, "audits", fmt.Sprintf("%s-%d.json", *id, time.Now().UTC().UnixNano()))

	audit := mutationAuditRecord{
		SchemaVersion: defaultSchema,
		ID:            *id,
		Stage:         *stage,
		Decision:      policy.Decision,
		ReviewedBy:    policy.Reviewer,
		ReviewedAt:    now,
		DecisionNotes: policy.Notes,
		Reason:        policy.Reason,
		Risk:          policy.Risk,
		ReworkNotes:   policy.ReworkNotes,
		ReReviewedBy:  policy.ReReviewedBy,
		ChangedFiles: []string{
			relContentPath,
			relMetaPath,
			"index.yaml",
		},
		Applied: policy.Decision == "approved",
	}
	if err := writeJSONAsYAML(auditPath, audit); err != nil {
		return err
	}
	if policy.Decision == "rejected" {
		return errors.New("ERR_MUTATION_REJECTED_PENDING_REWORK: write blocked until rework and re-review are completed")
	}

	if err := os.MkdirAll(filepath.Dir(contentPath), 0o755); err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(metaPath), 0o755); err != nil {
		return err
	}

	markdown := fmt.Sprintf("# %s\n\n%s\n", *title, entryBody)
	if err := os.WriteFile(contentPath, []byte(markdown), 0o644); err != nil {
		return err
	}

	meta := metadataFile{
		SchemaVersion: defaultSchema,
		ID:            *id,
		Title:         *title,
		Status:        "approved",
		UpdatedAt:     now,
		Review: reviewMeta{
			ReviewedBy:   policy.Reviewer,
			ReviewedAt:   now,
			Decision:     "approved",
			DecisionNote: policy.Notes,
		},
	}
	if err := writeJSONAsYAML(metaPath, meta); err != nil {
		return err
	}

	updated := false
	for i := range idx.Entries {
		if idx.Entries[i].ID == *id {
			idx.Entries[i] = indexEntry{
				ID:           *id,
				Type:         *typeValue,
				Domain:       *domain,
				Path:         relContentPath,
				MetadataPath: relMetaPath,
				Status:       "approved",
				UpdatedAt:    now,
				Title:        *title,
			}
			updated = true
			break
		}
	}
	if !updated {
		idx.Entries = append(idx.Entries, indexEntry{
			ID:           *id,
			Type:         *typeValue,
			Domain:       *domain,
			Path:         relContentPath,
			MetadataPath: relMetaPath,
			Status:       "approved",
			UpdatedAt:    now,
			Title:        *title,
		})
	}
	sort.Slice(idx.Entries, func(i, j int) bool { return idx.Entries[i].ID < idx.Entries[j].ID })
	idx.UpdatedAt = now
	if err := writeJSONAsYAML(filepath.Join(*root, "index.yaml"), idx); err != nil {
		return err
	}

	fmt.Printf("wrote entry %s at %s\n", *id, relContentPath)
	return nil
}

func runRetrieve(args []string) error {
	fs := flag.NewFlagSet("retrieve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "", "natural language query")
	domain := fs.String("domain", "", "optional domain filter")
	if err := fs.Parse(args); err != nil {
		return err
	}

	result, err := retrieve(*root, *query, *domain)
	if err != nil {
		return err
	}
	return printResult(result)
}

func runEvaluate(args []string) error {
	fs := flag.NewFlagSet("evaluate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	queryFile := fs.String("query-file", defaultEvaluationQuerySetPath, "path to evaluation query set JSON")
	corpusID := fs.String("corpus-id", defaultEvaluationCorpusID, "pinned corpus snapshot id")
	querySetID := fs.String("query-set-id", defaultEvaluationQuerySetID, "pinned query set id")
	configID := fs.String("config-id", defaultEvaluationConfigID, "retrieval configuration snapshot id")
	if err := fs.Parse(args); err != nil {
		return err
	}

	queries, err := loadEvaluationQueries(*queryFile)
	if err != nil {
		return err
	}

	report, err := evaluateRetrieval(*root, queries, *corpusID, *querySetID, *configID)
	if err != nil {
		return err
	}

	out, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func runServeReadGateway(args []string) error {
	fs := flag.NewFlagSet("serve-read-gateway", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	addr := fs.String("addr", "127.0.0.1:8788", "listen address")
	if err := fs.Parse(args); err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/memory/retrieve", readGatewayHandler(*root))
	return http.ListenAndServe(*addr, mux)
}

func runAPIRetrieve(args []string) error {
	fs := flag.NewFlagSet("api-retrieve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "", "natural language query")
	domain := fs.String("domain", "", "optional domain filter")
	sessionID := fs.String("session-id", "", "session identifier")
	gatewayURL := fs.String("gateway-url", "", "optional API gateway base URL")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*query) == "" {
		return errors.New("--query is required")
	}
	if strings.TrimSpace(*sessionID) == "" {
		return errors.New("--session-id is required")
	}

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	req := apiRetrieveRequest{
		Query:     *query,
		Domain:    *domain,
		SessionID: *sessionID,
	}
	resp, err := apiRetrieveWithFallback(*root, strings.TrimSpace(*gatewayURL), req, traceID, http.DefaultClient)
	if err != nil {
		return err
	}
	return printResult(resp)
}

func retrieve(root, query, domain string) (retrieveResult, error) {
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

func evaluateRetrieval(root string, queries []evaluationQuery, corpusID, querySetID, configID string) (evaluationReport, error) {
	report := evaluationReport{
		CorpusID:            corpusID,
		QuerySetID:          querySetID,
		ConfigID:            configID,
		Status:              "FAIL",
		FailingQueries:      []queryMiss{},
		DeterministicReplay: []deterministicReplay{},
	}
	if len(queries) == 0 {
		return report, errors.New("ERR_EVAL_QUERY_SET_INVALID: query set must contain at least one query")
	}

	total := len(queries)
	top1Useful := 0
	selectionModePresent := 0
	sourceTracePresent := 0
	fallbackChecks := 0
	fallbackStable := 0

	for _, q := range queries {
		result, err := retrieve(root, q.Query, q.Domain)
		if err != nil {
			return report, err
		}
		if result.SelectionMode != "" {
			selectionModePresent++
		}
		if result.SelectedID != "" && result.SourcePath != "" {
			sourceTracePresent++
		}
		if result.SelectionMode == "semantic" && q.ExpectedID != "" && result.SelectedID == q.ExpectedID {
			top1Useful++
		} else if q.ExpectedID != "" && result.SelectedID != q.ExpectedID {
			report.FailingQueries = append(report.FailingQueries, queryMiss{
				Query:      q.Query,
				ExpectedID: q.ExpectedID,
				ActualID:   result.SelectedID,
				Mode:       result.SelectionMode,
			})
		}

		if strings.HasPrefix(result.SelectionMode, "fallback_") {
			fallbackChecks++
			stable := true
			for i := 0; i < 4; i++ {
				again, err := retrieve(root, q.Query, q.Domain)
				if err != nil {
					return report, err
				}
				if again.SelectionMode != result.SelectionMode || again.SelectedID != result.SelectedID || again.SourcePath != result.SourcePath {
					stable = false
					break
				}
			}
			if stable {
				fallbackStable++
			}
			report.DeterministicReplay = append(report.DeterministicReplay, deterministicReplay{
				Query:      q.Query,
				Mode:       result.SelectionMode,
				SelectedID: result.SelectedID,
				SourcePath: result.SourcePath,
				StableRuns: 5,
			})
		}
	}

	report.Top1UsefulRate = metric(top1Useful, total)
	report.SelectionModeReporting = metric(selectionModePresent, total)
	report.SourceTraceCompleteness = metric(sourceTracePresent, total)
	if fallbackChecks == 0 {
		report.FallbackDeterminism = evaluationMetric{Numerator: 1, Denominator: 1, Rate: 1}
	} else {
		report.FallbackDeterminism = metric(fallbackStable, fallbackChecks)
	}

	pass := total >= 50 &&
		report.Top1UsefulRate.Rate >= 0.80 &&
		report.FallbackDeterminism.Rate == 1 &&
		report.SelectionModeReporting.Rate == 1 &&
		report.SourceTraceCompleteness.Rate == 1

	if pass {
		report.Status = "PASS"
		if report.Top1UsefulRate.Rate >= 0.80 && report.Top1UsefulRate.Rate <= 0.82 {
			report.Status = "WATCH"
		}
	}

	return report, nil
}

func metric(numerator, denominator int) evaluationMetric {
	if denominator == 0 {
		return evaluationMetric{Numerator: 0, Denominator: 0, Rate: 0}
	}
	return evaluationMetric{
		Numerator:   numerator,
		Denominator: denominator,
		Rate:        float64(numerator) / float64(denominator),
	}
}

func loadEvaluationQueries(path string) ([]evaluationQuery, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var queries []evaluationQuery
	if err := json.Unmarshal(data, &queries); err != nil {
		return nil, fmt.Errorf("ERR_EVAL_QUERY_SET_INVALID: cannot parse %s: %w", path, err)
	}
	for _, q := range queries {
		if strings.TrimSpace(q.Query) == "" {
			return nil, errors.New("ERR_EVAL_QUERY_SET_INVALID: each query must include query text")
		}
	}
	return queries, nil
}

func printResult(r any) error {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
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

func enforceWritePolicy(in writePolicyInput) (writePolicyDecision, error) {
	if isTrue(os.Getenv("AUTONOMOUS_RUN")) {
		return writePolicyDecision{}, errors.New("ERR_MUTATION_NOT_ALLOWED_DURING_AUTONOMOUS_RUN: writes are blocked during autonomous runs")
	}
	allowed := map[string]struct{}{"planning": {}, "architect": {}, "pm": {}}
	if _, ok := allowed[in.Stage]; !ok {
		return writePolicyDecision{}, errors.New("ERR_MUTATION_STAGE_INVALID: --stage must be planning, architect, or pm")
	}
	if strings.TrimSpace(in.Reviewer) == "" {
		return writePolicyDecision{}, errors.New("ERR_MUTATION_REVIEW_REQUIRED: --reviewer is required")
	}

	decision := strings.TrimSpace(strings.ToLower(in.Decision))
	if decision == "" {
		if in.ApprovedFlag {
			decision = "approved"
		} else {
			return writePolicyDecision{}, errors.New("ERR_MUTATION_REVIEW_REQUIRED: provide --decision=approved|rejected")
		}
	}
	if decision != "approved" && decision != "rejected" {
		return writePolicyDecision{}, errors.New("ERR_MUTATION_REVIEW_REQUIRED: --decision must be approved or rejected")
	}
	if strings.TrimSpace(in.Reason) == "" || strings.TrimSpace(in.Risk) == "" || strings.TrimSpace(in.Notes) == "" {
		return writePolicyDecision{}, errors.New("ERR_MUTATION_EVIDENCE_REQUIRED: --reason --risk and --notes are required")
	}
	if decision == "rejected" {
		if strings.TrimSpace(in.ReworkNotes) == "" || strings.TrimSpace(in.ReReviewedBy) == "" {
			return writePolicyDecision{}, errors.New("ERR_MUTATION_REJECTION_EVIDENCE_REQUIRED: rejected decisions require --rework-notes and --re-reviewed-by")
		}
	}

	return writePolicyDecision{
		Decision:     decision,
		Reviewer:     in.Reviewer,
		Notes:        in.Notes,
		Reason:       in.Reason,
		Risk:         in.Risk,
		ReworkNotes:  in.ReworkNotes,
		ReReviewedBy: in.ReReviewedBy,
	}, nil
}

func isTrue(v string) bool {
	v = strings.ToLower(strings.TrimSpace(v))
	return v == "1" || v == "true" || v == "yes"
}

func loadIndex(root string) (indexFile, error) {
	path := filepath.Join(root, "index.yaml")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			if err := os.MkdirAll(root, 0o755); err != nil {
				return indexFile{}, err
			}
			return indexFile{SchemaVersion: defaultSchema, UpdatedAt: time.Now().UTC().Format(time.RFC3339), Entries: []indexEntry{}}, nil
		}
		return indexFile{}, err
	}

	var idx indexFile
	if err := json.Unmarshal(data, &idx); err != nil {
		return indexFile{}, fmt.Errorf("ERR_SCHEMA_VERSION_INVALID: cannot parse %s: %w", path, err)
	}
	if strings.TrimSpace(idx.SchemaVersion) == "" {
		return indexFile{}, errors.New("ERR_SCHEMA_VERSION_INVALID: index schema_version is required")
	}
	if err := validateSchemaVersion(idx.SchemaVersion); err != nil {
		return indexFile{}, err
	}
	if err := validateIndex(idx, root); err != nil {
		return indexFile{}, err
	}
	return idx, nil
}

func validateIndex(idx indexFile, root string) error {
	if strings.TrimSpace(idx.UpdatedAt) == "" {
		return errors.New("ERR_SCHEMA_VALIDATION: index updated_at is required")
	}
	if _, err := time.Parse(time.RFC3339, idx.UpdatedAt); err != nil {
		return errors.New("ERR_SCHEMA_VALIDATION: index updated_at must be RFC3339")
	}
	seen := map[string]struct{}{}
	for _, e := range idx.Entries {
		if e.ID == "" || e.Type == "" || e.Domain == "" || e.Path == "" || e.MetadataPath == "" || e.Status == "" || e.UpdatedAt == "" {
			return errors.New("ERR_SCHEMA_VALIDATION: index entry missing required fields")
		}
		if _, ok := seen[e.ID]; ok {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: duplicate index entry id %s", e.ID)
		}
		seen[e.ID] = struct{}{}
		if e.Type != "prompt" && e.Type != "instruction" {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s has invalid type", e.ID)
		}
		if !isValidStatus(e.Status) {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s has invalid status", e.ID)
		}
		if _, err := time.Parse(time.RFC3339, e.UpdatedAt); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s has invalid updated_at", e.ID)
		}
		if e.Type == "prompt" && !strings.HasPrefix(e.Path, "prompts/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s path must be under prompts/", e.ID)
		}
		if e.Type == "instruction" && !strings.HasPrefix(e.Path, "instructions/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s path must be under instructions/", e.ID)
		}
		if !strings.HasPrefix(e.MetadataPath, "metadata/") {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s metadata path must be under metadata/", e.ID)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(e.Path))); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s content path does not exist", e.ID)
		}
		if _, err := os.Stat(filepath.Join(root, filepath.FromSlash(e.MetadataPath))); err != nil {
			return fmt.Errorf("ERR_SCHEMA_VALIDATION: index entry %s metadata path does not exist", e.ID)
		}
	}
	return nil
}

func validateSchemaVersion(version string) error {
	major, minor, err := parseMajorMinor(version)
	if err != nil {
		return fmt.Errorf("ERR_SCHEMA_VERSION_INVALID: %w", err)
	}
	if major > supportedMajor {
		return fmt.Errorf("ERR_SCHEMA_MAJOR_UNSUPPORTED: schema version %s is newer than supported major %d", version, supportedMajor)
	}
	if major == supportedMajor && minor > supportedMinor {
		fmt.Fprintf(os.Stderr, "WARN_SCHEMA_MINOR_NEWER_COMPAT: operating in compatibility mode for schema version %s\n", version)
	}
	return nil
}

func parseMajorMinor(v string) (int, int, error) {
	trimmed := strings.TrimSpace(v)
	parts := strings.Split(trimmed, ".")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("version must be MAJOR.MINOR")
	}
	var major, minor int
	_, err := fmt.Sscanf(trimmed, "%d.%d", &major, &minor)
	if err != nil {
		return 0, 0, fmt.Errorf("version must contain numeric MAJOR.MINOR")
	}
	return major, minor, nil
}

func isValidStatus(s string) bool {
	return s == "draft" || s == "approved"
}

func writeJSONAsYAML(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
