# Engineering Handoff: STORY-20260227-docs-workspace-surface-and-navigation-v1

## What Changed
- Added a dedicated `Docs Workspace` panel to the board UI with quick-jump doc cards.
- Added docs API read endpoints:
  - `GET /api/v1/docs/index` for document index metadata
  - `GET /api/v1/docs/view?id=<doc_id>` for inline document preview payload
- Linked core operational docs:
  - `products/athena-work/HUMANS.md`
  - `products/athena-work/DEVELOPMENT_CYCLE.md`
  - `knowledge-base/process/STAGE_EXIT_GATES.md`
  - `knowledge-base/operations/LOCAL_CONTROL_PLANE_QUICKSTART.md`
- Extended workspace UI doc test coverage for new docs routes and panel presence.

## Why It Changed
- Reduce navigation friction by making docs a first-class workspace surface.
- Keep policy/process guidance one click away from active board context.

## Test Updates Made
- Updated `products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh` with docs endpoint/panel assertions.

## Test Run Results
- `products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh` -> PASS
- `products/athena-work/tools/test_human_planning_workbench_v1.sh` -> PASS
- `python3 -m py_compile products/athena-work/ui/local_control_plane_api.py` -> PASS

## Open Risks/Questions
- Docs index is currently static allowlist-backed; future work may require dynamic relevance ranking.

## Recommended QA Focus Areas
- Verify docs panel usability on small laptop viewports.
- Verify keyboard-only docs navigation and preview reading flow.
- Verify operators can identify the right doc for common stage decisions in one click.

## New Gaps Discovered
- none
