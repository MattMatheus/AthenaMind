# Active Queue

Ordered execution queue for engineering implementation.

## Rules
- Top entry is the next story to execute.
- Exactly one story should be in active execution at a time.
- If no stories exist, implementation agents must report `no stories` and stop.
- Product management is triggered only when active is empty.

## Active Sequence
- (empty)
