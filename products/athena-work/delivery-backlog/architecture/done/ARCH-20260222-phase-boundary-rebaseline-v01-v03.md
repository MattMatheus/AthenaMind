# Architecture Story: Rebaseline Phase Boundaries and Ownership (v0.1-v0.3)

## Metadata
- `id`: ARCH-20260222-phase-boundary-rebaseline-v01-v03
- `owner_persona`: Software Architect - Ada.md
- `status`: qa
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0003, ADR-0007, ADR-0008]
- `decision_owner`: staff-personas/Software Architect - Ada.md
- `success_metric`: v0.1/v0.2/v0.3 boundary ambiguities reduced to zero unresolved conflicts in architecture baseline docs

## Decision Scope
- Produce an explicit phase-boundary contract mapping each module/domain from architecture docs to `v0.1 now`, `v0.2 next`, or `v0.3 later`.

## Problem Statement
- Current architecture artifacts mix memory-layer-only v0.1 constraints with broader module plans, creating ranking ambiguity when lanes are empty and backlog refill begins.

## Inputs
- ADRs:
  - `product-research/decisions/ADR-0003-v01-module-boundaries.md`
  - `product-research/decisions/ADR-0007-memory-layer-scope-refinement.md`
  - `product-research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`
- Architecture docs:
  - `product-research/architecture/MODULE_BOUNDARIES_V01.md`
  - `product-research/architecture/ARCHITECTURE_BASELINE_MAP_V01.md`
  - `product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`
- Constraints:
  - Preserve memory-layer-only scope for v0.1.
  - No runtime ownership assignments in v0.1 outputs.

## Outputs Required
- ADR updates:
  - ADR addendum or linked decision memo clarifying phase allocation and guardrails.
- Architecture artifacts:
  - Updated phase-boundary map with unambiguous module-to-phase ownership.
- Risk/tradeoff notes:
  - Throughput vs scope-control tradeoffs for aggressive parallelization.

## Acceptance Criteria
1. Each module in `MODULE_BOUNDARIES_V01.md` has one explicit phase assignment and ownership note.
2. No artifact claims runtime orchestration ownership in v0.1 outputs.
3. PM can rank backlog with dependency clarity across v0.1, v0.2, and v0.3.

## QA Focus
- Verify phase allocations are internally consistent across roadmap, architecture map, and ADR constraints.

## Intake Promotion Checklist (intake -> ready)
- [ ] Decision scope is explicit and bounded.
- [ ] Problem statement describes urgency and impact.
- [ ] Required inputs are listed (ADRs, architecture docs, constraints).
- [ ] Separation rule verified: architecture output, not implementation output.
- [ ] Required outputs are concrete and reviewable in QA handoff.
- [ ] Risks/tradeoffs include mitigation and owner.
