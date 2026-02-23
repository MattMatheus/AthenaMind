# ADR-0004: Adopt v0.1 MVP Constraints Framework and Deployment Gate

## Status
Accepted

## Context
Vision and pillars require bounded cost, reliable operation, and auditable memory behavior, but they do not provide full numeric thresholds.

## Decision
Adopt a mandatory constraints framework across four domains:
- Cost
- Latency
- Traceability
- Reliability

Adopt a deployment gate:
- If required thresholds are unset, autonomous memory-writing mode is disabled.
- Constraint keys and enforcement points follow `product-research/architecture/MVP_CONSTRAINTS_V01.md`.

## Consequences
- Positive:
  - Prevents unsafe autonomy without explicit limits
  - Converts qualitative goals into enforceable system behavior
  - Keeps threshold-setting as a deliberate product decision
- Negative:
  - Requires configuration and policy setup before full autonomy
  - Adds validation overhead in early iterations
- Neutral:
  - Numeric values remain intentionally open pending founder decisions

## Alternatives Considered
- Option A: Set default numeric thresholds now
  - Rejected because current source material does not define enough values to avoid guessing
- Option B: Rely on observability only, without hard enforcement
  - Rejected because it allows uncontrolled failure/cost modes

## Validation Plan
- Implement constraint validation checks in planning docs for each module story.
- Prove that autonomous mode refuses to start when required thresholds are absent.
