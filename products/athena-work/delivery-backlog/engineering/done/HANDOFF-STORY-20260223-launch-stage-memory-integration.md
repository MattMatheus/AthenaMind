# Engineering Handoff: STORY-20260223-launch-stage-memory-integration

## What Changed
- Integrated soft memory bootstrap into stage launching:
  - `tools/launch_stage.sh` now attempts `memory-cli bootstrap` for each launched stage (`engineering`, `qa`, `pm`, `planning`, `architect`, `cycle`).
  - Successful bootstrap payload is appended to launcher output under `memory_bootstrap_context:`.
  - Missing or failing `memory-cli` now emits warning-only behavior and does not block stage launch.
- Integrated soft episode write-back into observer cycle closure:
  - `tools/run_observer_cycle.sh` now attempts `memory-cli episode write` after observer report generation.
  - Episode write includes cycle metadata (`cycle-id`, derived `story-id`, summary/decisions text, and changed-files CSV).
  - Missing or failing `memory-cli` now emits warning-only behavior and does not block observer report generation.
- Added regression tests for both integration points:
  - `tools/test_launch_stage_memory_integration.sh`
  - `tools/test_observer_cycle_memory_integration.sh`
- Added the new tests to canonical docs test runner:
  - `tools/run_doc_tests.sh`
- Synced workflow docs to reflect the new soft integration behavior:
  - `HUMANS.md`
  - `DEVELOPMENT_CYCLE.md`
- Synced queue/state artifacts after transition:
  - moved story to `delivery-backlog/engineering/qa/`
  - updated `delivery-backlog/engineering/active/README.md`
  - updated `product-research/roadmap/PROGRAM_STATE_BOARD.md` queue counts and `Now` pointer

## Why It Changed
- The story requires bootstrap/write-back to be wired into normal stage execution so memory context is consumed automatically and cycle outcomes are persisted automatically, without creating a hard runtime dependency on `memory-cli`.

## Test Updates Made
- Added `tools/test_launch_stage_memory_integration.sh` with coverage for:
  - bootstrap invoked and appended when `memory-cli` is available
  - warning-only fallback when `memory-cli` is unavailable
- Added `tools/test_observer_cycle_memory_integration.sh` with coverage for:
  - episode write invoked with cycle metadata when `memory-cli` is available
  - warning-only fallback when `memory-cli` is unavailable

## Test Run Results
- `tools/test_launch_stage_memory_integration.sh` passed.
- `tools/test_observer_cycle_memory_integration.sh` passed.
- `tools/run_doc_tests.sh` passed.
- `go test ./...` passed.

## Open Risks / Questions
- `internal/retrieval/bootstrap.go` currently reads episode context from `episodes/<repo>/<scenario>/latest.json`, while `episode write` currently writes `episodes/<repo>/latest.json`; this pre-existing path mismatch can limit bootstrap episode context visibility until harmonized in follow-up work.

## Recommended QA Focus Areas
- Verify launch output contains `memory_bootstrap_context:` only when bootstrap succeeds and warns (without failure) when unavailable.
- Verify observer still writes report when `memory-cli` is unavailable or episode write fails.
- Verify `memory-cli episode write` receives expected cycle metadata (`cycle-id`, `story-id`, summary/decisions/files).
- Verify no regressions in queue ordering and stage launcher behavior.

## New Gaps Discovered
- Consider follow-up engineering bug to align episode latest-context path between bootstrap read (`episodes/<repo>/<scenario>/latest.json`) and episode write (`episodes/<repo>/latest.json`).
