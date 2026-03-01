# Story: Implement markdown sync worker and drift guard v1

## Metadata
- `id`: STORY-20260227-markdown-sync-worker-and-drift-guard-v1
- `owner_persona`: staff-personas/SRE - Nia.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0010, ADR-0013]
- `success_metric`: Drift detector reports zero unresolved critical drift incidents across pilot cycles.
- `release_checkpoint`: required

## Problem Statement
Backend and markdown can diverge without deterministic synchronization and guardrails, breaking operator trust.

## Scope
- In:
  - worker that exports canonical state to markdown artifacts
  - drift detection for lane order, missing artifacts, and stale revisions
  - guard script integration into stage/doc tests
- Out:
  - bidirectional free-form merge engine

## Acceptance Criteria
1. Canonical backend state exports deterministically to backlog and observer markdown artifacts.
2. Drift guard fails on critical divergence and reports clear remediation hints.
3. Sync output preserves AthenaWork artifact paths expected by launch/observer workflows.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- ARCH-20260227-markdown-sync-authority-and-conflict-policy
- STORY-20260227-workspace-api-state-machine-v1

## Notes
- Include dry-run mode for safe validation.
