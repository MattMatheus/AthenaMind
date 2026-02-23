# HUMANS

Quick operator guide for founders and new helpers.

## Summary
This repo runs a staged workflow: Planning (as needed) -> Architect (as needed) -> Engineering -> QA -> PM Refinement.
Use launcher commands to start each stage safely.
Backlog ranking is product-first unless process defects are actively blocking delivery (`knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`).
Docs policy: markdown in git is canonical, website docs are publish artifacts.

Commit policy is cycle-based:
- Do not commit at intermediate stage transitions.
- Run Observer after each completed cycle.
- Commit once per cycle with message format: `cycle-<cycle-id>`.

## Founder Control Center (Single Human Entry Point)
Use this section as the unified "humans area." It points only to canonical docs and should be checked first.

### Current Priorities (Canonical)
- Founder one-screen snapshot:
  - `product-research/roadmap/FOUNDER_SNAPSHOT.md`
- Program control plane status and operator priorities:
  - `product-research/roadmap/PROGRAM_STATE_BOARD.md`
- Research roadmap status and upcoming research work:
  - `product-research/roadmap/RESEARCH_BACKLOG.md`
- Phase-level plan and milestone sequencing:
  - `product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`
- v0.1 goals and progress thresholds (green/yellow/red):
  - `product-research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`
- Living memory-technique baseline vs experiment rubric:
  - `product-research/roadmap/MEMORY_BEST_PRACTICES_RUBRIC_V01.md`
- Dogfooding and KPI operating loop:
  - `product-research/roadmap/DOGFOODING_LOOP_AND_TELEMETRY_KPI_SET_V01.md`

### Release and Ship Readiness
- Stage exit gates (mandatory):
  - `knowledge-base/process/STAGE_EXIT_GATES.md`
- Backlog weighting policy (mandatory for PM ranking):
  - `knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`
- Release bundle policy/template (`done` != shipped):
  - `operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- Observer policy and artifacts:
  - `operating-system/observer/README.md`

### Quick Weekly Founder Loop
1. Check `product-research/roadmap/PROGRAM_STATE_BOARD.md` (`Now` and `Risks`).
2. Check `product-research/roadmap/RESEARCH_BACKLOG.md` (`Now` and `Next`).
3. Check KPI trend snapshot output from `operating-system/metrics/KPI_SNAPSHOT_TEMPLATE.md`.
4. Confirm release bundle status before declaring anything shipped.

## 60-Second Start
1. Ensure branch is `dev`.
2. Run `./tools/launch_stage.sh engineering`.
3. Follow the returned story and checklist.
4. Move completed work to QA (`delivery-backlog/engineering/qa/`), then run QA stage.
5. Run Observer (`tools/run_observer_cycle.sh --cycle-id <story-id>`).
6. Commit once for the cycle (`cycle-<cycle-id>`).
7. Repeat until engineering returns `no stories`.

## PM Intake Validation
- Before moving items from intake to active, run:
  - `./tools/validate_intake_items.sh`
- If validation fails:
  - fix missing metadata or invalid status values,
  - move misfiled `ARCH-*` items to `delivery-backlog/architecture/intake/`,
  - move misfiled `STORY-*`/`BUG-*` items to `delivery-backlog/engineering/intake/`,
  - rerun validation until it passes.

## Program Control Plane
- PM must keep `product-research/roadmap/PROGRAM_STATE_BOARD.md` current after each refinement cycle.
- PM must apply product-first backlog weighting during queue ranking (`knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`).
- New intake items must include traceability metadata:
  - `idea_id`, `phase`, `adr_refs`, and one metric field.
- `done` is not automatically `shipped`:
  - create/maintain release bundles using `operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`.
- Use `knowledge-base/process/STAGE_EXIT_GATES.md` as mandatory gate criteria for each stage.
- Do not use time estimates in planning/refinement artifacts; sequence by value/risk/dependencies.

## Stage Commands
- Planning: `./tools/launch_stage.sh planning`
- Engineering: `./tools/launch_stage.sh engineering`
- Architect: `./tools/launch_stage.sh architect`
- QA: `./tools/launch_stage.sh qa`
- PM: `./tools/launch_stage.sh pm`
- Continuous loop: `./tools/launch_stage.sh cycle`
- Observer: `./tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- Build docs site locally: `./tools/build_docs_site.sh`
- CI enforcement (Azure DevOps): `go test ./...` on push and PR (`azure-pipelines.yml`).
- CI docs publish (GitHub Actions): `.github/workflows/docs-publish.yml` deploys docs to `athena.teamorchestrator.com/docs/`.

