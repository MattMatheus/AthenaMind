package main

import (
	"bufio"
	"context"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"athenamind/internal/episode"
	"athenamind/internal/gateway"
	"athenamind/internal/governance"
	"athenamind/internal/index"
	"athenamind/internal/retrieval"
	"athenamind/internal/snapshot"
	"athenamind/internal/telemetry"
	"athenamind/internal/types"
	"go.opentelemetry.io/otel/attribute"
)

const (
	defaultEvaluationQuerySetPath = "cmd/memory-cli/testdata/eval-query-set-v1.json"
	defaultEvaluationCorpusID     = "memory-corpus-v1"
	defaultEvaluationQuerySetID   = "query-set-v1"
	defaultEvaluationConfigID     = "config-v1-confidence-0.34-margin-0.15"
)

func runWrite(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "write")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

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
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint for write-time indexing")

	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("memory.id", *id),
		attribute.String("memory.domain", *domain),
		attribute.String("memory.type", *typeValue),
		attribute.String("session.id", *sessionID),
	)
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
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("write", *sessionID, *scenarioID, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	var policy types.WritePolicyDecision
	_, policySpan := telemetry.StartSpan(ctx, "memory.policy.enforce")
	policy, err = governance.EnforceWritePolicy(types.WritePolicyInput{
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
	telemetry.EndSpan(policySpan, err)
	if err != nil {
		return err
	}

	_, upsertSpan := telemetry.StartSpan(ctx, "memory.index.upsert")
	err = index.UpsertEntry(*root, types.UpsertEntryInput{
		ID:       *id,
		Title:    *title,
		Type:     *typeValue,
		Domain:   *domain,
		Body:     *body,
		BodyFile: *bodyFile,
		Stage:    *stage,
	}, policy)
	telemetry.EndSpan(upsertSpan, err)
	if err != nil {
		return err
	}
	_, embedSpan := telemetry.StartSpan(ctx, "memory.embedding.index_entry")
	warning, err := retrieval.IndexEntryEmbedding(*root, *id, *embeddingEndpoint, *sessionID)
	telemetry.EndSpan(embedSpan, err)
	if err != nil {
		return err
	}
	if strings.TrimSpace(warning) != "" {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}

	fmt.Printf("wrote entry %s at %s\n", *id, fmt.Sprintf("%s/%s/%s.md", map[bool]string{true: "instructions", false: "prompts"}[*typeValue == "instruction"], *domain, *id))
	return nil
}

func runRetrieve(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "retrieve")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

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
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	mode := fs.String("mode", "classic", "retrieval mode: classic|hybrid")
	topK := fs.Int("top-k", 5, "number of candidate traces to return (1-50)")
	backend := fs.String("retrieval-backend", "sqlite", "retrieval backend: sqlite|qdrant|neo4j")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("memory.query", *query),
		attribute.String("memory.domain", *domain),
		attribute.String("session.id", *sessionID),
	)
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
		semanticHit := telemetryResult.SelectionMode == "semantic" ||
			telemetryResult.SelectionMode == "embedding_semantic" ||
			telemetryResult.SelectionMode == "hybrid_rrf"
		fallbackUsed := strings.HasPrefix(telemetryResult.SelectionMode, "fallback_")
		semanticRate := 0.0
		fallbackRate := 0.0
		if rate, ferr := telemetry.EmitRetrievalMetric(*root, telemetryResult); ferr == nil {
			semanticRate = rate.SemanticHitRate
			fallbackRate = rate.FallbackRate
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
			FallbackUsed:    fallbackUsed,
			SemanticHit:     semanticHit,
			PrecisionProxy:  telemetryResult.PrecisionHint,
			SemanticHitRate: semanticRate,
			FallbackRate:    fallbackRate,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("retrieve", *sessionID, *scenarioID, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	_, retrieveSpan := telemetry.StartSpan(ctx, "memory.retrieve.execute")
	result, warning, err := retrieval.RetrieveWithOptionsAndEndpointAndSession(
		*root,
		*query,
		*domain,
		*embeddingEndpoint,
		*sessionID,
		retrieval.RetrieveOptions{
			Mode:    *mode,
			TopK:    *topK,
			Backend: *backend,
		},
	)
	telemetry.EndSpan(retrieveSpan, err)
	if err != nil {
		return err
	}
	if strings.TrimSpace(warning) != "" {
		fmt.Fprintf(os.Stderr, "warning: %s\n", warning)
	}
	telemetryResult = result
	return printResult(result)
}

func runEvaluate(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "evaluate")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

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
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	mode := fs.String("mode", "classic", "retrieval mode under evaluation: classic|hybrid")
	topK := fs.Int("top-k", 5, "candidate trace size under evaluation (1-50)")
	backend := fs.String("retrieval-backend", "sqlite", "retrieval backend under evaluation: sqlite|qdrant|neo4j")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("session.id", *sessionID),
		attribute.String("scenario.id", *scenarioID),
		attribute.String("evaluation.query_file", *queryFile),
	)
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
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("evaluate", *sessionID, *scenarioID, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	_, loadQueriesSpan := telemetry.StartSpan(ctx, "memory.evaluation.load_queries")
	queries, err := retrieval.LoadEvaluationQueries(*queryFile)
	telemetry.EndSpan(loadQueriesSpan, err)
	if err != nil {
		return err
	}

	_, evaluateSpan := telemetry.StartSpan(ctx, "memory.evaluation.execute")
	report, err := retrieval.EvaluateRetrievalWithOptionsAndEmbeddingEndpoint(
		*root,
		queries,
		*corpusID,
		*querySetID,
		*configID,
		*embeddingEndpoint,
		retrieval.RetrieveOptions{
			Mode:    *mode,
			TopK:    *topK,
			Backend: *backend,
		},
	)
	telemetry.EndSpan(evaluateSpan, err)
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
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "bootstrap")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

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
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("repo.id", *repo),
		attribute.String("session.id", *sessionID),
		attribute.String("scenario.id", *scenario),
	)

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
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("retrieve", *sessionID, *scenario, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}

	_, bootstrapSpan := telemetry.StartSpan(ctx, "memory.bootstrap.execute")
	payload, err := retrieval.Bootstrap(*root, *repo, *sessionID, *scenario)
	telemetry.EndSpan(bootstrapSpan, err)
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
	mode := fs.String("mode", "classic", "retrieval mode: classic|hybrid")
	topK := fs.Int("top-k", 5, "number of candidate traces to request (1-50)")
	backend := fs.String("retrieval-backend", "sqlite", "retrieval backend: sqlite|qdrant|neo4j")
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
		Mode:      *mode,
		TopK:      *topK,
		Backend:   *backend,
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

