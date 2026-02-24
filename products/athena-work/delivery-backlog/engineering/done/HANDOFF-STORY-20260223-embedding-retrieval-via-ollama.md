# Engineering Handoff: STORY-20260223-embedding-retrieval-via-ollama

## What Changed
- Added embedding support with Ollama integration:
  - `internal/retrieval/embedding.go`
  - default endpoint constant: `http://localhost:11434`
  - `/api/embeddings` client with fast-fail cache and timeout
- Added embedding index/update flow on write:
  - `cmd/memory-cli/commands.go` adds `--embedding-endpoint` to `write`
  - after `index.UpsertEntry`, write path attempts `retrieval.IndexEntryEmbedding(...)`
  - graceful warning if embedding endpoint unavailable (write still succeeds)
- Added embedding storage in SQLite:
  - `internal/index/sqlite_store.go` creates `embeddings` table
  - added `index.UpsertEmbedding(...)` and `index.GetEmbeddings(...)`
- Updated retrieval pipeline to prefer embedding similarity:
  - `internal/retrieval/retrieval.go` adds `RetrieveWithEmbeddingEndpoint(...)`
  - embedding cosine similarity used as primary score when vectors are available
  - preserved confidence gate and deterministic fallback behavior
  - graceful fallback with warning to token-overlap when endpoint/store/vectors unavailable
- Updated evaluation path:
  - `internal/retrieval/eval.go` adds `EvaluateRetrievalWithEmbeddingEndpoint(...)`
  - `evaluate` command now accepts `--embedding-endpoint`
- Added compose file for local Ollama:
  - `docker-compose.ollama.yml`
- Updated CLI docs:
  - `knowledge-base/cli/commands.md`
  - `knowledge-base/cli/examples.md`

## Why It Changed
- Story requires embedding-based retrieval as the primary scorer while preserving existing governance and deterministic fallback behavior.
- Local Ollama integration provides model access without cloud dependency and supports v0.2 semantic quality improvement goals.

## Test Updates Made
- Updated `internal/retrieval/retrieval_test.go`:
  - `TestRetrieveUsesEmbeddingSimilarityWhenAvailable`
  - `TestRetrieveFallsBackWhenEmbeddingEndpointUnavailable`
  - `TestEmbeddingScoringImprovesOverTokenBaseline`
- Existing suite coverage retained for CLI/retrieval/snapshot/index paths.

## Test Run Results
- `go test ./internal/retrieval ./cmd/memory-cli` passed.
- `go test ./...` passed.
- `tools/run_doc_tests.sh` passed.

## Open Risks / Questions
- Embedding client currently uses `/api/embeddings` and model name `nomic-embed-text`; operators must ensure model availability in local Ollama runtime.
- Endpoint failure cache is process-local and time-based (30s); acceptable for CLI flow, but long-running daemon mode may want richer health-state handling.

## Recommended QA Focus Areas
- Verify `write` with reachable Ollama endpoint stores vectors and does not regress existing write semantics.
- Verify `retrieve` prioritizes embedding scoring when vectors exist and preserves deterministic fallback otherwise.
- Verify explicit warning output when embedding endpoint is unreachable.
- Verify `evaluate --embedding-endpoint` works with stubbed endpoint and preserves output schema.
- Validate `docker-compose.ollama.yml` service startup semantics and local endpoint availability.

## New Gaps Discovered
- Potential follow-up: add explicit CLI command for embedding backfill of pre-existing entries written before embedding endpoint availability.
