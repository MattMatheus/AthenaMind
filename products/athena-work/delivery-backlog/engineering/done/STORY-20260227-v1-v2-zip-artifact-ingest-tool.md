# Story: V1 -> V2 zip artifact ingest tool

## Metadata
- `id`: STORY-20260227-v1-v2-zip-artifact-ingest-tool
- `owner_persona`: staff-personas/Software Engineer - Alex.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Operator can import a provided `.zip` artifact bundle into the workspace in <=2 minutes with deterministic ingest results.
- `release_checkpoint`: required

## Problem Statement
- Current state seeding relies on direct local artifact edits and lacks a user-facing ingest path for packaged artifact objects.

## Scope
- In:
  - ingest command/tool that accepts user-provided `.zip` bundle of required artifact objects
  - validation of required artifact structure and schema compatibility for V1 -> V2 support
  - deterministic unpack + projection into local workspace state for UI consumption
  - clear ingest report with accepted, rejected, and migrated object counts
- Out:
  - remote artifact hosting or network fetch pipeline
  - multi-user concurrent ingest conflict resolution beyond local-first single operator flow

## Acceptance Criteria
1. Operator can provide a `.zip` bundle and run one command to ingest artifact objects into the local system.
2. Tool enforces required object schema validation and reports actionable errors for invalid bundles.
3. V1 artifacts are transformed or mapped to V2-compatible objects with explicit migration report output.
4. Ingested output is visible in workspace board/read-model without manual file editing.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-workspace-api-state-machine-v1
- STORY-20260227-markdown-sync-worker-and-drift-guard-v1

## Notes
- Keep ingest logic auditable and local-first by default.
