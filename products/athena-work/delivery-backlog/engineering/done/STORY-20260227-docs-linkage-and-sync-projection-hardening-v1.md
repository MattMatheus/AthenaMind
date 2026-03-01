# Story: Harden docs linkage and sync projection behavior for Kanban workspace

## Metadata
- `id`: STORY-20260227-docs-linkage-and-sync-projection-hardening-v1
- `owner_persona`: staff-personas/SRE - Nia.md
- `status`: done
- `idea_id`: PLAN-20260227-kanban-workspace-ideation
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Docs workspace linkage and markdown projection remain deterministic across dry-run and normal sync paths with zero unresolved critical drift events in regression suite.
- `release_checkpoint`: required

## Problem Statement
- Workspace delivery speed will degrade if docs linkage and sync projection are fragile or regress under frequent board updates.

## Scope
- In:
- Harden mapping from workspace entities to docs workspace references.
- Strengthen projection/sync guard behavior for ordering, missing artifact, and stale revision classes.
- Expand tests for deterministic projection output and drift failure hints.
- Out:
- Net-new docs IA redesign.
- External documentation platform integration.

## Acceptance Criteria
1. Docs workspace references are generated from canonical workspace state without manual correction.
2. Projection/sync guard emits deterministic error classes and actionable remediation hints.
3. Regression tests cover critical drift classes and pass in local control-plane workflow.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `products/athena-work/delivery-backlog/engineering/active/STORY-20260227-kanban-entity-relationship-implementation-v1.md`
- `products/athena-work/delivery-backlog/engineering/done/STORY-20260227-markdown-sync-worker-and-drift-guard-v1.md`

## Notes
- Keep output compatible with current launcher/observer/workspace board docs paths.
