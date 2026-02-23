# ADR-0010: Memory Schema Versioning Policy

## Status
Accepted

## Context
The v1 memory architecture defines layout and retrieval behavior but does not define explicit schema versioning and migration policy for `index.yaml`, metadata artifacts, and entry contracts.

Without explicit compatibility policy, CLI upgrades can introduce silent drift between old and new repositories.

## Decision
Adopt a schema versioning contract for v1 memory artifacts:

- Canonical version field: `schema_version` in:
  - `/memory/index.yaml`
  - `/memory/metadata/<entry-id>.yaml`
- Version format: semantic `MAJOR.MINOR` string (for example `1.0`).
- Compatibility rules:
  - Minor bump (`1.0` -> `1.1`): backward compatible additive fields only.
  - Major bump (`1.x` -> `2.0`): incompatible change allowed only with explicit migration path.
- CLI behavior:
  - If artifact `MAJOR` is greater than CLI-supported major: fail fast with explicit migration-needed error.
  - If artifact `MAJOR` matches and `MINOR` is newer: continue in compatibility mode, ignore unknown additive fields, and emit warning.
  - If artifact version is absent: treat as invalid schema and require migration/init repair.
- Field validation:
  - Required fields must be present and typed correctly before write/retrieve operations continue.
  - Unknown top-level fields are allowed only for minor-compatible extensions.

## Consequences
- Positive:
  - Reduces upgrade ambiguity and protects deterministic fallback behavior.
  - Creates clear contract for future schema evolution.
- Negative:
  - Adds validation and migration burden to CLI implementation.
- Neutral:
  - Some legacy artifacts may require one-time migration before use.

## Alternatives Considered
- Option A: no explicit schema versioning
  - Rejected due to high compatibility risk.
- Option B: strict lockstep version (no unknown fields)
  - Rejected because it slows additive evolution.

## Validation Plan
- Architecture docs define required fields and versioning rules.
- Engineering tests must cover:
  - unsupported major version fail-fast,
  - newer minor version compatibility mode,
  - invalid/missing `schema_version` handling.
