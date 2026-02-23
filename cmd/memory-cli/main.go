package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
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
	eventSchema    = "1.0"
	telemetryRel   = "telemetry/events.jsonl"

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

type telemetryEvent struct {
	EventName       string `json:"event_name"`
	EventVersion    string `json:"event_version"`
	TimestampUTC    string `json:"timestamp_utc"`
	SessionID       string `json:"session_id"`
	TraceID         string `json:"trace_id"`
	ScenarioID      string `json:"scenario_id"`
	Operation       string `json:"operation"`
	Result          string `json:"result"`
	PolicyGate      string `json:"policy_gate"`
	MemoryType      string `json:"memory_type"`
	LatencyMS       int64  `json:"latency_ms"`
	SelectedID      string `json:"selected_id,omitempty"`
	SelectionMode   string `json:"selection_mode,omitempty"`
	SourcePath      string `json:"source_path,omitempty"`
	OperatorVerdict string `json:"operator_verdict"`
	ErrorCode       string `json:"error_code,omitempty"`
	Reason          string `json:"reason,omitempty"`
}

type snapshotChecksum struct {
	Path   string `json:"path"`
	SHA256 string `json:"sha256"`
}

type snapshotManifest struct {
	SnapshotID    string             `json:"snapshot_id"`
	CreatedAt     string             `json:"created_at"`
	CreatedBy     string             `json:"created_by"`
	SchemaVersion string             `json:"schema_version"`
	IndexVersion  string             `json:"index_version"`
	Scope         string             `json:"scope"`
	Reason        string             `json:"reason"`
	Checksums     []snapshotChecksum `json:"checksums"`
	PayloadRefs   []string           `json:"payload_refs"`
}

type snapshotListRow struct {
	SnapshotID    string `json:"snapshot_id"`
	CreatedAt     string `json:"created_at"`
	CreatedBy     string `json:"created_by"`
	SchemaVersion string `json:"schema_version"`
	IndexVersion  string `json:"index_version"`
	Scope         string `json:"scope"`
	Reason        string `json:"reason"`
}

