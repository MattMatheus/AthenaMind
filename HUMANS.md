# HUMANS

Quick operator guide for founders and new helpers.

## Summary
This repo runs a staged workflow: Planning (as needed) -> Architect (as needed) -> Engineering -> QA -> PM Refinement.
Use launcher commands to start each stage safely.

Commit policy is cycle-based:
- Do not commit at intermediate stage transitions.
- Run Observer after each completed cycle.
- Commit once per cycle with message format: `cycle-<cycle-id>`.

## Founder Control Center (Single Human Entry Point)
Use this section as the unified "humans area." It points only to canonical docs and should be checked first.

### Current Priorities (Canonical)
- Founder one-screen snapshot:
  - `research/roadmap/FOUNDER_SNAPSHOT.md`
- Program control plane status and operator priorities:
  - `research/roadmap/PROGRAM_STATE_BOARD.md`
- Research roadmap status and upcoming research work:
  - `research/roadmap/RESEARCH_BACKLOG.md`
- Phase-level plan and milestone sequencing:
  - `research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`
- v0.1 goals and progress thresholds (green/yellow/red):
  - `research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`
- Living memory-technique baseline vs experiment rubric:
  - `research/roadmap/MEMORY_BEST_PRACTICES_RUBRIC_V01.md`
- Dogfooding and KPI operating loop:
  - `research/roadmap/DOGFOODING_LOOP_AND_TELEMETRY_KPI_SET_V01.md`

### Release and Ship Readiness
- Stage exit gates (mandatory):
  - `docs/process/STAGE_EXIT_GATES.md`
- Release bundle policy/template (`done` != shipped):
  - `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- Observer policy and artifacts:
  - `work-system/observer/README.md`

### Quick Weekly Founder Loop
1. Check `research/roadmap/PROGRAM_STATE_BOARD.md` (`Now` and `Risks`).
2. Check `research/roadmap/RESEARCH_BACKLOG.md` (`Now` and `Next`).
3. Check KPI trend snapshot output from `work-system/metrics/KPI_SNAPSHOT_TEMPLATE.md`.
4. Confirm release bundle status before declaring anything shipped.

## 60-Second Start
1. Ensure branch is `dev`.
2. Run `./scripts/launch_stage.sh engineering`.
3. Follow the returned story and checklist.
4. Move completed work to QA (`backlog/engineering/qa/`), then run QA stage.
5. Run Observer (`scripts/run_observer_cycle.sh --cycle-id <story-id>`).
6. Commit once for the cycle (`cycle-<cycle-id>`).
7. Repeat until engineering returns `no stories`.

## PM Intake Validation
- Before moving items from intake to active, run:
  - `./scripts/validate_intake_items.sh`
- If validation fails:
  - fix missing metadata or invalid status values,
  - move misfiled `ARCH-*` items to `backlog/architecture/intake/`,
  - move misfiled `STORY-*`/`BUG-*` items to `backlog/engineering/intake/`,
  - rerun validation until it passes.

## Program Control Plane
- PM must keep `research/roadmap/PROGRAM_STATE_BOARD.md` current after each refinement cycle.
- New intake items must include traceability metadata:
  - `idea_id`, `phase`, `adr_refs`, and one metric field.
- `done` is not automatically `shipped`:
  - create/maintain release bundles using `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`.
- Use `docs/process/STAGE_EXIT_GATES.md` as mandatory gate criteria for each stage.
- Do not use time estimates in planning/refinement artifacts; sequence by value/risk/dependencies.

## Stage Commands
- Planning: `./scripts/launch_stage.sh planning`
- Engineering: `./scripts/launch_stage.sh engineering`
- Architect: `./scripts/launch_stage.sh architect`
- QA: `./scripts/launch_stage.sh qa`
- PM: `./scripts/launch_stage.sh pm`
- Continuous loop: `./scripts/launch_stage.sh cycle`
- Observer: `./scripts/run_observer_cycle.sh --cycle-id <cycle-id>`

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
- Use one commit per completed cycle with format: `cycle-<cycle-id>`.
- Do not commit during intermediate stage transitions (`active -> qa`, `qa -> done`, etc.).
- Run Observer first and include `work-system/observer/OBSERVER-REPORT-<cycle-id>.md` in the cycle commit.
- Include all cycle artifacts in the same commit (handoff, QA result, queue/program updates).

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
