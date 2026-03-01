# Engineering Handoff: STORY-20260225-semantic-boundary-chunking-with-byte-cap-fallback

## What Changed
- Added semantic-boundary-first chunk planning with strict byte-cap enforcement:
  - `internal/retrieval/embedding.go`
  - New chunker: `chunkTextBySemanticBoundaries(...)`
  - New byte-safe splitter: `splitTextByBytes(...)` with rune-safe cutoff via `maxPrefixByBytes(...)`
- Updated embedding generation flow to chunk per input, generate per-chunk embeddings, and average chunk vectors:
  - `GenerateEmbeddings(...)` now applies semantic chunking for both Azure and Ollama paths.
  - Added `averageEmbeddings(...)` and `generateOllamaEmbeddings(...)`.
- Preserved provider safety behavior:
  - Every emitted chunk is hard-capped to `embeddingChunkMaxBytes` (2800 bytes) before provider calls.
  - Azure remains first-priority when configured; Ollama fallback behavior remains in place.

## Why It Changed
- The story requires improved semantic coherence in chunking while retaining strict request-size safety to avoid Azure/Ollama request failures.
- Prior behavior sent one full prompt per entry in this branch; large structured content could lose semantic locality or fail on provider limits.

## Test Updates Made
- Added `internal/retrieval/embedding_test.go` with:
  - `TestChunkTextBySemanticBoundariesPrefersHeadings`
  - `TestChunkTextBySemanticBoundariesRespectsByteCap`
  - `TestGenerateEmbeddingsChunksAndAveragesOllama`
  - `TestGenerateEmbeddingsChunksAzureInputsAndAverages`

## Test Run Results
- `go test ./internal/retrieval/...` passed.
- `go test ./...` passed.
- `./tools/run_stage_tests.sh` could not be run because the script does not exist in this repo.

## Open Risks / Questions
- Semantic boundary detection currently uses heuristic boundaries (markdown headings, blank lines, common function/class prefixes). It is intentionally conservative but not parser-accurate for every language.
- Averaging chunk embeddings is simple and stable; weighted averaging by chunk length could be considered if retrieval evals indicate benefit.

## Recommended QA Focus Areas
- Verify structured markdown and code-like inputs retrieve at least as well as before on representative queries.
- Validate Azure path still avoids oversize payload failures for long documents.
- Confirm no regression in embedding fallback behavior when Azure endpoint fails and Ollama is used.
