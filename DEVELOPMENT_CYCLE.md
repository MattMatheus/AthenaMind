# Development Cycle System

## Stage Launchers
- Engineering: `prompts/active/next-agent-seed-prompt.md`
- QA: `prompts/active/qa-agent-seed-prompt.md`
- PM Refinement: `prompts/active/pm-refinement-seed-prompt.md`

Quick launch helper:
- `scripts/launch_stage.sh engineering`
- `scripts/launch_stage.sh qa`
- `scripts/launch_stage.sh pm`

## Canonical Flow
1. PM ensures ranked stories exist in `backlog/active/`.
2. Engineering executes top story and moves it to `backlog/qa/` with handoff package.
3. QA validates and either:
   - moves story to `backlog/done/`, or
   - files prioritized bugs to `backlog/intake/` and returns story to `backlog/active/`.
4. PM refines intake bugs and re-ranks active queue.
5. Repeat until QA + Engineering are satisfied.

## Empty Backlog Rule
- If engineering launch is attempted with empty `backlog/active/`, agent must report:
  - `no stories`
- PM cycle is then responsible for creating/refining work.
