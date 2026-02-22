# ADR-0015: v0.1 Memory Telemetry Event Contract

## Status
Accepted

## Context
KPI and dogfooding docs define metrics and process expectations, but no formal telemetry contract existed for memory CLI producers and KPI consumers.

## Decision
Adopt a minimal telemetry contract for v0.1 memory flows with required fields, deterministic validation rules, and explicit KPI mapping.

Canonical artifact:
- `research/architecture/MEMORY_TELEMETRY_EVENT_CONTRACT_V01.md`

Contract principles:
- Minimal required fields only for v0.1.
- Deterministic required-field validation.
- Retrieval events always include selection provenance.
- Manual scoring support via `operator_verdict=not_scored`.

## Consequences
- Positive:
  - Enables consistent KPI generation and QA evidence.
  - Reduces implementation ambiguity for telemetry emission.
- Negative:
  - Adds strict schema compliance burden to CLI event emission paths.
  - Some fields require manual scoring until tooling matures.
- Neutral:
  - Backend transport/storage remains implementation-decoupled.

## Alternatives Considered
- Keep telemetry guidance informal in roadmap docs.
  - Rejected due to schema drift risk.
- Instrument a broad schema first, optimize later.
  - Rejected due to v0.1 scope and complexity limits.

## Follow-On Implementation Paths
- `backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- `backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
- `backlog/engineering/intake/STORY-20260222-dogfood-scenario-pack-v01.md`
