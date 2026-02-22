package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestWriteAndRetrieveSemantic(t *testing.T) {
	root := t.TempDir()

	err := runWrite([]string{
		"--root", root,
		"--id", "getting-started",
		"--title", "Getting Started Prompt",
		"--type", "prompt",
		"--domain", "onboarding",
		"--body", "Use this prompt to onboard a new engineer quickly",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "improve onboarding speed",
		"--risk", "low risk with rollback by git revert",
		"--notes", "approved in planning review",
	})
	if err != nil {
		t.Fatalf("runWrite failed: %v", err)
	}

	if _, err := os.Stat(filepath.Join(root, "index.yaml")); err != nil {
		t.Fatalf("index.yaml missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "prompts", "onboarding", "getting-started.md")); err != nil {
		t.Fatalf("content markdown missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "metadata", "getting-started.yaml")); err != nil {
		t.Fatalf("metadata missing: %v", err)
	}

	result, err := retrieve(root, "onboard engineer prompt", "")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	if result.SelectionMode == "" {
		t.Fatal("selection_mode should be populated")
	}
	if result.SelectedID == "" || result.SourcePath == "" {
		t.Fatal("selected_id and source_path should be populated")
	}
}

func TestWriteFailsDuringAutonomousRun(t *testing.T) {
	t.Setenv("AUTONOMOUS_RUN", "true")
	root := t.TempDir()
	err := runWrite([]string{
		"--root", root,
		"--id", "blocked",
		"--title", "Blocked",
		"--type", "prompt",
		"--domain", "security",
		"--body", "blocked",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "security policy",
		"--risk", "low",
		"--notes", "blocked in autonomous mode",
	})
	if err == nil {
		t.Fatal("expected write to fail during autonomous run")
	}
}

func TestWriteRejectsWithoutEvidence(t *testing.T) {
	root := t.TempDir()
	err := runWrite([]string{
		"--root", root,
		"--id", "missing-evidence",
		"--title", "Missing Evidence",
		"--type", "prompt",
		"--domain", "security",
		"--body", "blocked",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
	})
	if err == nil {
		t.Fatal("expected evidence enforcement error")
	}
}

func TestWriteRejectedDecisionBlocksApplyAndCreatesAudit(t *testing.T) {
	root := t.TempDir()
	err := runWrite([]string{
		"--root", root,
		"--id", "rejected-change",
		"--title", "Rejected Change",
		"--type", "prompt",
		"--domain", "security",
		"--body", "blocked",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "rejected",
		"--reason", "needs redesign",
		"--risk", "high rollback risk",
		"--notes", "rejected pending rework",
		"--rework-notes", "provide deterministic rollback plan",
		"--re-reviewed-by", "maya",
	})
	if err == nil {
		t.Fatal("expected rejected decision to block mutation")
	}
	if _, statErr := os.Stat(filepath.Join(root, "index.yaml")); !os.IsNotExist(statErr) {
		t.Fatalf("index.yaml should not exist after rejected decision, stat err: %v", statErr)
	}
	matches, globErr := filepath.Glob(filepath.Join(root, "audits", "rejected-change-*.json"))
	if globErr != nil {
		t.Fatalf("glob failed: %v", globErr)
	}
	if len(matches) != 1 {
		t.Fatalf("expected exactly one audit artifact, got %d", len(matches))
	}
}

func TestLoadIndexRejectsUnsupportedMajor(t *testing.T) {
	root := t.TempDir()
	idx := indexFile{
		SchemaVersion: "2.0",
		UpdatedAt:     "2026-02-22T00:00:00Z",
		Entries:       []indexEntry{},
	}
	data, _ := json.Marshal(idx)
	if err := os.WriteFile(filepath.Join(root, "index.yaml"), data, 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	_, err := loadIndex(root)
	if err == nil {
		t.Fatal("expected unsupported major schema to fail")
	}
}

func TestLoadIndexAllowsNewerMinorCompatibility(t *testing.T) {
	root := t.TempDir()
	now := time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(filepath.Join(root, "prompts", "ops"), 0o755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "metadata"), 0o755); err != nil {
		t.Fatalf("mkdir metadata: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "prompts", "ops", "entry.md"), []byte("# Entry\n\nbody\n"), 0o644); err != nil {
		t.Fatalf("write prompt: %v", err)
	}
	meta := metadataFile{
		SchemaVersion: "1.1",
		ID:            "entry",
		Title:         "Entry",
		Status:        "approved",
		UpdatedAt:     now,
		Review: reviewMeta{
			ReviewedBy:   "maya",
			ReviewedAt:   now,
			Decision:     "approved",
			DecisionNote: "compat mode",
		},
	}
	metaData, _ := json.Marshal(meta)
	if err := os.WriteFile(filepath.Join(root, "metadata", "entry.yaml"), metaData, 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}
	idx := indexFile{
		SchemaVersion: "1.1",
		UpdatedAt:     now,
		Entries: []indexEntry{{
			ID:           "entry",
			Type:         "prompt",
			Domain:       "ops",
			Path:         "prompts/ops/entry.md",
			MetadataPath: "metadata/entry.yaml",
			Status:       "approved",
			UpdatedAt:    now,
			Title:        "Entry",
		}},
	}
	idxData, _ := json.Marshal(idx)
	if err := os.WriteFile(filepath.Join(root, "index.yaml"), idxData, 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}
	if _, err := loadIndex(root); err != nil {
		t.Fatalf("expected newer minor compatibility mode, got error: %v", err)
	}
}

func TestSchemaValidationRejectsInvalidMetadata(t *testing.T) {
	root := t.TempDir()
	now := time.Now().UTC().Format(time.RFC3339)
	if err := os.MkdirAll(filepath.Join(root, "prompts", "ops"), 0o755); err != nil {
		t.Fatalf("mkdir prompts: %v", err)
	}
	if err := os.MkdirAll(filepath.Join(root, "metadata"), 0o755); err != nil {
		t.Fatalf("mkdir metadata: %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, "prompts", "ops", "entry.md"), []byte("# Entry\n\nbody\n"), 0o644); err != nil {
		t.Fatalf("write prompt: %v", err)
	}
	meta := metadataFile{
		SchemaVersion: "1.0",
		ID:            "entry",
		Title:         "Entry",
		Status:        "approved",
		UpdatedAt:     now,
		Review: reviewMeta{
			ReviewedBy:   "maya",
			ReviewedAt:   "",
			Decision:     "approved",
			DecisionNote: "missing reviewed_at",
		},
	}
	metaData, _ := json.Marshal(meta)
	if err := os.WriteFile(filepath.Join(root, "metadata", "entry.yaml"), metaData, 0o644); err != nil {
		t.Fatalf("write metadata: %v", err)
	}
	idx := indexFile{
		SchemaVersion: "1.0",
		UpdatedAt:     now,
		Entries: []indexEntry{{
			ID:           "entry",
			Type:         "prompt",
			Domain:       "ops",
			Path:         "prompts/ops/entry.md",
			MetadataPath: "metadata/entry.yaml",
			Status:       "approved",
			UpdatedAt:    now,
			Title:        "Entry",
		}},
	}
	idxData, _ := json.Marshal(idx)
	if err := os.WriteFile(filepath.Join(root, "index.yaml"), idxData, 0o644); err != nil {
		t.Fatalf("write index: %v", err)
	}

	_, err := retrieve(root, "entry", "")
	if err == nil {
		t.Fatal("expected metadata schema validation failure")
	}
}

func TestFallbackDeterminismAndMetadataCompleteness(t *testing.T) {
	root := t.TempDir()
	for _, tc := range []struct {
		id     string
		title  string
		domain string
		body   string
	}{
		{id: "alpha", title: "Alpha", domain: "ops", body: "shared terms"},
		{id: "beta", title: "Beta", domain: "ops", body: "shared terms"},
	} {
		err := runWrite([]string{
			"--root", root,
			"--id", tc.id,
			"--title", tc.title,
			"--type", "prompt",
			"--domain", tc.domain,
			"--body", tc.body,
			"--stage", "planning",
			"--reviewer", "maya",
			"--decision", "approved",
			"--reason", "seed corpus",
			"--risk", "low",
			"--notes", "approved",
		})
		if err != nil {
			t.Fatalf("runWrite failed for %s: %v", tc.id, err)
		}
	}

	first, err := retrieve(root, "non matching query", "")
	if err != nil {
		t.Fatalf("retrieve failed: %v", err)
	}
	for i := 0; i < 5; i++ {
		again, err := retrieve(root, "non matching query", "")
		if err != nil {
			t.Fatalf("retrieve retry failed: %v", err)
		}
		if again.SelectedID != first.SelectedID || again.SelectionMode != first.SelectionMode || again.SourcePath != first.SourcePath {
			t.Fatal("fallback result should be deterministic across runs")
		}
		if again.SelectionMode == "" || again.SelectedID == "" || again.SourcePath == "" {
			t.Fatal("response metadata should always include mode, id, and source path")
		}
	}
}

func TestEvaluateRetrievalHarnessThresholds(t *testing.T) {
	root := t.TempDir()
	queries := make([]evaluationQuery, 0, 50)
	for i := 0; i < 50; i++ {
		id := fmt.Sprintf("q-entry-%02d", i)
		err := runWrite([]string{
			"--root", root,
			"--id", id,
			"--title", "Title " + id,
			"--type", "prompt",
			"--domain", "bench",
			"--body", "query token " + id,
			"--stage", "planning",
			"--reviewer", "maya",
			"--decision", "approved",
			"--reason", "seed benchmark",
			"--risk", "low",
			"--notes", "approved",
		})
		if err != nil {
			t.Fatalf("seed write failed: %v", err)
		}
		queries = append(queries, evaluationQuery{Query: id, Domain: "bench", ExpectedID: id})
	}

	report, err := evaluateRetrieval(root, queries, "corpus-v1", "query-set-v1", "config-v1")
	if err != nil {
		t.Fatalf("evaluateRetrieval failed: %v", err)
	}
	if report.SelectionModeReporting.Rate != 1 {
		t.Fatalf("expected 100%% selection mode reporting, got %f", report.SelectionModeReporting.Rate)
	}
	if report.SourceTraceCompleteness.Rate != 1 {
		t.Fatalf("expected 100%% source trace completeness, got %f", report.SourceTraceCompleteness.Rate)
	}
	if report.FallbackDeterminism.Rate != 1 {
		t.Fatalf("expected 100%% fallback determinism, got %f", report.FallbackDeterminism.Rate)
	}
}

func TestSemanticConfidenceGate(t *testing.T) {
	if isSemanticConfident(0.90, 0.82) {
		t.Fatal("expected low margin to fail confidence gate")
	}
	if !isSemanticConfident(0.90, 0.60) {
		t.Fatal("expected clear top score to pass confidence gate")
	}
}

func TestAPIRetrieveParityWithCLI(t *testing.T) {
	root := t.TempDir()
	err := runWrite([]string{
		"--root", root,
		"--id", "api-parity-entry",
		"--title", "API Parity Entry",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "Retrieve this entry for parity testing",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "seed parity scenario",
		"--risk", "low",
		"--notes", "approved",
	})
	if err != nil {
		t.Fatalf("seed write failed: %v", err)
	}

	server := httptest.NewServer(readGatewayHandler(root))
	defer server.Close()

	req := apiRetrieveRequest{
		Query:     "api parity entry",
		Domain:    "ops",
		SessionID: "session-parity",
	}
	got, err := apiRetrieveWithFallback(root, server.URL, req, "trace-parity", server.Client())
	if err != nil {
		t.Fatalf("api retrieve with gateway failed: %v", err)
	}
	want, err := retrieve(root, req.Query, req.Domain)
	if err != nil {
		t.Fatalf("cli retrieve failed: %v", err)
	}

	if got.SelectedID != want.SelectedID || got.SelectionMode != want.SelectionMode || got.SourcePath != want.SourcePath {
		t.Fatalf("expected parity with CLI, got=%+v want=%+v", got, want)
	}
	if got.FallbackUsed {
		t.Fatal("did not expect fallback for healthy gateway path")
	}
	if got.TraceID == "" {
		t.Fatal("expected trace_id in gateway response")
	}
}

func TestAPIRetrieveFallbackWhenGatewayUnavailable(t *testing.T) {
	root := t.TempDir()
	err := runWrite([]string{
		"--root", root,
		"--id", "api-fallback-entry",
		"--title", "API Fallback Entry",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "Retrieve this entry through fallback path",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "seed fallback scenario",
		"--risk", "low",
		"--notes", "approved",
	})
	if err != nil {
		t.Fatalf("seed write failed: %v", err)
	}

	client := &http.Client{Timeout: 200 * time.Millisecond}
	req := apiRetrieveRequest{
		Query:     "api fallback entry",
		Domain:    "ops",
		SessionID: "session-fallback",
	}
	got, err := apiRetrieveWithFallback(root, "http://127.0.0.1:1", req, "trace-fallback", client)
	if err != nil {
		t.Fatalf("api retrieve fallback failed: %v", err)
	}
	want, err := retrieve(root, req.Query, req.Domain)
	if err != nil {
		t.Fatalf("cli retrieve failed: %v", err)
	}

	if !got.FallbackUsed {
		t.Fatal("expected fallback to be used when gateway is unavailable")
	}
	if got.FallbackCode != "ERR_API_WRAPPER_UNAVAILABLE" {
		t.Fatalf("expected fallback code ERR_API_WRAPPER_UNAVAILABLE, got %s", got.FallbackCode)
	}
	if got.SelectedID != want.SelectedID || got.SelectionMode != want.SelectionMode || got.SourcePath != want.SourcePath {
		t.Fatalf("fallback output must match CLI retrieve semantics, got=%+v want=%+v", got, want)
	}
}

func TestTelemetryEmissionForWriteRetrieveEvaluateAndFailure(t *testing.T) {
	root := t.TempDir()
	telemetryPath := filepath.Join(root, "events.jsonl")

	if err := runWrite([]string{
		"--root", root,
		"--id", "telemetry-entry",
		"--title", "Telemetry Entry",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "telemetry body",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "seed telemetry flow",
		"--risk", "low",
		"--notes", "approved",
		"--session-id", "sess-1",
		"--scenario-id", "scn-1",
		"--memory-type", "procedural",
		"--telemetry-file", telemetryPath,
	}); err != nil {
		t.Fatalf("runWrite failed: %v", err)
	}

	if err := runRetrieve([]string{
		"--root", root,
		"--query", "telemetry entry",
		"--domain", "ops",
		"--session-id", "sess-1",
		"--scenario-id", "scn-1",
		"--memory-type", "semantic",
		"--telemetry-file", telemetryPath,
	}); err != nil {
		t.Fatalf("runRetrieve failed: %v", err)
	}

	querySetPath := filepath.Join(root, "queries.json")
	querySet := []evaluationQuery{{
		Query:      "telemetry-entry",
		Domain:     "ops",
		ExpectedID: "telemetry-entry",
	}}
	data, _ := json.Marshal(querySet)
	if err := os.WriteFile(querySetPath, data, 0o644); err != nil {
		t.Fatalf("write query set: %v", err)
	}
	if err := runEvaluate([]string{
		"--root", root,
		"--query-file", querySetPath,
		"--session-id", "sess-1",
		"--scenario-id", "scn-1",
		"--memory-type", "state",
		"--telemetry-file", telemetryPath,
	}); err != nil {
		t.Fatalf("runEvaluate failed: %v", err)
	}

	t.Setenv("AUTONOMOUS_RUN", "true")
	if err := runWrite([]string{
		"--root", root,
		"--id", "blocked-telemetry",
		"--title", "Blocked",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "blocked",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "blocked in autonomous mode",
		"--risk", "low",
		"--notes", "blocked",
		"--session-id", "sess-2",
		"--scenario-id", "scn-2",
		"--memory-type", "procedural",
		"--telemetry-file", telemetryPath,
	}); err == nil {
		t.Fatal("expected autonomous write to fail")
	}

	events := readTelemetryEvents(t, telemetryPath)
	if len(events) != 4 {
		t.Fatalf("expected 4 telemetry events, got %d", len(events))
	}

	seenFailure := false
	seenRetrieve := false
	for _, ev := range events {
		if ev.EventName == "" || ev.EventVersion == "" || ev.TimestampUTC == "" || ev.SessionID == "" || ev.TraceID == "" || ev.ScenarioID == "" || ev.Operation == "" || ev.Result == "" || ev.PolicyGate == "" || ev.MemoryType == "" || ev.OperatorVerdict == "" {
			t.Fatalf("telemetry event missing required fields: %+v", ev)
		}
		if ev.LatencyMS < 0 {
			t.Fatalf("telemetry latency must be non-negative: %+v", ev)
		}
		if ev.Operation == "retrieve" {
			seenRetrieve = true
			if ev.SelectedID == "" || ev.SelectionMode == "" || ev.SourcePath == "" {
				t.Fatalf("retrieve telemetry missing required retrieval fields: %+v", ev)
			}
		}
		if ev.Result == "fail" {
			seenFailure = true
			if ev.ErrorCode == "" {
				t.Fatalf("failed telemetry event must include error_code: %+v", ev)
			}
		}
	}
	if !seenRetrieve {
		t.Fatal("expected at least one retrieve telemetry event")
	}
	if !seenFailure {
		t.Fatal("expected at least one failed telemetry event")
	}
}

func readTelemetryEvents(t *testing.T, path string) []telemetryEvent {
	t.Helper()

	file, err := os.Open(path)
	if err != nil {
		t.Fatalf("open telemetry file: %v", err)
	}
	defer file.Close()

	events := []telemetryEvent{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev telemetryEvent
		if err := json.Unmarshal(line, &ev); err != nil {
			t.Fatalf("parse telemetry line: %v", err)
		}
		events = append(events, ev)
	}
	if err := scanner.Err(); err != nil {
		t.Fatalf("scan telemetry file: %v", err)
	}
	return events
}

func TestSnapshotCreateListRestoreLifecycle(t *testing.T) {
	root := t.TempDir()
	if err := runWrite([]string{
		"--root", root,
		"--id", "snap-base",
		"--title", "Snapshot Base",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "base",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "seed",
		"--risk", "low",
		"--notes", "approved",
	}); err != nil {
		t.Fatalf("seed write failed: %v", err)
	}

	if err := runSnapshotCreate([]string{
		"--root", root,
		"--created-by", "tester",
		"--reason", "checkpoint",
		"--session-id", "sess-snap",
	}); err != nil {
		t.Fatalf("snapshot create failed: %v", err)
	}

	rows, err := listSnapshots(root)
	if err != nil {
		t.Fatalf("snapshot list failed: %v", err)
	}
	if len(rows) != 1 {
		t.Fatalf("expected 1 snapshot, got %d", len(rows))
	}
	snapshotID := rows[0].SnapshotID

	if err := runWrite([]string{
		"--root", root,
		"--id", "snap-new",
		"--title", "Snapshot New",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "new data",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "mutate",
		"--risk", "low",
		"--notes", "approved",
	}); err != nil {
		t.Fatalf("second write failed: %v", err)
	}
	idx, err := loadIndex(root)
	if err != nil {
		t.Fatalf("load index after mutation: %v", err)
	}
	if len(idx.Entries) != 2 {
		t.Fatalf("expected 2 entries after mutation, got %d", len(idx.Entries))
	}

	if err := runSnapshotRestore([]string{
		"--root", root,
		"--snapshot-id", snapshotID,
		"--session-id", "sess-snap",
		"--stage", "pm",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "rollback",
		"--risk", "low",
		"--notes", "approved restore",
	}); err != nil {
		t.Fatalf("snapshot restore failed: %v", err)
	}

	idx, err = loadIndex(root)
	if err != nil {
		t.Fatalf("load index after restore: %v", err)
	}
	if len(idx.Entries) != 1 {
		t.Fatalf("expected 1 entry after restore, got %d", len(idx.Entries))
	}
}

func TestSnapshotRestoreRejectsIncompatibleManifest(t *testing.T) {
	root := t.TempDir()
	if err := runWrite([]string{
		"--root", root,
		"--id", "snap-compat",
		"--title", "Snapshot Compat",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "compat",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "seed",
		"--risk", "low",
		"--notes", "approved",
	}); err != nil {
		t.Fatalf("seed write failed: %v", err)
	}
	if err := runSnapshotCreate([]string{
		"--root", root,
		"--created-by", "tester",
		"--reason", "checkpoint",
	}); err != nil {
		t.Fatalf("snapshot create failed: %v", err)
	}
	rows, err := listSnapshots(root)
	if err != nil || len(rows) != 1 {
		t.Fatalf("snapshot list failed: %v len=%d", err, len(rows))
	}
	manifestPath := filepath.Join(root, "snapshots", rows[0].SnapshotID, "manifest.json")
	m, err := loadSnapshotManifest(root, rows[0].SnapshotID)
	if err != nil {
		t.Fatalf("load manifest: %v", err)
	}
	m.SchemaVersion = "2.0"
	data, _ := json.MarshalIndent(m, "", "  ")
	if err := os.WriteFile(manifestPath, append(data, '\n'), 0o644); err != nil {
		t.Fatalf("write manifest: %v", err)
	}

	err = runSnapshotRestore([]string{
		"--root", root,
		"--snapshot-id", rows[0].SnapshotID,
		"--stage", "pm",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "restore",
		"--risk", "low",
		"--notes", "approved",
	})
	if err == nil || !strings.Contains(err.Error(), "ERR_SNAPSHOT_COMPATIBILITY_BLOCKED") {
		t.Fatalf("expected compatibility blocked error, got %v", err)
	}
}

func TestSnapshotAuditEventChain(t *testing.T) {
	root := t.TempDir()
	if err := runWrite([]string{
		"--root", root,
		"--id", "snap-audit",
		"--title", "Snapshot Audit",
		"--type", "prompt",
		"--domain", "ops",
		"--body", "audit",
		"--stage", "planning",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "seed",
		"--risk", "low",
		"--notes", "approved",
	}); err != nil {
		t.Fatalf("seed write failed: %v", err)
	}
	if err := runSnapshotCreate([]string{
		"--root", root,
		"--created-by", "tester",
		"--reason", "checkpoint",
		"--session-id", "sess-audit",
	}); err != nil {
		t.Fatalf("snapshot create failed: %v", err)
	}
	rows, _ := listSnapshots(root)
	if err := runSnapshotRestore([]string{
		"--root", root,
		"--snapshot-id", rows[0].SnapshotID,
		"--session-id", "sess-audit",
		"--stage", "pm",
		"--reviewer", "maya",
		"--decision", "approved",
		"--reason", "restore",
		"--risk", "low",
		"--notes", "approved",
	}); err != nil {
		t.Fatalf("snapshot restore failed: %v", err)
	}

	matches, err := filepath.Glob(filepath.Join(root, "audits", "*.json"))
	if err != nil {
		t.Fatalf("glob audits failed: %v", err)
	}
	found := map[string]bool{}
	for _, p := range matches {
		data, _ := os.ReadFile(p)
		var ev snapshotAuditEvent
		if json.Unmarshal(data, &ev) == nil && ev.EventName != "" {
			found[ev.EventName] = true
		}
	}
	required := []string{
		"snapshot.create.requested",
		"snapshot.create.completed",
		"snapshot.restore.requested",
		"snapshot.restore.policy_decision",
		"snapshot.restore.completed",
	}
	for _, name := range required {
		if !found[name] {
			t.Fatalf("missing required snapshot audit event: %s", name)
		}
	}
}
