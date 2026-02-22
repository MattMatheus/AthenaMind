# Architecture Active Queue

Ordered execution queue for architecture stories.

## Rules
- Top entry is the next architecture story to execute.
- If no stories exist, architect launcher must report `no stories`.
- Architecture outputs move through `backlog/architecture/qa` before `done`.

## Active Sequence
1. `ARCH-20260222-adr-intake-workflow-bootstrap.md`
