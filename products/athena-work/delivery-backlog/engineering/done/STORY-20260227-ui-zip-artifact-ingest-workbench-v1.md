# Story: UI zip artifact ingest workbench v1

## Metadata
- `id`: STORY-20260227-ui-zip-artifact-ingest-workbench-v1
- `owner_persona`: staff-personas/Product Designer - Clara.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Operator can ingest a zip artifact bundle from the UI in <=3 interactions and immediately preview updated board state.
- `release_checkpoint`: required

## Problem Statement
- Zip ingest currently requires terminal usage, which adds context switching and slows non-terminal operator workflows.

## Scope
- In:
  - add UI ingest panel for artifact zip upload + execute action
  - API endpoint to receive upload and invoke ingest tool safely
  - inline ingest result panel with accepted/rejected/migrated counts
  - post-ingest board refresh path so operators can verify changes immediately
- Out:
  - remote storage of uploaded bundles
  - multi-user conflict resolution for concurrent ingests

## Acceptance Criteria
1. Operator can upload a `.zip` bundle from UI and run ingest without leaving workspace board.
2. UI shows deterministic ingest report including accepted/rejected/migrated counts and error reason on failure.
3. Board/timeline refresh immediately after successful ingest and reflect new state.
4. Accessibility and keyboard usage remain low-vision-friendly and operationally clear.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-v1-v2-zip-artifact-ingest-tool

## Notes
- Keep API local-first and avoid external upload dependencies.
