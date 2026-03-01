# Story: Add research communication exception policy and audit trail

## Metadata
- `id`: STORY-20260227-research-communication-exception-audit
- `owner_persona`: staff-personas/Security Engineer - Sora.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002]
- `success_metric`: All research-mode undocumented agent communications are logged with policy reason and zero leakage into governed lanes.
- `release_checkpoint`: required

## Problem Statement
Research needs flexibility, but undocumented agent communication in governed lanes undermines policy controls unless tightly scoped and auditable.

## Scope
- In:
  - explicit research-mode flagging and policy gate
  - audit event emission for every exception path
  - hard block for exception path outside research context
- Out:
  - generalized unrestricted communication channels

## Acceptance Criteria
1. Undocumented agent-to-agent communication is permitted only in research context with explicit flag and audit event.
2. Governed lanes reject undocumented communication attempts with clear policy error.
3. Audit events include `cycle_id`, `story_id`, `session_id`, `source_agent`, `target_agent`, `reason`.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-workspace-api-state-machine-v1

## Notes
- Ensure parity between CLI and API policy enforcement.
