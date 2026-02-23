# Memory API Wrapper Contract (v0.3)

## Purpose
Define a reversible API wrapper contract for memory operations that preserves CLI-first behavior as the authoritative fallback path.

## Scope
- In scope:
  - API contracts for `retrieve`, `write`, and `evaluate`.
  - Parity and fallback behavior relative to CLI.
  - Error semantics and compatibility guarantees.
- Out of scope:
  - Full cloud deployment architecture.
  - Provider-specific backend adapters beyond contract boundaries.

## Design Constraints
- Local-first remains default.
- Wrapper is optional and reversible.
- CLI remains canonical execution path for behavior parity.

## API Surface (v0.3 baseline)
### `POST /memory/retrieve`
Request:
- `query` (required)
- `domain` (optional)
- `session_id` (required)

Response:
- `selected_id`
- `selection_mode`
- `source_path`
- `confidence`
- `reason`
- `trace_id`

### `POST /memory/write`
Request:
- `id`, `title`, `type`, `domain`, `body`
- governance evidence fields (`reviewer`, `decision`, `reason`, `risk`, `notes`)
- `session_id`

Response:
- `status`
- `entry_id`
- `audit_record_path`
- `trace_id`

### `POST /memory/evaluate`
Request:
- `query_set` or `query_set_path`
- `corpus_id`, `query_set_id`, `config_id`
- `session_id`

Response:
- evaluation report payload matching CLI `evaluate` shape
- `trace_id`

## Parity Contract
1. For equivalent inputs/corpus, API outputs must match CLI response semantics.
2. `selection_mode`, `selected_id`, and `source_path` must be preserved for retrieve parity.
3. API wrapper must not introduce hidden policy bypass paths.

## Fallback and Rollback Contract
- If wrapper service is unavailable or parity checks fail:
  - Route to CLI path.
  - Emit explicit fallback event and reason code.
- Rollback to CLI-only mode:
  - Must not require data migration.
  - Must preserve request/response contracts at integration boundaries.

## Error Semantics
- `ERR_API_WRAPPER_UNAVAILABLE`
- `ERR_API_CLI_PARITY_MISMATCH`
- `ERR_API_POLICY_GATE_BLOCKED`
- `ERR_API_INPUT_INVALID`

## QA Validation Requirements
- Contract tests compare API and CLI output parity on pinned query sets.
- Failure-path tests verify deterministic fallback to CLI.
- Governance tests verify review gates remain enforced through API wrapper path.

## Follow-On Implementation
- `delivery-backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md`

## References
- `product-research/decisions/ADR-0017-memory-api-wrapper-contract-v03.md`
- `product-research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`
