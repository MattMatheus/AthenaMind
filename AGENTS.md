# AGENTS

Navigation and operating guide for agents working in AthenaMind + AthenaWork.

## Mission Context

- `AthenaMind`: governed memory runtime for agentic coding workflows.
- `AthenaWork`: staged workflow that constrains and steers agent execution.

## Repository Product Boundaries

- AthenaMind: `products/athena-mind/`
- AthenaWork: `products/athena-work/`

Use product skills when applicable:
- `skills/athena-mind/SKILL.md`
- `skills/athena-work/SKILL.md`

Root compatibility links (`tools`, `stage-prompts`, `staff-personas`, `delivery-backlog`, `operating-system`, `HUMANS.md`, `DEVELOPMENT_CYCLE.md`) are canonical operator entry points.

## First 5 Minutes

1. Read `HUMANS.md`.
2. Read `DEVELOPMENT_CYCLE.md`.
3. Read `delivery-backlog/engineering/active/README.md`.
4. Launch requested stage with `tools/launch_stage.sh <stage>`.

## Canonical Stage Prompts

- Planning: `stage-prompts/active/planning-seed-prompt.md`
- Engineering: `stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `stage-prompts/active/qa-agent-seed-prompt.md`
- PM: `stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `stage-prompts/active/cycle-seed-prompt.md`

## Source-of-Truth Maps

- Staff specialist directory: `staff-personas/STAFF_DIRECTORY.md`
- Architecture lane: `delivery-backlog/architecture/README.md`
- Process-improvement system: `operating-system/README.md`
- Observer reports and resume context: `operating-system/observer/README.md`
- User/operator docs map: `knowledge-base/INDEX.md`

## Mandatory Behavioral Rules

- Branch must match `ATHENA_REQUIRED_BRANCH` (default `dev`; launcher enforces this).
- Respect backlog lanes and state transitions; do not skip stages.
- Use architecture lane for architecture work (`delivery-backlog/architecture/`).
- If engineering returns `no stories`, stop and hand off to PM refinement.
- Do not commit during intermediate stage transitions.
- Run observer after each completed cycle: `tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
- Commit once per cycle using `cycle-<cycle-id>`.
- Include observer report in cycle commit: `operating-system/observer/OBSERVER-REPORT-<cycle-id>.md`.
- Keep resume context current: `operating-system/observer/RESUME_CONTEXT.md`.
- Apply stage exits: `knowledge-base/process/stage-exit-gates.md`.
- Apply backlog weighting policy: `knowledge-base/process/backlog-weighting-policy.md`.
- Treat `done` as QA-complete, not auto-shipped; use release bundles in `operating-system/handoff/`.
- Prioritize by value/risk/dependency; do not add time estimates.
- Repository markdown is source of truth.
- Use `tools/run_stage_tests.sh` before handoff.

## Documentation Sync Rule

When work-system behavior changes (stage flow, handoff rules, commit conventions, state transitions, launch commands):
1. Update `HUMANS.md`.
2. Update `DEVELOPMENT_CYCLE.md`.
3. Update affected stage prompt(s).
4. Update `knowledge-base/process/README.md` and related docs.
