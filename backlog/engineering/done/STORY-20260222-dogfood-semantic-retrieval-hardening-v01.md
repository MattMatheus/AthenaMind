# Story: Harden Semantic Retrieval Precision and Trace Completeness from Dogfood Baseline

## Metadata
- `id`: STORY-20260222-dogfood-semantic-retrieval-hardening-v01
- `owner_persona`: QA Engineer - Iris.md
- `status`: done
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0002, ADR-0008, ADR-0012, ADR-0015]
- `success_metric`: semantic scenario precision_at_3 improves from 66.7% to >= target band and trace completeness reaches >= 95%
- `release_checkpoint`: required

## Problem Statement
- Initial dogfood run indicates weak semantic retrieval precision and incomplete trace fields on one scenario path, reducing trust and auditability.

## Scope
- In:
  - Improve top-3 retrieval relevance for semantic scenario queries.
  - Ensure retrieval trace evidence fields are consistently populated.
  - Re-run `SCN-SEM-01` and capture KPI delta.
- Out:
  - Broad v0.2 retrieval algorithm redesign.
  - Non-memory runtime orchestration work.

## Acceptance Criteria
1. `SCN-SEM-01` precision_at_3 improves and no materially incorrect recall is observed.
2. Trace completeness for semantic retrieval events is >= 95%.
3. Updated dogfood run artifact records before/after KPI comparison.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `work-system/metrics/DOGFOOD_SCENARIO_PACK_V01.md`
- `work-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22.md`
- `research/architecture/MEMORY_TELEMETRY_EVENT_CONTRACT_V01.md`

## Notes
- Prioritized from first dogfood baseline run as highest-impact weakness for trusted memory behavior.
