# Architecture Handoff: ARCH-20260222-file-memory-cli-v1-architecture

## Decision(s) Made
- Accepted file-based memory CLI v1 architecture with repo-local root at `/memory`.
- Fixed v1 scope to write/retrieve only.
- Set retrieval design to semantic primary with deterministic fallback order:
  1. exact-key lookup
  2. path-priority lookup
- Established pre-MVP governance rule: memory mutations only outside autonomous runs and under human review.

## Alternatives Considered
- Centralized API-first memory before CLI: rejected for pre-MVP complexity.
- Deterministic lookup only without semantic retrieval: rejected for weaker natural-language usability.

## Risks and Mitigations
- Semantic ranking drift:
  - Mitigate with confidence gating and explicit fallback mode signaling.
- Layout sprawl:
  - Mitigate with strict `/memory` conventions plus canonical `index.yaml`.
- Mutation policy bypass:
  - Mitigate with workflow rules and QA verification of mutation provenance.

## Updated Artifacts
- `research/decisions/ADR-0009-file-memory-cli-v1-architecture.md`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `backlog/architecture/qa/ARCH-20260222-file-memory-cli-v1-architecture.md`

## Validation Commands and Results
- `./scripts/run_doc_tests.sh` -> PASS

## Open Questions for QA Focus
- Does the architecture artifact define enough deterministic behavior for implementation testing?
- Are fallback and governance requirements clear enough to prevent scope drift in engineering?
