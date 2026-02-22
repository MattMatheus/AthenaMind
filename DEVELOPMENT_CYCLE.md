# Development Cycle System

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
3. Engineering executes top story, runs tests, commits with a story-linked message, and moves it to `backlog/engineering/qa/` with handoff package.
4. QA validates and either:
   - moves story to `backlog/engineering/done/`, or
   - files prioritized bugs to `backlog/engineering/intake/` and returns story to `backlog/engineering/active/`.
   - QA then commits QA artifacts and backlog state transitions as `qa-<story-id>`.
5. PM refines intake bugs/stories, re-ranks active queue, and commits refinement/state updates.
6. Repeat until QA + Engineering are satisfied.

PM intake validation requirement:
- Run `scripts/validate_intake_items.sh` before promoting intake items to active queues.
- Validation failures must be fixed (metadata/status) and lane crossover must be corrected before ranking.

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
  - Repeat until active queue is drained.

## QA Commit Format
- QA commits must use: `qa-<story-id>`.
- Omit story title in QA commit messages to keep commit subjects short and uniform.

## Planning Commit Format
- Planning commits must use: `plan-<plan-id>`.
- `plan-id` should match the planning session file id (`PLAN-YYYYMMDD-<slug>`).

## Work-System Doc Sync Rule
- If any change modifies workflow behavior (stage flow, launch commands, commit conventions, state transitions, handoff requirements), update:
  - `HUMANS.md`
  - `AGENTS.md`
  - relevant prompt files under `prompts/active/`

## Empty Backlog Rule
- If engineering launch is attempted with empty `backlog/engineering/active/`, agent must report:
  - `no stories`
- PM cycle is then responsible for creating/refining work.
