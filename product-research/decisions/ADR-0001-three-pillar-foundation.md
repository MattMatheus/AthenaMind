# ADR-0001: Adopt Three-Pillar Memory Foundation

## Status
Accepted

## Context
The pivot requires a durable architecture that separates long-term instructions, persistent semantic/state memory, and short-term execution context. The `product-research/ingest/archive/pillars.md` note defines this model and emphasizes secure sandboxed execution, repository-local directives, and semantic code navigation.

## Decision
AthenaMind will adopt the following foundational architecture:
1. Procedural Memory in versioned Markdown
2. State & Semantic Memory in persistent storage with retrievable episodes/preferences
3. Working Environment Memory in isolated ephemeral runtimes

All new design artifacts should map explicitly to one or more pillars and their cross-pillar contracts.

## Consequences
- Positive:
  - Clear boundary between instruction memory, persistent learned state, and execution context
  - Better safety and observability via isolated runtime and auditable memory flows
  - Better product clarity for roadmap and positioning
- Negative:
  - Higher upfront design and integration complexity
  - Requires explicit governance for memory promotion/conflict resolution
- Neutral:
  - Storage and runtime technology choices remain open as implementation details

## Alternatives Considered
- Option A: Single unified memory store for all contexts
  - Rejected due to weaker separation of concerns and harder auditing
- Option B: Procedural memory only (Markdown) without persistent episodic state
  - Rejected due to poor cross-session compounding

## Validation Plan
- Create v0.1 architecture and API contracts that demonstrate one end-to-end flow across all three pillars.
- Confirm each run can answer: what was loaded, what was recalled, what changed, and why.
