# Story: Build human planning workbench v1 (simple, direct, low-vision-first)

## Metadata
- `id`: STORY-20260227-human-planning-workbench-v1
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0004, ADR-0013]
- `success_metric`: Humans can capture direction and produce next-stage recommendation in <=10 minutes with <=1 navigation context switch.
- `release_checkpoint`: required

## Problem Statement
Planning currently requires reading and correlating many artifacts. This increases cognitive load and contributes to state drift when humans miss critical fields or confirmations.

## Scope
- In:
  - single-screen planning workbench showing goals, constraints, risks, and next-stage recommendation
  - structured direction confirmation action for human approval
  - plain-language summaries and explicit "what happens next" section
- Out:
  - long-form rich text editor
  - advanced collaboration/comment threads

## Acceptance Criteria
1. Planning UI captures required fields (`direction`, `constraints`, `risks`, `next_stage`) in one screen.
2. Human can explicitly confirm direction and see confirmation status immediately.
3. Layout is low-vision-friendly by default: large text, high contrast, simple hierarchy, and direct labels.
4. Output exports to workflow artifacts without requiring manual copying across files.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-workspace-ui-read-only-board-v1
- STORY-20260227-human-direction-confirmation-gate

## Notes
- Optimize for first-time operator comprehension over feature breadth.
