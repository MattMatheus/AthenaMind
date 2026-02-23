# ADR-0006: Governance and Human-in-the-Loop Policy

## Status
Accepted

## Context
AthenaMind requires explicit governance and human-review controls for trust and safety. Scope is memory-layer-only for v0.1/v1.0, with local-first single-operator usage.

## Decision
Adopt governance framework components:
- Policy decision points for promotion/conflict/autonomous side effects
- Risk-class model and review trigger matrix
- Mandatory policy decision audit events

Adopt the following v0.1 policy values:
- `autonomy.default_mode`: `review-first`
- `autonomy.user_override`: allowed to accelerate onboarding
- `memory_budget.policy`: `strict`
  - hard caps on memory operations (retrieval/query/write/context budget)
  - abort/deny on budget breach
- `review.required_risk_classes`: `medium` and `high`
- `override.roles`: `operator`
  - local-first single-user model: the operator is the only override authority in v0.1 and likely v1.0

## Consequences
- Positive:
  - Strong safety baseline for early product trust
  - Fast onboarding path via explicit operator override
  - Clear alignment with local single-operator product reality
- Negative:
  - Review-first and strict budgets can slow some workflows
  - Higher chance of manual interventions in early tuning
- Neutral:
  - Multi-operator/centralized memory governance deferred to v2.0 direction

## Alternatives Considered
- Option A: `auto-first` default
  - Rejected for v0.1 due to higher early governance risk
- Option B: soft memory budgets
  - Rejected due to weak protection against context/bloat failure modes
- Option C: delegated override roles in v0.1
  - Rejected because local single-operator model is current product scope

## Validation Plan
- Add policy simulation tests for allow/deny/review + operator override paths.
- Verify strict memory-budget enforcement with audit traces.
- Revisit role model when multi-user/swarm memory is introduced (v2.0 planning).
