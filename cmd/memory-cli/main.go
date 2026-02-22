package main

import (
	"bytes"
	"crypto/sha256"
	"encoding/json"
	"encoding/hex"
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
	EventName      string `json:"event_name"`
	EventVersion   string `json:"event_version"`
	TimestampUTC   string `json:"timestamp_utc"`
	SessionID      string `json:"session_id"`
	TraceID        string `json:"trace_id"`
	ScenarioID     string `json:"scenario_id"`
	Operation      string `json:"operation"`
	Result         string `json:"result"`
	PolicyGate     string `json:"policy_gate"`
	MemoryType     string `json:"memory_type"`
	LatencyMS      int64  `json:"latency_ms"`
	SelectedID     string `json:"selected_id,omitempty"`
	SelectionMode  string `json:"selection_mode,omitempty"`
	SourcePath     string `json:"source_path,omitempty"`
	OperatorVerdict string `json:"operator_verdict"`
	ErrorCode      string `json:"error_code,omitempty"`
	Reason         string `json:"reason,omitempty"`
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

func emitTelemetry(root, telemetryPath string, ev telemetryEvent) error {
	path := strings.TrimSpace(telemetryPath)
	if path == "" {
		path = filepath.Join(root, filepath.FromSlash(telemetryRel))
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err := f.Write(append(data, '\n')); err != nil {
		return err
	}
	return nil
}

func telemetryErrorCode(err error) string {
	msg := strings.TrimSpace(err.Error())
	if msg == "" {
		return "ERR_UNKNOWN"
	}
	first := strings.Fields(msg)
	code := strings.TrimSuffix(first[0], ":")
	code = strings.TrimSpace(code)
	if strings.HasPrefix(code, "ERR_") {
		return code
	}
	return "ERR_COMMAND_FAILED"
}

func normalizeMemoryType(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "procedural", "state", "semantic":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "semantic"
	}
}

func normalizeOperatorVerdict(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "correct", "partially_correct", "incorrect", "not_scored":
		return strings.ToLower(strings.TrimSpace(v))
	default:
		return "not_scored"
	}
}

func normalizeTelemetryValue(v, fallback string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return fallback
	}
	return v
}

func enforceConstraintChecks(operation, sessionID, scenarioID, traceID string) error {
	if err := enforceCostConstraint(operation); err != nil {
		return err
	}
	if err := enforceTraceabilityConstraint(sessionID, scenarioID, traceID); err != nil {
		return err
	}
	if err := enforceReliabilityConstraint(); err != nil {
		return err
	}
	return nil
}

func enforceCostConstraint(operation string) error {
	maxPerRun := 0.50
	if v := strings.TrimSpace(os.Getenv("MEMORY_CONSTRAINT_COST_MAX_PER_RUN_USD")); v != "" {
		if _, err := fmt.Sscanf(v, "%f", &maxPerRun); err != nil {
			return errors.New("ERR_CONSTRAINT_COST_CONFIG_INVALID: MEMORY_CONSTRAINT_COST_MAX_PER_RUN_USD must be numeric")
		}
	}
	estimated := map[string]float64{
		"write":    0.08,
		"retrieve": 0.02,
		"evaluate": 0.30,
	}[operation]
	if estimated > maxPerRun {
		return fmt.Errorf("ERR_CONSTRAINT_COST_BUDGET_EXCEEDED: estimated_%s_cost_usd=%.2f exceeds max_per_run_usd=%.2f", operation, estimated, maxPerRun)
	}
	return nil
}

func enforceTraceabilityConstraint(sessionID, scenarioID, traceID string) error {
	if isTrue(os.Getenv("MEMORY_CONSTRAINT_FORCE_TRACE_MISSING")) {
		return errors.New("ERR_CONSTRAINT_TRACEABILITY_INCOMPLETE: trace policy forced missing required fields")
	}
	if strings.TrimSpace(sessionID) == "" || strings.TrimSpace(scenarioID) == "" || strings.TrimSpace(traceID) == "" {
		return errors.New("ERR_CONSTRAINT_TRACEABILITY_INCOMPLETE: session_id, scenario_id, and trace_id are required")
	}
	return nil
}

