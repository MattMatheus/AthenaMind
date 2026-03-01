# Engineering Handoff: STORY-20260227-docs-linkage-and-sync-projection-hardening-v1

## Summary
- Hardened docs workspace linkage so docs references are generated from canonical workspace state (including queue/story-linked docs).
- Strengthened markdown drift guard output with deterministic drift classes and remediation ids.
- Extended markdown projection to publish a deterministic docs index read model artifact.
- Added regression tests for canonical docs linkage and drift class/hint behavior.

## Changes Implemented
- `products/athena-work/ui/local_control_plane_api.py`
  - Added workspace path resolution across `/workspace` and local repo roots.
  - Replaced static docs index behavior with canonical-state-generated docs linkage (`build_docs_index(state)`).
  - Added deterministic auto-linking for `next_story`, engineering `active`, and engineering `qa` stories when files exist.
  - Updated doc view resolution to use canonical resolved workspace paths.
- `products/athena-work/tools/markdown_sync_worker.sh`
  - Added projected docs artifact output:
    - `products/athena-work/operating-system/observer/LATEST_DOCS_INDEX_READ_MODEL.md`
  - Included docs projection in dry-run drift checks and sync writes.
- `products/athena-work/tools/check_markdown_drift.sh`
  - Added deterministic drift metadata:
    - `drift_class: <class>`
    - `remediation_id: <id>`
  - Added env override inputs for regression harness:
    - `ATHENA_MD_STATE_FILE`
    - `ATHENA_MD_ACTIVE_README`
    - `ATHENA_MD_ACTIVE_STORY_DIR`
- `products/athena-work/tools/test_markdown_sync_worker_and_drift_guard_v1.sh`
  - Added checks for deterministic remediation ids and projected docs index artifact.
- `products/athena-work/tools/test_docs_workspace_linkage_v1.sh`
  - New regression test proving docs links are generated from canonical queue/story state.
- `products/athena-work/tools/test_markdown_drift_guard_regression_v1.sh`
  - New regression test covering deterministic ordering/stale/missing drift classes and remediation ids.
- `products/athena-work/tools/run_doc_tests.sh`
  - Added both new regression tests to canonical suite.

## Validation
- `./tools/test_docs_workspace_linkage_v1.sh` -> PASS
- `./tools/test_markdown_drift_guard_regression_v1.sh` -> PASS
- `./tools/test_markdown_sync_worker_and_drift_guard_v1.sh` -> PASS
- `./tools/run_doc_tests.sh` -> PASS

## Risks / Follow-ups
- Runtime state overrides can still affect live UI docs linkage if stale runtime files are intentionally loaded; keep canonical and runtime state synchronized during demo runs.
- QA can additionally verify docs linkage behavior through UI interaction against current queue state.
