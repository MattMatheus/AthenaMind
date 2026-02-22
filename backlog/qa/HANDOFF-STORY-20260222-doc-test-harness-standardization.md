# Handoff: STORY-20260222-doc-test-harness-standardization

## What Changed
- Added shared docs test harness library:
  - `scripts/lib/doc_test_harness.sh`
- Added canonical docs test runner:
  - `scripts/run_doc_tests.sh`
- Added standard guidance and reusable template/checklist:
  - `research/roadmap/DOC_TEST_HARNESS_STANDARD.md`
- Updated engineering directive to reference canonical command:
  - `prompts/active/next-agent-seed-prompt.md`
- Added lifecycle guidance reference:
  - `DEVELOPMENT_CYCLE.md`
- Migrated existing doc-focused tests to shared harness pattern:
  - `scripts/test_goals_scorecard_v01.sh`
  - `scripts/test_phased_plan_v01_v03.sh`
  - `scripts/test_coding_readiness_gate.sh`
- Added story validation test:
  - `scripts/test_doc_test_harness_standardization.sh`
- Moved story from `backlog/active/` to `backlog/qa/`, updated status, and updated active queue ordering.

## Why It Changed
- Repo had ad hoc doc validation scripts with duplicated assertion logic and no canonical entrypoint.
- This implementation creates one reusable harness pattern and one default command path for engineering runs.

## Test Updates Made
- Added:
  - `scripts/test_doc_test_harness_standardization.sh`
- Updated:
  - `scripts/run_doc_tests.sh` to include all doc-focused tests
  - existing story tests to source shared harness library

## Test Run Results
- Command: `scripts/run_doc_tests.sh`
- Result: PASS
- Additional command: `scripts/test_doc_test_harness_standardization.sh`
- Result: PASS

## Open Risks/Questions
- Launcher currently may not respect `backlog/active/README.md` sequence due portable sed issue; tracked as intake bug.

## Recommended QA Focus Areas
- Verify canonical docs test command path is referenced in engineering prompt.
- Verify shared harness functions are used by migrated tests.
- Verify reusable template and pattern guidance are clear and actionable.
- Verify canonical runner executes all current doc-driven tests.

## New Gaps Discovered During Implementation
- `backlog/intake/BUG-20260222-launch-stage-readme-parse-fallback.md`
