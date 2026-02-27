# Story: Add CLI adapter for launcher and observer compatibility

## Metadata
- `id`: STORY-20260227-cli-adapter-for-launcher-observer
- `owner_persona`: staff-personas/Software Engineer - Max.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0004, ADR-0013]
- `success_metric`: Existing `launch_stage.sh` and `run_observer_cycle.sh` workflows operate via API adapter with no regression in stage outputs.
- `release_checkpoint`: required

## Problem Statement
Operators rely on current scripts. A backend migration must preserve this interface to avoid disrupting daily workflow.

## Scope
- In:
  - adapter calls from launcher/observer scripts to workspace API
  - fallback mode when API unavailable
  - direction confirmation prompts/flags aligned with API `confirm_direction`
- Out:
  - full script replacement

## Acceptance Criteria
1. Stage launcher and observer flows function with backend integration enabled.
2. Direction-changing actions require explicit human confirmation path in script workflow.
3. Script output remains concise and backward-compatible for operator usage.
4. Agent-facing output remains deterministic and minimal (`status`, `action`, `why`, `next`) to reduce parsing overhead and prompt/context bloat.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-workspace-api-state-machine-v1

## Notes
- Preserve existing exit code expectations for automation.
