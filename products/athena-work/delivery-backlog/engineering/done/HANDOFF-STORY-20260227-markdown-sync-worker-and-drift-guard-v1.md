# Engineering Handoff: STORY-20260227-markdown-sync-worker-and-drift-guard-v1

## What Changed
- Added canonical backend read-model state artifact:
  - `products/athena-work/operating-system/state/backend_read_model_v1.json`
- Added deterministic markdown projection worker with dry-run support:
  - `tools/markdown_sync_worker.sh`
- Added drift guard with deterministic conflict classes and remediation hints:
  - `tools/check_markdown_drift.sh`
- Worker projections now include:
  - engineering active queue projection with `projection_version`
  - `operating-system/observer/LATEST_BOARD_READ_MODEL.md`
  - `operating-system/observer/LATEST_TIMELINE_READ_MODEL.md`
- API read model now loads canonical state artifact:
  - `products/athena-work/ui/local_control_plane_api.py`
- Added drift/sync doc test coverage:
  - `tools/test_markdown_sync_worker_and_drift_guard_v1.sh`
  - `tools/run_doc_tests.sh`

## Why It Changed
- Enforce backend-authoritative, deterministic markdown projection.
- Detect critical divergence (`ordering_conflict`, `missing_artifact`, `stale_revision`) with clear operator remediation.
- Preserve backlog/observer artifact paths used by launch and observer workflows.

## Test Updates Made
- Added sync worker and drift guard contract checks.
- Added dry-run status assertions and post-projection pass assertions.
- Added canonical state and projected artifact existence checks.

## Test Run Results
- `tools/run_doc_tests.sh` -> PASS
- `go test ./...` (repo root) -> PASS

## Open Risks/Questions
- Current canonical state file is local artifact and needs ownership discipline to prevent manual drift.
- Future backend runtime should become source of truth for this state file projection.

## Recommended QA Focus Areas
- Verify drift guard failure messaging is actionable in real operator drift scenarios.
- Verify projection output remains stable across repeated runs.
- Verify queue-empty state behaves correctly in stage launch flow.

## New Gaps Discovered
- none