func enforceReliabilityConstraint() error {
	if isTrue(os.Getenv("MEMORY_CONSTRAINT_RELIABILITY_FREEZE")) {
		return errors.New("ERR_CONSTRAINT_RELIABILITY_FREEZE_ACTIVE: autonomous promotion paths are frozen")
	}
	return nil
}

func isLatencyDegraded(elapsedMs int64) bool {
	if isTrue(os.Getenv("MEMORY_CONSTRAINT_FORCE_LATENCY_DEGRADED")) {
		return true
	}
	threshold := int64(700)
	if v := strings.TrimSpace(os.Getenv("MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS")); v != "" {
		var parsed int64
		if _, err := fmt.Sscanf(v, "%d", &parsed); err == nil && parsed > 0 {
			threshold = parsed
		}
	}
	return elapsedMs > threshold
}

func createSnapshot(root, createdBy, reason string) (snapshotManifest, error) {
	idx, err := loadIndex(root)
	if err != nil {
		return snapshotManifest{}, err
	}

	now := time.Now().UTC()
	snapshotID := fmt.Sprintf("snapshot-%d", now.UnixNano())
	baseDir := filepath.Join(root, "snapshots", snapshotID)
	payloadDir := filepath.Join(baseDir, "payload")

	refs := collectSnapshotRefs(idx)
	if len(refs) == 0 {
		return snapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: no payload references found for snapshot")
	}

	checksums := make([]snapshotChecksum, 0, len(refs))
	for _, rel := range refs {
		src := filepath.Join(root, filepath.FromSlash(rel))
		dst := filepath.Join(payloadDir, filepath.FromSlash(rel))
		if err := copyFile(src, dst); err != nil {
			return snapshotManifest{}, err
		}
		sum, err := fileSHA256(dst)
		if err != nil {
			return snapshotManifest{}, err
		}
		checksums = append(checksums, snapshotChecksum{Path: rel, SHA256: sum})
	}

	manifest := snapshotManifest{
		SnapshotID:    snapshotID,
		CreatedAt:     now.Format(time.RFC3339),
		CreatedBy:     strings.TrimSpace(createdBy),
		SchemaVersion: defaultSchema,
		IndexVersion:  idx.SchemaVersion,
		Scope:         "full",
		Reason:        strings.TrimSpace(reason),
		Checksums:     checksums,
		PayloadRefs:   refs,
	}
	if err := validateSnapshotManifest(manifest); err != nil {
		return snapshotManifest{}, err
	}
	if err := writeJSONAsYAML(filepath.Join(baseDir, "manifest.json"), manifest); err != nil {
		return snapshotManifest{}, err
	}
	return manifest, nil
}

func listSnapshots(root string) ([]snapshotListRow, error) {
	snapRoot := filepath.Join(root, "snapshots")
	if _, err := os.Stat(snapRoot); errors.Is(err, os.ErrNotExist) {
		return []snapshotListRow{}, nil
	}

	entries, err := os.ReadDir(snapRoot)
	if err != nil {
		return nil, err
	}
	rows := make([]snapshotListRow, 0, len(entries))
	for _, ent := range entries {
		if !ent.IsDir() {
			continue
		}
		manifest, err := loadSnapshotManifest(root, ent.Name())
		if err != nil {
			return nil, err
		}
		rows = append(rows, snapshotListRow{
			SnapshotID:    manifest.SnapshotID,
			CreatedAt:     manifest.CreatedAt,
			CreatedBy:     manifest.CreatedBy,
			SchemaVersion: manifest.SchemaVersion,
			IndexVersion:  manifest.IndexVersion,
			Scope:         manifest.Scope,
			Reason:        manifest.Reason,
		})
	}
	sort.Slice(rows, func(i, j int) bool { return rows[i].CreatedAt > rows[j].CreatedAt })
	return rows, nil
}

func restoreSnapshot(root, snapshotID string) error {
	manifest, err := loadSnapshotManifest(root, snapshotID)
	if err != nil {
		return err
	}
	if err := checkSnapshotCompatibility(root, manifest); err != nil {
		return err
	}
	if err := verifySnapshotChecksums(root, manifest); err != nil {
		return err
	}

	// restore as forward revision by replacing active payload paths while preserving audits/snapshots history
	for _, rel := range manifest.PayloadRefs {
		src := filepath.Join(root, "snapshots", snapshotID, "payload", filepath.FromSlash(rel))
		dst := filepath.Join(root, filepath.FromSlash(rel))
		if err := copyFile(src, dst); err != nil {
			return err
		}
	}
	return nil
}

