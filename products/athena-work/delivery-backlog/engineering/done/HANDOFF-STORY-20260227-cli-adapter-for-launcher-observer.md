# Engineering Handoff: STORY-20260227-cli-adapter-for-launcher-observer

## What Changed
- Added workspace API adapter shell library:
  - `products/athena-work/tools/lib/workspace_api_adapter.sh`
- Integrated adapter status and direction-confirmation gate into launcher:
  - `products/athena-work/tools/launch_stage.sh`
- Integrated adapter status and direction-confirmation gate into observer:
  - `products/athena-work/tools/run_observer_cycle.sh`
- Added adapter-focused tests and wired them into doc test runner:
  - `products/athena-work/tools/test_launch_stage_workspace_api_adapter.sh`
  - `products/athena-work/tools/test_observer_workspace_api_adapter.sh`
  - `products/athena-work/tools/run_doc_tests.sh`

## Why It Changed
- Preserve current CLI workflow while enabling optional workspace API mode.
- Enforce explicit human confirmation for direction-changing actions via script gate.
- Keep agent-facing output deterministic and minimal in fallback or connected mode.

## Test Updates Made
- Added launcher adapter test coverage for:
  - fallback status output when API is unavailable
  - direction-change confirmation enforcement
- Added observer adapter test coverage for:
  - fallback status output when API is unavailable
  - direction-change confirmation enforcement

## Test Run Results
- `tools/run_doc_tests.sh` -> PASS
- `go test ./...` (repo root) -> PASS

## Open Risks/Questions
- `confirm_direction` API call currently soft-fails to warning when endpoint is unavailable; strict runtime enforcement remains a control-plane responsibility.
- Adapter health check currently assumes `/health` contract stability.

## Recommended QA Focus Areas
- Verify adapter output fields remain stable for parser/agent consumers.
- Verify direction-change gate behavior across all launch stages and observer path.
- Verify no behavior regression when `WORKSPACE_API_ENABLED` is unset.

## New Gaps Discovered
- none
