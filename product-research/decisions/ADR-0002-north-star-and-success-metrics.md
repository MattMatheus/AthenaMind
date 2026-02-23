# ADR-0002: Adopt North Star and Success Metrics Framework

## Status
Accepted

## Context
AthenaMind has a foundational three-pillar architecture but lacks a product-level objective function to drive prioritization and evaluate tradeoffs.

## Decision
Adopt the following north star:
- AthenaMind is the memory-and-execution substrate that improves contextual accuracy, auditability, and cost-bounded autonomy across sessions.

Adopt the following metric hierarchy:
- Primary outcomes:
  - Task continuity success rate
  - Rework reduction
  - Trusted autonomy rate
- Guardrails:
  - Memory precision
  - Memory error rate
  - Cost predictability
  - Trace completeness
- Adoption:
  - Weekly active memory-enabled workflows
  - Multi-session workflow completion rate

Roadmap and architecture decisions should state which primary metric they improve and which guardrails they affect.

## Consequences
- Positive:
  - Prioritization becomes outcome-driven
  - Easier to reject low-leverage feature work
  - Stronger alignment between product and architecture
- Negative:
  - Requires instrumentation and baseline measurement work early
  - Some metrics need human-in-the-loop evaluation at first
- Neutral:
  - Specific tooling/telemetry implementation remains open

## Alternatives Considered
- Option A: Feature-count roadmap without explicit scorecard
  - Rejected due to weak signal on actual value delivered
- Option B: Cost-only optimization
  - Rejected because low cost without trust/accuracy is not useful

## Validation Plan
- Define baseline values for all primary and guardrail metrics during v0.1.
- Require each major roadmap item to include expected metric deltas.
- Review scorecard each iteration and adjust backlog ordering accordingly.
