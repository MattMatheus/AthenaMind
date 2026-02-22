# Story: Define Backlog State Transition Checklist

## Metadata
- `id`: STORY-20260222-state-transition-checklist
- `owner_persona`: QA Engineer - Iris.md
- `status`: qa

## Problem Statement
Story movement across `active`, `qa`, `done`, and `blocked` needs checklist gates to avoid inconsistent transitions.

## Scope
- In:
  - Define required evidence for each state transition
  - Define who can approve each transition
  - Add checklist artifact referenced by prompts
- Out:
  - Automating transitions

## Acceptance Criteria
1. A state transition checklist file exists and is referenced by stage prompts.
2. Required artifacts are defined for `active -> qa` and `qa -> done`.
3. Failure path requirements are defined for `qa -> active` with defects.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `DEVELOPMENT_CYCLE.md`
- prompt files under `prompts/active/`

## Notes
- Recommended hardening step already identified.
