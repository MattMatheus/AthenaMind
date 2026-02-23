package main

import (
	"athenamind/internal/governance"
	"athenamind/internal/index"
	"athenamind/internal/telemetry"
)

func emitTelemetry(root, telemetryPath string, ev telemetryEvent) error {
	return telemetry.Emit(root, telemetryPath, ev)
}

func telemetryErrorCode(err error) string {
	return telemetry.TelemetryErrorCode(err)
}

func normalizeMemoryType(v string) string {
	return telemetry.NormalizeMemoryType(v)
}

func normalizeOperatorVerdict(v string) string {
	return telemetry.NormalizeOperatorVerdict(v)
}

func normalizeTelemetryValue(v, fallback string) string {
	return telemetry.NormalizeTelemetryValue(v, fallback)
}

func enforceConstraintChecks(operation, sessionID, scenarioID, traceID string) error {
	return governance.EnforceConstraintChecks(operation, sessionID, scenarioID, traceID)
}

func enforceWritePolicy(in writePolicyInput) (writePolicyDecision, error) {
	return governance.EnforceWritePolicy(in)
}

func isLatencyDegraded(elapsedMs int64) bool {
	return governance.IsLatencyDegraded(elapsedMs)
}

func loadIndex(root string) (indexFile, error) {
	return index.LoadIndex(root)
}

func validateSchemaVersion(version string) error {
	return index.ValidateSchemaVersion(version)
}

func parseMajorMinor(v string) (int, int, error) {
	return index.ParseMajorMinor(v)
}

func isValidStatus(s string) bool {
	return index.IsValidStatus(s)
}

func writeJSONAsYAML(path string, v any) error {
	return index.WriteJSONAsYAML(path, v)
}