func runEpisode(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: memory-cli episode <write|list> [flags]")
	}
	switch args[0] {
	case "write":
		return runEpisodeWrite(args[1:])
	case "list":
		return runEpisodeList(args[1:])
	default:
		return fmt.Errorf("unknown episode subcommand: %s", args[0])
	}
}

func runVerify(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: memory-cli verify <embeddings|health> [flags]")
	}
	switch args[0] {
	case "embeddings":
		return runVerifyEmbeddings(args[1:])
	case "health":
		return runVerifyHealth(args[1:])
	default:
		return fmt.Errorf("unknown verify subcommand: %s", args[0])
	}
}

func runTelemetry(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: memory-cli telemetry <tail> [flags]")
	}
	switch args[0] {
	case "tail":
		return runTelemetryTail(args[1:])
	default:
		return fmt.Errorf("unknown telemetry subcommand: %s", args[0])
	}
}

func runTelemetryTail(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "telemetry.tail")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("telemetry tail", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	lines := fs.Int("lines", 20, "number of recent records to show per source")
	source := fs.String("source", "events", "telemetry source: events|retrieval|both")
	follow := fs.Bool("follow", false, "stream new telemetry records after printing recent lines")
	followPollMS := fs.Int("follow-poll-ms", 500, "follow polling interval in milliseconds")
	followSeconds := fs.Int("follow-seconds", 0, "optional follow duration in seconds (0 = until interrupted)")
	operation := fs.String("operation", "", "optional filter: operation value")
	result := fs.String("result", "", "optional filter: result value")
	sessionID := fs.String("session-id", "", "optional filter: session_id value")
	eventsFile := fs.String("telemetry-file", "", "optional events jsonl path (default: <root>/telemetry/events.jsonl)")
	retrievalFile := fs.String("retrieval-metrics-file", "", "optional retrieval metrics jsonl path (default: <root>/telemetry/retrieval-metrics.jsonl)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *lines < 1 {
		return errors.New("--lines must be >= 1")
	}
	if *followPollMS < 50 {
		return errors.New("--follow-poll-ms must be >= 50")
	}
	if *followSeconds < 0 {
		return errors.New("--follow-seconds must be >= 0")
	}

	mode := strings.ToLower(strings.TrimSpace(*source))
	if mode != "events" && mode != "retrieval" && mode != "both" {
		return errors.New("--source must be one of: events|retrieval|both")
	}
	filters := telemetryRecordFilters{
		Operation: strings.TrimSpace(*operation),
		Result:    strings.TrimSpace(*result),
		SessionID: strings.TrimSpace(*sessionID),
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("telemetry.source", mode),
		attribute.Int("telemetry.lines", *lines),
		attribute.Bool("telemetry.follow", *follow),
		attribute.String("telemetry.filter.operation", filters.Operation),
		attribute.String("telemetry.filter.result", filters.Result),
		attribute.String("telemetry.filter.session_id", filters.SessionID),
	)

	out := map[string]any{
		"root":   *root,
		"source": mode,
		"lines":  *lines,
		"follow": *follow,
	}
	out["filters"] = map[string]string{
		"operation":  filters.Operation,
		"result":     filters.Result,
		"session_id": filters.SessionID,
	}

	streamPaths := map[string]string{}
	if mode == "events" || mode == "both" {
		path := strings.TrimSpace(*eventsFile)
		if path == "" {
			path = filepath.Join(*root, filepath.FromSlash(telemetry.TelemetryRel))
		}
		_, eventsSpan := telemetry.StartSpan(ctx, "memory.telemetry.tail.events")
		rows, tailErr := tailJSONL(path, *lines)
		telemetry.EndSpan(eventsSpan, tailErr)
		if tailErr != nil {
			return tailErr
		}
		rows = filterTelemetryRows(rows, filters)
		out["events_path"] = path
		out["events"] = rows
		out["events_count"] = len(rows)
		streamPaths["events"] = path
	}

	if mode == "retrieval" || mode == "both" {
		path := strings.TrimSpace(*retrievalFile)
		if path == "" {
			path = filepath.Join(*root, filepath.FromSlash(telemetry.RetrievalMetricsRel))
		}
		_, retrievalSpan := telemetry.StartSpan(ctx, "memory.telemetry.tail.retrieval")
		rows, tailErr := tailJSONL(path, *lines)
		telemetry.EndSpan(retrievalSpan, tailErr)
		if tailErr != nil {
			return tailErr
		}
		rows = filterTelemetryRows(rows, filters)
		out["retrieval_metrics_path"] = path
		out["retrieval_metrics"] = rows
		out["retrieval_metrics_count"] = len(rows)
		streamPaths["retrieval_metrics"] = path
	}
	if err := printResult(out); err != nil {
		return err
	}
	if !*follow {
		return nil
	}
	return followTelemetryTail(ctx, streamPaths, filters, time.Duration(*followPollMS)*time.Millisecond, time.Duration(*followSeconds)*time.Second)
}

