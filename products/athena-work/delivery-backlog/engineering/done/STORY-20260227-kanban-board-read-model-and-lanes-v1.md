# Story: Build Kanban board read model and lane rendering v1

## Metadata
- `id`: STORY-20260227-kanban-board-read-model-and-lanes-v1
- `owner_persona`: staff-personas/Product Designer - Clara.md
- `status`: done
- `idea_id`: PLAN-20260227-kanban-workspace-ideation
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Operator can identify lane status, next task, blockers, and current cycle state in <=60 seconds on first screen.
- `release_checkpoint`: required

## Problem Statement
- Current workspace UI is informative but not yet a full Kanban lane/card operational board.

## Scope
- In:
  - read-model driven lane columns with cards and blocker markers
  - cycle summary and queue-state visualization aligned to canonical state model
  - deterministic card ordering from authoritative projection
- Out:
  - write/edit lane operations
  - multi-user realtime collaboration

## Acceptance Criteria
1. Board presents lane/card model consistent with canonical workspace read model.
2. Card ordering and blocker markers are deterministic and reproducible after refresh.
3. First-screen view exposes stage, next task, blockers, and required confirmation without scroll on common laptop resolution.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- ARCH-20260227-kanban-workspace-information-architecture-v1
- STORY-20260227-markdown-sync-worker-and-drift-guard-v1

## Notes
- Preserve low-vision defaults as baseline.
