# Memory CLI v1 Architecture

## Purpose
Define an implementation-ready architecture for a repo-local, file-based memory layer accessed through CLI workflows.

## Scope
- In scope:
  - Repo-local memory storage rooted at `/memory`.
  - Go-based CLI write/retrieve behavior distributed as a stable portable binary.
  - Semantic retrieval pipeline with deterministic fallback.
  - Mutation governance for reviewed, out-of-run updates.
  - Non-shell primary implementation architecture.
- Out of scope:
  - Centralized memory API/service.
  - Multi-backend abstraction layers beyond files.
  - UI layer for memory inspection.
  - Shell-script-only runtime architecture.

## Canonical Layout
```text
/memory/
  AGENTS.md
  index.yaml
  prompts/
    <domain>/
      <entry-id>.md
  instructions/
    <domain>/
      <entry-id>.md
  metadata/
    <entry-id>.yaml
```

Layout notes:
- `AGENTS.md`: lane-local behavior guidance for autonomous agent navigation and constraints.
- `index.yaml`: deterministic registry of known entries, ids, and canonical paths.
- `prompts/` and `instructions/`: primary human-readable memory artifacts.
- `metadata/`: optional structured metadata for retrieval tuning and governance state.

## Memory Entry Model
Each memory entry must define:
- `id`: stable key (for exact-key fallback).
- `title`: human label.
- `type`: `prompt` or `instruction`.
- `domain`: functional category (for path-priority fallback).
- `updated_at`: ISO timestamp.
- `status`: `draft` or `approved`.
- `body`: markdown content (for human-readable artifacts).

Schema/versioning contract:
- `index.yaml` and metadata files must include `schema_version`.
- Version compatibility and migration behavior is defined in:
  - `research/architecture/MEMORY_SCHEMA_AND_VERSIONING_POLICY.md`

## CLI Interaction Model (v1)
Required behavior-level commands:
- Binary requirement:
  - Provide a stable portable Go binary as the production CLI interface.
- Write:
  - Create or update memory entry by `id`.
  - Persist markdown artifact under canonical folder.
  - Update `index.yaml` and structured metadata.
- Retrieve:
  - Accept natural language query and optional domain hints.
  - Return ranked results with selected winner plus rationale metadata.
  - Execute deterministic fallback chain when semantic result confidence is insufficient.

## Retrieval Pipeline
1. Load candidate set from `/memory/index.yaml` and relevant domains.
2. Execute semantic scoring over candidate content/metadata.
3. If confidence and ambiguity checks pass, return top semantic result.
4. Fallback (deterministic):
   - Step 1: exact-key lookup against `id`.
   - Step 2: path-priority lookup by domain/type path order.
5. Return predictable output format with:
   - `selected_id`
   - `selection_mode` (`semantic` | `fallback_exact_key` | `fallback_path_priority`)
   - `source_path`

## Quality Bar (v1)
- For common queries, first result should usually be useful.
- When semantic confidence is low or ambiguous, fallback must produce deterministic and explainable output.

## Mutation Governance (Pre-MVP)
- Memory mutation is disallowed during autonomous agent execution cycles.
- Allowed mutation workflows:
  - planning stage,
  - architect stage,
  - PM refinement stage.
- All mutations require human review before acceptance.
- Governance state for entries should be represented in metadata (`draft` vs `approved`).

## Risks and Mitigations
- Risk: semantic retrieval drift produces inconsistent answers.
  - Mitigation: explicit confidence gates + deterministic fallback chain + selection mode telemetry.
- Risk: file layout sprawl reduces discoverability.
  - Mitigation: strict `/memory` root conventions + `index.yaml` as canonical registry.
- Risk: memory edits bypass review.
  - Mitigation: workflow rule enforcement and QA checks for mutation provenance.

## Roadmap Boundary
Future phases can wrap CLI functionality behind an API contract and later add a UI, but v1 must not depend on those components.
Work-system scripts are allowed for repository operations and testing, but production memory runtime behavior must remain in the Go CLI binary.

## References
- `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- `research/decisions/ADR-0009-file-memory-cli-v1-architecture.md`
- `research/decisions/ADR-0010-memory-schema-versioning-policy.md`
- `research/architecture/MEMORY_SCHEMA_AND_VERSIONING_POLICY.md`
- `research/planning/sessions/PLAN-20260222-idea-generation-session.md`
