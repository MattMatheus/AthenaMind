# Architecture Active Queue

Ordered execution queue for architecture stories.

## Rules
- Top entry is the next architecture story to execute.
- If no stories exist, architect launcher must report `no stories`.
- Architecture outputs move through `delivery-backlog/architecture/qa` before `done`.

## Active Sequence
1. `ARCH-20260227-shared-workspace-control-plane-contract.md`
2. `ARCH-20260227-markdown-sync-authority-and-conflict-policy.md`
