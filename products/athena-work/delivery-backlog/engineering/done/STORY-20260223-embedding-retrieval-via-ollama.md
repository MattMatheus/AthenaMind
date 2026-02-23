# Story: Add embedding-based retrieval via local Ollama

## Metadata
- `id`: STORY-20260223-embedding-retrieval-via-ollama
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: qa
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0005, ADR-0012]
- `success_metric`: embedding-based retrieval achieves higher top-1 useful rate than token-overlap baseline on the existing 50-query evaluation benchmark
- `release_checkpoint`: required

## Problem Statement
- The current semantic retrieval uses token-overlap scoring, which limits recall quality for natural language queries. ADR-0005 specifies FAISS exact-first indexing and ADR-0012 defines retrieval quality gates. Embedding-based retrieval is the single largest quality improvement available for the retrieval pipeline. Using Ollama in Docker provides local embedding generation at zero cloud cost.

## Scope
- In:
  - Integrate with Ollama API for local embedding generation (e.g., `nomic-embed-text` or similar small model)
  - Store embeddings in SQLite (or simple flat vector store) alongside index entries
  - Update retrieval scoring pipeline: embedding similarity → confidence gate → existing fallback chain
  - Add `--embedding-endpoint` flag to CLI for Ollama URL (default: `http://localhost:11434`)
  - Graceful fallback: if Ollama is unavailable, fall back to token-overlap scoring with warning
  - Add Docker Compose file for Ollama service
  - Run evaluation harness and compare against token-overlap baseline
- Out:
  - Cloud embedding services
  - FAISS ANN structures (exact/flat search is sufficient for v0.2 scale)
  - Changes to governance or write path

## Acceptance Criteria
1. Embeddings generated via Ollama API for memory entries on write.
2. Retrieval uses embedding similarity as primary scoring when embeddings are available.
3. Fallback to token-overlap scoring when Ollama is unavailable.
4. Top-1 useful rate on 50-query benchmark improves over token-overlap baseline.
5. Confidence gating and deterministic fallback chain behavior preserved.
6. Docker Compose file provided for local Ollama service.
7. `go test ./...` passes (tests use mock/stub for embedding endpoint).

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260223-sqlite-index-store (embeddings stored in SQLite)
- STORY-20260223-split-main-go-into-packages (retrieval logic should be in internal/retrieval)

## Notes
- Start with a small embedding model to keep Docker resource usage low.
- The evaluation harness already exists — use it as the quality gate for this change.
- This is the story that transitions from "token matching" to "semantic understanding."
