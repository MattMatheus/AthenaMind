# ADR-0012: Semantic Retrieval Quality Gates (v1)

## Status
Accepted

## Context
The v1 memory CLI architecture requires semantic retrieval with deterministic fallback, but quality expectations were qualitative and not directly testable in QA.

## Decision
Adopt quantitative retrieval quality gates for v1:

Primary semantic quality gate:
- `Top-1 Useful Rate >= 80%` on a representative fixed query set (minimum 50 queries).
  - "Useful" means the top selected result satisfies query intent without requiring fallback.

Fallback interaction gates:
- When semantic confidence is below threshold or ambiguity is detected:
  - fallback must execute deterministically in required order.
- `Fallback Determinism = 100%` for repeated identical inputs on same corpus/config.

Reliability gates:
- `Selection Mode Reporting = 100%`:
  - each retrieval response must include `selection_mode`.
- `Source Trace Completeness = 100%`:
  - each response includes stable `source_path` and `selected_id`.

Operational method:
- Evaluate on pinned benchmark corpus and query set version.
- Compute metrics in CI/QA harness with reproducible config.

## Consequences
- Positive:
  - QA can make consistent pass/fail decisions.
  - Engineering can tune retrieval against explicit targets.
- Negative:
  - Building/maintaining benchmark harness adds overhead.
- Neutral:
  - Thresholds can be revised only through follow-on ADR.

## Alternatives Considered
- Option A: qualitative "usually useful" only
  - Rejected due to non-deterministic QA outcomes.
- Option B: much stricter initial threshold (`>= 90%`)
  - Rejected as high schedule risk for initial v1 rollout.

## Validation Plan
- Architecture artifact defines benchmark and harness contract.
- QA validates threshold results and fallback determinism before acceptance.
