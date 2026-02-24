# Architecture Story: Define File-Based Memory CLI v1 Architecture

## Metadata
- `id`: ARCH-20260222-file-memory-cli-v1-architecture
- `owner_persona`: Software Architect - Ada.md
- `status`: done

## Decision Scope
- Define architecture for a repo-local, file-based memory layer accessed through a CLI executable, with v1 scope limited to write/retrieve flows.

## Problem Statement
- The memory system needs an explicit architecture before engineering begins so CLI behavior, file layout, retrieval semantics, and governance rules are stable and reviewable.

## Inputs
- ADRs:
- Existing accepted ADRs in `product-research/decisions/`.
- Architecture docs:
- `DEVELOPMENT_CYCLE.md`
- `HUMANS.md`
- `product-research/planning/sessions/PLAN-20260222-idea-generation-session.md`
- Constraints:
- V1 supports write + retrieve only (no autonomous mutation in agent runs).
- Memory artifacts must be Markdown or machine-compliant JSON/YAML.
- Memory is repo-local by default.
- Canonical memory root path is `/memory` at repository root.
- Semantic lookup is required in v1.
- Memory updates require review pre-MVP.

## Outputs Required
- ADR updates:
- Decision record for memory CLI v1 contract and lifecycle rules.
- Architecture artifacts:
- Canonical memory folder layout and naming strategy.
- Retrieval pipeline design (semantic lookup + deterministic fallback path).
- Deterministic fallback order specification: exact-key lookup first, then path-priority lookup.
- Mutation governance model for reviewed updates outside agent runs.
- Risk/tradeoff notes:
- Tradeoffs between Markdown-first vs structured index metadata.
- Failure modes for semantic retrieval and mitigation strategy.
- Deferred roadmap boundary (API wrapper, multi-backend, UI).

## Acceptance Criteria
1. Architecture defines CLI interaction model for v1 write/retrieve operations with concrete command-level behaviors.
2. Architecture defines memory file layout, artifact types, and ownership boundaries clearly enough for implementation stories.
3. Architecture defines semantic retrieval strategy, deterministic fallback behavior, and review gates for memory mutation.
4. Architecture defines a v1 retrieval quality bar: first result is usually useful for common queries and fallback output remains predictable.

## QA Focus
- Validate that architecture output is implementation-ready, testable, and constrained to v1 scope without premature inclusion of post-MVP roadmap items.
