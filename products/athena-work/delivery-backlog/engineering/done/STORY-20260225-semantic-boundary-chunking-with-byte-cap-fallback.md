# Story: Add semantic-boundary chunking with byte-cap fallback for embeddings

## Metadata
- `id`: STORY-20260225-semantic-boundary-chunking-with-byte-cap-fallback
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0005, ADR-0012]
- `success_metric`: retrieval embeddings for structured markdown and code inputs use semantic-boundary chunks first while preserving zero 400/413 failures on Azure embedding calls in regression tests
- `release_checkpoint`: required

## Problem Statement
- Current embedding chunking in `internal/retrieval/embedding.go` is byte-cap-first. It protects Azure/Ollama limits, but it can split across headings/functions and degrade semantic retrieval quality. We need semantic-boundary-first chunking while keeping strict byte-cap fallback safety.

## Scope
- In:
  - Implement semantic-boundary-first chunk planning for embedding inputs.
  - Preserve hard byte-cap fallback behavior for provider safety.
  - Add unit tests for markdown heading boundaries, code block/function boundaries, and long unstructured text fallback.
  - Ensure existing embedding error handling/fallback behavior remains intact.
- Out:
  - Changes to retrieval ranking, confidence gates, recency boosting, or fallback chain.
  - Changes to vector backend providers.

## Acceptance Criteria
1. Embedding request chunking prefers semantic boundaries (markdown headings, paragraph/function-like block boundaries) before byte-size splitting.
2. Every emitted chunk remains within configured max byte size.
3. Existing provider-safety behavior remains: long inputs are still accepted and chunked to avoid endpoint failures.
4. Existing retrieval tests remain green; new tests cover semantic chunk behavior and byte-cap guarantees.
5. `tools/run_stage_tests.sh` passes for touched scope.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- None

## Notes
- This story intentionally keeps the byte-cap limit as a non-negotiable hard constraint and improves chunk quality before that limit is applied.
