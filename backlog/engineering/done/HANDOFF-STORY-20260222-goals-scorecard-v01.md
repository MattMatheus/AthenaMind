# Handoff: STORY-20260222-goals-scorecard-v01

## What Changed
- Added accepted decision record: `research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`.
- Defined frozen v0.1 goals and numeric target bands for all ADR-0002 metrics.
- Added explicit goal-to-metric mapping in ADR-0008.
- Updated planning references:
  - `PM-TODO.md` now marks goals/criteria finalization complete and links ADR-0008.
  - `research/roadmap/RESEARCH_BACKLOG.md` now records ADR-0008 as completed.
- Added story validation test: `scripts/test_goals_scorecard_v01.sh`.
- Moved story from `backlog/active/` to `backlog/qa/` and updated status to `qa`.
- Updated active queue order in `backlog/active/README.md`.

## Why It Changed
- The story required explicit v0.1 goals and scorecard target bands tied to ADR-0002 metrics.
- Planning docs needed an explicit decision reference so sprint planning and QA can use measurable thresholds instead of subjective interpretation.

## Test Updates Made
- Added `scripts/test_goals_scorecard_v01.sh` to verify:
  - ADR-0008 exists.
  - ADR-0008 includes all ADR-0002 metrics.
  - `PM-TODO.md` references ADR-0008.
  - `research/roadmap/RESEARCH_BACKLOG.md` references ADR-0008.

## Test Run Results
- Command: `AthenaMind/scripts/test_goals_scorecard_v01.sh`
- Result: PASS (11 checks passed, 0 failed).

## Open Risks/Questions
- Numeric target bands are fixed for v0.1 but instrumentation/baseline data quality may vary by metric.
- No centralized shared test harness exists yet for doc-driven stories; this can create inconsistent ad hoc tests.

## Recommended QA Focus Areas
- Confirm ADR-0008 target bands are sensible for v0.1 memory-layer-only scope (`ADR-0007` alignment).
- Verify all ADR-0002 metrics are present and unambiguous in ADR-0008.
- Verify planning references exist and are usable during prioritization (`PM-TODO.md`, `research/roadmap/RESEARCH_BACKLOG.md`).
- Validate the new test script catches deliberate negative edits (missing metric/reference).