func loadSnapshotManifest(root, snapshotID string) (snapshotManifest, error) {
	path := filepath.Join(root, "snapshots", snapshotID, "manifest.json")
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return snapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: snapshot manifest not found")
		}
		return snapshotManifest{}, err
	}
	var m snapshotManifest
	if err := json.Unmarshal(data, &m); err != nil {
		return snapshotManifest{}, errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: cannot parse manifest")
	}
	if err := validateSnapshotManifest(m); err != nil {
		return snapshotManifest{}, err
	}
	return m, nil
}

func validateSnapshotManifest(m snapshotManifest) error {
	if strings.TrimSpace(m.SnapshotID) == "" || strings.TrimSpace(m.CreatedAt) == "" || strings.TrimSpace(m.CreatedBy) == "" ||
		strings.TrimSpace(m.SchemaVersion) == "" || strings.TrimSpace(m.IndexVersion) == "" || strings.TrimSpace(m.Scope) == "" || strings.TrimSpace(m.Reason) == "" {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: missing required manifest fields")
	}
	if m.Scope != "full" {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: unsupported snapshot scope")
	}
	if _, err := time.Parse(time.RFC3339, m.CreatedAt); err != nil {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: created_at must be RFC3339")
	}
	if len(m.PayloadRefs) == 0 || len(m.Checksums) == 0 {
		return errors.New("ERR_SNAPSHOT_MANIFEST_INVALID: payload references and checksums are required")
	}
	return nil
}

func checkSnapshotCompatibility(root string, m snapshotManifest) error {
	if err := validateSchemaVersion(m.SchemaVersion); err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	if err := validateSchemaVersion(m.IndexVersion); err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	cur, err := loadIndex(root)
	if err != nil {
		return err
	}
	snapMajor, _, err := parseMajorMinor(m.IndexVersion)
	if err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	curMajor, _, err := parseMajorMinor(cur.SchemaVersion)
	if err != nil {
		return fmt.Errorf("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: %w", err)
	}
	if snapMajor != curMajor {
		return errors.New("ERR_SNAPSHOT_COMPATIBILITY_BLOCKED: snapshot index major version does not match current index major")
	}
	return nil
}

func verifySnapshotChecksums(root string, m snapshotManifest) error {
	expected := map[string]string{}
	for _, c := range m.Checksums {
		expected[c.Path] = c.SHA256
	}
	for _, rel := range m.PayloadRefs {
		sum, err := fileSHA256(filepath.Join(root, "snapshots", m.SnapshotID, "payload", filepath.FromSlash(rel)))
		if err != nil {
			return fmt.Errorf("ERR_SNAPSHOT_INTEGRITY_CHECK_FAILED: %w", err)
		}
		if expected[rel] == "" || expected[rel] != sum {
			return errors.New("ERR_SNAPSHOT_INTEGRITY_CHECK_FAILED: checksum mismatch")
		}
	}
	return nil
}

func collectSnapshotRefs(idx indexFile) []string {
	refs := []string{"index.yaml"}
	for _, e := range idx.Entries {
		refs = append(refs, e.Path, e.MetadataPath)
	}
	seen := map[string]struct{}{}
	out := make([]string, 0, len(refs))
	for _, r := range refs {
		r = filepath.ToSlash(strings.TrimSpace(r))
		if r == "" {
			continue
		}
		if _, ok := seen[r]; ok {
			continue
		}
		seen[r] = struct{}{}
		out = append(out, r)
	}
	sort.Strings(out)
	return out
}

func writeSnapshotAudit(root string, ev snapshotAuditEvent) error {
	name := strings.ReplaceAll(ev.EventName, ".", "-")
	path := filepath.Join(root, "audits", fmt.Sprintf("%s-%d.json", name, time.Now().UTC().UnixNano()))
	return writeJSONAsYAML(path, ev)
}

func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0o644)
}

func fileSHA256(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:]), nil
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
