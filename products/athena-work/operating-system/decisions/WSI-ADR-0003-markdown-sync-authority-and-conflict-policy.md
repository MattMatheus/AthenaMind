# WSI-ADR-0003: Markdown Sync Authority and Conflict Policy

## Status
Accepted

## Context
AthenaWork requires repo-readable markdown artifacts while transitioning to backend-authoritative workflow state. Without explicit sync authority and conflict policy, concurrent edits can create drift and hidden mutation.

## Decision
Adopt a single-writer sync model with deterministic conflict handling.

1. Authority
- Control-plane backend is the only writer of canonical workflow state.
- Markdown lanes are derived read models.
- Manual markdown edits are allowed only through documented exception hooks and must be reconciled by policy.

2. Sync direction
- Primary path: backend -> markdown projection.
- Reverse path (markdown -> backend) is disabled by default and allowed only for approved operator override events.

3. Conflict policy
- Required conflict taxonomy: `content_conflict`, `ordering_conflict`, `missing_artifact`, `stale_revision`.
- Resolution rules are deterministic and emit immutable audit events.

4. Drift and blocking policy
- Stage-critical artifact drift triggers blocking conditions for transition completion.
- Observer-linked alarms are mandatory for unresolved drift.

## Consequences
- Positive:
  - Preserves deterministic state authority under concurrency.
  - Maintains auditable, repo-readable artifacts for human planning.
- Negative:
  - Reduces ad-hoc markdown editing flexibility.
  - Adds reconciliation workflow overhead.

## Validation Plan
- Verify conflict classes are emitted deterministically for equivalent inputs.
- Verify unresolved stage-critical drift blocks stage transitions.
- Verify override hooks require explicit operator decision capture.
