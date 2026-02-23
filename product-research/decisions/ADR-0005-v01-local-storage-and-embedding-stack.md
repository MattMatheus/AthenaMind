# ADR-0005: Adopt v0.1 Local-First Storage and Embedding Stack

## Status
Accepted

## Context
The project vision and pillars emphasize local-first operation with optional cloud expansion. v0.1 requires concrete defaults to prevent architecture drift.

## Decision
For v0.1, adopt:
- State storage: SQLite (WAL mode)
- Semantic retrieval: FAISS exact-first index strategy
- Abstraction: adapter interfaces that permit future cloud-backed providers (for v0.3)

Additional policy:
- Use advanced FAISS ANN structures only after measured need from scorecard constraints.
- Keep provider-specific cloud logic outside core memory modules.

## Consequences
- Positive:
  - Fastest path to local reliability and low operational overhead
  - Strong alignment with local-first/offline-capable goals
  - Preserves cloud optionality without early lock-in
- Negative:
  - Potential scale ceiling until ANN/cloud modes are introduced
  - Requires careful packaging/testing of local vector dependencies
- Neutral:
  - Cloud provider choice remains open under adapter contracts

## Alternatives Considered
- Option A: Cloud-first storage/retrieval in v0.1
  - Rejected due to local-first requirement and higher operational burden
- Option B: ANN-heavy index default from day one
  - Rejected due to insufficient workload evidence

## Validation Plan
- Instrument retrieval latency/quality metrics under v0.1 workloads.
- Define explicit trigger criteria for ANN or cloud migration in a follow-up ADR.
