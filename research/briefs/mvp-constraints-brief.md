# Research Brief: MVP Constraints (Cost, Latency, Traceability, Reliability)

## Source Inputs
- `research/ingest/pillars.md`
- `README.md` (vision prompt)
- `research/decisions/ADR-0002-north-star-and-success-metrics.md`
- `research/architecture/MODULE_BOUNDARIES_V01.md`
- `research/references/external-sources-2026-02-22-mvp-constraints.md`

## Problem Statement
AthenaMind needs enforceable MVP constraints so autonomy remains trustworthy, bounded, and operationally predictable.

## Key Insights (80/20)
- Constraint categories are clear from source docs, but numeric thresholds are not fully specified.
- Hard enforcement at runtime/API is required; dashboards alone are insufficient.
- Without an explicit error-budget policy, reliability metrics will not drive decisions.
- Traceability must be complete enough to explain every memory-affecting action.

## MVP Constraint Classes
- Cost constraints:
  - Pre-run estimate required
  - Hard abort on budget breach
  - Periodic spend ceiling support (vision references avoiding large burn events)

- Latency constraints:
  - p95 targets required per key flow (bootstrap, retrieval, writeback)
  - Degradation mode if latency budgets are threatened

- Traceability constraints:
  - Mandatory audit chain for memory read/write/promotion decisions
  - Run/session correlation IDs required

- Reliability constraints:
  - Documented SLOs and error-budget policy required
  - Explicit freeze/escalation behavior on budget exhaustion

## What Can Be Decided Now (No Guessing)
- Constraint framework and enforcement points can be locked now.
- Numeric thresholds should remain configurable and gated by founder-approved values.

## What Cannot Be Decided Yet From Current Inputs
- Exact numeric SLO targets (availability/latency) for each flow.
- Spend window definition for the "$2k burn" ceiling (per run/day/week/month).
- Confidence thresholds that trigger human review for memory promotion.

## Recommended Decisions
- Adopt a strict "no threshold, no autonomous mode" policy.
- Require a `constraints.yaml` (or equivalent) as a deployment gate for v0.1.
