# Story: Implement workspace API state machine v1 with policy gates

## Metadata
- `id`: STORY-20260227-workspace-api-state-machine-v1
- `owner_persona`: staff-personas/Software Engineer - Max.md
- `status`: qa
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0004, ADR-0013, WSI-ADR-0002]
- `success_metric`: Core transitions are enforced by API with >=99% valid transition success and deterministic rejection codes.
- `release_checkpoint`: required

## Problem Statement
Without a transactional transition API, concurrent human+agent operations can cause state drift and ambiguous stage outcomes.

## Scope
- In:
  - API endpoints for story/cycle transitions (`promote`, `fail`, `close_cycle`, `confirm_direction`)
  - transition precondition engine aligned to AthenaWork gates
  - immutable transition event write on every accepted transition
  - policy handling for research-mode non-canonical communication
- Out:
  - advanced analytics
  - external auth federation

## Acceptance Criteria
1. All canonical transitions are enforced through API preconditions and produce deterministic error codes.
2. Research-mode exception exists for undocumented agent-to-agent communication only when explicitly flagged as research and fully logged.
3. Non-research lanes reject undocumented agent-to-agent communication attempts.
4. `confirm_direction` checkpoint must be recorded for human-confirmed direction changes before transition execution.
5. API returns concise, machine-readable responses optimized for agent execution (`code`, `reason`, `next_action`, `correlation_id`) with p95 transition response <= 300ms in local mode.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- ARCH-20260227-shared-workspace-control-plane-contract
- ARCH-20260227-markdown-sync-authority-and-conflict-policy

## Notes
- Keep rejection reasons human-readable for UI display.
