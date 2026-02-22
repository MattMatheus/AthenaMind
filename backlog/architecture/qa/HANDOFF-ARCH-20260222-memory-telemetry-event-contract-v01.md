# Architecture Handoff: ARCH-20260222-memory-telemetry-event-contract-v01

## Decision(s) Made
- Accepted v0.1 telemetry contract decision in:
  - `research/decisions/ADR-0015-v01-memory-telemetry-event-contract.md`
- Defined canonical architecture artifact with required fields, validation rules, producer/consumer map, and KPI traceability mapping:
  - `research/architecture/MEMORY_TELEMETRY_EVENT_CONTRACT_V01.md`
- Linked roadmap and architecture references to canonical telemetry contract artifacts.

## Alternatives Considered
- Keep telemetry schema guidance only in roadmap text.
  - Rejected due to drift and implementation ambiguity.
- Start with broad instrumentation schema.
  - Rejected due to v0.1 minimal-scope constraints.

## Risks and Mitigations
- Risk: strict required fields increase implementation overhead.
  - Mitigation: scope fixed to KPI-critical minimum and manual scoring fallback.
- Risk: KPI miscalculation if producers diverge from schema.
  - Mitigation: explicit validation rules and contract tests in follow-on engineering stories.
- Risk: retention/redaction non-compliance.
  - Mitigation: baseline retention and redaction policy are now explicit in contract.

## Updated Artifacts
- `research/decisions/ADR-0015-v01-memory-telemetry-event-contract.md`
- `research/architecture/MEMORY_TELEMETRY_EVENT_CONTRACT_V01.md`
- `research/roadmap/DOGFOODING_LOOP_AND_TELEMETRY_KPI_SET_V01.md`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `backlog/architecture/qa/ARCH-20260222-memory-telemetry-event-contract-v01.md`

## Follow-On Implementation Story Paths
- `backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- `backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
- `backlog/engineering/intake/STORY-20260222-dogfood-scenario-pack-v01.md`

## Validation Commands and Results
- `scripts/validate_intake_items.sh` -> PASS
- `scripts/run_doc_tests.sh` -> PASS

## Open Questions for QA Focus
- Are all KPI template fields computable from required event fields?
- Are retrieval provenance requirements (`selected_id`, `selection_mode`, `source_path`) enforced clearly enough for tests?
- Does the contract remain minimal without missing mandatory governance data?
