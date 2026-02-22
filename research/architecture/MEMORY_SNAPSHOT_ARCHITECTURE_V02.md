# Memory Snapshot Architecture (v0.2)

## Purpose
Define implementation-ready contracts for snapshot lifecycle operations (`create`, `list`, `restore`) in post-v0.1 execution.

## Scope
- In scope:
  - Local-first snapshot manifests and payload references.
  - Full-restore flow with governance and compatibility gates.
  - Deterministic audit events for lifecycle operations.
- Out of scope:
  - Partial restore execution (planned next increment).
  - Cloud snapshot backends and cross-repo restore orchestration.

## Snapshot Lifecycle Contract
### 1) Create
- Operation: `snapshot create`
- Required inputs:
  - `scope=full`
  - `reason`
  - `created_by`
- Behavior:
  - Capture point-in-time references for procedural/state/semantic artifacts.
  - Generate immutable `snapshot_id`.
  - Persist manifest with checksums.

### 2) List
- Operation: `snapshot list`
- Required outputs per row:
  - `snapshot_id`
  - `created_at`
  - `created_by`
  - `schema_version`
  - `index_version`
  - `scope`
  - `reason`

### 3) Restore (full)
- Operation: `snapshot restore --snapshot-id <id>`
- Required gates:
  - Compatibility check for `schema_version` and `index_version`.
  - Governance review policy check.
- Behavior:
  - Restore creates a new forward revision.
  - Previous state and snapshot artifacts remain auditable.

## Data Model
- Manifest path:
  - `/memory/snapshots/<snapshot_id>/manifest.json`
- Manifest required fields:
  - `snapshot_id`
  - `created_at`
  - `created_by`
  - `schema_version`
  - `index_version`
  - `scope`
  - `reason`
  - `checksums[]`
  - `payload_refs[]`

## Compatibility Rules
1. Restore rejects snapshots with unsupported major schema/index versions.
2. Restore allows newer minor versions only in explicit compatibility mode.
3. Missing or invalid version fields are treated as schema violations.

## Governance and Audit Rules
- `snapshot restore` requires policy-gate evaluation and reviewer evidence.
- Required audit event chain:
  - `snapshot.create.requested`
  - `snapshot.create.completed`
  - `snapshot.restore.requested`
  - `snapshot.restore.policy_decision`
  - `snapshot.restore.completed|failed`
- Each event must include `snapshot_id`, `session_id`, `trace_id`, and `result`.

## Failure Modes and Deterministic Responses
- Compatibility failure:
  - `ERR_SNAPSHOT_COMPATIBILITY_BLOCKED`
- Missing manifest field:
  - `ERR_SNAPSHOT_MANIFEST_INVALID`
- Policy gate failure:
  - `ERR_SNAPSHOT_RESTORE_POLICY_BLOCKED`
- Checksum mismatch:
  - `ERR_SNAPSHOT_INTEGRITY_CHECK_FAILED`

## Rollout Sequence
1. Implement create/list/full-restore contract.
2. Validate governance + compatibility behavior in QA.
3. Add partial restore and retention optimization as a separate follow-on slice.

## Follow-On Implementation
- `backlog/engineering/intake/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`

## References
- `research/decisions/ADR-0016-memory-snapshot-lifecycle-contract-v02.md`
- `research/roadmap/MEMORY_SNAPSHOT_DESIGN_BRIEF_POST_V01.md`
- `research/architecture/MEMORY_SCHEMA_AND_VERSIONING_POLICY.md`
- `research/architecture/MEMORY_MUTATION_REVIEW_WORKFLOW.md`
