# AthenaMind v0.1 MVP Constraints

## Objective
Define enforceable constraint contracts for cost, latency, traceability, and reliability.

## Constraint Policy
- Constraints are mandatory runtime controls, not optional telemetry.
- If required thresholds are unset, autonomous memory-writing mode must be disabled.

## Required Constraint Domains

### 1) Cost
Enforcement points:
- Pre-run estimator in `orchestrator-api`
- In-run budget checks in `workspace-runtime`
- Post-run spend accounting in `audit-telemetry`

Required config keys (must be explicitly set):
- `cost.max_per_run_usd`
- `cost.max_period_spend_usd`
- `cost.period_window`
- `cost.abort_policy`

Default v0.1 targets:
- `cost.max_per_run_usd`: `0.50`
- `cost.max_period_spend_usd`: `25.00`
- `cost.period_window`: `7d`
- `cost.abort_policy`: `fail_closed`

Owner:
- `personas/SRE - Nia.md`

### 2) Latency
Enforcement points:
- Request-level timers in `orchestrator-api`
- Retrieval timers in `semantic-memory`
- Runtime execution timers in `workspace-runtime`

Required config keys:
- `latency.p95_bootstrap_ms`
- `latency.p95_retrieval_ms`
- `latency.p95_writeback_ms`
- `latency.degradation_policy`

Default v0.1 targets:
- `latency.p95_bootstrap_ms`: `1500`
- `latency.p95_retrieval_ms`: `700`
- `latency.p95_writeback_ms`: `1200`
- `latency.degradation_policy`: `fallback_then_block`

Owner:
- `personas/SRE - Nia.md`

### 3) Traceability
Enforcement points:
- Mandatory event emission from each module to `audit-telemetry`
- Run/session/memory-object correlation IDs at API boundary

Required config keys:
- `trace.required_events`
- `trace.retention_days`
- `trace.redaction_policy`

Default v0.1 targets:
- `trace.required_events`: `retrieve,write,evaluate,policy_gate_decision`
- `trace.retention_days`: `30`
- `trace.redaction_policy`: `strict_no_secret_payloads`

Owner:
- `personas/SRE - Nia.md`

### 4) Reliability
Enforcement points:
- SLO computation in `audit-telemetry`
- Error-budget policy actions orchestrated by `orchestrator-api`

Required config keys:
- `reliability.slo_targets`
- `reliability.error_budget_window`
- `reliability.freeze_policy`
- `reliability.override_roles`

Default v0.1 targets:
- `reliability.slo_targets`: `retrieve_success_rate>=99.0%, write_success_rate>=99.5%`
- `reliability.error_budget_window`: `30d`
- `reliability.freeze_policy`: `freeze_on_budget_exhaustion`
- `reliability.override_roles`: `Product Manager - Maya, SRE - Nia, Founder/Operator`

Owner:
- `personas/SRE - Nia.md`

## v0.1 Deployment Gate
Autonomous mode (`memory promotion` + `unattended writeback`) is allowed only when:
- All required config keys are set
- Validation checks pass
- Error-budget policy document exists and is approved

## Policy Actions (Testable)
- Cost budget breach:
  - Action: block write-affecting operations and return explicit budget-exceeded error.
- Latency degradation:
  - Action: force deterministic fallback mode and mark session for review.
- Traceability missing required fields:
  - Action: reject operation as non-compliant.
- Reliability budget exhaustion:
  - Action: freeze autonomous promotion paths until override by allowed role.

## Founder Input Status
- This baseline freezes default target values for immediate implementation and QA.
- Founder/operator may revise targets through a follow-on ADR update, but values are no longer unspecified.
