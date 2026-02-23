# Handoff: STORY-20260222-memory-api-read-gateway-v03

## What Changed
- Added read-only API gateway server path:
  - `memory-cli serve-read-gateway`
  - endpoint: `POST /memory/retrieve`
- Added API client retrieval path with explicit fallback behavior:
  - `memory-cli api-retrieve`
  - falls back to CLI retrieve with `ERR_API_WRAPPER_UNAVAILABLE` when gateway path is unreachable.
- Added parity enforcement between gateway response and CLI retrieve contract:
  - mismatch returns `ERR_API_CLI_PARITY_MISMATCH`.
- Added API request/response structs with rationale and trace metadata:
  - includes `selected_id`, `selection_mode`, `source_path`, `confidence`, `reason`, `trace_id`.
- Added contract tests:
  - `TestAPIRetrieveParityWithCLI`
  - `TestAPIRetrieveFallbackWhenGatewayUnavailable`

## Why It Changed
- v0.3 requires a reversible API read wrapper that preserves CLI semantics and deterministic fallback behavior.

## Test Updates Made
- Updated `cmd/memory-cli/main_test.go` with API/CLI parity and fallback contract tests.

## Test Run Results
- Command: `go test ./cmd/memory-cli`
- Result: PASS
- Command: `tools/run_doc_tests.sh`
- Result: PASS

## Acceptance Criteria Trace
- AC1 (schema-consistent retrieve with rationale metadata): PASS
  - API responses preserve CLI retrieval contract fields and include `trace_id` + `reason`.
- AC2 (contract tests for parity + fallback): PASS
  - Added and passing parity/fallback tests in `main_test.go`.
- AC3 (reversible CLI-only mode): PASS
  - `api-retrieve` supports direct CLI path and explicit fallback behavior without data migration.

## Open Risks/Questions
- Current gateway surface covers read path only; write/evaluate API parity remains future work.

## Recommended QA Focus Areas
- Verify parity check behavior on representative query sets.
- Verify fallback metadata values when gateway endpoint is unavailable.
- Verify API handler rejects malformed payloads with deterministic error semantics.

## New Gaps Discovered During Implementation
- None beyond planned follow-on API write/evaluate scope.
