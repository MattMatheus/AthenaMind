# Common Questions

## What is in scope for v0.1?
Memory-layer behavior only: file-backed write/retrieve, governance, and quality gates.

Reference: `research/decisions/ADR-0007-memory-layer-scope-refinement.md`

## Does `done` mean shipped?
No. `done` is QA-complete. Shipping requires release checkpoint bundle approval.

Reference: `docs/process/STAGE_EXIT_GATES.md`

## Can autonomous agents write memory directly?
No. Mutation is blocked in autonomous runs and requires reviewer evidence.

Reference: `docs/concepts/governance.md`

## How do I recover from a bad corpus update?
Use snapshot restore with policy approval and compatibility checks.

Reference: `docs/workflows/snapshot-recovery.md`

## Why did retrieval return a fallback result?
Semantic confidence gate likely failed, so deterministic fallback was applied.

Reference: `docs/concepts/retrieval-and-quality.md`

## How do I report a docs or behavior issue?
Create an intake bug with reproducible commands and outputs.

Reference: `backlog/engineering/intake/BUG_TEMPLATE.md`
