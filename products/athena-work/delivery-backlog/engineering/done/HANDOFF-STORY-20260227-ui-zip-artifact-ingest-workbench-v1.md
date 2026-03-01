# Engineering Handoff: STORY-20260227-ui-zip-artifact-ingest-workbench-v1

## What Changed
- Added UI artifact-ingest workbench panel:
  - zip file chooser
  - ingest action
  - status and report output
- Added API endpoint:
  - `POST /api/v1/artifacts/ingest`
  - accepts `filename` + `bundle_base64`
  - validates `.zip` input, 10MB cap, and base64 integrity
  - invokes local ingest tool and returns structured ingest result
- Triggered board refresh after successful ingest so updated cards/timeline appear immediately.
- Expanded workspace UI doc-test coverage for ingest endpoint/panel assertions.

## Why It Changed
- Remove terminal-only requirement for artifact ingest and keep operator workflow inside the board UI.

## Test Updates Made
- Updated `products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh` with ingest endpoint and UI checks.

## Test Run Results
- `products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh` -> PASS
- `products/athena-work/tools/test_v1_v2_zip_artifact_ingest.sh` -> PASS
- `python3 -m py_compile products/athena-work/ui/local_control_plane_api.py` -> PASS

## Open Risks/Questions
- Upload path currently uses base64 JSON body; this is acceptable for current bundle sizes but may be tuned to multipart if payload sizes grow.

## Recommended QA Focus Areas
- Verify ingest failure messaging for invalid/non-zip uploads.
- Verify mobile viewport usability of ingest panel controls.
- Verify post-ingest board refresh under repeated runs.

## New Gaps Discovered
- none
