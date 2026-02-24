# Dogfood Scenario Run: 2026-02-22 (Semantic Hardening)

## Run Metadata
- `run_id`: `DOGFOOD-RUN-2026-02-22-02`
- `pack_version`: `v0.1`
- `scenarios_executed`: 4
- `operator`: `QA Engineer - Iris.md`
- `focus`: `semantic retrieval precision and trace completeness hardening`

## Scenario Outcomes

| scenario_id | memory_type | session_id | result | operator_verdict | latency_ms (retrieve p95) | failure_class |
| --- | --- | --- | --- | --- | --- | --- |
| SCN-PROC-01 | procedural | cyc-proc-20260222-03 | success | correct | 175 | none |
| SCN-STATE-01 | state | cyc-state-20260222-02 | success | correct | 225 | none |
| SCN-SEM-01 | semantic | cyc-sem-20260222-02 | success | correct | 355 | none |
| SCN-PROC-02 | procedural | cyc-proc-20260222-04 | success | correct | 205 | none |

## KPI-Relevant Snapshot Annotations
- `resume_success_rate`: `1/1 = 100%` for explicitly resumed state scenario.
- `preference_adherence_rate`: `4/4 = 100%`.
- `precision_at_3`: `3/3 = 100%` on semantic retrieval checks.
- `wrong_memory_recall_rate`: `0/4 = 0%`.
- `trace_completeness_rate`: `4/4 = 100%`.
- `review_gate_bypass_rate`: `0%` (no medium/high-risk bypass observed).
- `p95_retrieval_latency_ms`: `355` (semantic scenario improved from baseline).

## Before/After Delta vs Baseline Run
- Baseline reference:
  - `operating-system/metrics/dogfood-scenario-run-2026-02-22.md`
- `precision_at_3`: `66.7% -> 100%` (`+33.3pp`)
- `trace_completeness_rate`: `75% -> 100%` (`+25pp`)
- `wrong_memory_recall_rate`: `25% -> 0%` (`-25pp`)
- `semantic_retrieval_p95_latency_ms`: `470 -> 355` (`-115ms`)

## Failure Classification Summary
- `retrieval`: 0
- `operator-ux`: 0
- `policy`: 0
- `memory-quality`: 0

## Notes
- This hardening follow-up run closes the weak-signal path identified in the baseline semantic scenario and provides delta evidence for the next KPI snapshot update.
