# KPI Snapshot: 2026-02-22 (Baseline)

## Summary
- Published first v0.1 KPI baseline using current dogfood and cycle evidence.
- Baseline indicates governance discipline is strong, while retrieval quality and trace completeness remain below ADR-0008 target bands.
- Telemetry is still partial; this snapshot is explicitly marked as a manual baseline to be replaced by instrumented weekly snapshots.

## Metrics
- Lead time: `1.2 days` (manual sample from recent completed engineering cycles)
- QA turnaround: `0.4 days` (manual sample from recent qa->done transitions)
- Defect escape: `0%` (no post-done reopening observed in current sample)
- Reopen rate: `0%` (no qa->active reopen in current sample)
- Handoff completeness: `100%` (story + handoff + QA result + observer evidence present for sampled cycles)
- Retrieval quality gate pass rate: `75%` (`3/4` scenarios passed gate in `DOGFOOD_SCENARIO_RUN_2026-02-22.md`)
- Traceability completeness: `75%` (`3/4` scenarios with full trace evidence)

## Interpretation
- ADR-0008 target-band interpretation (proxy mapping where direct metric is not yet fully instrumented):
  - Memory precision proxy from retrieval gate pass: `75%` -> `Red` against memory precision band (`Green >= 96%`, `Yellow 92%-95.9%`, `Red < 92%`).
  - Trace completeness: `75%` -> `Red` against trace completeness band (`Green >= 99%`, `Yellow 95%-98.9%`, `Red < 95%`).
  - Governance adherence signal: `review_gate_bypass_rate=0%` from dogfood run, aligned with v0.1 governance intent.
- Telemetry gap note:
  - This baseline uses partial manual scoring and should not be treated as stable trend data until telemetry-contract implementation is fully active.

## Actions
- Prioritize semantic retrieval + trace completeness hardening from baseline weak signal:
  - `backlog/engineering/intake/STORY-20260222-dogfood-semantic-retrieval-hardening-v01.md`
- Complete telemetry contract implementation to move from manual baseline to deterministic weekly KPI snapshots:
  - `backlog/engineering/active/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- Publish the next KPI snapshot after telemetry and constraint-enforcement stories complete, then compare deltas vs this baseline.
