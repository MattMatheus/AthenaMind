# Work System Improvement

This area manages improvements to the development process itself.

## Goal
Treat process quality as a product: measurable, iterative, and testable.

## Core Loop
1. Capture process gaps in `work-system/backlog/intake/`.
2. Refine and rank into `work-system/backlog/active/`.
3. Implement process changes.
4. QA the process change impact.
5. Measure KPI deltas in `work-system/metrics/`.
6. Keep/adjust/revert via decision records.
7. Maintain explicit shipped checkpoint bundles before declaring release completion.

## Documentation Maintenance Rule
Any accepted process change must include documentation sync:
- `HUMANS.md` (operator quick guide)
- `AGENTS.md` (agent discovery/operating guide)
- `DEVELOPMENT_CYCLE.md` (canonical stage behavior)

## Structure
- `backlog/`: state-machine for process-improvement stories
- `experiments/`: hypothesis-driven trials
- `decisions/`: process ADR-style records
- `metrics/`: process KPIs and trend snapshots
- `playbooks/`: operational runbooks for cycle stages
- `retros/`: sprint and incident retrospectives
- `handoff/`: current process handoff status
