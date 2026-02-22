# Handoff: STORY-20260222-qa-regression-rubric

## What Changed
- Added deterministic QA rubric:
  - `backlog/QA_REGRESSION_RUBRIC.md`
- Rubric includes:
  - pass/fail gates for `backlog/qa` stories
  - P0/P1/P2/P3 mapping aligned to `BUG_TEMPLATE.md`
  - minimum evidence requirements for bug filing
  - pass/fail handoff examples
- Updated QA seed prompt to explicitly apply rubric:
  - `prompts/active/qa-agent-seed-prompt.md`
- Added story validation test:
  - `scripts/test_qa_regression_rubric.sh`
- Updated canonical docs test runner:
  - `scripts/run_doc_tests.sh`
- Moved story from `backlog/active/` to `backlog/qa/` and updated status to `qa`.
- Updated active queue:
  - `backlog/active/README.md`

## Why It Changed
- QA decisions needed deterministic gates and consistent severity mapping to reduce pass/fail variability across cycles.

## Test Updates Made
- Added `scripts/test_qa_regression_rubric.sh` to verify rubric existence/content, prompt reference, and priority mapping linkage.

## Test Run Results
- Command: `scripts/run_doc_tests.sh`
- Result: PASS
- Additional command: `scripts/test_qa_regression_rubric.sh`
- Result: PASS

## Open Risks/Questions
- Rubric remains documentation-driven; consistent operator use is required for full benefit.

## Recommended QA Focus Areas
- Verify pass/fail gates are strict enough to prevent ambiguous QA outcomes.
- Verify severity mapping language matches bug filing practice.
- Verify pass/fail examples are usable as templates for real reviews.

## New Gaps Discovered During Implementation
- None.
