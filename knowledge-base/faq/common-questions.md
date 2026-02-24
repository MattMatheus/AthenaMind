# Common Questions

## What is AthenaMind focused on?

A governable memory layer for agentic coding workflows: write/retrieve/evaluate/snapshot/episode flows with deterministic safety behavior.

## Is AthenaMind a runtime orchestration platform?

No. It is a memory layer. Runtime orchestration is intentionally out of scope for v0.1 product focus.

## What is AthenaWork in this slim repo?

AthenaWork is the staged operational workflow system. In slim mode, the full operator pack is archived at `/Users/foundry/Experiments/Archived/AthenaMind-internal-2026-02-24/products/athena-work`.

## Does OTLP require a secret?

Sometimes. Local/open collectors may not require auth; managed remote collectors usually require token headers or mTLS.

## Why did retrieval return fallback mode?

Semantic confidence was insufficient or latency degradation policy forced deterministic fallback.

## How do I tune retrieval behavior?

Adjust mode/backend/top-k, embedding endpoint, and quality query sets; then re-run `evaluate`.
