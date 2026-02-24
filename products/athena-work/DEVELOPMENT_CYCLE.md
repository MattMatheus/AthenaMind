# Development Cycle System

## Purpose

Define the enforced AthenaWork loop used to execute governed delivery with AthenaMind.

## Stage Launchers

- Planning: `stage-prompts/active/planning-seed-prompt.md`
- Engineering: `stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `stage-prompts/active/qa-agent-seed-prompt.md`
- PM: `stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `stage-prompts/active/cycle-seed-prompt.md`

Helper commands:
- `tools/launch_stage.sh <planning|engineering|architect|qa|pm|cycle>`
- `tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- `tools/run_stage_tests.sh`

## Canonical Flow

1. Planning (optional): shape ambiguous requests into intake artifacts.
2. Architect (optional): resolve architecture items from architecture lane.
3. Engineering: execute top-ranked story from engineering active lane.
4. QA: validate story outcome, regressions, and acceptance criteria.
5. Observer: generate cycle evidence and resume context.
6. PM: refine intake and reprioritize active queue.
7. Commit once for completed cycle (`cycle-<cycle-id>`).

## Backlog Model

- Engineering lane:
  - `delivery-backlog/engineering/intake/`
  - `delivery-backlog/engineering/active/`
  - `delivery-backlog/engineering/qa/`
  - `delivery-backlog/engineering/done/`
  - `delivery-backlog/engineering/blocked/`
- Architecture lane:
  - `delivery-backlog/architecture/intake/`
  - `delivery-backlog/architecture/ready/`
  - `delivery-backlog/architecture/active/`
  - `delivery-backlog/architecture/qa/`
  - `delivery-backlog/architecture/done/`

## Guardrails

- Branch safety enforced by launcher (`ATHENA_REQUIRED_BRANCH`, default `dev`).
- Preflight readiness enforced before engineering/architecture execution.
- `no stories` from engineering is a valid stop signal.
- No intermediate commits during stage transitions.
- `done` means QA-complete; shipping is controlled by release handoff artifacts.

## Quality Gates

- Run `tools/run_stage_tests.sh` before handoff.
- Push/non-PR scope: docs + targeted Go tests.
- PR scope: docs + `go test ./...`.
- Apply stage exit criteria in `knowledge-base/process/stage-exit-gates.md`.

## Documentation Sync Rule

If stage behavior changes (commands, state transitions, commit policy, handoff requirements), update:
1. `HUMANS.md`
2. `AGENTS.md`
3. relevant `stage-prompts/active/*`
4. `knowledge-base/process/README.md`
