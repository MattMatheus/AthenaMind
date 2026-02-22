# AGENTS

Navigation and operating guide for new or lost agents in this repository.

## Mission Context
AthenaMind v0.1 is a memory-layer product with a strict staged development workflow.

## First 5 Minutes
1. Read `HUMANS.md` for operator expectations.
2. Read `DEVELOPMENT_CYCLE.md` for stage rules.
3. Read `backlog/engineering/active/README.md` for queue order.
4. Launch the requested stage with `scripts/launch_stage.sh <stage>`.

## Canonical Stage Prompts
- Engineering: `prompts/active/next-agent-seed-prompt.md`
- Architect: `prompts/active/architect-agent-seed-prompt.md`
- QA: `prompts/active/qa-agent-seed-prompt.md`
- PM: `prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `prompts/active/cycle-seed-prompt.md`

## Source-of-Truth Maps
- Staff roles: `personas/STAFF_DIRECTORY.md`
- Product/architecture decisions: `research/decisions/`
- Architecture backlog lane: `backlog/architecture/README.md`
- Process-improvement system: `work-system/README.md`
- Wiki docs root: `docs/INDEX.md`

## Mandatory Behavioral Rules
- Branch must be `dev` (launcher enforces this).
- Respect backlog state model and do not skip stages.
- Use `backlog/architecture/` for architecture item types (do not place architecture stories in `backlog/engineering/active/`).
- Do not fabricate work when engineering reports `no stories`.
- Commit after each stage transition.
- Architect commit format: `arch-<story-id>`.
- QA commit format: `qa-<story-id>`.

## Documentation Sync Rule
When work-system behavior changes (stage flow, handoff rules, commit conventions, state transitions, launch commands):
1. Update `HUMANS.md`.
2. Update `DEVELOPMENT_CYCLE.md`.
3. Update affected stage prompt(s).
4. Update `docs/process/README.md` or linked process docs as needed.

## Progressive Disclosure Rule
Prefer shallow-to-deep guidance:
1. Summary and action steps first.
2. Detailed workflow second.
3. Edge cases and references last.
