# ADR-0008: Freeze v0.1 Goals and Scorecard Target Bands

## Status
Accepted

## Context
ADR-0002 established the north star and scorecard metric hierarchy, but v0.1 planning still lacked explicit numeric targets. Without fixed target bands, backlog prioritization and QA outcomes cannot be tied to measurable thresholds.

ADR-0007 refined scope to memory-layer outcomes for v0.1, so targets must emphasize memory continuity, trust, and predictable cost within that boundary.

## Decision
Freeze the following v0.1 product goals and target bands. These targets remain fixed for v0.1 unless superseded by a later ADR.

v0.1 goals:
- Goal 1: Improve multi-session continuity so users can resume work with less re-briefing.
- Goal 2: Increase trusted autonomy with audit-complete and low-error memory behavior.
- Goal 3: Improve memory quality (precision up, error rate down) under real workflows.
- Goal 4: Keep outcomes cost-bounded and show repeat weekly adoption.

Band definitions:
- Green: at or above v0.1 target threshold.
- Yellow: watch zone; below target but above minimum acceptable floor.
- Red: below minimum acceptable floor.

Scorecard target bands mapped to ADR-0002 metrics:
- Task continuity success rate:
  - Red: < 70%
  - Yellow: 70% to 84%
  - Green: >= 85%
- Rework reduction (vs baseline):
  - Red: < 20%
  - Yellow: 20% to 34%
  - Green: >= 35%
- Trusted autonomy rate:
  - Red: < 60%
  - Yellow: 60% to 74%
  - Green: >= 75%
- Memory precision:
  - Red: < 92%
  - Yellow: 92% to 95.9%
  - Green: >= 96%
- Memory error rate:
  - Red: > 5%
  - Yellow: 2% to 5%
  - Green: < 2%
- Cost predictability (absolute variance from estimate):
  - Red: > 35%
  - Yellow: 20% to 35%
  - Green: <= 20%
- Trace completeness:
  - Red: < 95%
  - Yellow: 95% to 98.9%
  - Green: >= 99%
- Weekly active memory-enabled workflows:
  - Red: < 5 workflows/week
  - Yellow: 5 to 9 workflows/week
  - Green: >= 10 workflows/week
- Multi-session workflow completion rate:
  - Red: < 55%
  - Yellow: 55% to 69%
  - Green: >= 70%

Goal-to-metric mapping:
- Goal 1: Task continuity success rate, Rework reduction, Multi-session workflow completion rate
- Goal 2: Trusted autonomy rate, Trace completeness, Memory error rate
- Goal 3: Memory precision, Memory error rate
- Goal 4: Cost predictability, Weekly active memory-enabled workflows

## Consequences
- Positive:
  - Sprint planning can rank work against explicit thresholds.
  - QA can evaluate outcomes as pass/watch/fail against fixed bands.
  - Decision trail now ties v0.1 priorities directly to ADR-0002 metrics.
- Negative:
  - Some targets may require baseline/instrumentation updates to measure reliably.
  - Thresholds may feel conservative for high-performing scenarios.
- Neutral:
  - Metrics remain unchanged; only target bands and goal mappings are frozen.

## Alternatives Considered
- Option A: Keep qualitative goals without numeric bands.
  - Rejected because it preserves subjective prioritization.
- Option B: Delay target setting until end of v0.1.
  - Rejected because planning and QA need immediate measurable gates.

## Validation Plan
- Reference this ADR in planning artifacts and story-level acceptance checks.
- Review scorecard weekly during v0.1 using these bands.
- Any target adjustment requires a follow-on ADR update, not ad hoc edits.
