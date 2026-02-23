# Planning Session: Architecture Gap to Execution Backlog Expansion

## Metadata
- `id`: PLAN-20260222-architecture-gap-to-execution
- `date`: 2026-02-22
- `facilitator`: staff-personas/Product Manager - Maya.md
- `status`: finalized

## Session Goal
- Convert documented architecture/roadmap gaps into ranked intake artifacts so empty architecture and engineering lanes can be refilled immediately.

## Problem Framing
- Users affected:
  - Founder/operator coordinating a multi-agent delivery system.
  - Architecture and engineering agents waiting on ranked executable work.
- Jobs-to-be-done:
  - Reconcile planned architecture scope with implemented v0.1 reality.
  - Promote remaining planned ideas into concrete, traceable execution artifacts.
- Pain points:
  - Both `delivery-backlog/architecture/active` and `delivery-backlog/engineering/active` are empty.
  - Program board still marks v0.1 in progress with unstarted release/KPI operational work.
  - Several architecture artifacts describe future contracts that are not yet converted into active intake.

## Constraints
- Technical:
  - Preserve ADR-0007 memory-layer-only scope guardrail for v0.1.
  - Keep production runtime behavior anchored to Go CLI for memory flows.
- Product/business:
  - No time estimates; rank by value, risk, and dependency.
  - `done` is not equal to shipped; release bundle evidence is mandatory.
- Timeline/dependencies:
  - Architecture decisions must precede dependent implementation work.
  - KPI and release evidence should be re-established before broad v0.2 expansion.

## Idea Notes
1. Idea: Rebaseline phase boundaries and architecture maps before scaling the backlog.
   - Value: prevents mixed-scope execution and false completion signals.
   - Risks: over-documenting if not tied to intake and ranking outcomes.
   - Unknowns: exact v0.1 closure criteria versus v0.2 entry threshold.
2. Idea: Operationalize telemetry/KPI loop defined in roadmap docs.
   - Value: enables evidence-based ranking for v0.2/v0.3 work.
   - Risks: metric noise if event contract is underspecified.
   - Unknowns: minimal telemetry payload required for first reliable snapshot.
3. Idea: Promote post-v0.1 snapshot and API-wrapper ideas into explicit architecture + engineering tracks.
   - Value: turns parked strategy into executable backlog.
   - Risks: parallel expansion can dilute v0.1 closeout focus.
   - Unknowns: how aggressively to start v0.3 adapter work before v0.2 stabilizes.

## Decision Split
- Engineering intake candidates:
  - `delivery-backlog/engineering/intake/STORY-20260222-release-checkpoint-bundle-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-architecture-baseline-map-drift-repair.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-dogfood-scenario-pack-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md`
- Architecture/ADR intake candidates:
  - `delivery-backlog/architecture/intake/ARCH-20260222-phase-boundary-rebaseline-v01-v03.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-mvp-constraint-contract-hardening.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-memory-telemetry-event-contract-v01.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-memory-snapshot-architecture-v02.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-memory-api-wrapper-contract-v03.md`

## Intake Artifacts Created
- Engineering stories:
  - `delivery-backlog/engineering/intake/STORY-20260222-release-checkpoint-bundle-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-architecture-baseline-map-drift-repair.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-dogfood-scenario-pack-v01.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
  - `delivery-backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md`
- Architecture stories:
  - `delivery-backlog/architecture/intake/ARCH-20260222-phase-boundary-rebaseline-v01-v03.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-mvp-constraint-contract-hardening.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-memory-telemetry-event-contract-v01.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-memory-snapshot-architecture-v02.md`
  - `delivery-backlog/architecture/intake/ARCH-20260222-memory-api-wrapper-contract-v03.md`
- Traceability defaults:
  - `idea_id`: PLAN-20260222-architecture-gap-to-execution
  - `phase`: v0.1 | v0.2 | v0.3
  - `adr_refs`: ADR-0002, ADR-0006, ADR-0007, ADR-0008, ADR-0009, ADR-0010, ADR-0011, ADR-0012
  - metric field: success_metric / impact metric on every intake artifact

## Recommended Next Stage
- `architect`
- Rationale:
  - Multiple high-priority items require architecture contracts first (phase boundaries, constraints hardening, telemetry contract, snapshot contract, API wrapper contract).
  - Running architect first creates deterministic inputs for PM ranking and prevents engineering from starting under ambiguous cross-phase assumptions.

## Open Questions
- Should v0.3 API wrapper work begin as a design-only track until first KPI baseline is published?
- Should snapshot implementation stay fully blocked behind v0.2 architecture acceptance, or can scaffolding start in parallel?
- Which founder-approved numeric targets should be locked first for constraint enforcement hardening?

## Ranked Execution Order (Dependency-First)
1. `ARCH-20260222-phase-boundary-rebaseline-v01-v03.md`
2. `ARCH-20260222-mvp-constraint-contract-hardening.md`
3. `ARCH-20260222-memory-telemetry-event-contract-v01.md`
4. `STORY-20260222-release-checkpoint-bundle-v01.md`
5. `STORY-20260222-kpi-snapshot-baseline-v01.md`
6. `STORY-20260222-architecture-baseline-map-drift-repair.md`
7. `STORY-20260222-memory-cli-telemetry-contract-v01.md`
8. `STORY-20260222-dogfood-scenario-pack-v01.md`
9. `ARCH-20260222-memory-snapshot-architecture-v02.md`
10. `STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
11. `ARCH-20260222-memory-api-wrapper-contract-v03.md`
12. `STORY-20260222-memory-api-read-gateway-v03.md`
