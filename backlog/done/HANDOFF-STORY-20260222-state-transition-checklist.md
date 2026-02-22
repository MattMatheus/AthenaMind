# Handoff: STORY-20260222-state-transition-checklist

## What Changed
- Added canonical transition checklist artifact:
  - `backlog/STATE_TRANSITION_CHECKLIST.md`
- Checklist defines required artifacts, approvals, and gate decisions for:
  - `active -> qa`
  - `qa -> done`
  - `qa -> active` (defect return path)
- Updated stage prompts to reference checklist:
  - `prompts/active/next-agent-seed-prompt.md`
  - `prompts/active/qa-agent-seed-prompt.md`
  - `prompts/active/pm-refinement-seed-prompt.md`
- Added story validation test:
  - `scripts/test_state_transition_checklist.sh`
- Updated canonical docs test runner:
  - `scripts/run_doc_tests.sh`
- Moved story from `backlog/active/` to `backlog/qa/` and updated status to `qa`.
- Updated active queue (`backlog/active/README.md`) to empty.

## Why It Changed
- State movement rules needed explicit gates to reduce transition inconsistency and enforce required evidence across engineering, QA, and PM cycles.

## Test Updates Made
- Added `scripts/test_state_transition_checklist.sh` to verify checklist existence/content and prompt references.

## Test Run Results
- Command: `scripts/run_doc_tests.sh`
- Result: PASS
- Additional command: `scripts/test_state_transition_checklist.sh`
- Result: PASS

## Open Risks/Questions
- Checklist enforcement is process-driven and depends on operators following prompt constraints.

## Recommended QA Focus Areas
- Verify transition requirements are complete for pass and fail paths.
- Verify prompt references are clear and non-ambiguous for each stage.
- Verify empty active queue behavior remains correct with updated `backlog/active/README.md`.

## New Gaps Discovered During Implementation
- None.
