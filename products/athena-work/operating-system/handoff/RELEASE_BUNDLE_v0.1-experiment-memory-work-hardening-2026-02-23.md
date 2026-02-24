# Release Bundle: v0.1-experiment-memory-work-hardening-2026-02-23

## Decision
- `ship` (candidate for merge to `main`)

## Scope
- Included stories:
  - Memory retrieval resilience hardening (skip invalid candidates, preserve deterministic retrieval with explicit warnings).
  - Embedding store resilience hardening (skip corrupt embedding rows without failing full reads).
  - Work-system reliability hardening (stage inference in observer write-back, duplicate intake id detection, branch guard override for isolated branch workflows).
  - Doc/test harness reliability hardening (lane-absence tolerant backlog consistency test).
- Included bugs:
  - Retrieval path failed whole query on partial metadata corruption.
  - Embedding record read failed whole result set on a single malformed vector row.
  - Program state consistency test failed hard when optional backlog lane directories were absent.
- Excluded deferred items:
  - Podman-backed advanced storage backend experiments (explicitly deferred in this cycle to keep scope low-risk and no-install).

## Evidence
- QA result artifacts:
  - `go test ./...` passed on branch `experiment`.
  - `products/athena-work/tools/run_doc_tests.sh` passed.
  - Focused integration checks passed:
    - `products/athena-work/tools/test_launch_stage_readme_queue.sh`
    - `products/athena-work/tools/test_launch_stage_memory_integration.sh`
    - `products/athena-work/tools/test_observer_cycle_memory_integration.sh`
    - `products/athena-work/tools/test_validate_intake_duplicate_ids.sh`
- Validation commands/results:
  - `go test ./internal/retrieval` -> pass
  - `go test ./internal/index` -> pass
  - `go test ./...` -> pass
  - `products/athena-work/tools/run_doc_tests.sh` -> pass
- Program board snapshot reference:
  - `product-research/roadmap/PROGRAM_STATE_BOARD.md`

## Risk and Rollback
- Known risks:
  - Retrieval and embedding reads now skip corrupt rows/files; this improves availability but can mask data-quality drift if warnings are ignored.
  - Work-system stage inference uses naming heuristics and may need explicit overrides if future cycle-id conventions change.
- Rollback direction:
  - Revert commits `5561e7e` and `3fca513` (or revert this merge commit) to restore previous behavior.

## Outcome Signals
- Baseline metric snapshot:
  - `products/athena-work/operating-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md`
- Expected trend direction:
  - Lower retrieval failure rate under partial corruption.
  - Lower operator friction in branch-isolated execution and observer cycle closure.
  - Higher intake hygiene due to duplicate-id rejection.

## Notes
- Decision log and milestone trace are captured in `PROGRESS.md`.
