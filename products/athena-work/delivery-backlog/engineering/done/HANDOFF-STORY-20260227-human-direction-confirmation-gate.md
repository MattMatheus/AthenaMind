# Engineering Handoff: STORY-20260227-human-direction-confirmation-gate

## What Changed
- Enforced a full direction confirmation data model in workflow adapter:
  - `confirmed_by`
  - `confirmed_at`
  - `scope`
  - `expiry`
  - deterministic rejection for missing/expired/superseded confirmations
  - `products/athena-work/tools/lib/workspace_api_adapter.sh`
- Added direction confirmation status visibility in CLI outputs and planning summary:
  - direction confirmation status block for direction-changing operations
  - one-screen planning direction summary card
  - `products/athena-work/tools/launch_stage.sh`
- Added direction confirmation evidence export to observer markdown report:
  - `products/athena-work/tools/run_observer_cycle.sh`
- Updated control-plane contract and stage prompt guidance:
  - `products/athena-work/operating-system/contracts/WORKSPACE_API_STATE_MACHINE_V1.md`
  - `products/athena-work/stage-prompts/active/planning-seed-prompt.md`
- Expanded test coverage for new gate semantics and visibility:
  - `products/athena-work/tools/test_launch_stage_workspace_api_adapter.sh`
  - `products/athena-work/tools/test_observer_workspace_api_adapter.sh`
  - `products/athena-work/tools/test_workspace_api_state_machine_v1.sh`
  - `products/athena-work/tools/test_stage_exit_pipeline.sh`

## Why It Changed
- Prevent direction-changing operations from running on incomplete, stale, or superseded human confirmations.
- Provide high-clarity confirmation visibility for both CLI workflows and markdown evidence consumers.

## Test Updates Made
- Added rejection-path assertions:
  - `ERR_CONFIRM_DIRECTION_REQUIRED`
  - `ERR_CONFIRM_DIRECTION_EXPIRED`
  - `ERR_CONFIRM_DIRECTION_SUPERSEDED`
- Added acceptance-path assertions for valid confirmation state + output visibility.
- Added planning summary-card checks and observer report evidence checks.

## Test Run Results
- `tools/run_doc_tests.sh` -> PASS
- `go test ./...` (repo root) -> PASS

## Open Risks/Questions
- Expiry check assumes ISO8601 UTC lexical ordering in shell environment.
- API-side persistence/validation remains required for full runtime parity.

## Recommended QA Focus Areas
- Validate direction confirmation status output remains stable for UI/agent parsers.
- Validate report export fields for all direction confirmation states.
- Validate planning summary readability with low-vision defaults.

## New Gaps Discovered
- none
