# Story: Define Founder-Operator Daily Workflow Script

## Metadata
- `id`: STORY-20260222-founder-operator-workflow
- `owner_persona`: Product Manager - Maya.md
- `status`: done

## Problem Statement
A repeatable day-to-day operator workflow is needed so the system is used consistently inside Codex app sessions.

## Scope
- In:
  - Define startup, stage launch, QA loop, and shutdown routine
  - Include branch discipline and escalation rules
- Out:
  - Team/multi-operator workflows (v2.0)

## Acceptance Criteria
1. A concise operator workflow doc exists.
2. It references `launch_stage.sh` and required prompts.
3. It includes explicit “if X then Y” instructions for empty backlog and QA failures.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `tools/launch_stage.sh`
- `DEVELOPMENT_CYCLE.md`

## Notes
- Supports local single-operator model.
