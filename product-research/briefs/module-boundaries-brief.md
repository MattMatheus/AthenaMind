# Research Brief: 3-Pillar Module Boundaries (v0.1)

## Source Inputs
- `product-research/architecture/FOUNDATION_PILLARS.md`
- `product-research/decisions/ADR-0001-three-pillar-foundation.md`
- `product-research/decisions/ADR-0002-north-star-and-success-metrics.md`
- `product-research/references/external-sources-2026-02-22.md`

## Problem Statement
The three-pillar architecture is defined conceptually, but AthenaMind still needs concrete module boundaries to support parallel development, enforce governance, and avoid cross-cutting sprawl.

## Key Insights (80/20)
- The main failure mode is boundary leakage: memory, policy, runtime, and telemetry logic blending into one service.
- We need one explicit "memory plane" and one explicit "execution plane" with shared contracts.
- Auditability should be an independent concern with strict write-once trace records.
- Retrieval/index strategy must be pluggable from day one to avoid early lock-in.

## Proposed v0.1 Module Boundaries
1. `procedural-memory` (Pillar 1)
- Load/validate/version global and repo-local Markdown directives.
- Resolve precedence (`global` vs `local`) and emit merged startup context.

2. `state-memory` (Pillar 2)
- Persist sessions, episodes, preferences, and checkpointable workflow state.
- Own structured schemas and state migration rules.

3. `semantic-memory` (Pillar 2)
- Embedding generation, vector index operations, hybrid retrieval orchestration.
- Abstract provider/index implementations (local FAISS first; cloud optional).

4. `workspace-runtime` (Pillar 3)
- Isolated run lifecycle (pod/container create, clone, execute, teardown).
- Enforce runtime budgets and isolation profiles.

5. `semantic-navigation` (Pillar 3 support)
- LSP/symbolic tools and snippet extraction from workspace context.
- Return minimal relevant context for tasks.

6. `memory-governance`
- Promotion rules (local -> global), confidence thresholds, conflict policy.
- Human review gates for sensitive/low-confidence writes.

7. `audit-telemetry`
- Trace every memory read/write/decision with run/session linkage.
- Emit OTEL spans/attributes and store immutable audit events.

8. `orchestrator-api`
- Unified public interface for memory and runtime operations.
- Request validation, policy enforcement, and cross-module workflow coordination.

## Boundary Contracts (Critical)
- `procedural-memory` can read Markdown only; it cannot persist semantic state.
- `state-memory` is source-of-truth for structured state; retrieval modules never mutate state directly.
- `semantic-memory` returns candidates + scores + provenance, not final decisions.
- `memory-governance` is the only module that authorizes preference promotion.
- `workspace-runtime` cannot write global memory directly; writes must flow through governance + state modules.
- `audit-telemetry` receives events from all modules and cannot be bypassed.

## Implications for AthenaMind
- Enables parallel implementation by ownership boundary.
- Supports product guardrails from ADR-0002 through measurable module responsibilities.
- Reduces regression risk by isolating policy and observability from business logic.

## Open Questions
- What is the minimal event schema for end-to-end memory auditability?
- Should semantic retrieval use one module or split query planner/index adapter now?
- Do we need a dedicated `cost-estimator` module in v0.1 or keep inside orchestrator-api?

## Recommended Decisions
- Adopt these eight module boundaries for v0.1 architecture planning.
- Require every implementation story to declare module ownership and boundary constraints.
