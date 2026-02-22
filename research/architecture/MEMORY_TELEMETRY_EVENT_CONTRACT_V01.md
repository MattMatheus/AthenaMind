# Memory Telemetry Event Contract (v0.1)

## Purpose
Define the minimal, implementation-ready telemetry event schema required to compute v0.1 KPI snapshots and enforce memory quality/governance checks.

## Scope
- In scope:
  - Memory CLI events for `write`, `retrieve`, and `evaluate`.
  - KPI-critical event fields only.
  - Deterministic validation rules for required fields.
- Out of scope:
  - External telemetry backend choice.
  - Non-memory runtime event schemas.

## Event Types
- `memory.retrieve`
- `memory.write`
- `memory.evaluate`
- `memory.policy_gate_decision`

## Required Fields
- `event_name`: one of the event types above.
- `event_version`: `1.0`
- `timestamp_utc`: RFC3339 UTC timestamp.
- `session_id`: stable session identifier.
- `trace_id`: correlation identifier for event chain.
- `scenario_id`: dogfooding scenario identifier.
- `operation`: `retrieve|write|evaluate|promote|reject|review`
- `result`: `success|fail`
- `policy_gate`: `none|low|medium|high`
- `memory_type`: `procedural|state|semantic`
- `latency_ms`: integer, `>= 0`
- `selected_id`: required when `operation=retrieve`.
- `selection_mode`: required when `operation=retrieve`; `semantic|fallback_exact_key|fallback_path_priority`
- `source_path`: required when `operation=retrieve`.
- `operator_verdict`: `correct|partially_correct|incorrect|not_scored`

## Optional Fields
- `error_code`: deterministic reason code on failures.
- `domain`: memory domain hint.
- `reviewer`: reviewer identifier for governed write decisions.
- `decision`: `approved|rejected` for review events.
- `reason`: short human-readable reason for decision/failure.

## Producer and Consumer Map
- Producers:
  - `cmd/memory-cli` command paths (`write`, `retrieve`, `evaluate`).
- Primary consumers:
  - KPI snapshot workflow (`work-system/metrics/KPI_SNAPSHOT_TEMPLATE.md`).
  - QA verification and release evidence bundles.

## Validation Rules
1. Events missing required fields are invalid and must be rejected by emitter-side validation.
2. `latency_ms` must be present and non-negative.
3. Retrieval events must always include `selected_id`, `selection_mode`, and `source_path`.
4. `operator_verdict` may be `not_scored` for initial manual scoring mode.
5. `event_version` increments only for schema-breaking changes.

## KPI Traceability Mapping
- `trace_completeness_rate`:
  - Uses completeness of `trace_id`, `operation`, `result`, and retrieval selection fields.
- `review_gate_bypass_rate`:
  - Uses `policy_gate` and `memory.policy_gate_decision` event streams.
- `p95_retrieval_latency_ms`:
  - Uses `memory.retrieve.latency_ms`.
- Retrieval quality gate pass evidence:
  - Uses `selection_mode`, `selected_id`, `source_path`, `operator_verdict`.

## Retention and Redaction Baseline
- Minimum retention: `30` days.
- Redaction policy: never include secret/token payloads; redact free-text fields when needed.

## References
- `research/decisions/ADR-0015-v01-memory-telemetry-event-contract.md`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `research/roadmap/DOGFOODING_LOOP_AND_TELEMETRY_KPI_SET_V01.md`
