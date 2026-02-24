# Release Bundle: v0.1-initial-2026-02-22

## Decision
- `ship`

## Scope
- Included stories:
  - `delivery-backlog/engineering/done/STORY-20260222-architecture-baseline-map-drift-repair.md`
  - `delivery-backlog/engineering/done/STORY-20260222-dogfood-scenario-pack-v01.md`
  - `delivery-backlog/engineering/done/STORY-20260222-dogfood-semantic-retrieval-hardening-v01.md`
  - `delivery-backlog/engineering/done/STORY-20260222-kpi-snapshot-baseline-v01.md`
  - `delivery-backlog/engineering/done/STORY-20260222-kpi-delta-snapshot-after-semantic-hardening-v01.md`
  - `delivery-backlog/engineering/done/STORY-20260222-memory-cli-telemetry-contract-v01.md`
  - `delivery-backlog/engineering/done/STORY-20260222-mvp-constraint-enforcement-v01.md`
  - `delivery-backlog/engineering/done/STORY-20260222-release-checkpoint-bundle-v01.md`
- Included bugs:
  - `delivery-backlog/engineering/done/BUG-20260222-launch-stage-readme-parse-fallback.md`
- Excluded deferred items:
  - `delivery-backlog/engineering/done/STORY-20260222-memory-api-read-gateway-v03.md` (`release_checkpoint: deferred`)
  - `delivery-backlog/engineering/done/STORY-20260222-memory-snapshot-mvp-implementation-v02.md` (`release_checkpoint: deferred`)

## Evidence
- QA result artifacts:
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-architecture-baseline-map-drift-repair.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-dogfood-scenario-pack-v01.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-dogfood-semantic-retrieval-hardening-v01.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-kpi-snapshot-baseline-v01.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-kpi-delta-snapshot-after-semantic-hardening-v01.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-memory-cli-telemetry-contract-v01.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-mvp-constraint-enforcement-v01.md`
  - `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-release-checkpoint-bundle-v01.md`
- Validation commands/results:
  - `go test ./cmd/memory-cli` -> `PASS`
  - `tools/run_doc_tests.sh` -> `PASS`
- Program board snapshot reference:
  - `product-research/roadmap/PROGRAM_STATE_BOARD.md`

## Risk and Rollback
- Known risks:
  - Release decision is constrained to founder/operator early-alpha usage; broader rollout and v0.2 scale claims remain out of scope.
  - Deferred items (`memory-api-read-gateway-v03`, `memory-snapshot-mvp-implementation-v02`) are excluded and must stay outside this ship boundary.
- Rollback direction:
  - Revert release decision to `hold` if KPI deltas regress below ADR-0008 target-band intent in the next dogfood cycle.
  - Roll back by operating at checkpoint boundary: pause ship claims, keep backlog in quality-gated mode, and prioritize corrective hardening stories.

## Outcome Signals
- Baseline metric snapshot:
  - `operating-system/metrics/kpi-snapshot-2026-02-22-baseline.md`
- Post-hardening delta snapshot:
  - `operating-system/metrics/kpi-snapshot-2026-02-22-delta-post-hardening.md`
- Expected trend direction:
  - Sustain `Green` proxy status for memory precision and trace completeness while maintaining governance adherence and bounded latency.

## Notes
- This bundle keeps the `done` vs `shipped` boundary explicit and records the founder-early-alpha ship decision for v0.1 scope only.
