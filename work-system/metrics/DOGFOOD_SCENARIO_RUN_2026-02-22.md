# Dogfood Scenario Run: 2026-02-22

## Run Metadata
- `run_id`: `DOGFOOD-RUN-2026-02-22-01`
- `pack_version`: `v0.1`
- `scenarios_executed`: 4
- `operator`: `QA Engineer - Iris.md`

## Scenario Outcomes

| scenario_id | memory_type | session_id | result | operator_verdict | latency_ms (retrieve p95) | failure_class |
| --- | --- | --- | --- | --- | --- | --- |
| SCN-PROC-01 | procedural | cyc-proc-20260222-01 | success | correct | 180 | none |
| SCN-STATE-01 | state | cyc-state-20260222-01 | success | partially_correct | 240 | operator-ux |
| SCN-SEM-01 | semantic | cyc-sem-20260222-01 | fail | partially_correct | 470 | retrieval |
| SCN-PROC-02 | procedural | cyc-proc-20260222-02 | success | correct | 210 | none |

## KPI-Relevant Snapshot Annotations
- `resume_success_rate`: `1/1 = 100%` for explicitly resumed state scenario.
- `preference_adherence_rate`: `3/4 = 75%` (one partial adherence in semantic scenario).
- `precision_at_3`: `2/3 = 66.7%` on semantic retrieval checks (below expected target band intent).
- `wrong_memory_recall_rate`: `1/4 = 25%` (semantic scenario produced one materially incorrect recall).
- `trace_completeness_rate`: `3/4 = 75%` (one session missing full retrieval selection evidence).
- `review_gate_bypass_rate`: `0%` (no medium/high-risk bypass observed).
- `p95_retrieval_latency_ms`: `470` (from semantic scenario; high relative to other runs).

## Failure Classification Summary
- `retrieval`: 1 (SCN-SEM-01)
- `operator-ux`: 1 (SCN-STATE-01)
- `policy`: 0
- `memory-quality`: 0

## Prioritized Follow-On Action
1. `P1` retrieval precision and trace completeness hardening for semantic scenario path.
   - Created intake story:
     - `backlog/engineering/intake/STORY-20260222-dogfood-semantic-retrieval-hardening-v01.md`

## Notes
- This run is a manual baseline and intentionally keeps scoring lightweight while telemetry emission is still being completed.
