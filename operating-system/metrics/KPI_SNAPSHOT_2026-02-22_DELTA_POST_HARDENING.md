# KPI Snapshot: 2026-02-22 (Post-Semantic-Hardening Delta)

## Summary
- Published follow-up KPI evidence after semantic hardening to compare against the baseline snapshot.
- Memory precision and trace completeness proxy metrics improved from `Red` to `Green` bands against ADR-0008 thresholds.
- Governance adherence remained stable (`review_gate_bypass_rate=0%`) while semantic retrieval latency improved.

## Baseline Reference
- Baseline KPI snapshot:
  - `operating-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md`
- Hardening run evidence:
  - `operating-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22-HARDENING.md`

## Before/After Deltas (Impacted Metrics)
| metric | baseline | post-hardening | delta |
| --- | --- | --- | --- |
| memory_precision_proxy (`precision_at_3`) | `66.7%` | `100%` | `+33.3pp` |
| trace_completeness | `75%` | `100%` | `+25pp` |
| wrong_memory_recall_rate | `25%` | `0%` | `-25pp` |
| semantic_retrieval_p95_latency_ms | `470` | `355` | `-115ms` |

## Updated ADR-0008 Band Interpretation
- Memory precision proxy (`precision_at_3 = 100%`):
  - Baseline: `Red`
  - Post-hardening: `Green` (`>= 96%`)
- Trace completeness (`100%`):
  - Baseline: `Red`
  - Post-hardening: `Green` (`>= 99%`)
- Memory error proxy (`wrong_memory_recall_rate = 0%`):
  - Baseline: `Red` proxy signal (`25%`)
  - Post-hardening: `Green` proxy signal (`0%`, below `<2%` error-rate intent)

## Interpretation
- Post-hardening evidence closes the highest-risk v0.1 memory-quality signal from the baseline run.
- Remaining release decision work is control-plane: refresh release bundle scope/evidence and record an explicit ship/hold decision from current data.

## Actions
- Execute `STORY-20260222-release-checkpoint-refresh-v01` to update release bundle scope and decision with current KPI evidence.
- Keep founder/operator messaging aligned to release checkpoint state (`done` vs `shipped` boundary remains explicit).
