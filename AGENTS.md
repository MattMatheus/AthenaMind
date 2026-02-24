# AGENTS

Navigation and operating guide for new or lost agents in this repository.

## Mission Context
AthenaMind v0.1 is a memory-layer product with a strict staged development workflow.

## Repository Split Context
This repository now contains:
- AthenaMind product boundary: `products/athena-mind/`
- AthenaWork product boundary: `products/athena-work/`

Use product-specific skills:
- AthenaMind skill: `skills/athena-mind/SKILL.md`
- AthenaWork skill: `skills/athena-work/SKILL.md`

Compatibility links are kept at repo root for existing commands.

## First 5 Minutes
1. Read `HUMANS.md` for operator expectations.
2. Read `DEVELOPMENT_CYCLE.md` for stage rules.
3. Read `delivery-backlog/engineering/active/README.md` for queue order.
4. Launch the requested stage with `tools/launch_stage.sh <stage>`.

## Canonical Stage Prompts
- Planning: `stage-prompts/active/planning-seed-prompt.md`
- Engineering: `stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `stage-prompts/active/qa-agent-seed-prompt.md`
- PM: `stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle: `stage-prompts/active/cycle-seed-prompt.md`

## Source-of-Truth Maps
- Staff roles: `staff-personas/STAFF_DIRECTORY.md`
- Product/architecture decisions: `product-research/decisions/`
- Architecture backlog lane: `delivery-backlog/architecture/README.md`
- Process-improvement system: `operating-system/README.md`
- Observer artifacts: `operating-system/observer/README.md`
- Wiki docs root: `knowledge-base/INDEX.md`

## Mandatory Behavioral Rules
- Branch must match `ATHENA_REQUIRED_BRANCH` (default `dev`; launcher enforces this).
- Respect backlog state model and do not skip stages.
- Use `delivery-backlog/architecture/` for architecture item types (do not place architecture stories in `delivery-backlog/engineering/active/`).
- Do not fabricate work when engineering reports `no stories`.
- Do not commit during intermediate stage transitions.
- Run observer after each completed cycle using `tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
- Commit once per cycle with format `cycle-<cycle-id>`.
- Include observer report artifact in the cycle commit (`operating-system/observer/OBSERVER-REPORT-<cycle-id>.md`).
- Keep resume handoff context current in `operating-system/observer/RESUME_CONTEXT.md` (updated by observer script).
- Apply stage exits from `knowledge-base/process/STAGE_EXIT_GATES.md`.
- Respect product-first backlog weighting (`knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`).
- Keep `product-research/roadmap/PROGRAM_STATE_BOARD.md` in sync during PM refinement.
- Treat `done` as QA-complete, not automatically shipped; use release checkpoint bundles.
- Do not add time estimates; prioritize by value/risk/dependency sequence.
- Keep Azure DevOps CI gate green (targeted tests on push/non-PR; full `go test ./...` on PR via `azure-pipelines.yml`).
- Treat repository markdown as the documentation source of truth; do not edit published docs directly.
- Ensure docs changes remain publishable through `.github/workflows/docs-publish.yml`.
- Default memory storage is external (`~/.athena/memory/<repo-id>`); do not rely on in-repo memory unless explicitly configured.

## Documentation Sync Rule
When work-system behavior changes (stage flow, handoff rules, commit conventions, state transitions, launch commands):
1. Update `HUMANS.md`.
2. Update `DEVELOPMENT_CYCLE.md`.
3. Update affected stage prompt(s).
4. Update `knowledge-base/process/README.md` or linked process docs as needed.

## Progressive Disclosure Rule
Prefer shallow-to-deep guidance:
1. Summary and action steps first.
2. Detailed workflow second.
3. Edge cases and references last.
