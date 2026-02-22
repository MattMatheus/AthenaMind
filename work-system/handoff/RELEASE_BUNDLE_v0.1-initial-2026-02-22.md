# Release Bundle: v0.1-initial-2026-02-22

## Decision
- `hold`

## Scope
- Included stories:
  - `backlog/engineering/done/STORY-20260222-architecture-baseline-map-drift-repair.md`
  - `backlog/engineering/done/STORY-20260222-dogfood-scenario-pack-v01.md`
  - `backlog/engineering/done/STORY-20260222-kpi-snapshot-baseline-v01.md`
  - `backlog/engineering/done/STORY-20260222-memory-cli-telemetry-contract-v01.md`
  - `backlog/engineering/done/STORY-20260222-mvp-constraint-enforcement-v01.md`
- Included bugs:
  - `backlog/engineering/done/BUG-20260222-launch-stage-readme-parse-fallback.md`
- Excluded deferred items:
  - `backlog/engineering/done/STORY-20260222-memory-api-read-gateway-v03.md` (`release_checkpoint: deferred`)
  - `backlog/engineering/done/STORY-20260222-memory-snapshot-mvp-implementation-v02.md` (`release_checkpoint: deferred`)
  - `backlog/engineering/active/STORY-20260222-release-checkpoint-bundle-v01.md` (control-plane story still in progress)

## Evidence
- QA result artifacts:
  - `backlog/engineering/done/QA-RESULT-STORY-20260222-architecture-baseline-map-drift-repair.md`
  - `backlog/engineering/done/QA-RESULT-STORY-20260222-dogfood-scenario-pack-v01.md`
  - `backlog/engineering/done/QA-RESULT-STORY-20260222-kpi-snapshot-baseline-v01.md`
  - `backlog/engineering/done/QA-RESULT-STORY-20260222-memory-cli-telemetry-contract-v01.md`
  - `backlog/engineering/done/QA-RESULT-STORY-20260222-mvp-constraint-enforcement-v01.md`
- Validation commands/results:
  - `go test ./cmd/memory-cli` -> `PASS`
  - `scripts/run_doc_tests.sh` -> `PASS`
- Program board snapshot reference:
  - `research/roadmap/PROGRAM_STATE_BOARD.md`

## Risk and Rollback
- Known risks:
  - One required control-plane story (`release-checkpoint-bundle-v01`) is still active, so ship boundary is not yet closed.
  - Retrieval quality and trace completeness remain below target band in baseline KPI evidence.
- Rollback direction:
  - Keep decision at `hold`; continue active queue until checkpoint story is done and KPI deltas improve.

## Outcome Signals
- Baseline metric snapshot:
  - `work-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md`
- Expected trend direction:
  - Move retrieval quality and traceability from `Red` toward `Yellow/Green` via telemetry + constraint enforcement follow-ons.

## Notes
- This bundle establishes the initial `done` vs `shipped` boundary and must be updated whenever release-bound scope changes.
