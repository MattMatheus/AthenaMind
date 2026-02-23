# ADR-0017: Memory API Wrapper Contract (v0.3)

## Status
Accepted

## Context
v1 memory architecture explicitly allows future API wrapping, but no accepted contract existed for parity, fallback, or rollback behavior.

## Decision
Adopt a v0.3 API wrapper contract with reversible behavior and CLI parity requirements.

Canonical artifact:
- `product-research/architecture/MEMORY_API_WRAPPER_CONTRACT_V03.md`

Decision highlights:
- API operations cover retrieve/write/evaluate.
- CLI is authoritative fallback path.
- Parity checks are mandatory for wrapper confidence.
- Rollback to CLI-only requires no migration.

## Consequences
- Positive:
  - Enables controlled API expansion without violating local-first defaults.
  - Provides explicit parity and fallback test expectations.
- Negative:
  - Adds contract-test and parity-validation overhead.
  - Introduces potential wrapper/CLI drift risk if governance is weak.
- Neutral:
  - Cloud/provider adapters remain optional follow-on work.

## Alternatives Considered
- Build API wrapper without explicit parity contract.
  - Rejected due to regression and hidden-divergence risk.
- Keep CLI-only forever with no API contract.
  - Rejected to preserve v0.3 expansion path.

## Follow-On Implementation Paths
- `delivery-backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md`
