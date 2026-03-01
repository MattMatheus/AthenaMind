# Memory Snapshot Design Brief post-v0.1

## Snapshot Use Cases
- Deterministic checkpointing for reproducible debugging.
- Fast rollback to known-good memory states.

## Restore Semantics
- Restore rehydrates indexed memory and metadata as an atomic revision.

## Data Model and Versioning Implications
- Snapshot manifests must include schema version and compatibility markers.

## Integration Points (Current Modules)
- `cmd/memory-cli`
- `internal/memory`
- observer write-back integration paths

## Scope
- This is post-v0.1 work and remains aligned with ADR-0007 constraints.
