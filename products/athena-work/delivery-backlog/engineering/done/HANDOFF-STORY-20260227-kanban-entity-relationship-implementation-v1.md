# Engineering Handoff: STORY-20260227-kanban-entity-relationship-implementation-v1

## Summary
- Implemented a canonical workspace entity model in the local control-plane API so board and timeline are derived from one deterministic relationship graph.
- Added relationship drift checks for contradictory lane membership, invalid `next_story` linkage, and incomplete timeline metadata.
- Added regression coverage proving consistent state passes and contradictory state fails.

## Changes Implemented
- `products/athena-work/ui/local_control_plane_api.py`
  - Added `build_workspace_entity_model(state)`.
  - Added deterministic normalizers for card/timeline inputs.
  - Board payload now derives lane/card data from the canonical model.
  - Timeline payload now derives events from the same canonical model.
  - Added `entity_model` metadata and relationship error count in board/timeline responses.
- `products/athena-work/tools/test_workspace_entity_relationship_model_v1.sh`
  - New regression test validating both clean and contradictory relationship scenarios.
- `products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh`
  - Extended checks for canonical entity model and relationship drift fields.
- `products/athena-work/tools/run_doc_tests.sh`
  - Added canonical entity model regression test to the standard suite.
- `products/athena-work/operating-system/state/backend_read_model_v1.json`
  - Synced canonical state to current queue projection for deterministic drift-guard behavior.

## Validation
- `./tools/test_workspace_entity_relationship_model_v1.sh` -> PASS
- `./tools/test_workspace_ui_read_only_board_v1.sh` -> PASS
- `./tools/test_markdown_sync_worker_and_drift_guard_v1.sh` -> PASS
- `./tools/run_doc_tests.sh` -> PASS

## Risks / Follow-ups
- Runtime state file (`state/runtime/backend_read_model_v1.local.json`) can still override canonical defaults for live sessions; keep projected queue and runtime state aligned during active demo cycles.
- Next story should continue left-anchored layout delivery and preserve entity-model compatibility.
