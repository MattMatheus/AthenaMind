# Story: Azure OpenAI Foundry Integration

## Metadata
- `id`: STORY-20260223-azure-openai-integration
- `owner_persona`: Software Architect - Ada.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0005, ADR-0017]
- `success_metric`: Sub-200ms semantic retrieval via Azure OpenAI
- `release_checkpoint`: required

## Problem Statement
Local embedding models (Ollama) on older hardware trigger latency degradation gates. We need a high-performance external embedding provider.

## Scope
- In: Native Azure OpenAI support in `GenerateEmbedding`, Entra ID authentication, batch support.
- Out: Other cloud providers (AWS, GCP).

## Acceptance Criteria
1. `GenerateEmbedding` uses Azure OpenAI if configured.
2. Entra ID Service Principal authentication is supported.
3. Batch embedding is supported for large operations.
4. Retrieval latency stays below the 700ms P95 threshold (conceptually).

## QA Checks
- Verified Entra ID token fetching and caching.
* Verified semantic retrieval with Azure-generated vectors.
* Verified batch re-indexing.

## Dependencies
- None

## Notes
- Integrated into `internal/retrieval/embedding.go`.
