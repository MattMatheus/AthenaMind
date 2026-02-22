# ADR-0003: Adopt v0.1 Module Boundaries for Three-Pillar Architecture

## Status
Accepted

## Context
The three-pillar model is accepted (ADR-0001), and product scorecard metrics are accepted (ADR-0002). To execute reliably, AthenaMind needs concrete module boundaries and contracts.

## Decision
Adopt the following v0.1 modules:
- `procedural-memory`
- `state-memory`
- `semantic-memory`
- `workspace-runtime`
- `semantic-navigation`
- `memory-governance`
- `audit-telemetry`
- `orchestrator-api`

Adopt the boundary contracts in `research/architecture/MODULE_BOUNDARIES_V01.md`, including:
- governance-gated memory promotions
- mandatory audit emission for recall/write decisions
- no direct global-memory writes from runtime modules

## Consequences
- Positive:
  - Enables parallel workstreams with reduced coupling
  - Improves policy/audit enforceability
  - Aligns implementation to north-star metrics
- Negative:
  - Increased interface design overhead early
  - Requires disciplined contract tests to prevent boundary erosion
- Neutral:
  - Storage/runtime vendor choices remain implementation details

## Alternatives Considered
- Option A: Single service with internal namespaces
  - Rejected due to weak isolation and governance bypass risk
- Option B: Split only by three pillars (one module per pillar)
  - Rejected because policy and telemetry need explicit independent ownership

## Validation Plan
- Add one integration scenario proving all mandatory boundaries are respected.
- Track module-level scorecard impact assumptions in planning docs.
- Reject stories that do not declare module ownership and interface changes.
