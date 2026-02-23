package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

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
