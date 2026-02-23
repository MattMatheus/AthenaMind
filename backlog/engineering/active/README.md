# Engineering Active Queue

Ordered execution queue for engineering implementation.

## Rules
- Top entry is the next engineering story to execute.
- If no engineering stories exist, launcher returns `no stories`.
- Queue ordering is product-first; process stories only outrank product work when process defects are blocking delivery (`docs/process/BACKLOG_WEIGHTING_POLICY.md`).

## Active Sequence
No active engineering stories.
