# ADR-0013: Phase Boundary Contract and Ownership (v0.1-v0.3)

## Status
Accepted

## Context
Architecture artifacts included valid module boundaries but did not provide one explicit phase-allocation contract for backlog sequencing across `v0.1`, `v0.2`, and `v0.3`.

This ambiguity created ranking risk when lanes were repopulated from planning outputs.

## Decision
Adopt a canonical phase boundary contract:

- `v0.1` (now): memory-layer reliability and governance in local-first mode.
  - Primary modules: `procedural-memory`, `state-memory`, `semantic-memory` (CLI/file-based), `memory-governance`, `audit-telemetry` (memory events).
- `v0.2` (next): memory quality maturity and operational instrumentation depth.
  - Primary modules: `semantic-memory` quality upgrades, `semantic-navigation` integration, `memory-governance` hardening, `audit-telemetry` KPI loop maturity.
- `v0.3` (later): cloud-ready and API-wrapper extensions that remain optional and reversible.
  - Primary modules: `orchestrator-api` as memory API wrapper/adapters, with CLI-local fallback retained.
- `workspace-runtime`:
  - Not owned by AthenaMind core in `v0.1`.
  - Treated as an external integration surface in later phases, not a core runtime-ownership commitment.

## Guardrails
- No architecture artifact may assign runtime-execution ownership to AthenaMind core in `v0.1`.
- Engineering ranking must preserve dependency order:
  - architecture contract first,
  - then implementation stories that depend on that contract.
- `done` remains separate from `shipped`; release checkpoint evidence is mandatory.

## Consequences
- Positive:
  - Removes phase ambiguity for PM ranking and architecture-to-engineering handoffs.
  - Reinforces ADR-0007 scope control while keeping expansion paths explicit.
- Negative:
  - Defers some high-throughput implementation until architecture contracts are accepted.
  - Requires periodic contract maintenance as phase goals evolve.
- Neutral:
  - Existing ADRs remain valid; this ADR consolidates allocation and sequencing rules.

## Alternatives Considered
- Option A: Keep phase allocation implicit in multiple docs.
  - Rejected due to recurring drift and ranking ambiguity.
- Option B: Push all queued ideas directly to engineering and resolve phase ownership ad hoc.
  - Rejected due to avoidable rework and scope leakage risk.

## Follow-On Implementation Paths
- `backlog/engineering/intake/STORY-20260222-release-checkpoint-bundle-v01.md`
- `backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
- `backlog/engineering/intake/STORY-20260222-architecture-baseline-map-drift-repair.md`
- `backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- `backlog/engineering/intake/STORY-20260222-dogfood-scenario-pack-v01.md`
- `backlog/engineering/intake/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
- `backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md`
