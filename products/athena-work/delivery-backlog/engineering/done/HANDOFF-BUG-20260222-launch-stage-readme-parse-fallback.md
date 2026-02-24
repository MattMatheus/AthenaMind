# Handoff: BUG-20260222-launch-stage-readme-parse-fallback

## What Changed
- Fixed Active README sequence parsing in launcher to be portable across BSD and GNU sed:
  - `<repo>/tools/launch_stage.sh`
- Fixed path resolution for README entries that are bare filenames so they resolve under `delivery-backlog/active/` instead of repository root.
- Added regression test that verifies engineering launch respects README queue order:
  - `<repo>/tools/test_launch_stage_readme_queue.sh`
- Updated canonical docs test runner to include the new regression test:
  - `<repo>/tools/run_doc_tests.sh`
- Moved bug card from `delivery-backlog/active/` to `delivery-backlog/qa/` and updated status to `qa`.
- Updated active queue ordering in `<repo>/delivery-backlog/active/README.md`.

## Why It Changed
- Engineering stage launch could ignore PM-ranked queue ordering because the previous sed expression was not portable in BSD sed.
- After enabling README parsing, bare filename queue entries also needed correct path resolution to avoid false "story not found" failures.

## Test Updates Made
- Added: `<repo>/tools/test_launch_stage_readme_queue.sh`
- Updated: `<repo>/tools/run_doc_tests.sh`

## Test Run Results
- Command: `tools/run_doc_tests.sh`
- Result: PASS
- Additional command: `tools/launch_stage.sh engineering`
- Result: PASS (`story: delivery-backlog/active/BUG-20260222-launch-stage-readme-parse-fallback.md` before state move)

## Open Risks/Questions
- The regression test currently uses known story filenames from this workspace to force non-alphabetical ordering. If these names change, test setup will need to be updated.

## Recommended QA Focus Areas
- Confirm README ordering is honored when queue order differs from alphabetical order.
- Confirm launcher works with bare filename entries in `delivery-backlog/active/README.md`.
- Confirm fallback behavior still works if README sequence is missing or malformed.

## New Gaps Discovered During Implementation
- None.
