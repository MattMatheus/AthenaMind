# Research Backlog States

This backlog uses folder states so agents can move a story file through a clear lifecycle.

## States
- `intake/`: raw ideas not yet refined into actionable stories.
- `ready/`: refined and ready to be picked up.
- `active/`: currently in execution by the active agent.
- `qa/`: implementation/research complete, waiting for QA review.
- `blocked/`: cannot proceed without specific dependency/input.
- `done/`: passed QA and accepted.
- `archive/`: historical closed items no longer part of active planning horizon.

## Core Flow
`intake -> ready -> active -> qa -> done`

## Exceptions
- Any state can move to `blocked` when a blocker appears.
- `done` items can be moved to `archive` during periodic cleanup.
