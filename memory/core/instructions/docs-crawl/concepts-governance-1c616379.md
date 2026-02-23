# Concepts Governance 1c616379

# Governance

## Summary
Memory mutation is controlled by policy gates, reviewer evidence, and audit artifacts.

## Mutation Rules
- Writes are blocked during autonomous runs (`AUTONOMOUS_RUN=true`).
- Allowed mutation stages:
  - `planning`
  - `architect`
  - `pm`
- Required review evidence:
  - `--reviewer`
  - `--decision approved|rejected`
  - `--reason`
  - `--risk`
  - `--notes`

## Rejection Rules
- Rejected decisions require:
  - `--rework-notes`
  - `--re-reviewed-by`
- Rejected writes create audit artifacts and do not apply content changes.

## Policy/Error Signals
- `ERR_MUTATION_NOT_ALLOWED_DURING_AUTONOMOUS_RUN`
- `ERR_MUTATION_STAGE_INVALID`
- `ERR_MUTATION_EVIDENCE_REQUIRED`
- `ERR_MUTATION_REJECTED_PENDING_REWORK`

## References
- `product-research/decisions/ADR-0011-memory-mutation-review-workflow-contract.md`
- `cmd/memory-cli/main.go`
