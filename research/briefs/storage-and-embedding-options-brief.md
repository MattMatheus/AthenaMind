# Research Brief: Storage and Embedding Options (Local-First)

## Source Inputs
- `README.md` (vision prompt)
- `research/ingest/archive/pillars.md`
- `research/architecture/MODULE_BOUNDARIES_V01.md`
- `research/decisions/ADR-0001-three-pillar-foundation.md`
- `research/references/external-sources-2026-02-22-storage-embedding.md`

## Problem Statement
AthenaMind needs concrete v0.1 choices for structured state storage and semantic retrieval to avoid implementation drift.

## Key Insights (80/20)
- Vision already points to `SQLite + FAISS` for local-first mode and optional cloud later.
- The biggest risk is over-optimizing ANN/index complexity before baseline workload data.
- A clean adapter boundary preserves optional migration to Cosmos/cloud retrieval services.

## Option Summary
- Option A: SQLite (WAL) + FAISS exact first (`Flat`/IDMap)
  - Strong local-first fit; lowest operational complexity.

- Option B: SQLite + more complex FAISS IVF/PQ at v0.1
  - Premature without corpus scale/latency evidence.

- Option C: Cloud-first Cosmos + managed vector from v0.1
  - Conflicts with local-first and offline-friendly goals.

## Recommended v0.1 Stack
- Structured state: SQLite in WAL mode
- Semantic index: FAISS exact-first baseline with adapter to evolve index strategy
- Cloud optionality: adapter contract for Cosmos-backed state and remote embedding service in v0.3

## Open Questions
- Trigger thresholds for switching from exact to ANN index strategy
- Default embedding model/provider for local/offline vs cloud burst modes
- Cross-platform packaging plan for FAISS dependencies

## Recommended Decisions
- Adopt Option A for v0.1.
- Defer ANN complexity until metric evidence (latency/cost/recall quality) indicates need.
