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

### 3) Traceability
Enforcement points:
- Mandatory event emission from each module to `audit-telemetry`
- Run/session/memory-object correlation IDs at API boundary

Required config keys:
- `trace.required_events`
- `trace.retention_days`
- `trace.redaction_policy`

### 4) Reliability
Enforcement points:
- SLO computation in `audit-telemetry`
- Error-budget policy actions orchestrated by `orchestrator-api`

Required config keys:
- `reliability.slo_targets`
- `reliability.error_budget_window`
- `reliability.freeze_policy`
- `reliability.override_roles`

## v0.1 Deployment Gate
Autonomous mode (`memory promotion` + `unattended writeback`) is allowed only when:
- All required config keys are set
- Validation checks pass
- Error-budget policy document exists and is approved

## Open Inputs Required From Founder
- Numeric targets for SLOs/latency and spend windows
- Aggressiveness of freeze policy when reliability budget is exhausted
- Acceptable risk level for automated promotions before human review
