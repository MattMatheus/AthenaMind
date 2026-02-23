package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"athenamind/internal/gateway"
	"athenamind/internal/governance"
	"athenamind/internal/index"
	"athenamind/internal/retrieval"
	"athenamind/internal/snapshot"
	"athenamind/internal/telemetry"
	"athenamind/internal/types"
)

const (
	defaultEvaluationQuerySetPath = "cmd/memory-cli/testdata/eval-query-set-v1.json"
	defaultEvaluationCorpusID     = "memory-corpus-v1"
	defaultEvaluationQuerySetID   = "query-set-v1"
	defaultEvaluationConfigID     = "config-v1-confidence-0.34-margin-0.15"
)

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
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.write",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "write",
			Result:          result,
			PolicyGate:      "medium",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := governance.EnforceConstraintChecks("write", *sessionID, *scenarioID, traceID); err != nil {
		return err
	}

	policy, err := governance.EnforceWritePolicy(types.WritePolicyInput{
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

	if err := index.UpsertEntry(*root, types.UpsertEntryInput{
		ID:       *id,
		Title:    *title,
		Type:     *typeValue,
		Domain:   *domain,
		Body:     *body,
		BodyFile: *bodyFile,
		Stage:    *stage,
	}, policy); err != nil {
		return err
	}

	fmt.Printf("wrote entry %s at %s\n", *id, fmt.Sprintf("%s/%s/%s.md", map[bool]string{true: "instructions", false: "prompts"}[*typeValue == "instruction"], *domain, *id))
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
	telemetryResult := types.RetrieveResult{
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
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.retrieve",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "retrieve",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			SelectedID:      telemetryResult.SelectedID,
			SelectionMode:   telemetryResult.SelectionMode,
			SourcePath:      telemetryResult.SourcePath,
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := governance.EnforceConstraintChecks("retrieve", *sessionID, *scenarioID, traceID); err != nil {
		return err
	}

	result, err := retrieval.Retrieve(*root, *query, *domain)
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
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.evaluate",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenarioID, "scenario-manual"),
			Operation:       "evaluate",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := governance.EnforceConstraintChecks("evaluate", *sessionID, *scenarioID, traceID); err != nil {
		return err
	}

	queries, err := retrieval.LoadEvaluationQueries(*queryFile)
	if err != nil {
		return err
	}

	report, err := retrieval.EvaluateRetrieval(*root, queries, *corpusID, *querySetID, *configID)
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

func runBootstrap(args []string) (err error) {
	fs := flag.NewFlagSet("bootstrap", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	repo := fs.String("repo", "", "repository identifier")
	sessionID := fs.String("session-id", "", "telemetry session identifier")
	scenario := fs.String("scenario", "", "scenario identifier")
	memoryType := fs.String("memory-type", "state", "telemetry memory type: procedural|state|semantic")
	operatorVerdict := fs.String("operator-verdict", "not_scored", "telemetry operator verdict")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*repo) == "" {
		return errors.New("--repo is required")
	}
	if strings.TrimSpace(*sessionID) == "" {
		return errors.New("--session-id is required")
	}
	if strings.TrimSpace(*scenario) == "" {
		return errors.New("--scenario is required")
	}

	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	telemetryResult := types.BootstrapPayload{}
	defer func() {
		result := "success"
		errorCode := ""
		reason := ""
		selectedID := "__none__"
		selectionMode := "bootstrap_empty"
		sourcePath := "__none__"
		if len(telemetryResult.MemoryEntries) > 0 {
			selectedID = telemetryResult.MemoryEntries[0].ID
			selectionMode = telemetryResult.MemoryEntries[0].SelectionMode
			sourcePath = telemetryResult.MemoryEntries[0].SourcePath
		}
		if err != nil {
			result = "fail"
			errorCode = telemetry.TelemetryErrorCode(err)
			reason = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.bootstrap",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*scenario, "scenario-bootstrap"),
			Operation:       "bootstrap",
			Result:          result,
			PolicyGate:      "none",
			MemoryType:      telemetry.NormalizeMemoryType(*memoryType),
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			SelectedID:      selectedID,
			SelectionMode:   selectionMode,
			SourcePath:      sourcePath,
			OperatorVerdict: telemetry.NormalizeOperatorVerdict(*operatorVerdict),
			ErrorCode:       errorCode,
			Reason:          reason,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	if err := governance.EnforceConstraintChecks("retrieve", *sessionID, *scenario, traceID); err != nil {
		return err
	}

	payload, err := retrieval.Bootstrap(*root, *repo, *sessionID, *scenario)
	if err != nil {
		return err
	}
	telemetryResult = payload
	return printResult(payload)
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
	mux.Handle("/memory/retrieve", gateway.ReadGatewayHandler(*root))
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
	req := types.APIRetrieveRequest{
		Query:     *query,
		Domain:    *domain,
		SessionID: *sessionID,
	}
	resp, err := gateway.APIRetrieveWithFallback(*root, strings.TrimSpace(*gatewayURL), req, traceID, http.DefaultClient)
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
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.create.requested",
		SnapshotID: "pending",
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	manifest, err := snapshot.CreateSnapshot(*root, *createdBy, *reason)
	if err != nil {
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.create.completed",
			SnapshotID:  "pending",
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.create.completed",
		SnapshotID: manifest.SnapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
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
	rows, err := snapshot.ListSnapshots(*root)
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
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.restore.requested",
		SnapshotID: *snapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	policy, policyErr := governance.EnforceWritePolicy(types.WritePolicyInput{
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
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.restore.policy_decision",
			SnapshotID:  *snapshotID,
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.restore.failed",
			SnapshotID:  *snapshotID,
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:   "snapshot.restore.policy_decision",
		SnapshotID:  *snapshotID,
		SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:     traceID,
		Result:      "success",
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		Description: policy.Decision,
	})

	if err := snapshot.RestoreSnapshot(*root, *snapshotID); err != nil {
		_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
			EventName:   "snapshot.restore.failed",
			SnapshotID:  *snapshotID,
			SessionID:   telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:     traceID,
			Result:      "fail",
			Timestamp:   time.Now().UTC().Format(time.RFC3339),
			ErrorCode:   telemetry.TelemetryErrorCode(err),
			Description: err.Error(),
		})
		return err
	}
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.restore.completed",
		SnapshotID: *snapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	return nil
}

func printResult(r any) error {
	out, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(out))
	return nil
}

func exitErr(err error) {
	fmt.Fprintln(os.Stderr, err.Error())
	os.Exit(1)
}
