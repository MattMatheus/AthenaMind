# HUMANS

Operator guide for founders, non-technical users, and collaborators running AthenaWork.

## Summary

AthenaWork is a constrained stage workflow around AthenaMind:
- Planning (optional)
- Architect (optional)
- Engineering
- QA
- PM
- Cycle (engineering + QA loop)

Use launcher scripts so operators do not need to memorize deep process details.

## 60-Second Start

1. Confirm branch policy for your run (`ATHENA_REQUIRED_BRANCH`, default `dev`).
2. Start with `./tools/launch_stage.sh planning` for ambiguous requests, otherwise `engineering`.
3. Follow stage checklist from launcher output.
4. Move items through queue states (`active -> qa -> done` or back to `active` with defects).
5. Run observer after each completed cycle: `./tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
6. Commit once per cycle: `cycle-<cycle-id>`.

## Founder Control Center

- Current execution queue: `delivery-backlog/engineering/active/README.md`
- Architecture queue: `delivery-backlog/architecture/README.md`
- Observer/report evidence: `operating-system/observer/README.md`
- Release handoff and ship decision templates: `operating-system/handoff/README.md`
- KPI and dogfood metrics: `operating-system/metrics/README.md`

## Stage Commands

- Planning: `./tools/launch_stage.sh planning`
- Architect: `./tools/launch_stage.sh architect`
- Engineering: `./tools/launch_stage.sh engineering`
- QA: `./tools/launch_stage.sh qa`
- PM: `./tools/launch_stage.sh pm`
- Continuous loop: `./tools/launch_stage.sh cycle`
- Observer: `./tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- Stage tests: `./tools/run_stage_tests.sh`

## Core Rules

- Do not fabricate work when engineering reports `no stories`.
- Do not commit during stage transitions.
- One commit per completed cycle.
- Run observer before cycle commit.
- Architecture work stays in architecture lane, not engineering lane.
- Apply process gates from `knowledge-base/process/STAGE_EXIT_GATES.md`.
- Apply backlog weighting from `knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`.

## Memory Integration

- AthenaWork can call AthenaMind memory flows during launch and observer steps.
- Default memory root is external: `~/.athena/memory/<repo-id>`.
- Override with `ATHENA_MEMORY_ROOT=<path>`.
- Use `ATHENA_MEMORY_IN_REPO=1` only when repo-local memory is explicitly desired.

## If You Are Stuck

1. Read `DEVELOPMENT_CYCLE.md`.
2. Read `knowledge-base/process/README.md`.
3. Check queue order in `delivery-backlog/engineering/active/README.md`.
4. Re-run stage launcher and follow checklist output exactly.
