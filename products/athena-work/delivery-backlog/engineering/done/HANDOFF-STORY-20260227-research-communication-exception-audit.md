# Engineering Handoff: STORY-20260227-research-communication-exception-audit

## What Changed
- Extended workspace API adapter with research exception policy enforcement + audit logging:
  - `products/athena-work/tools/lib/workspace_api_adapter.sh`
- Integrated research exception policy checks into stage launcher and observer flows:
  - `products/athena-work/tools/launch_stage.sh`
  - `products/athena-work/tools/run_observer_cycle.sh`
- Expanded adapter tests to cover:
  - rejection outside research lane
  - explicit flag requirement
  - required-field audit logging
  - `products/athena-work/tools/test_launch_stage_workspace_api_adapter.sh`
  - `products/athena-work/tools/test_observer_workspace_api_adapter.sh`
- Updated API state-machine contract and test coverage for required audit fields:
  - `products/athena-work/operating-system/contracts/WORKSPACE_API_STATE_MACHINE_V1.md`
  - `products/athena-work/tools/test_workspace_api_state_machine_v1.sh`

## Why It Changed
- Enforce research-only exception path for undocumented agent-to-agent communication.
- Guarantee auditable exception events with required metadata fields.
- Preserve CLI/API policy parity during control-plane migration.

## Test Updates Made
- Added research exception block/allow path checks to launcher adapter integration test.
- Added research exception block/allow path checks to observer adapter integration test.
- Added contract assertions for `source_agent`, `target_agent`, and `reason`.

## Test Run Results
- `tools/run_doc_tests.sh` -> PASS
- `go test ./...` (repo root) -> PASS

## Open Risks/Questions
- CLI audit log currently writes JSONL local file and does not call backend event API directly.
- Additional policy tests may be needed when runtime API handlers are introduced.

## Recommended QA Focus Areas
- Validate non-research hard block behavior remains stable across all stage invocations.
- Validate research exception log entries always contain required identifiers.
- Validate no regression in existing launcher and observer outputs for non-exception flows.

## New Gaps Discovered
- none
