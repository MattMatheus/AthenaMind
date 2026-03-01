# Story: Build docs workspace surface and navigation v1

## Metadata
- `id`: STORY-20260227-docs-workspace-surface-and-navigation-v1
- `owner_persona`: staff-personas/Technical Writer - Clara.md
- `status`: done
- `idea_id`: PLAN-20260227-kanban-workspace-ideation
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Operator can open relevant docs for active work in <=1 navigation context switch from board.
- `release_checkpoint`: required

## Problem Statement
- Docs are present but not surfaced as a first-class workspace experience connected to current work context.

## Scope
- In:
  - docs workspace panel/index integrated into Kanban workspace
  - task-linked doc references and quick jump navigation
  - plain-language summaries for key workflow/process docs
- Out:
  - full-text semantic search engine
  - external docs publishing pipeline

## Acceptance Criteria
1. Workspace UI shows docs surface with direct links relevant to active tasks/lane context.
2. Operator can jump from board to targeted docs with one interaction.
3. Docs surface remains readable and low-vision-friendly with clear hierarchy and focus states.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- ARCH-20260227-kanban-workspace-information-architecture-v1
- STORY-20260227-kanban-board-read-model-and-lanes-v1

## Notes
- Prioritize docs that reduce workflow ambiguity for first-time and returning operators.
