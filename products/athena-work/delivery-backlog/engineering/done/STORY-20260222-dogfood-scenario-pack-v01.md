# Story: Build Initial Dogfooding Scenario Pack and Evaluation Loop

## Metadata
- `id`: STORY-20260222-dogfood-scenario-pack-v01
- `owner_persona`: QA Engineer - Iris.md
- `status`: qa
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0002, ADR-0008, ADR-0012]
- `success_metric`: at least 3 recurring memory scenarios run through a repeatable scoring loop with recorded outcomes
- `release_checkpoint`: required

## Problem Statement
- Dogfooding loop requires scenario-based evidence, but there is no canonical v0.1 scenario pack wired into repeatable evaluation.

## Scope
- In:
  - Define 3 to 5 recurring operator scenarios across procedural/state/semantic memory.
  - Create repeatable execution and scoring instructions.
  - Record first run outputs and classify failures.
- Out:
  - Full automation of scenario orchestration.
  - v0.2 retrieval tuning changes.

## Acceptance Criteria
1. Scenario pack is versioned and covers multiple memory types.
2. First run is executed and results are captured with KPI-relevant annotations.
3. At least one prioritized follow-on action is generated from observed failures or weak signals.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `delivery-backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- `operating-system/metrics/KPI_SNAPSHOT_TEMPLATE.md`

## Notes
- This story converts KPI intent into actionable weekly operations.
