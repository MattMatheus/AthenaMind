<!-- AUDIENCE: Internal/Technical -->

# Next Research Agent Directive

Your task is to begin the next research cycle by executing the top item in the active research ingest.

## Primary Task

- **Backlog Item:** `research/roadmap/RESEARCH_BACKLOG.md` -> `Now` item 1: Program state consistency and release-checkpoint adoption.

## Scope Guardrail (Mandatory)

AthenaMind v0.1 is memory-layer only:
- In scope: procedural memory, state/semantic memory, explainability/auditability, governance/HITL.
- Out of scope: runtime execution orchestration, pod/container lifecycle ownership.

## Current Canonical Context
- `research/decisions/ADR-0001-three-pillar-foundation.md`
- `research/decisions/ADR-0002-north-star-and-success-metrics.md`
- `research/decisions/ADR-0003-v01-module-boundaries.md`
- `research/decisions/ADR-0004-v01-mvp-constraints-framework.md`
- `research/decisions/ADR-0005-v01-local-storage-and-embedding-stack.md`
- `research/decisions/ADR-0006-governance-and-hitl-policy.md`
- `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- `research/product/VISION_WORKSHOP_2026-02-22.md`

## Rules of Engagement
- Favor 80/20 insights with high leverage.
- Keep claims traceable to source notes or explicit assumptions.
- Prefer reversible decisions when uncertainty is high.
- Research external primary sources when assumptions need verification.

## Handoff Contract (Mandatory)

Before ending your session:
1. Truncate and rewrite `research/handoff.md` with only:
   - What was produced
   - Why those outputs matter now
   - Validation performed
   - Open questions for next cycle
   - Recommended next backlog item
2. Update `research/roadmap/RESEARCH_BACKLOG.md` to reflect completion/progression.
3. Update this file (`research/research-agent-seed-prompt.md`) so `Primary Task` points to the next research cycle item.
