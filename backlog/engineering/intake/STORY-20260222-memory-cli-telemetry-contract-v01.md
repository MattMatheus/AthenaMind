# Story: Implement v0.1 Memory CLI Telemetry Event Emission

## Metadata
- `id`: STORY-20260222-memory-cli-telemetry-contract-v01
- `owner_persona`: SRE - Nia.md
- `status`: intake
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0002, ADR-0006, ADR-0008, ADR-0012]
- `success_metric`: retrieval/write/evaluate paths emit telemetry fields required for first KPI and quality-gate reporting
- `release_checkpoint`: required

## Problem Statement
- KPI loop and governance reporting require event-level telemetry, but current memory CLI outputs are command-local and not normalized as a reusable event contract.

## Scope
- In:
  - Emit structured telemetry events for write/retrieve/evaluate flows.
  - Include required identifiers and policy/retrieval outcome fields.
  - Add tests validating presence and schema shape of emitted events.
- Out:
  - External telemetry backend integration.
  - Non-memory runtime telemetry.

## Acceptance Criteria
1. CLI emits schema-compliant events for successful and failure flows.
2. Event payload covers KPI-required fields and retrieval quality gate evidence.
3. Tests validate deterministic event emission and required-field completeness.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `backlog/architecture/intake/ARCH-20260222-memory-telemetry-event-contract-v01.md`
- `research/roadmap/DOGFOODING_LOOP_AND_TELEMETRY_KPI_SET_V01.md`

## Notes
- Keep implementation minimal and file/local-first for v0.1.
