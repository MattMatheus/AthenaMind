# Story: Add Cycle Navigation Index for Operators and Agents

## Metadata
- `id`: STORY-20260222-docs-navigation-hardening
- `owner_persona`: Technical Writer - Clara.md
- `status`: qa

## Problem Statement
Cycle artifacts are now distributed across backlog, prompts, personas, and research docs; we need a single navigation entrypoint.

## Scope
- In:
  - Create a concise cycle index doc linking canonical files and stage launch points
  - Include branch rule and `no stories` behavior
- Out:
  - Rewriting all existing docs

## Acceptance Criteria
1. A cycle index doc exists at repo root.
2. It links backlog states, prompts, launcher script, personas directory, and handoff docs.
3. It includes “first 5 minutes” operator instructions.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `DEVELOPMENT_CYCLE.md`
- `personas/STAFF_DIRECTORY.md`

## Notes
- Operator onboarding accelerator.
