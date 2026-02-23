# Memory Lifecycle

## Summary
AthenaMind stores approved prompt/instruction entries in a file-backed memory root with explicit schema and status contracts.

## Core Model
- Memory root: `memory/`
- Registry: `memory/index.yaml`
- Entry content:
  - `memory/stage-prompts/<domain>/<id>.md`
  - `memory/instructions/<domain>/<id>.md`
- Entry metadata: `memory/metadata/<id>.yaml`

## Lifecycle States
- Entry status: `draft` or `approved`.
- Runtime retrieval only uses approved entries.
- Snapshot lifecycle:
  - `snapshot create`
  - `snapshot list`
  - `snapshot restore`

## Constraints
- Schema is `MAJOR.MINOR`; unsupported major versions are blocked.
- Snapshot restore requires compatibility checks and checksum validation.

## References
- `product-research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `product-research/decisions/ADR-0010-memory-schema-versioning-policy.md`
- `product-research/decisions/ADR-0016-memory-snapshot-lifecycle-contract-v02.md`
