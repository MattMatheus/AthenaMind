# HUMANS

Quick operator guide for founders and new helpers.

## Summary
This repo runs a staged workflow: Planning (as needed) -> Architect (as needed) -> Engineering -> QA -> PM Refinement.
Use launcher commands to start each stage safely.

## 60-Second Start
1. Ensure branch is `dev`.
2. Run `./scripts/launch_stage.sh engineering`.
3. Follow the returned story and checklist.
4. Move completed work to QA (`backlog/engineering/qa/`), then run QA stage.
5. Repeat until engineering returns `no stories`.

## Stage Commands
- Planning: `./scripts/launch_stage.sh planning`
- Engineering: `./scripts/launch_stage.sh engineering`
- Architect: `./scripts/launch_stage.sh architect`
- QA: `./scripts/launch_stage.sh qa`
- PM: `./scripts/launch_stage.sh pm`
- Continuous loop: `./scripts/launch_stage.sh cycle`

## Backlog State Model
- `backlog/engineering/intake/`: raw new engineering work (stories/bugs)
- `backlog/engineering/active/`: ranked engineering execution queue
- `backlog/engineering/qa/`: engineering items awaiting QA decision
- `backlog/engineering/done/`: accepted engineering work
- `backlog/engineering/blocked/`: blocked engineering work
- `backlog/architecture/`: separate architecture item type and lifecycle

## How Ideas Enter The System
0. Optional planning session (recommended for ambiguous/new ideas):
   - Run `./scripts/launch_stage.sh planning`.
   - Capture session notes in `research/planning/sessions/` using `research/planning/PLANNING_SESSION_TEMPLATE.md`.
   - Convert outputs into intake artifacts for engineering and/or architecture lanes.
1. New product/engineering idea:
   - Create a story in `backlog/engineering/intake/` using `backlog/engineering/intake/STORY_TEMPLATE.md`.
1a. New architecture/ADR idea:
   - Create an architecture story in `backlog/architecture/intake/` using `backlog/architecture/intake/ARCH_STORY_TEMPLATE.md`.
2. New defect/regression from QA:
   - Create a bug in `backlog/engineering/intake/` using `backlog/engineering/intake/BUG_TEMPLATE.md`.
3. New process/workflow improvement:
   - Create a process story in `work-system/backlog/engineering/intake/` using `work-system/backlog/engineering/intake/PROCESS_STORY_TEMPLATE.md`.
4. PM refinement then decides what moves into `backlog/engineering/active/` (or stays parked/deferred).

## Commit Rules
- Planning: commit format must be `plan-<plan-id>`.
- Engineering: include story id in commit message.
- Architect: commit format must be `arch-<story-id>`.
- QA: commit format must be `qa-<story-id>`.
- PM: commit refinement/state changes after queue updates.

## If You Forgot What To Do
1. Read `DEVELOPMENT_CYCLE.md`.
2. Read `CYCLE_INDEX.md`.
3. Check `backlog/engineering/active/README.md` for current order.
4. Run stage launcher and follow prompt output.

## If Something Looks Wrong
- Wrong branch? Launcher aborts if not `dev`.
- No active stories? Engineering returns `no stories`; run PM refinement.
- QA found issues? Bugs go to `backlog/engineering/intake/`, story returns to `backlog/engineering/active/`.

## Progressive Disclosure
- Human quick path: this file.
- Operational details: `DEVELOPMENT_CYCLE.md`, `OPERATOR_DAILY_WORKFLOW.md`.
- Deep docs/wiki: `docs/INDEX.md`.
