# WSI-ADR-0001: Add Architect Launch Stage

## Status
Accepted

## Context
Architecture and ADR work was competing with PM scoping and engineering cycles. The workflow needed a dedicated stage to execute architecture-owned stories predictably.

## Decision
Add a first-class launch stage:
- `scripts/launch_stage.sh architect`
- Prompt: `prompts/active/architect-agent-seed-prompt.md`

Selection rule:
- Launcher selects the top active story owned by `Software Architect - Ada.md`.
- If none exist, return `no stories`.

Commit rule:
- Architect commits must use `arch-<story-id>`.

## Consequences
- Positive:
  - Clear separation between architecture decisions and PM scoping.
  - Better queue discipline for ADR/architecture updates.
- Negative:
  - Adds one more stage command to operator workflow.

## Validation Plan
- Verify launcher returns architect prompt + story when architecture-owned stories exist.
- Verify launcher returns `no stories` when none exist.
- Verify docs/guides include architect stage and commit rule.
