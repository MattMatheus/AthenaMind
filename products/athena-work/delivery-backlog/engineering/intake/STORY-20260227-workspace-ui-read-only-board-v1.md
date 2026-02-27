# Story: Build read-only workspace UI board v1

## Metadata
- `id`: STORY-20260227-workspace-ui-read-only-board-v1
- `owner_persona`: staff-personas/Product Designer - Clara.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0004, ADR-0013]
- `success_metric`: Human operators can identify current stage, next story, and blocking issues within 60 seconds.
- `release_checkpoint`: required

## Problem Statement
Humans currently need to navigate many files to understand workflow state, causing navigation overhead and misinterpretation.

## Scope
- In:
  - local web UI showing lane board, active queue, cycle timeline, and drift alerts
  - visual indicator for direction-confirmed vs unconfirmed actions
  - explicit marker when research-mode communication exception is active
- Out:
  - write/edit actions in v1
  - role-based permission management

## Acceptance Criteria
1. Board view renders engineering and architecture lanes from backend read model.
2. Timeline shows recent transitions and observer cycle closures with correlation IDs.
3. UI clearly marks actions requiring human direction confirmation.
4. UI clearly marks research-mode communication exception events.
5. UI defaults are low-vision-friendly: minimum 18px body text, high contrast, clear focus state, generous spacing, and no dense multi-column text blocks.
6. Key human tasks (identify stage, next story, blockers, required confirmation) are visible on first screen without scrolling on common laptop resolution.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-workspace-api-state-machine-v1

## Notes
- Keep mobile readability acceptable for quick operator checks.
