# Story: Consolidate Architecture Baseline and ADR Cross-Reference Map

## Metadata
- `id`: STORY-20260222-architecture-baseline-consolidation
- `owner_persona`: Software Architect - Ada.md
- `status`: done

## Problem Statement
Architecture decisions and design artifacts are spread across multiple ADR and research files. We need one consolidated architecture baseline map so engineering and PM can scope work against a stable reference.

## Scope
- In:
  - Create a consolidated architecture baseline map doc
  - Cross-reference ADR-0001 through ADR-0008 and current architecture docs
  - Identify unresolved architectural gaps that require follow-up stories
- Out:
  - Implementing memory-layer runtime code

## Acceptance Criteria
1. A consolidated architecture baseline map exists under `research/architecture/`.
2. The map links each major architecture area to governing ADR(s).
3. At least three unresolved architecture gaps are captured as intake stories if discovered.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `research/architecture/FOUNDATION_PILLARS.md`
- `research/architecture/MODULE_BOUNDARIES_V01.md`
- `research/architecture/MVP_CONSTRAINTS_V01.md`
- `research/decisions/ADR-0001-three-pillar-foundation.md`
- `research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`

## Notes
- This story is intended for architect stage execution.
