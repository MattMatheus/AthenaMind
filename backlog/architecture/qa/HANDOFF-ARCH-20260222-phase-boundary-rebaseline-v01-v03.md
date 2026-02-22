# Architecture Handoff: ARCH-20260222-phase-boundary-rebaseline-v01-v03

## Decision(s) Made
- Accepted an explicit phase-boundary contract across `v0.1`, `v0.2`, and `v0.3` in:
  - `research/decisions/ADR-0013-phase-boundary-contract-v01-v03.md`
- Formalized module-to-phase allocation and ownership notes in:
  - `research/architecture/MODULE_BOUNDARIES_V01.md`
- Rebased architecture baseline map to current resolved and open work, removing stale unresolved references:
  - `research/architecture/ARCHITECTURE_BASELINE_MAP_V01.md`

## Alternatives Considered
- Keep phase ownership implicit across multiple docs.
  - Rejected due to repeated ranking ambiguity and drift.
- Push all refilled backlog directly to engineering before architecture rebaseline.
  - Rejected due to dependency/order risk and cross-phase scope leakage.

## Risks and Mitigations
- Risk: broad backlog refill can reduce focus on v0.1 closeout signals.
  - Mitigation: contract locks `v0.1` priority to release checkpoint and KPI baseline work before major v0.2 expansion.
- Risk: phase map can drift again as new ideas arrive.
  - Mitigation: baseline map now carries explicit open-work intake references and ADR-0013 as canonical decision source.
- Risk: confusion around `workspace-runtime` ownership.
  - Mitigation: explicit note that it is an external integration surface and not v0.1 core ownership.

## Updated Artifacts
- `research/decisions/ADR-0013-phase-boundary-contract-v01-v03.md`
- `research/architecture/MODULE_BOUNDARIES_V01.md`
- `research/architecture/ARCHITECTURE_BASELINE_MAP_V01.md`
- `backlog/architecture/qa/ARCH-20260222-phase-boundary-rebaseline-v01-v03.md`

## Follow-On Implementation Story Paths
- `backlog/engineering/intake/STORY-20260222-release-checkpoint-bundle-v01.md`
- `backlog/engineering/intake/STORY-20260222-kpi-snapshot-baseline-v01.md`
- `backlog/engineering/intake/STORY-20260222-architecture-baseline-map-drift-repair.md`
- `backlog/engineering/intake/STORY-20260222-memory-cli-telemetry-contract-v01.md`
- `backlog/engineering/intake/STORY-20260222-dogfood-scenario-pack-v01.md`
- `backlog/engineering/intake/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
- `backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md`

## Validation Commands and Results
- `scripts/run_doc_tests.sh` -> FAIL (program board queue counts stale before state sync)
- `scripts/validate_intake_items.sh` -> PASS
- `scripts/run_doc_tests.sh` -> PASS (after queue-count sync and state transition updates)

## Open Questions for QA Focus
- Are phase allocations internally consistent across ADR-0013, module boundaries, phase plan, and program board status?
- Does the updated baseline map fully eliminate stale unresolved-gap references?
- Is `workspace-runtime` treatment clearly scoped as non-v0.1 core ownership in all updated artifacts?
