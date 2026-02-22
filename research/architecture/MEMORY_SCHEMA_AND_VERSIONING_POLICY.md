# Memory Schema and Versioning Policy (v1)

## Purpose
Define enforceable schema and compatibility rules for `/memory` artifacts used by the CLI.

## Artifact Coverage
- `/memory/index.yaml`
- `/memory/metadata/<entry-id>.yaml`
- `/memory/prompts/<domain>/<entry-id>.md`
- `/memory/instructions/<domain>/<entry-id>.md`

## Versioning Contract
- `schema_version` is required in `index.yaml` and each metadata file.
- Format: `MAJOR.MINOR` (string).
- v1 baseline: `1.0`.

Compatibility policy:
1. Same major, newer minor: compatible additive changes only.
2. Higher major than supported by CLI: reject operation and require migration.
3. Missing/invalid version: reject operation and report schema violation.

## index.yaml Canonical Schema (v1.0)
Required top-level fields:
- `schema_version`: string (`1.0`)
- `updated_at`: ISO-8601 timestamp
- `entries`: array of entry descriptors

Required entry descriptor fields:
- `id`: string, stable unique key
- `type`: enum (`prompt` | `instruction`)
- `domain`: string
- `path`: repo-relative path under `/memory/prompts` or `/memory/instructions`
- `metadata_path`: repo-relative path under `/memory/metadata`
- `status`: enum (`draft` | `approved`)
- `updated_at`: ISO-8601 timestamp

Validation rules:
- `id` values must be unique.
- `path` and `metadata_path` must point to existing files.
- `type` and `path` folder must match (`prompt` -> `/prompts/`, `instruction` -> `/instructions/`).

## Metadata Schema (v1.0)
Required fields:
- `schema_version`: string (`1.0`)
- `id`: string; must equal filename stem and index entry id
- `title`: string
- `status`: enum (`draft` | `approved`)
- `updated_at`: ISO-8601 timestamp
- `review`: object with:
  - `reviewed_by`: string
  - `reviewed_at`: ISO-8601 timestamp or empty for draft

Optional additive fields:
- `tags`: string array
- `embedding_model`: string
- `source_refs`: string array
- `notes`: string

## Markdown Entry Contract (v1.0)
Each markdown entry file must contain:
- Title heading (`# <title>`)
- Body content
- No runtime-only metadata; authoritative structured state remains in metadata YAML and index.

## Migration Policy
- Minor migrations (`1.0` -> `1.1`):
  - additive-only fields,
  - no required field removals or type changes.
- Major migrations (`1.x` -> `2.0`):
  - must provide deterministic migration procedure,
  - CLI rejects write/retrieve until migration is completed.

## CLI Error Handling Expectations
- Unsupported major: `ERR_SCHEMA_MAJOR_UNSUPPORTED`.
- Missing/invalid version: `ERR_SCHEMA_VERSION_INVALID`.
- Required field violation: `ERR_SCHEMA_VALIDATION`.

## References
- `research/decisions/ADR-0010-memory-schema-versioning-policy.md`
- `research/decisions/ADR-0009-file-memory-cli-v1-architecture.md`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
