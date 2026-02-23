# Development Cycle System

## Human Operator Entry Point
- `HUMANS.md` is the canonical founder/operator navigation page.
- It links to roadmap status, release checkpoints, and KPI tracking docs.

## Stage Launchers
- Planning: `prompts/active/planning-seed-prompt.md`
- Engineering: `prompts/active/next-agent-seed-prompt.md`
- Architect: `prompts/active/architect-agent-seed-prompt.md`
- QA: `prompts/active/qa-agent-seed-prompt.md`
- PM Refinement: `prompts/active/pm-refinement-seed-prompt.md`
- Cycle Loop: `prompts/active/cycle-seed-prompt.md`

Quick launch helper:
- `scripts/launch_stage.sh planning`
- `scripts/launch_stage.sh engineering`
- `scripts/launch_stage.sh architect`
- `scripts/launch_stage.sh qa`
- `scripts/launch_stage.sh pm`
- `scripts/launch_stage.sh cycle`
- `scripts/run_observer_cycle.sh --cycle-id <cycle-id>`
- CI gate (Azure DevOps): `go test ./...` runs on push and PR via `azure-pipelines.yml`.

Memory integration (soft dependency):
- `launch_stage.sh` invokes `memory-cli bootstrap` and appends bootstrap context to stage output when available.
- `run_observer_cycle.sh` invokes `memory-cli episode write` after observer report generation when available.
- Missing/failed `memory-cli` calls must not block stage flow; scripts emit warnings and proceed.

## Doc Validation Standard
- Canonical docs validation command: `scripts/run_doc_tests.sh`
- Story-specific doc tests live under `scripts/test_*.sh`.
- Shared assertion helpers live under `scripts/lib/doc_test_harness.sh`.
- Standard and template: `research/roadmap/DOC_TEST_HARNESS_STANDARD.md`.

## Branch Safety Rule
- All stage launches require the current git branch to be `dev`.
- If branch is not `dev`, launcher aborts with:
  - `abort: active branch is '<branch>'; expected 'dev'`
- This is intentional to prevent accidental execution from the wrong branch.

## Canonical Flow
0. Planning session (optional, recommended for new/ambiguous ideas) captures interactive notes and creates intake artifacts.
1. PM ensures ranked stories exist in `backlog/engineering/active/`.
2. Architect executes top architecture story in `backlog/architecture/active/`.
3. Engineering executes top story, runs tests, prepares handoff package, and moves story to `backlog/engineering/qa/`.
4. QA validates and either:
   - moves story to `backlog/engineering/done/`, or
   - files prioritized bugs to `backlog/engineering/intake/` and returns story to `backlog/engineering/active/`.
5. Observer runs at cycle boundary, captures deterministic diff/memory deltas, and writes `work-system/observer/OBSERVER-REPORT-<cycle-id>.md`.
6. Commit once for the full cycle using `cycle-<cycle-id>`.
7. PM refines intake bugs/stories, re-ranks active queue, and updates control-plane artifacts.
8. Repeat until QA + Engineering are satisfied.
9. Shipping checkpoint (separate from `done`):
   - PM/operator prepares release bundle and explicit ship/hold decision.

PM intake validation requirement:
- Run `scripts/validate_intake_items.sh` before promoting intake items to active queues.
- Validation failures must be fixed (metadata/status) and lane crossover must be corrected before ranking.
- Ranking must apply product-first weighting per `docs/process/BACKLOG_WEIGHTING_POLICY.md`.

Traceability requirement:
- New stories and bugs must include `phase`, ADR references, and metric traceability metadata.
- Planning-originated work must include `idea_id`.

Program board requirement:
- PM refinement must update `research/roadmap/PROGRAM_STATE_BOARD.md` with queue counts and Now/Next priorities.

Stage exit gates:
- Use `docs/process/STAGE_EXIT_GATES.md` as mandatory acceptance gate for stage transitions and cycle closure.

Backlog weighting policy:
- PM ranks product stories above process stories by default.
- Process stories may outrank product work only when a broken process is blocking delivery or stage-gate enforcement.
- Record the override reason in queue notes when process work is elevated.

No-time-estimate rule:
- Pipeline sequencing is value/risk/dependency based; do not require duration estimates in stage artifacts.

## Architecture Item Type
- Architecture work uses a separate lane: `backlog/architecture/`.
- Architecture lifecycle: `backlog/architecture/intake -> ready -> active -> qa -> done`.
- This keeps architecture standards independent from engineering story standards.

## Cycle Mode
- `scripts/launch_stage.sh cycle` emits the seed for an engineering+QA loop.
- Loop behavior:
  - Run engineering stage.
  - If result is `no stories`, stop.
  - Run QA stage for completed story.
  - Run observer and generate cycle report.
  - Commit once for that cycle.
  - Repeat until active backlog is drained.

## Commit Convention
- Commit exactly once per completed cycle.
- Commit format: `cycle-<cycle-id>`.
- Include observer report and all cycle artifacts in the same commit.
- Do not commit intermediate stage transitions.

## Work-System Doc Sync Rule
- If any change modifies workflow behavior (stage flow, launch commands, commit conventions, state transitions, handoff requirements), update:
  - `HUMANS.md`
  - `AGENTS.md`
  - relevant prompt files under `prompts/active/`

## Empty Backlog Rule
- If engineering launch is attempted with empty `backlog/engineering/active/`, agent must report:
  - `no stories`
- PM cycle is then responsible for creating/refining work.
