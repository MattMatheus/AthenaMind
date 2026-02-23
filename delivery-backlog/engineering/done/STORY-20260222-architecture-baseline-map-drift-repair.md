# Story: Repair Architecture Baseline Map Drift

## Metadata
- `id`: STORY-20260222-architecture-baseline-map-drift-repair
- `owner_persona`: Technical Writer - Clara.md
- `status`: qa
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0007, ADR-0009, ADR-0010, ADR-0011, ADR-0012]
- `success_metric`: baseline map has zero stale intake references and accurate planned-vs-implemented status mapping
- `release_checkpoint`: required

## Problem Statement
- `ARCHITECTURE_BASELINE_MAP_V01.md` still reports unresolved gaps as intake links that no longer exist, creating false backlog signals.

## Scope
- In:
  - Update baseline map unresolved section to reflect current done/intake status.
  - Add explicit planned-vs-implemented matrix for key architecture domains.
  - Cross-link to current intake artifacts for true open items.
- Out:
  - New runtime feature development.
  - ADR policy changes.

## Acceptance Criteria
1. No references in architecture baseline map point to missing intake files.
2. Planned vs implemented status is explicit for each baseline domain.
3. Document tests remain green after updates.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `product-research/architecture/ARCHITECTURE_BASELINE_MAP_V01.md`
- Current backlog state in `delivery-backlog/architecture/` and `delivery-backlog/engineering/`

## Notes
- This story reduces planning errors by aligning architecture map output with real queue state.
