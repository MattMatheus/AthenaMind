# Backlog State Transition Checklist

## Purpose
Define required artifacts, evidence, and approvals for consistent backlog movement between `active`, `qa`, `done`, and return-to-active failure paths.

## Transition: `active -> qa` (Engineering Handoff)
Required artifacts:
- updated story file in `backlog/qa/` with `status: qa`
- handoff package at `backlog/qa/HANDOFF-<story>.md`
- passing test output for required commands
- commit including story id

Approval:
- Engineering operator executing the story

Gate decision:
- Allow transition only if tests pass and handoff package is complete.

## Transition: `qa -> done` (QA Pass)
Required artifacts:
- QA result artifact with explicit pass verdict at `backlog/done/QA-RESULT-<story>.md`
- evidence that acceptance criteria are satisfied
- regression risk check result
- commit including story id

Approval:
- QA operator for current cycle

Gate decision:
- Allow transition only if rubric gates pass and no blocking defects are found.

## Transition: `qa -> active` (QA Fail With Defects)
Required artifacts:
- bug files in `backlog/intake/` using `BUG_TEMPLATE.md`
- each bug includes priority (`P0-P3`) and required evidence
- QA result or notes include fail verdict and linked bug path(s)
- story moved back to `backlog/active/`
- commit including story id

Approval:
- QA operator for current cycle

Gate decision:
- Required when any blocking acceptance criteria or regression defect exists.

## Operational Rules
- Never move directly `active -> done`.
- Never move state forward with failing tests.
- Keep queue order in `backlog/active/README.md` aligned after every transition.