func tailJSONL(path string, maxLines int) ([]map[string]any, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []map[string]any{}, nil
		}
		return nil, err
	}
	defer f.Close()

	if maxLines < 1 {
		return []map[string]any{}, nil
	}
	lines := make([]string, 0, maxLines)
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		if len(lines) == maxLines {
			lines = lines[1:]
		}
		lines = append(lines, line)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(lines) == 0 {
		return []map[string]any{}, nil
	}
	out := make([]map[string]any, 0, len(lines))
	for _, line := range lines {
		row, err := parseJSONLine(line)
		if err != nil {
			return nil, err
		}
		out = append(out, row)
	}
	return out, nil
}

func parseJSONLine(line string) (map[string]any, error) {
	var row map[string]any
	if err := json.Unmarshal([]byte(line), &row); err != nil {
		return nil, fmt.Errorf("parse jsonl line: %w", err)
	}
	return row, nil
}

type telemetryRecordFilters struct {
	Operation string
	Result    string
	SessionID string
}

func filterTelemetryRows(rows []map[string]any, filters telemetryRecordFilters) []map[string]any {
	if filters.Operation == "" && filters.Result == "" && filters.SessionID == "" {
		return rows
	}
	out := make([]map[string]any, 0, len(rows))
	for _, row := range rows {
		if telemetryRowMatches(row, filters) {
			out = append(out, row)
		}
	}
	return out
}

func telemetryRowMatches(row map[string]any, filters telemetryRecordFilters) bool {
	if filters.Operation != "" && telemetryRowString(row, "operation") != filters.Operation {
		return false
	}
	if filters.Result != "" && telemetryRowString(row, "result") != filters.Result {
		return false
	}
	if filters.SessionID != "" && telemetryRowString(row, "session_id") != filters.SessionID {
		return false
	}
	return true
}

func telemetryRowString(row map[string]any, key string) string {
	v, ok := row[key]
	if !ok || v == nil {
		return ""
	}
	s, ok := v.(string)
	if ok {
		return s
	}
	return strings.TrimSpace(fmt.Sprint(v))
}

