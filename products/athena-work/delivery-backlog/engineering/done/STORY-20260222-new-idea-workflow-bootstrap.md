# Story: Bootstrap New Idea Intake Workflow

## Metadata
- `id`: STORY-20260222-new-idea-workflow-bootstrap
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: done

## Problem Statement
- Intake rules now exist, but engineering lacks lightweight automation/validation to prevent malformed story creation and lane crossover during fast backlog expansion.

## Scope
- In:
- Add validation checks/scripts for required story metadata fields and valid status values.
- Add lane-separation checks to flag architecture-story ids in engineering intake (and vice versa).
- Add operator-facing guidance for handling validation failures during PM refinement.
- Out:
- Full workflow orchestration platform changes.
- Stage-order or commit-convention changes.

## Acceptance Criteria
1. Validation commands reliably detect missing required metadata and invalid status values in intake stories.
2. Lane-separation checks detect obvious architecture-vs-engineering story misfiling.
3. PM refinement documentation includes failure-handling instructions that are clear enough for single-operator execution.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `HUMANS.md` intake rule section
- `DEVELOPMENT_CYCLE.md` architecture lane and refinement flow
- `delivery-backlog/architecture/INTAKE_REFINEMENT_GUIDE.md`

## Notes
- Keep implementation simple: shell/doc-test helpers are acceptable for intake validation.
