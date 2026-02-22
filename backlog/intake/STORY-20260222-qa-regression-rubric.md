# Story: Define QA Regression Rubric and Defect Triage Rules

## Metadata
- `id`: STORY-20260222-qa-regression-rubric
- `owner_persona`: QA Engineer - Iris.md
- `status`: intake

## Problem Statement
QA decisions need a standardized pass/fail rubric so story outcomes and bug priorities are consistent across cycles.

## Scope
- In:
  - Define pass/fail gates for `backlog/qa` stories
  - Define how to map findings to P0/P1/P2/P3
  - Define minimum evidence requirements for QA bug filing
- Out:
  - Implementing automation tooling

## Acceptance Criteria
1. A QA rubric doc exists with deterministic pass/fail criteria.
2. Priority mapping rules align with `BUG_TEMPLATE.md`.
3. QA handoff examples are included for at least one pass and one fail case.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `backlog/intake/BUG_TEMPLATE.md`
- `prompts/active/qa-agent-seed-prompt.md`

## Notes
- Needed to reduce QA variability during first sprint.