type snapshotAuditEvent struct {
	EventName   string `json:"event_name"`
	SnapshotID  string `json:"snapshot_id"`
	SessionID   string `json:"session_id"`
	TraceID     string `json:"trace_id"`
	Result      string `json:"result"`
	Timestamp   string `json:"timestamp_utc"`
	ErrorCode   string `json:"error_code,omitempty"`
	Description string `json:"description,omitempty"`
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
	case "snapshot":
		err = runSnapshot(os.Args[2:])
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

func runWrite(args []string) (err error) {
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
	sessionID := fs.String("session-id", "session-local", "telemetry session identifier")
	scenarioID := fs.String("scenario-id", "scenario-manual", "telemetry scenario identifier")
	memoryType := fs.String("memory-type", "semantic", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
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
	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	defer func() {
		result := "success"
		errorCode := ""
		reason := ""
		if err != nil {
			result = "fail"
			errorCode = telemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := emitTelemetry(*root, *telemetryFile, telemetryEvent{
			EventName:       "memory.write",
			EventVersion:    eventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      normalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "write",
			Result:          result,
			PolicyGate:      "medium",
			MemoryType:      normalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			OperatorVerdict: normalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := enforceConstraintChecks("write", *sessionID, *scenarioID, traceID); err != nil {
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

func runRetrieve(args []string) (err error) {
	fs := flag.NewFlagSet("retrieve", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "", "natural language query")
	domain := fs.String("domain", "", "optional domain filter")
	sessionID := fs.String("session-id", "session-local", "telemetry session identifier")
	scenarioID := fs.String("scenario-id", "scenario-manual", "telemetry scenario identifier")
	memoryType := fs.String("memory-type", "semantic", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	telemetryResult := retrieveResult{
		SelectedID:    "__none__",
		SelectionMode: "fallback_path_priority",
		SourcePath:    "__none__",
	}
	defer func() {
		result := "success"
		errorCode := ""
		reason := telemetryResult.Reason
		if err != nil {
			result = "fail"
			errorCode = telemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := emitTelemetry(*root, *telemetryFile, telemetryEvent{
			EventName:       "memory.retrieve",
			EventVersion:    eventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      normalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "retrieve",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      normalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			SelectedID:      telemetryResult.SelectedID,
			SelectionMode:   telemetryResult.SelectionMode,
			SourcePath:      telemetryResult.SourcePath,
			OperatorVerdict: normalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := enforceConstraintChecks("retrieve", *sessionID, *scenarioID, traceID); err != nil {
		return err
	}

	result, err := retrieve(*root, *query, *domain)
	if err != nil {
		return err
	}
	telemetryResult = result
	return printResult(result)
}

func runEvaluate(args []string) (err error) {
	fs := flag.NewFlagSet("evaluate", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	sessionID := fs.String("session-id", "session-local", "telemetry session identifier")
	scenarioID := fs.String("scenario-id", "scenario-manual", "telemetry scenario identifier")
	memoryType := fs.String("memory-type", "semantic", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	queryFile := fs.String("query-file", defaultEvaluationQuerySetPath, "path to evaluation query set JSON")
	corpusID := fs.String("corpus-id", defaultEvaluationCorpusID, "pinned corpus snapshot id")
	querySetID := fs.String("query-set-id", defaultEvaluationQuerySetID, "pinned query set id")
	configID := fs.String("config-id", defaultEvaluationConfigID, "retrieval configuration snapshot id")
	if err := fs.Parse(args); err != nil {
		return err
	}
	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	defer func() {
		result := "success"
		errorCode := ""
		reason := ""
		if err != nil {
			result = "fail"
			errorCode = telemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := emitTelemetry(*root, *telemetryFile, telemetryEvent{
			EventName:       "memory.evaluate",
			EventVersion:    eventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      normalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "evaluate",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      normalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			OperatorVerdict: normalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := enforceConstraintChecks("evaluate", *sessionID, *scenarioID, traceID); err != nil {
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

func runSnapshot(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: memory-cli snapshot <create|list|restore> [flags]")
	}
	switch args[0] {
	case "create":
		return runSnapshotCreate(args[1:])
	case "list":
		return runSnapshotList(args[1:])
	case "restore":
		return runSnapshotRestore(args[1:])
	default:
		return fmt.Errorf("unknown snapshot subcommand: %s", args[0])
	}
}

func runSnapshotCreate(args []string) (err error) {
	fs := flag.NewFlagSet("snapshot create", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	createdBy := fs.String("created-by", "", "snapshot actor")
	reason := fs.String("reason", "", "snapshot rationale")
	scope := fs.String("scope", "full", "snapshot scope")
	sessionID := fs.String("session-id", "session-local", "session identifier")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*createdBy) == "" {
		return errors.New("--created-by is required")
	}
	if strings.TrimSpace(*reason) == "" {
		return errors.New("--reason is required")
	}
	if strings.TrimSpace(*scope) != "full" {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: only scope=full is supported in v0.2 MVP")
	}

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	_ = writeSnapshotAudit(*root, snapshotAuditEvent{
		EventName:  "snapshot.create.requested",
		SnapshotID: "pending",
		SessionID:  normalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	manifest, err := createSnapshot(*root, *createdBy, *reason)
	if err != nil {
		_ = writeSnapshotAudit(*root, snapshotAuditEvent{
			EventName:   "snapshot.create.completed",
			SnapshotID:  "pending",
			SessionID:   normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = writeSnapshotAudit(*root, snapshotAuditEvent{
		EventName:  "snapshot.create.completed",
		SnapshotID: manifest.SnapshotID,
		SessionID:  normalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})
	return printResult(manifest)
}

func runSnapshotList(args []string) error {
	fs := flag.NewFlagSet("snapshot list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	rows, err := listSnapshots(*root)
	if err != nil {
		return err
	}
	return printResult(rows)
}

func runSnapshotRestore(args []string) (err error) {
	fs := flag.NewFlagSet("snapshot restore", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	snapshotID := fs.String("snapshot-id", "", "snapshot id")
	sessionID := fs.String("session-id", "session-local", "session identifier")
	stage := fs.String("stage", "pm", "policy stage: planning|architect|pm")
	reviewer := fs.String("reviewer", "", "reviewer identity")
	approved := fs.Bool("approved", false, "legacy flag: equivalent to --decision=approved")
	decision := fs.String("decision", "", "review decision: approved|rejected")
	notes := fs.String("notes", "", "decision notes")
	reason := fs.String("reason", "", "reason for restore")
	risk := fs.String("risk", "", "risk and mitigation note")
	reworkNotes := fs.String("rework-notes", "", "required when --decision=rejected")
	reReviewedBy := fs.String("re-reviewed-by", "", "required when --decision=rejected")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*snapshotID) == "" {
		return errors.New("--snapshot-id is required")
	}

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	_ = writeSnapshotAudit(*root, snapshotAuditEvent{
		EventName:  "snapshot.restore.requested",
		SnapshotID: *snapshotID,
		SessionID:  normalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	policy, policyErr := enforceWritePolicy(writePolicyInput{
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
	if policyErr != nil {
		err = fmt.Errorf("ERR_SNAPSHOT_RESTORE_POLICY_BLOCKED: %w", policyErr)
		_ = writeSnapshotAudit(*root, snapshotAuditEvent{
			EventName:   "snapshot.restore.policy_decision",
			SnapshotID:  *snapshotID,
			SessionID:   normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetryErrorCode(err),
			Description: err.Error(),
		})
		_ = writeSnapshotAudit(*root, snapshotAuditEvent{
			EventName:   "snapshot.restore.failed",
			SnapshotID:  *snapshotID,
			SessionID:   normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = writeSnapshotAudit(*root, snapshotAuditEvent{
		EventName:   "snapshot.restore.policy_decision",
		SnapshotID:  *snapshotID,
		SessionID:   normalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:     traceID,
		Result:      "success",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Description: policy.Decision,
	})

	if err := restoreSnapshot(*root, *snapshotID); err != nil {
		_ = writeSnapshotAudit(*root, snapshotAuditEvent{
			EventName:   "snapshot.restore.failed",
			SnapshotID:  *snapshotID,
			SessionID:   normalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = writeSnapshotAudit(*root, snapshotAuditEvent{
		EventName:  "snapshot.restore.completed",
		SnapshotID: *snapshotID,
		SessionID:  normalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	return nil
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
