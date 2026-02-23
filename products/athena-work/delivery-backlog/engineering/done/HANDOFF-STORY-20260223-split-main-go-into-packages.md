# Engineering Handoff: STORY-20260223-split-main-go-into-packages

## What Changed
- Refactored CLI domain logic into internal packages:
  - `internal/types`
  - `internal/index`
  - `internal/governance`
  - `internal/retrieval`
  - `internal/telemetry`
  - `internal/snapshot`
  - `internal/gateway`
- Reduced `cmd/memory-cli/main.go` to subcommand wiring only.
- Moved command execution and flag parsing into `cmd/memory-cli/commands.go`.
- Added compatibility wrappers/type aliases in `cmd/memory-cli/*` so existing CLI tests and symbols remain stable.

## Why It Changed
- Story requires decomposing monolithic CLI logic into module boundaries aligned to ADR-0003/ADR-0009 while preserving behavior.
- This unblocks future memory/bootstrap/episode work by reducing coupling and improving testability.

## Test Updates Made
- Added independent unit tests for each extracted package:
  - `internal/index/index_test.go`
  - `internal/governance/governance_test.go`
  - `internal/retrieval/retrieval_test.go`
  - `internal/telemetry/telemetry_test.go`
  - `internal/snapshot/snapshot_test.go`
  - `internal/gateway/gateway_test.go`

## Test Run Results
- `go test ./...` passed.
- `tools/run_doc_tests.sh` passed.

## Open Risks / Questions
- `runWrite` output path formatting is preserved semantically; QA should spot-check command output text parity on prompt/instruction writes.

## Recommended QA Focus Areas
- CLI behavior parity for:
  - `write`, `retrieve`, `evaluate`, `snapshot`, `api-retrieve`, `serve-read-gateway`
- Telemetry event fields and error code behavior parity.
- Snapshot create/list/restore compatibility checks.
- Retrieval fallback determinism and API parity checks.

## New Gaps Discovered
- None.
