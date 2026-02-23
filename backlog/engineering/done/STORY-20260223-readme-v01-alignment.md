# Story: Align README with v0.1 reality

## Metadata
- `id`: STORY-20260223-readme-v01-alignment
- `owner_persona`: personas/Technical Writer - Clara.md
- `status`: qa
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0007, ADR-0009]
- `success_metric`: README accurately describes what v0.1 delivers today with no aspirational claims about unimplemented features
- `release_checkpoint`: deferred

## Problem Statement
- The current README.md describes the full 3-pillar vision including FAISS, Podman pods, SQLite, and cloud integration — none of which are implemented. Since agents read the README as context, an aspirational README can cause agents to over-scope or build toward v0.3 when they should stay in v0.1 lane. The README should describe the current product accurately.

## Scope
- In:
  - Rewrite README.md to describe what v0.1 actually delivers: a Go CLI for file-based memory with write, retrieve, evaluate, snapshot, and read gateway commands, governance-gated mutations, telemetry, and constraint enforcement.
  - Move the original vision content to `docs/product/VISION.md` (preserve, do not delete).
  - Link from README to the phased plan (`research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`) for future direction.
  - Include quick-start usage examples for the CLI.
- Out:
  - Code changes
  - Changes to ADRs or architecture docs

## Acceptance Criteria
1. README.md describes only implemented v0.1 features.
2. No references to FAISS, Podman, SQLite, embeddings, or cloud services in README unless clearly marked as future/planned.
3. Original vision content preserved in `docs/product/VISION.md`.
4. README includes basic usage examples for write and retrieve commands.
5. README links to phased implementation plan for roadmap context.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- None

## Notes
- This is a documentation-only story. No code changes.
