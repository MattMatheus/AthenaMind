# HUMANS

Quick operator guide for founders and new helpers.

## Summary
This repo runs a staged workflow: Engineering -> QA -> PM Refinement.
Use launcher commands to start each stage safely.

## 60-Second Start
1. Ensure branch is `dev`.
2. Run `./scripts/launch_stage.sh engineering`.
3. Follow the returned story and checklist.
4. Move completed work to QA (`backlog/qa/`), then run QA stage.
5. Repeat until engineering returns `no stories`.

## Stage Commands
- Engineering: `./scripts/launch_stage.sh engineering`
- QA: `./scripts/launch_stage.sh qa`
- PM: `./scripts/launch_stage.sh pm`
- Continuous loop: `./scripts/launch_stage.sh cycle`

## Backlog State Model
- `backlog/intake/`: raw new work (stories/bugs)
- `backlog/active/`: ranked execution queue
- `backlog/qa/`: awaiting QA decision
- `backlog/done/`: accepted work
- `backlog/blocked/`: waiting on dependency

## Commit Rules
- Engineering: include story id in commit message.
- QA: commit format must be `qa-<story-id>`.
- PM: commit refinement/state changes after queue updates.

## If You Forgot What To Do
1. Read `DEVELOPMENT_CYCLE.md`.
2. Read `CYCLE_INDEX.md`.
3. Check `backlog/active/README.md` for current order.
4. Run stage launcher and follow prompt output.

## If Something Looks Wrong
- Wrong branch? Launcher aborts if not `dev`.
- No active stories? Engineering returns `no stories`; run PM refinement.
- QA found issues? Bugs go to `backlog/intake/`, story returns to `backlog/active/`.

## Progressive Disclosure
- Human quick path: this file.
- Operational details: `DEVELOPMENT_CYCLE.md`, `OPERATOR_DAILY_WORKFLOW.md`.
- Deep docs/wiki: `docs/INDEX.md`.
