# Story: Enforce human direction confirmation gate in workflow

## Metadata
- `id`: STORY-20260227-human-direction-confirmation-gate
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002]
- `success_metric`: 100% of direction-changing transitions include a recorded human confirmation artifact.
- `release_checkpoint`: required

## Problem Statement
Confirmed direction can still be lost or misapplied if not explicitly captured in workflow state. This causes human-agent misalignment and state drift.

## Scope
- In:
  - direction confirmation data model (`confirmed_by`, `confirmed_at`, `scope`, `expiry`)
  - workflow policy requiring confirmation before direction-changing actions
  - UI and CLI visibility for confirmation status
  - planning-stage confirmation summary card (`direction`, `constraints`, `next_stage`, `confirmed_by`)
- Out:
  - legal-signature workflow

## Acceptance Criteria
1. Direction-changing operations are blocked unless a valid human confirmation exists.
2. Confirmation artifacts are visible in UI timeline and exported markdown evidence.
3. Expired or superseded confirmations are rejected with deterministic reasons.
4. Planning stage can produce and display a one-screen confirmation summary that a human can validate quickly.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- ARCH-20260227-shared-workspace-control-plane-contract
- STORY-20260227-workspace-api-state-machine-v1

## Notes
- Keep confirmation UX lightweight for routine daily cycles.
