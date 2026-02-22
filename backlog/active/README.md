# Active Queue

Ordered execution queue for engineering implementation.

## Rules
- Top entry is the next story to execute.
- Exactly one story should be in active execution at a time.
- If no stories exist, implementation agents must report `no stories` and stop.
- Product management is triggered only when active is empty.

## Active Sequence
1. `backlog/active/STORY-20260222-phased-plan-v01-v03.md`
2. `backlog/active/STORY-20260222-state-transition-checklist.md`
