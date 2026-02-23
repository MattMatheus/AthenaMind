# Concepts

## Summary
Core concepts that explain how AthenaMind stores, governs, and retrieves memory artifacts.

## Intended Audience
- Product users who need mental models before running workflows.
- Engineers validating behavior against architecture decisions.

## Preconditions
- Familiarity with basic CLI execution.
- Awareness that AthenaMind v0.1 scope is memory-layer only.

## Main Flow
1. Read memory model and lifecycle in `memory-lifecycle.md`.
2. Read retrieval quality contract in `retrieval-and-quality.md`.
3. Read governance constraints in `governance.md`.
4. Map each concept to linked ADRs.

## Failure Modes
- Confusing implementation detail with product contract behavior.
- Interpreting `done` backlog state as shipped product behavior.
- Using unsupported runtime orchestration assumptions.

## References
- `product-research/decisions/ADR-0001-three-pillar-foundation.md`
- `product-research/decisions/ADR-0011-memory-mutation-review-workflow-contract.md`
- `product-research/decisions/ADR-0012-semantic-retrieval-quality-gates-v1.md`
- `product-research/decisions/ADR-0016-memory-snapshot-lifecycle-contract-v02.md`

## Pages
- `memory-lifecycle.md`
- `retrieval-and-quality.md`
- `governance.md`
