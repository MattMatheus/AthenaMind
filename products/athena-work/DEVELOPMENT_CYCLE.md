# Development Cycle System

## Human Operator Entry Point
- `HUMANS.md` is the canonical founder/operator navigation page.
- It links to roadmap status, release checkpoints, and KPI tracking docs.

## Stage Launchers
- Planning: `stage-prompts/active/planning-seed-prompt.md`
- Engineering: `stage-prompts/active/next-agent-seed-prompt.md`
- Architect: `stage-prompts/active/architect-agent-seed-prompt.md`
- QA: `stage-prompts/active/qa-agent-seed-prompt.md`
- PM Refinement: `stage-prompts/active/pm-refinement-seed-prompt.md`
- Cycle Loop: `stage-prompts/active/cycle-seed-prompt.md`

Quick launch helper:
- `tools/launch_stage.sh planning`
- `tools/launch_stage.sh engineering`
- `tools/launch_stage.sh architect`
- `tools/launch_stage.sh qa`
- `tools/launch_stage.sh pm`
- `tools/launch_stage.sh cycle`
- `tools/run_stage_tests.sh` (auto: docs + targeted on push/non-PR, docs + full suite on PR)
- `tools/run_observer_cycle.sh --cycle-id <cycle-id>`
- `tools/build_docs_site.sh` (local docs-site build)
- CI gate (Azure DevOps):
  - Push/non-PR: targeted tests (`go test ./internal/governance ./internal/episode ./cmd/memory-cli`) + docs tests.
  - Pull request: full suite (`go test ./...`) + docs tests.
- Docs publish gate (GitHub Actions): `.github/workflows/docs-publish.yml` publishes docs artifacts.

Documentation source-of-truth rule:
- Edit docs only in repository markdown; published site content is generated output.

Memory integration (soft dependency):
- `launch_stage.sh` invokes `memory-cli bootstrap` and appends bootstrap context to stage output when available.
- `run_observer_cycle.sh` invokes `memory-cli episode write` after observer report generation when available.
- Default memory root is external to the repo: `~/.athena/memory/<repo-id>`.
- Set `ATHENA_MEMORY_IN_REPO=1` to use `<repo>/memory` for local compatibility.
- Set `ATHENA_MEMORY_ROOT` to fully override memory storage location.
- Missing/failed `memory-cli` calls must not block stage flow; scripts emit warnings and proceed.

Preflight gate enforcement:
- `launch_stage.sh` enforces story readiness for engineering/architecture before launch output:
  - `## Metadata` section must exist.
  - metadata must include `idea_id`, `phase`, and `adr_refs`.
  - `## Acceptance Criteria` section must exist.
- If prior observer reports exist, `operating-system/observer/RESUME_CONTEXT.md` must exist before stage launch.
- Default mode is hard gate (`ATHENA_PREFLIGHT_MODE=enforce`); use `ATHENA_PREFLIGHT_MODE=warn` for warning-only behavior.

## Doc Validation Standard
- Canonical docs validation command: `tools/run_doc_tests.sh`
- Story-specific doc tests live under `tools/test_*.sh`.
- Shared assertion helpers live under `tools/lib/doc_test_harness.sh`.
- Standard and template: `product-research/roadmap/DOC_TEST_HARNESS_STANDARD.md`.

## Branch Safety Rule
- All stage launches require the current git branch to match `ATHENA_REQUIRED_BRANCH` (default `dev`).
- If branch does not match, launcher aborts with:
  - `abort: active branch is '<branch>'; expected '<required-branch>'`
- This is intentional to prevent accidental execution from the wrong branch.

## Canonical Flow
0. Planning session (optional, recommended for new/ambiguous ideas) captures interactive notes and creates intake artifacts.
1. PM ensures ranked stories exist in `delivery-backlog/engineering/active/`.
2. Architect executes top architecture story in `delivery-backlog/architecture/active/`.
3. Engineering executes top story, runs tests, prepares handoff package, and moves story to `delivery-backlog/engineering/qa/`.
   - Required during stage execution: docs tests + targeted tests for touched scope.
   - Required during PR validation: full suite `go test ./...`.
4. QA validates and either:
   - moves story to `delivery-backlog/engineering/done/`, or
   - files prioritized bugs to `delivery-backlog/engineering/intake/` and returns story to `delivery-backlog/engineering/active/`.
5. Observer runs at cycle boundary, captures deterministic diff/memory deltas, and writes `operating-system/observer/OBSERVER-REPORT-<cycle-id>.md`.
6. Commit once for the full cycle using `cycle-<cycle-id>`.
7. PM refines intake bugs/stories, re-ranks active queue, and updates control-plane artifacts.
8. Repeat until QA + Engineering are satisfied.
9. Shipping checkpoint (separate from `done`):
   - PM/operator prepares release bundle and explicit ship/hold decision.

PM intake validation requirement:
- Run `tools/validate_intake_items.sh` before promoting intake items to active queues.
- Validation failures must be fixed (metadata/status) and lane crossover must be corrected before ranking.
- Ranking must apply product-first weighting per `knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`.

Traceability requirement:
- New stories and bugs must include `phase`, ADR references, and metric traceability metadata.
- Planning-originated work must include `idea_id`.

Program board requirement:
- PM refinement must update `product-research/roadmap/PROGRAM_STATE_BOARD.md` with queue counts and Now/Next priorities.

Stage exit gates:
- Use `knowledge-base/process/STAGE_EXIT_GATES.md` as mandatory acceptance gate for stage transitions and cycle closure.

Backlog weighting policy:
- PM ranks product stories above process stories by default.
- Process stories may outrank product work only when a broken process is blocking delivery or stage-gate enforcement.
- Record the override reason in queue notes when process work is elevated.

No-time-estimate rule:
- Pipeline sequencing is value/risk/dependency based; do not require duration estimates in stage artifacts.

## Architecture Item Type
- Architecture work uses a separate lane: `delivery-backlog/architecture/`.
- Architecture lifecycle: `delivery-backlog/architecture/intake -> ready -> active -> qa -> done`.
- This keeps architecture standards independent from engineering story standards.

## Cycle Mode
- `tools/launch_stage.sh cycle` emits the seed for an engineering+QA loop.
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
  - relevant prompt files under `stage-prompts/active/`

## Empty Backlog Rule
- If engineering launch is attempted with empty `delivery-backlog/engineering/active/`, agent must report:
  - `no stories`
- PM cycle is then responsible for creating/refining work.