### Memory Integration (Soft Dependency)
- `launch_stage.sh` attempts `memory-cli bootstrap` and appends returned context to stage output.
- `run_observer_cycle.sh` attempts `memory-cli episode write` after report generation.
- If `memory-cli` is unavailable or returns an error, workflows continue with warning-only behavior.

## Backlog State Model
- `delivery-backlog/engineering/intake/`: raw new engineering work (stories/bugs)
- `delivery-backlog/engineering/active/`: ranked engineering execution queue
- `delivery-backlog/engineering/qa/`: engineering items awaiting QA decision
- `delivery-backlog/engineering/done/`: accepted engineering work
- `delivery-backlog/engineering/blocked/`: blocked engineering work
- `delivery-backlog/architecture/`: separate architecture item type and lifecycle

## How Ideas Enter The System
0. Optional planning session (recommended for ambiguous/new ideas):
   - Run `./tools/launch_stage.sh planning`.
   - Capture session notes in `product-research/planning/sessions/` using `product-research/planning/PLANNING_SESSION_TEMPLATE.md`.
   - Convert outputs into intake artifacts for engineering and/or architecture lanes.
1. New product/engineering idea:
   - Create a story in `delivery-backlog/engineering/intake/` using `delivery-backlog/engineering/intake/STORY_TEMPLATE.md`.
1a. New architecture/ADR idea:
   - Create an architecture story in `delivery-backlog/architecture/intake/` using `delivery-backlog/architecture/intake/ARCH_STORY_TEMPLATE.md`.
2. New defect/regression from QA:
   - Create a bug in `delivery-backlog/engineering/intake/` using `delivery-backlog/engineering/intake/BUG_TEMPLATE.md`.
3. New process/workflow improvement:
   - Create a process story in `operating-system/backlog/engineering/intake/` using `operating-system/backlog/engineering/intake/PROCESS_STORY_TEMPLATE.md`.
4. PM refinement then decides what moves into `delivery-backlog/engineering/active/` (or stays parked/deferred).

## Commit Rules
- Use one commit per completed cycle with format: `cycle-<cycle-id>`.
- Do not commit during intermediate stage transitions (`active -> qa`, `qa -> done`, etc.).
- Run Observer first and include `operating-system/observer/OBSERVER-REPORT-<cycle-id>.md` in the cycle commit.
- Include all cycle artifacts in the same commit (handoff, QA result, queue/program updates).

## If You Forgot What To Do
1. Read `DEVELOPMENT_CYCLE.md`.
2. Read `knowledge-base/process/CYCLE_INDEX.md`.
3. Check `delivery-backlog/engineering/active/README.md` for current order.
4. Run stage launcher and follow prompt output.

## If Something Looks Wrong
- Wrong branch? Launcher aborts if not `dev`.
- No active stories? Engineering returns `no stories`; run PM refinement.
- QA found issues? Bugs go to `delivery-backlog/engineering/intake/`, story returns to `delivery-backlog/engineering/active/`.

## Progressive Disclosure
- Human quick path: this file.
- Operational details: `DEVELOPMENT_CYCLE.md`, `knowledge-base/process/OPERATOR_DAILY_WORKFLOW.md`.
- Deep knowledge-base/wiki: `knowledge-base/INDEX.md`.
