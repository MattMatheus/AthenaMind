# ADR-0014: v0.1 Constraint Targets and Freeze Policy

## Status
Accepted

## Context
`product-research/architecture/MVP_CONSTRAINTS_V01.md` defined constraint domains and required keys but left operational targets open, which blocked deterministic implementation sequencing and QA.

## Decision
Freeze v0.1 baseline targets and owners for cost, latency, traceability, and reliability controls.

Baseline values:
- Cost:
  - `cost.max_per_run_usd`: `0.50`
  - `cost.max_period_spend_usd`: `25.00`
  - `cost.period_window`: `7d`
  - `cost.abort_policy`: `fail_closed`
- Latency:
  - `latency.p95_bootstrap_ms`: `1500`
  - `latency.p95_retrieval_ms`: `700`
  - `latency.p95_writeback_ms`: `1200`
  - `latency.degradation_policy`: `fallback_then_block`
- Traceability:
  - `trace.required_events`: `retrieve,write,evaluate,policy_gate_decision`
  - `trace.retention_days`: `30`
  - `trace.redaction_policy`: `strict_no_secret_payloads`
- Reliability:
  - `reliability.slo_targets`: `retrieve_success_rate>=99.0%, write_success_rate>=99.5%`
  - `reliability.error_budget_window`: `30d`
  - `reliability.freeze_policy`: `freeze_on_budget_exhaustion`
  - `reliability.override_roles`: `Product Manager - Maya, SRE - Nia, Founder/Operator`

Owner:
- Operational owner for constraint policy: `staff-personas/SRE - Nia.md`

## Policy Actions
- Cost budget breach: block write-affecting operations (`fail_closed`).
- Latency degradation: force deterministic fallback and mark for review.
- Missing required traceability fields: reject operation.
- Reliability budget exhaustion: freeze autonomous promotion paths until approved override.

## Consequences
- Positive:
  - Enables implementation and QA without unresolved numeric/policy ambiguity.
  - Makes freeze behavior explicit and auditable.
- Negative:
  - Initial targets may need recalibration after KPI baselines are published.
  - Tight fail-closed settings can reduce throughput early.
- Neutral:
  - Target changes remain allowed via follow-on ADR revision.

## Alternatives Considered
- Keep founder-input placeholders open until later.
  - Rejected due to sequencing and validation blockage.
- Use permissive defaults to avoid early freezes.
  - Rejected due to governance drift risk.

## Follow-On Implementation Paths
- `delivery-backlog/engineering/intake/STORY-20260222-mvp-constraint-enforcement-v01.md`
- `delivery-backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- `delivery-backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