func followTelemetryTail(ctx context.Context, paths map[string]string, filters telemetryRecordFilters, poll time.Duration, duration time.Duration) error {
	notifyCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	runCtx := notifyCtx
	if duration > 0 {
		var cancel context.CancelFunc
		runCtx, cancel = context.WithTimeout(notifyCtx, duration)
		defer cancel()
	}

	offsets := map[string]int64{}
	for source, path := range paths {
		size, err := fileSize(path)
		if err != nil {
			return err
		}
		offsets[source] = size
	}

	ticker := time.NewTicker(poll)
	defer ticker.Stop()

	for {
		select {
		case <-runCtx.Done():
			return nil
		case <-ticker.C:
			sources := sortedKeys(paths)
			for _, source := range sources {
				path := paths[source]
				rows, nextOffset, err := readJSONLFromOffset(path, offsets[source])
				if err != nil {
					return err
				}
				offsets[source] = nextOffset
				for _, row := range rows {
					if !telemetryRowMatches(row, filters) {
						continue
					}
					event := map[string]any{
						"source": source,
						"record": row,
					}
					data, err := json.Marshal(event)
					if err != nil {
						return err
					}
					fmt.Println(string(data))
				}
			}
		}
	}
}

func fileSize(path string) (int64, error) {
	info, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return 0, nil
		}
		return 0, err
	}
	return info.Size(), nil
}

func readJSONLFromOffset(path string, offset int64) ([]map[string]any, int64, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return []map[string]any{}, 0, nil
		}
		return nil, offset, err
	}
	defer f.Close()

	info, err := f.Stat()
	if err != nil {
		return nil, offset, err
	}
	if info.Size() < offset {
		offset = 0
	}
	if _, err := f.Seek(offset, 0); err != nil {
		return nil, offset, err
	}

	rows := []map[string]any{}
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" {
			continue
		}
		row, err := parseJSONLine(line)
		if err != nil {
			return nil, offset, err
		}
		rows = append(rows, row)
	}
	if err := sc.Err(); err != nil {
		return nil, offset, err
	}
	nextOffset, err := f.Seek(0, 1)
	if err != nil {
		return nil, offset, err
	}
	return rows, nextOffset, nil
}

