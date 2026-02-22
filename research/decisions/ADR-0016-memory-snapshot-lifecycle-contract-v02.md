# ADR-0016: Memory Snapshot Lifecycle Contract (v0.2)

## Status
Accepted

## Context
Snapshot capability was defined in a design brief but lacked accepted lifecycle contracts and deterministic error/governance behavior for implementation sequencing.

## Decision
Adopt a v0.2 snapshot contract with operations:
- `create`
- `list`
- `restore` (full scope only for first increment)

Canonical architecture artifact:
- `research/architecture/MEMORY_SNAPSHOT_ARCHITECTURE_V02.md`

Core contract decisions:
- Local-first manifest + payload references.
- Immutable `snapshot_id`.
- Restore creates forward revision and preserves audit trail.
- Compatibility and policy gates are mandatory before restore.

## Consequences
- Positive:
  - Converts snapshot from concept to implementation-ready architecture.
  - Establishes deterministic safety and governance behavior.
- Negative:
  - Adds compatibility and integrity validation complexity.
  - First increment excludes partial restore functionality.
- Neutral:
  - Future cloud adapters remain optional and phase-gated.

## Alternatives Considered
- Implement snapshot immediately without accepted contract.
  - Rejected due to restore safety risk.
- Include partial restore in first implementation slice.
  - Rejected to minimize first-slice risk and complexity.

## Follow-On Implementation Paths
- `backlog/engineering/intake/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
- `backlog/engineering/intake/STORY-20260222-memory-api-read-gateway-v03.md` (for later API exposure path)
