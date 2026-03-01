# Engineering Handoff: STORY-20260227-human-planning-workbench-v1

## What Changed
- Expanded workspace UI to include single-screen planning workbench:
  - fields: `direction`, `constraints`, `risks`, `next_stage`
  - explicit `Confirm direction` action with immediate status feedback
  - explicit plain-language `What happens next` summary
  - `products/athena-work/ui/index.html`
- Added planning export API endpoint to eliminate manual copy/paste:
  - `POST /api/v1/planning/export`
  - writes markdown artifact to `product-research/planning/sessions/PLAN-WORKBENCH-*.md`
  - `products/athena-work/ui/local_control_plane_api.py`
- Added planning-workbench doc test and wired into canonical runner:
  - `tools/test_human_planning_workbench_v1.sh`
  - `tools/run_doc_tests.sh`

## Why It Changed
- Reduce planning cognitive load and context switching for human operators.
- Ensure planning outputs are directly exportable into workflow artifacts.
- Preserve low-vision-first readability and direct UX.

## Test Updates Made
- Added assertions for:
  - required one-screen planning fields
  - explicit confirmation action and immediate status
  - explicit "what happens next" section
  - export endpoint contract and validation rules

## Test Run Results
- `tools/run_doc_tests.sh` -> PASS
- `go test ./...` (repo root) -> PASS

## Open Risks/Questions
- Export endpoint currently appends timestamped artifacts and does not deduplicate semantically equivalent exports.
- Server-side artifact path is local-first and assumes compose-mounted workspace.

## Recommended QA Focus Areas
- Verify first-time operator completion of planning flow in <=10 minutes.
- Verify export artifact structure is directly usable by downstream workflow steps.
- Verify confirmation status remains clear after repeated edits/exports.

## New Gaps Discovered
- none
