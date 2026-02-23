# Research Brief: North Star and Success Metrics

## Source Inputs
- `product-research/ingest/archive/pillars.md`
- `product-research/architecture/FOUNDATION_PILLARS.md`
- `product-research/decisions/ADR-0001-three-pillar-foundation.md`

## Problem Statement
AthenaMind needs a clear product north star and measurable success criteria so architecture choices optimize for real operator outcomes, not feature volume.

## Key Insights (80/20)
- The core value is reduced context rot with trustworthy, bounded autonomy.
- Trust requires explainable memory behavior: what was loaded, recalled, and written.
- Safety and cost boundaries are product features, not implementation details.
- The best early signal is cycle-time and rework reduction on repeat workflows.

## North Star
AthenaMind becomes the default memory-and-execution substrate for agentic work by making each run more contextually accurate, auditable, and cost-bounded than stateless alternatives.

## Success Metric Framework
Primary outcome metrics:
- Task continuity success rate: percent of resumed tasks that complete without re-briefing from scratch.
- Rework reduction: percent reduction in duplicate clarification and repeated setup steps.
- Trusted autonomy rate: percent of runs accepted without manual rollback due to memory/execution errors.

Guardrail metrics:
- Memory precision: percent of recalled memory judged relevant by user or evaluator.
- Memory error rate: percent of runs with incorrect memory recall influencing output.
- Cost predictability: percent of runs that stay within estimated token/runtime budget bands.
- Trace completeness: percent of runs with complete recall/write audit records.

Adoption metrics:
- Weekly active workflows using memory retrieval.
- Multi-session workflow completion rate.

## Implications for AthenaMind
- Telemetry and audit fields must be first-class in v0.1.
- Memory promotion rules (local -> global) need explicit policy before broad rollout.
- MVP scope should prioritize reliability and auditability over broad integrations.

## Open Questions
- What threshold defines acceptable memory precision for autonomous actions?
- How strict should budget aborts be versus graceful degradation?
- Which metrics require human labeling versus automated proxies?

## Recommended Decisions
- Adopt the north star and metric framework as default product scorecard.
- Gate roadmap priorities by expected improvement in primary outcome metrics.
