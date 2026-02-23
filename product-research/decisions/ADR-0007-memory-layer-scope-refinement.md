# ADR-0007: Refine Scope to Memory Layer Only (No Runtime Execution in v0.1)

## Status
Accepted

## Context
Vision workshop feedback clarified that AthenaMind v0.1 should be a memory layer consumed by agents, not an execution runtime platform. Earlier artifacts inherited runtime-heavy assumptions from Athena migration context.

## Decision
For v0.1, AthenaMind scope is refined to:
- In scope:
  - Procedural memory (global + repo-local directives)
  - State/semantic memory (episodes, preferences, retrieval)
  - Explainability/auditability for memory-driven agent actions
  - Instant resume and long-term preference adherence improvements
- Out of scope:
  - Runtime execution orchestration
  - Pod/container lifecycle management
  - Sandboxed code execution responsibilities

Implementation constraints:
- Local-first default (local disk storage)
- Cloud-ready extension model for later phases
- Composable architecture that allows intelligence-layer evolution without hard coupling

Planned (post-v0.1) enhancement:
- Memory snapshot capability to guard against memory rot through versioned restore points.

## Consequences
- Positive:
  - Tighter product focus and reduced v0.1 complexity
  - Faster path to proving memory value
  - Cleaner boundaries for future integrations with external runtimes/agents
- Negative:
  - Existing runtime-oriented artifacts must be reinterpreted or superseded
  - Some previously planned modules become deferred or adapter-only
- Neutral:
  - Governance and traceability remain core requirements

## Alternatives Considered
- Option A: Keep full memory + runtime scope in v0.1
  - Rejected due to founder-directed product focus and scope risk
- Option B: Memory-only but postpone semantic memory
  - Rejected due to reduced ability to prove adherence/continuity improvements

## Validation Plan
- Future ADRs must explicitly avoid assigning execution-runtime ownership to AthenaMind core.
- Re-baseline metrics and governance docs to memory-layer responsibilities.
