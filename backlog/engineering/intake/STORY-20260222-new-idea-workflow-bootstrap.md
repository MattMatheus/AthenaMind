# Story: Bootstrap New Idea Intake Workflow

## Metadata
- `id`: STORY-20260222-new-idea-workflow-bootstrap
- `owner_persona`: personas/Product Manager - Maya.md
- `status`: intake

## Problem Statement
- We need a repeatable intake path so new product/engineering ideas and architecture/ADR ideas are consistently captured in the correct backlog lane with the correct template.

## Scope
- In:
- Define and document the PM intake triage checklist for new idea submissions.
- Ensure engineering ideas are filed as stories in `backlog/engineering/intake/` from `STORY_TEMPLATE.md`.
- Ensure architecture/ADR ideas are filed as architecture stories in `backlog/architecture/intake/` from `ARCH_STORY_TEMPLATE.md`.
- Out:
- Implementing specific product features or architecture changes from intake items.
- Changing stage order or commit conventions.

## Acceptance Criteria
1. A PM-ready intake checklist exists in this story notes and is usable without additional clarification.
2. At least one engineering intake story and one architecture intake story are created using the canonical templates.
3. PM refinement can pull these items into `ready` with clear problem statements, scope, and acceptance criteria.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `HUMANS.md` intake rule section
- `DEVELOPMENT_CYCLE.md` architecture lane and refinement flow

## Notes
- PM Intake Checklist:
- Confirm whether the idea is implementation work (`engineering`) or decision/design work (`architecture/ADR`).
- Create the item in the corresponding `intake` folder using the canonical template.
- Ensure ID format and status are valid (`intake`).
- Add enough detail for PM refinement to rank and move to `ready`.
