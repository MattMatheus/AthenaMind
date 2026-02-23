# Memory Snapshot Design Brief (Post-v0.1)

## Purpose
Define a post-v0.1 snapshot capability that mitigates memory rot while preserving current v0.1 scope boundaries.

## Scope Position
- Priority: post-v0.1 only.
- This brief defines behavior and interfaces, not implementation.
- Runtime execution ownership remains out of scope per `ADR-0007`.

## Snapshot Use Cases
1. Operator rollback after a bad memory promotion or policy misconfiguration.
2. Reproducible debugging of recall/write behavior at a known memory state.
3. Controlled migration checkpoints for schema upgrades and index rebuilds.
4. Fast recovery from accidental preference drift or low-quality semantic writes.

## Restore Semantics
- Snapshot is a point-in-time, immutable restore point for memory state.
- Restore operation creates a new forward revision and does not erase audit history.
- Partial restore modes are supported by memory area:
  - procedural directives
  - state memory (episodes/preferences)
  - semantic index state
- Default restore mode is full-memory restore with operator confirmation.
- Restore must be policy-gated by `memory-governance` and fully audit-logged.

## Data Model and Versioning Implications
- Snapshot metadata:
  - `snapshot_id` (stable identifier)
  - `created_at` (timestamp)
  - `created_by` (operator/system actor)
  - `schema_version` (state-memory schema marker)
  - `index_version` (semantic-memory index marker)
  - `scope` (`full|procedural|state|semantic`)
  - `reason` (human-readable rationale)
- Storage behavior:
  - local-first snapshot manifest and payload references
  - retention policy with size and count caps
  - checksum validation for payload integrity
- Compatibility rules:
  - restore is blocked when snapshot schema/index versions are incompatible
  - upgrade adapters may provide forward-compat restore when explicitly approved

## Integration Points (Current Modules)
- `state-memory`:
  - source of truth for episodic/preference data serialization and restore
  - schema compatibility checks and migration hooks
- `semantic-memory`:
  - index checkpoint export/import and version guardrails
  - retrieval integrity verification post-restore
- `procedural-memory`:
  - directive snapshot capture and precedence replay checks
- `memory-governance`:
  - authorization policy for snapshot create/restore/delete actions
  - mandatory review path for medium/high-risk restore operations
- `audit-telemetry`:
  - event chain for snapshot creation, validation, restore, and rollback outcomes
- `orchestrator-api`:
  - API contracts for snapshot lifecycle operations and status reporting

## Rollout Recommendation
1. Defer implementation from v0.1 critical path.
2. Start in post-v0.1 planning window as a constrained feature slice:
   - create snapshot
   - list snapshot metadata
   - full restore with governance gates
3. Add partial restore and advanced retention optimization in a later increment.

## Risks and Constraints
- Snapshot frequency can increase storage cost and retention pressure.
- Restore safety depends on strict compatibility validation and policy gates.
- Semantic index drift can cause retrieval variance unless post-restore checks run.

## Success Criteria for Future Delivery
- Restore succeeds with complete traceability and no silent data loss.
- Recovery time from a bad write/promotion is materially reduced.
- No bypass of governance approval classes during restore.
