# Engineering Handoff: STORY-20260222-new-idea-workflow-bootstrap

## What Changed
- Added intake validation command `/Users/foundry/Source/orchestrator/AthenaMind/tools/validate_intake_items.sh`.
- Added doc-test wrapper `/Users/foundry/Source/orchestrator/AthenaMind/tools/test_intake_validation.sh` and wired it into `/Users/foundry/Source/orchestrator/AthenaMind/tools/run_doc_tests.sh`.
- Implemented checks for:
  - required metadata fields,
  - valid status values,
  - lane-separation violations (`ARCH-*` in engineering intake and `STORY-*`/`BUG-*` in architecture intake).
- Added PM failure-handling instructions in:
  - `/Users/foundry/Source/orchestrator/AthenaMind/stage-prompts/active/pm-refinement-seed-prompt.md`
  - `/Users/foundry/Source/orchestrator/AthenaMind/delivery-backlog/architecture/INTAKE_REFINEMENT_GUIDE.md`
  - `/Users/foundry/Source/orchestrator/AthenaMind/HUMANS.md`
  - `/Users/foundry/Source/orchestrator/AthenaMind/DEVELOPMENT_CYCLE.md`

## Why It Changed
- Implements intake workflow guardrails so malformed cards and lane crossover are detected early in PM refinement.

## Test Updates Made
- Added `tools/test_intake_validation.sh` and integrated it into canonical docs test run.

## Test Run Results
- `go test ./...` -> PASS
- `tools/run_doc_tests.sh` -> PASS

## Open Risks/Questions
- Current lane checks are pattern-based and catch obvious crossover; nuanced misclassification still depends on PM judgment.

## Recommended QA Focus Areas
- Validate script behavior against malformed metadata and misfiled IDs.
- Confirm PM guidance is clear and executable for failure handling.

## New Gaps Discovered During Implementation
- None.
