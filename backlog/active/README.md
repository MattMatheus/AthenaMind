# Active Queue

Ordered execution queue for engineering implementation.

## Rules
- Top entry is the next story to execute.
- Exactly one story should be in active execution at a time.
- If no stories exist, implementation agents must report `no stories` and stop.
- Product management is triggered only when active is empty.

## Active Sequence
1. `BUG-20260222-launch-stage-readme-parse-fallback.md`
2. `STORY-20260222-memory-snapshot-plan.md`
