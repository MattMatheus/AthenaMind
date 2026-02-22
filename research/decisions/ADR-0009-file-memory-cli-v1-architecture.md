# ADR-0009: File-Based Memory CLI v1 Architecture

## Status
Accepted

## Context
AthenaMind v0.1 is scoped to the memory layer and must deliver an initial working memory system before multi-backend or centralized service expansion. Planning outputs defined the first push as write/retrieve only, repo-local by default, with human-review mutation controls and semantic retrieval support.

## Decision
Adopt a repo-local file-based memory architecture for v1 with these constraints:
- Canonical memory root: `/memory` (repo root).
- Interface: CLI-first contract (runtime/language implementation details remain open for engineering).
- v1 scope: write + retrieve only.
- Memory mutation governance: updates are allowed only outside autonomous agent runs and require human review pre-MVP.
- Artifact formats:
  - Markdown for human-readable prompt/instruction memory.
  - JSON/YAML for structured index/metadata where needed.
  - Code examples remain in actual code files (not in memory prose artifacts).
- Retrieval behavior:
  - Primary path: semantic lookup for relevance.
  - Deterministic fallback order: exact-key lookup first, then path-priority lookup.
- v1 quality bar:
  - First result should usually be useful for common queries.
  - Fallback path must always provide predictable output.

Deferred from v1:
- Multi-backend layers 2/3.
- Centralized memory API wrapper and container-hosted API integration.
- UI memory exploration layer.

## Consequences
- Positive:
  - Provides a narrow, implementable foundation aligned to v0.1 scope.
  - Preserves a clean path from CLI to future API wrapper without locking backend strategy too early.
  - Keeps memory artifacts inspectable and auditable for founder-operator workflows.
- Negative:
  - Semantic retrieval quality may vary before mature ranking/tuning exists.
  - Human-reviewed mutation flow adds operational overhead in early iterations.
- Neutral:
  - Detailed CLI command syntax can evolve as long as architecture invariants remain stable.

## Alternatives Considered
- Option A: Build centralized API-first memory before CLI
  - Rejected due to unnecessary pre-MVP complexity and infrastructure coupling.
- Option B: Deterministic lookup only (no semantic retrieval in v1)
  - Rejected because it weakens usability for natural-language memory access.

## Validation Plan
- Architecture artifact must define:
  - canonical `/memory` layout and ownership boundaries,
  - retrieval pipeline including fallback order,
  - mutation governance gates.
- Engineering implementation must include tests for:
  - write/retrieve flows,
  - semantic retrieval with deterministic fallback behavior,
  - mutation restriction outside autonomous agent runs.