func sortedKeys(m map[string]string) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func runVerifyEmbeddings(args []string) error {
	fs := flag.NewFlagSet("verify embeddings", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	showMissing := fs.Bool("show-missing", true, "include missing entry ids")
	if err := fs.Parse(args); err != nil {
		return err
	}

	idx, err := index.LoadIndex(*root)
	if err != nil {
		return err
	}
	embeddings, err := index.GetEmbeddingRecords(*root, nil)
	if err != nil {
		return err
	}

	missing := make([]string, 0)
	for _, e := range idx.Entries {
		if rec, ok := embeddings[e.ID]; !ok || len(rec.Vector) == 0 {
			missing = append(missing, e.ID)
		}
	}
	sort.Strings(missing)

	provider := "ollama_or_custom_endpoint"
	if strings.TrimSpace(os.Getenv("AZURE_OPENAI_ENDPOINT")) != "" {
		provider = "azure_openai"
	}

	report := struct {
		Root              string   `json:"root"`
		ProviderHint      string   `json:"provider_hint"`
		IndexedEntries    int      `json:"indexed_entries"`
		StoredEmbeddings  int      `json:"stored_embeddings"`
		MissingEmbeddings int      `json:"missing_embeddings"`
		MissingIDs        []string `json:"missing_ids,omitempty"`
	}{
		Root:              *root,
		ProviderHint:      provider,
		IndexedEntries:    len(idx.Entries),
		StoredEmbeddings:  len(embeddings),
		MissingEmbeddings: len(missing),
	}
	if *showMissing {
		report.MissingIDs = missing
	}

	return printResult(report)
}

func runVerifyHealth(args []string) error {
	fs := flag.NewFlagSet("verify health", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	query := fs.String("query", "memory lifecycle", "health-check query")
	domain := fs.String("domain", "", "optional domain filter")
	sessionID := fs.String("session-id", "session-local", "session identifier")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	if err := fs.Parse(args); err != nil {
		return err
	}

	report, err := retrieval.EvaluateSemanticHealth(*root, *query, *domain, *embeddingEndpoint, *sessionID)
	if err != nil {
		return err
	}
	health := struct {
		Root              string  `json:"root"`
		Query             string  `json:"query"`
		IndexedEntries    int     `json:"indexed_entries"`
		StoredEmbeddings  int     `json:"stored_embeddings"`
		MissingEmbeddings int     `json:"missing_embeddings"`
		CoverageOK        bool    `json:"coverage_ok"`
		SelectionMode     string  `json:"selection_mode"`
		SemanticAvailable bool    `json:"semantic_available"`
		Warning           string  `json:"warning,omitempty"`
		SelectedID        string  `json:"selected_id"`
		Confidence        float64 `json:"confidence"`
		Pass              bool    `json:"pass"`
	}{
		Root:              *root,
		Query:             *query,
		IndexedEntries:    report.IndexedEntries,
		StoredEmbeddings:  report.StoredEmbeddings,
		MissingEmbeddings: report.MissingEmbeddings,
		CoverageOK:        report.CoverageOK,
		SelectionMode:     report.SelectionMode,
		SemanticAvailable: report.SemanticAvailable,
		Warning:           report.Warning,
		SelectedID:        report.SelectedID,
		Confidence:        report.Confidence,
		Pass:              report.Pass,
	}
	return printResult(health)
}

func runEpisodeWrite(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "episode.write")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("episode write", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	repo := fs.String("repo", "", "repository identifier")
	sessionID := fs.String("session-id", "", "session identifier")
	cycleID := fs.String("cycle-id", "", "cycle id")
	storyID := fs.String("story-id", "", "story id")
	outcome := fs.String("outcome", "", "episode outcome: success|partial|blocked")
	summary := fs.String("summary", "", "episode summary")
	summaryFile := fs.String("summary-file", "", "summary file path")
	filesChanged := fs.String("files-changed", "", "comma-separated changed files")
	decisions := fs.String("decisions", "", "episode decisions")
	decisionsFile := fs.String("decisions-file", "", "decisions file path")
	stage := fs.String("stage", "pm", "policy stage: planning|architect|pm")
	reviewer := fs.String("reviewer", "", "reviewer identity")
	approved := fs.Bool("approved", false, "legacy flag: equivalent to --decision=approved")
	decision := fs.String("decision", "", "review decision: approved|rejected")
	notes := fs.String("notes", "", "decision notes")
	reason := fs.String("reason", "", "reason for write")
	risk := fs.String("risk", "", "risk and mitigation note")
	reworkNotes := fs.String("rework-notes", "", "required when --decision=rejected")
	reReviewedBy := fs.String("re-reviewed-by", "", "required when --decision=rejected")
	telemetryFile := fs.String("telemetry-file", "", "optional telemetry output file (default: <root>/telemetry/events.jsonl)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("repo.id", *repo),
		attribute.String("session.id", *sessionID),
		attribute.String("cycle.id", *cycleID),
		attribute.String("story.id", *storyID),
	)

	startedAt := time.Now().UTC()
	traceID := fmt.Sprintf("trace-%d", startedAt.UnixNano())
	defer func() {
		result := "success"
		errorCode := ""
		reasonMsg := ""
		if err != nil {
			result = "fail"
			errorCode = telemetry.TelemetryErrorCode(err)
			reasonMsg = err.Error()
		}
		emitErr := telemetry.Emit(*root, *telemetryFile, types.TelemetryEvent{
			EventName:       "memory.episode.write",
			EventVersion:    telemetry.EventSchema,
			TimestampUTC:    time.Now().UTC().Format(time.RFC3339),
			SessionID:       telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
			TraceID:         traceID,
			ScenarioID:      telemetry.NormalizeTelemetryValue(*cycleID, "episode"),
			Operation:       "episode_write",
			Result:          result,
			PolicyGate:      "medium",
			MemoryType:      "state",
			LatencyMS:       time.Since(startedAt).Milliseconds(),
			OperatorVerdict: "not_scored",
			ErrorCode:       errorCode,
			Reason:          reasonMsg,
		})
		if err == nil && emitErr != nil {
			err = emitErr
		}
	}()
	_, enforceSpan := telemetry.StartSpan(ctx, "memory.constraint_check")
	err = governance.EnforceConstraintChecks("write", *sessionID, *cycleID, traceID)
	telemetry.EndSpan(enforceSpan, err)
	if err != nil {
		return err
	}
	var policy types.WritePolicyDecision
	_, policySpan := telemetry.StartSpan(ctx, "memory.policy.enforce")
	policy, err = governance.EnforceWritePolicy(types.WritePolicyInput{
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
	telemetry.EndSpan(policySpan, err)
	if err != nil {
		return err
	}
	_, writeSpan := telemetry.StartSpan(ctx, "memory.episode.write_record")
	record, err := episode.Write(*root, types.WriteEpisodeInput{
		Repo:          *repo,
		SessionID:     *sessionID,
		CycleID:       *cycleID,
		StoryID:       *storyID,
		Outcome:       *outcome,
		Summary:       *summary,
		SummaryFile:   *summaryFile,
		FilesChanged:  *filesChanged,
		Decisions:     *decisions,
		DecisionsFile: *decisionsFile,
		Stage:         *stage,
	}, policy)
	telemetry.EndSpan(writeSpan, err)
	if err != nil {
		return err
	}
	return printResult(record)
}

func runEpisodeList(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "episode.list")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("episode list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)

	root := fs.String("root", "memory", "memory root path")
	repo := fs.String("repo", "", "repository identifier")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("repo.id", *repo),
	)
	_, listSpan := telemetry.StartSpan(ctx, "memory.episode.list")
	rows, err := episode.List(*root, *repo)
	telemetry.EndSpan(listSpan, err)
	if err != nil {
		return err
	}
	return printResult(rows)
}

func runSnapshotCreate(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "snapshot.create")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

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
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("session.id", *sessionID),
		attribute.String("snapshot.scope", *scope),
	)

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.create.requested",
		SnapshotID: "pending",
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	_, snapshotSpan := telemetry.StartSpan(ctx, "memory.snapshot.create")
	manifest, err := snapshot.CreateSnapshot(*root, *createdBy, *reason)
	telemetry.EndSpan(snapshotSpan, err)
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

func runSnapshotList(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "snapshot.list")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("snapshot list", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(attribute.String("memory.root", *root))
	_, listSpan := telemetry.StartSpan(ctx, "memory.snapshot.list")
	rows, err := snapshot.ListSnapshots(*root)
	telemetry.EndSpan(listSpan, err)
	if err != nil {
		return err
	}
	return printResult(rows)
}

func runSnapshotRestore(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "snapshot.restore")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

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
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("snapshot.id", *snapshotID),
		attribute.String("session.id", *sessionID),
	)

	traceID := fmt.Sprintf("trace-%d", time.Now().UTC().UnixNano())
	_ = snapshot.WriteSnapshotAudit(*root, types.SnapshotAuditEvent{
		EventName:  "snapshot.restore.requested",
		SnapshotID: *snapshotID,
		SessionID:  telemetry.NormalizeTelemetryValue(*sessionID, "session-local"),
		TraceID:    traceID,
		Result:     "success",
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	})

	_, policySpan := telemetry.StartSpan(ctx, "memory.policy.enforce")
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
	telemetry.EndSpan(policySpan, policyErr)
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

	_, restoreSpan := telemetry.StartSpan(ctx, "memory.snapshot.restore")
	err = snapshot.RestoreSnapshot(*root, *snapshotID)
	telemetry.EndSpan(restoreSpan, err)
	if err != nil {
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

func runReindexAll(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "reindex-all")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("reindex-all", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	if err := fs.Parse(args); err != nil {
		return err
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("embedding.endpoint", *embeddingEndpoint),
	)

	_, loadIndexSpan := telemetry.StartSpan(ctx, "memory.index.load")
	idx, err := index.LoadIndex(*root)
	telemetry.EndSpan(loadIndexSpan, err)
	if err != nil {
		return err
	}

	_, loadEmbeddingSpan := telemetry.StartSpan(ctx, "memory.embeddings.load")
	embeddings, err := index.GetEmbeddingRecords(*root, nil)
	telemetry.EndSpan(loadEmbeddingSpan, err)
	if err != nil {
		return err
	}

	var missing []string
	for _, e := range idx.Entries {
		if rec, ok := embeddings[e.ID]; !ok || len(rec.Vector) == 0 {
			missing = append(missing, e.ID)
		}
	}

	if len(missing) == 0 {
		fmt.Println("No entries missing embeddings.")
		return nil
	}

	fmt.Printf("Reindexing %d entries...\n", len(missing))
	_, batchSpan := telemetry.StartSpan(ctx, "memory.embedding.batch_index")
	warnings, err := retrieval.IndexEntriesEmbeddingBatch(*root, missing, *embeddingEndpoint, "")
	telemetry.EndSpan(batchSpan, err)
	if err != nil {
		return err
	}

	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}

	fmt.Printf("Successfully processed %d entries.\n", len(missing))
	return nil
}

func runSyncQdrant(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "sync-qdrant")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("sync-qdrant", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	qdrantURL := fs.String("qdrant-url", "", "qdrant base URL (default env ATHENA_QDRANT_URL or http://localhost:6333)")
	collection := fs.String("collection", "", "qdrant collection name (default env ATHENA_QDRANT_COLLECTION or athena_memories)")
	batchSize := fs.Int("batch-size", 128, "upsert batch size")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if *batchSize < 1 {
		return errors.New("--batch-size must be >= 1")
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("qdrant.url", *qdrantURL),
		attribute.String("qdrant.collection", *collection),
		attribute.Int("qdrant.batch_size", *batchSize),
	)

	_, syncSpan := telemetry.StartSpan(ctx, "memory.qdrant.sync")
	report, err := retrieval.SyncQdrantCollection(*root, *qdrantURL, *collection, *batchSize)
	telemetry.EndSpan(syncSpan, err)
	if err != nil {
		return err
	}
	return printResult(report)
}

func runCrawl(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "crawl")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("crawl", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	dir := fs.String("dir", "", "directory to crawl for markdown files")
	domain := fs.String("domain", "auto-crawled", "domain for crawled entries")
	reviewer := fs.String("reviewer", "system", "reviewer identity")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")

	if err := fs.Parse(args); err != nil {
		return err
	}
	if *dir == "" {
		return errors.New("--dir is required")
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("crawl.dir", *dir),
		attribute.String("crawl.domain", *domain),
	)

	var mdFiles []string
	_, walkSpan := telemetry.StartSpan(ctx, "memory.crawl.walk_markdown")
	err = filepath.Walk(*dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			// Skip files in the memory root itself
			absPath, _ := filepath.Abs(path)
			absRoot, _ := filepath.Abs(*root)
			if !strings.HasPrefix(absPath, absRoot) {
				mdFiles = append(mdFiles, path)
			}
		}
		return nil
	})
	telemetry.EndSpan(walkSpan, err)
	if err != nil {
		return err
	}

	if len(mdFiles) == 0 {
		fmt.Printf("No markdown files found in %s\n", *dir)
		return nil
	}

	fmt.Printf("Found %d markdown files. Indexing...\n", len(mdFiles))

	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: *reviewer,
		Reason:   "bulk crawl indexing",
		Notes:    "automated crawl",
		Risk:     "none",
	}

	var processedIDs []string
	_, upsertSpan := telemetry.StartSpan(ctx, "memory.crawl.upsert_entries")
	for _, f := range mdFiles {
		id := buildCrawlEntryID(*dir, f)
		title := strings.Title(strings.ReplaceAll(id, "-", " "))

		upsertIn := types.UpsertEntryInput{
			ID:       id,
			Title:    title,
			Type:     "instruction",
			Domain:   *domain,
			BodyFile: f,
			Stage:    "pm",
		}

		if err := index.UpsertEntry(*root, upsertIn, policy); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to index %s: %v\n", f, err)
			continue
		}
		processedIDs = append(processedIDs, id)
	}
	telemetry.EndSpan(upsertSpan, nil)

	if len(processedIDs) > 0 {
		fmt.Printf("Batch embedding %d entries...\n", len(processedIDs))
		_, embedSpan := telemetry.StartSpan(ctx, "memory.embedding.batch_index")
		warnings, err := retrieval.IndexEntriesEmbeddingBatch(*root, processedIDs, *embeddingEndpoint, "")
		telemetry.EndSpan(embedSpan, err)
		if err != nil {
			return err
		}
		for _, w := range warnings {
			fmt.Fprintf(os.Stderr, "warning: %s\n", w)
		}
	}

	fmt.Printf("Successfully crawled and indexed %d files.\n", len(processedIDs))
	return nil
}

