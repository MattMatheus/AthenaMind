# Handoff: STORY-20260222-memory-cli-telemetry-contract-v01

## What Changed
- Added local-first telemetry event emission for all CLI flows:
  - `write` -> `memory.write`
  - `retrieve` -> `memory.retrieve`
  - `evaluate` -> `memory.evaluate`
- Added telemetry flags to command paths:
  - `--session-id`
  - `--scenario-id`
  - `--memory-type`
  - `--operator-verdict`
  - `--telemetry-file` (optional; defaults to `<root>/telemetry/events.jsonl`)
- Added schema-oriented telemetry event model with KPI-required identifiers and retrieval evidence fields.
- Added failure-path telemetry coverage with deterministic error codes.
- Added telemetry contract test coverage:
  - `TestTelemetryEmissionForWriteRetrieveEvaluateAndFailure`

## Why It Changed
- KPI and governance reporting needed reusable, normalized event payloads emitted directly from memory CLI flows without introducing external telemetry dependencies.

## Test Updates Made
- Added telemetry emission test in `cmd/memory-cli/main_test.go` validating:
  - event emission across write/retrieve/evaluate flows
  - required field completeness
  - retrieval evidence fields (`selected_id`, `selection_mode`, `source_path`)
  - failure-path emission with `error_code`

## Test Run Results
- Command: `go test ./cmd/memory-cli`
- Result: PASS
- Command: `scripts/run_doc_tests.sh`
- Result: PASS

## Acceptance Criteria Trace
- AC1 (schema-compliant events for success and failure): PASS
  - success events emitted for write/retrieve/evaluate; fail event emitted for policy-blocked write.
- AC2 (KPI-required and retrieval evidence fields): PASS
  - event payload includes required identifiers and retrieve evidence fields.
- AC3 (deterministic emission + required-field tests): PASS
  - telemetry tests assert event count, required field presence, and failure error code population.

## Open Risks/Questions
- Event validation is emitter-side by convention; strict schema validator command is still a follow-on hardening option.

## Recommended QA Focus Areas
- Verify retrieval failure path emits complete required retrieval evidence placeholders.
- Verify telemetry file path behavior when root is custom.
- Verify added flags remain backward-compatible for existing CLI usage.

## New Gaps Discovered During Implementation
- None.
