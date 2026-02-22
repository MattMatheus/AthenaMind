# Development Cycle System

## Stage Launchers
- Engineering: `prompts/active/next-agent-seed-prompt.md`
- QA: `prompts/active/qa-agent-seed-prompt.md`
- PM Refinement: `prompts/active/pm-refinement-seed-prompt.md`

Quick launch helper:
- `scripts/launch_stage.sh engineering`
- `scripts/launch_stage.sh qa`
- `scripts/launch_stage.sh pm`

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
1. PM ensures ranked stories exist in `backlog/active/`.
2. Engineering executes top story, runs tests, commits with a story-linked message, and moves it to `backlog/qa/` with handoff package.
3. QA validates and either:
   - moves story to `backlog/done/`, or
   - files prioritized bugs to `backlog/intake/` and returns story to `backlog/active/`.
   - QA then commits QA artifacts and backlog state transitions with a story-linked message.
4. PM refines intake bugs, re-ranks active queue, and commits refinement/state updates.
5. Repeat until QA + Engineering are satisfied.

## Empty Backlog Rule
- If engineering launch is attempted with empty `backlog/active/`, agent must report:
  - `no stories`
- PM cycle is then responsible for creating/refining work.
