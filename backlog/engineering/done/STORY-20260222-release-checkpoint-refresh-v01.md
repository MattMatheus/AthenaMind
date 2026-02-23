# Story: Refresh v0.1 Release Checkpoint Bundle After Hardening Cycle

## Metadata
- `id`: STORY-20260222-release-checkpoint-refresh-v01
- `owner_persona`: Product Manager - Maya.md
- `status`: done
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0002, ADR-0008]
- `success_metric`: release bundle is refreshed with current scope/evidence and explicit ship/hold decision grounded in post-hardening KPI evidence
- `release_checkpoint`: required

## Problem Statement
- The existing release bundle remains on `hold` and contains stale scope references; a refreshed checkpoint decision is required after hardening and KPI delta publication.

## Scope
- In:
  - Refresh `work-system/handoff/RELEASE_BUNDLE_v0.1-initial-2026-02-22.md` with current included/excluded scope.
  - Link post-hardening KPI and QA evidence.
  - Record explicit ship/hold decision with updated rationale and rollback direction.
- Out:
  - New runtime code changes.
  - Non-v0.1 release planning.

## Acceptance Criteria
1. Release bundle scope and evidence links match current backlog states and completed artifacts.
2. Bundle includes post-hardening KPI evidence and updated risk assessment.
3. Bundle decision is explicit (`ship` or `hold`) with rationale tied to ADR-0008 target bands.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `backlog/engineering/active/STORY-20260222-kpi-delta-snapshot-after-semantic-hardening-v01.md`
- `work-system/handoff/RELEASE_BUNDLE_v0.1-initial-2026-02-22.md`
- `docs/process/STAGE_EXIT_GATES.md`

## Notes
- Keep `done` vs `shipped` boundary explicit; this story is the v0.1 release-decision control-plane gate.