func runReembedChanged(args []string) (err error) {
	ctx, commandSpan := telemetry.StartCommandSpan(context.Background(), "reembed-changed")
	defer func() {
		telemetry.EndSpan(commandSpan, err)
	}()

	fs := flag.NewFlagSet("reembed-changed", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	root := fs.String("root", "memory", "memory root path")
	repoRoot := fs.String("repo-root", ".", "repository root used to resolve changed files")
	filesChanged := fs.String("files-changed", "", "comma-separated changed file paths")
	domain := fs.String("domain", "auto-crawled", "default domain for newly discovered markdown files")
	reviewer := fs.String("reviewer", "system", "reviewer identity for upsert policy")
	sessionID := fs.String("session-id", "", "embedding session identifier")
	embeddingEndpoint := fs.String("embedding-endpoint", retrieval.DefaultEmbeddingEndpoint, "embedding service endpoint")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if strings.TrimSpace(*filesChanged) == "" {
		fmt.Println("No changed files provided. Nothing to re-embed.")
		return nil
	}
	commandSpan.SetAttributes(
		attribute.String("memory.root", *root),
		attribute.String("repo.root", *repoRoot),
		attribute.String("session.id", *sessionID),
	)

	_, rootSpan := telemetry.StartSpan(ctx, "memory.reembed.resolve_paths")
	absRepoRoot, err := filepath.Abs(*repoRoot)
	if err != nil {
		telemetry.EndSpan(rootSpan, err)
		return err
	}
	absMemoryRoot, err := filepath.Abs(*root)
	if err != nil {
		telemetry.EndSpan(rootSpan, err)
		return err
	}
	telemetry.EndSpan(rootSpan, nil)

	existing := map[string]types.IndexEntry{}
	_, loadIndexSpan := telemetry.StartSpan(ctx, "memory.index.load")
	if idx, err := index.LoadIndex(*root); err == nil {
		for _, e := range idx.Entries {
			existing[e.ID] = e
		}
	}
	telemetry.EndSpan(loadIndexSpan, nil)

	policy := types.WritePolicyDecision{
		Decision: "approved",
		Reviewer: *reviewer,
		Reason:   "re-embed changed markdown files",
		Notes:    "observer-triggered consistency update",
		Risk:     "low",
	}

	seen := map[string]struct{}{}
	var targetIDs []string
	_, updateSpan := telemetry.StartSpan(ctx, "memory.reembed.upsert_changed")
	for _, raw := range strings.Split(*filesChanged, ",") {
		changed := strings.TrimSpace(raw)
		if changed == "" {
			continue
		}
		clean := filepath.Clean(changed)
		absPath := clean
		if !filepath.IsAbs(absPath) {
			absPath = filepath.Join(absRepoRoot, clean)
		}
		absPath, err = filepath.Abs(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", changed, err)
			continue
		}
		if strings.HasPrefix(absPath, absMemoryRoot+string(os.PathSeparator)) || absPath == absMemoryRoot {
			continue
		}
		info, statErr := os.Stat(absPath)
		if statErr != nil {
			fmt.Fprintf(os.Stderr, "warning: skipping %s: %v\n", changed, statErr)
			continue
		}
		if info.IsDir() || !strings.HasSuffix(strings.ToLower(info.Name()), ".md") {
			continue
		}

		entryID := buildCrawlEntryID(absRepoRoot, absPath)
		if _, ok := seen[entryID]; ok {
			continue
		}
		seen[entryID] = struct{}{}

		entryType := "instruction"
		entryDomain := *domain
		entryTitle := strings.Title(strings.ReplaceAll(entryID, "-", " "))
		if old, ok := existing[entryID]; ok {
			if strings.TrimSpace(old.Type) != "" {
				entryType = old.Type
			}
			if strings.TrimSpace(old.Domain) != "" {
				entryDomain = old.Domain
			}
			if strings.TrimSpace(old.Title) != "" {
				entryTitle = old.Title
			}
		}

		if err := index.UpsertEntry(*root, types.UpsertEntryInput{
			ID:       entryID,
			Title:    entryTitle,
			Type:     entryType,
			Domain:   entryDomain,
			BodyFile: absPath,
			Stage:    "pm",
		}, policy); err != nil {
			fmt.Fprintf(os.Stderr, "warning: failed to update %s from %s: %v\n", entryID, changed, err)
			continue
		}
		targetIDs = append(targetIDs, entryID)
	}
	telemetry.EndSpan(updateSpan, nil)

	if len(targetIDs) == 0 {
		fmt.Println("No changed markdown files eligible for re-embedding.")
		return nil
	}

	_, embedSpan := telemetry.StartSpan(ctx, "memory.embedding.batch_index")
	warnings, err := retrieval.IndexEntriesEmbeddingBatch(*root, targetIDs, *embeddingEndpoint, *sessionID)
	telemetry.EndSpan(embedSpan, err)
	if err != nil {
		return err
	}
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "warning: %s\n", w)
	}
	fmt.Printf("Re-embedded %d changed markdown entries.\n", len(targetIDs))
	return nil
}

func buildCrawlEntryID(crawlRoot, filePath string) string {
	rel, err := filepath.Rel(crawlRoot, filePath)
	if err != nil || strings.HasPrefix(rel, "..") {
		rel = filePath
	}
	rel = filepath.ToSlash(rel)
	base := strings.TrimSuffix(rel, filepath.Ext(rel))
	slug := slugify(base)
	if slug == "" {
		slug = "doc"
	}
	sum := sha1.Sum([]byte(rel))
	suffix := hex.EncodeToString(sum[:4])
	const maxSlugLen = 72
	if len(slug) > maxSlugLen {
		slug = slug[:maxSlugLen]
	}
	return slug + "-" + suffix
}

func slugify(s string) string {
	var b strings.Builder
	lastDash := false
	for _, r := range strings.ToLower(s) {
		isAlpha := r >= 'a' && r <= 'z'
		isDigit := r >= '0' && r <= '9'
		if isAlpha || isDigit {
			b.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash {
			b.WriteByte('-')
			lastDash = true
		}
	}
	out := strings.Trim(b.String(), "-")
	return out
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
