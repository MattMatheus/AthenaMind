# AthenaMind Product Guide

## What AthenaMind Is

AthenaMind is a local-first memory layer for agentic coding workflows. It focuses on durable, governable memory operations rather than runtime execution orchestration.

## Product Scope

In scope:
- Procedural memory (instructions/prompts).
- State and semantic memory (episodes, retrieval).
- Governance and auditability (review evidence, deterministic behavior).
- Observability (telemetry events and OTel traces).

Out of scope for v0.1:
- Owning code execution runtime orchestration.
- Pod/container lifecycle management as a product responsibility.

## User Personas

- `Engineer`: wants reliable recall and deterministic retrieval behavior.
- `Data Scientist`: wants tunable retrieval/backends and measurable quality gates.
- `Operator`: wants policy controls, audit evidence, and telemetry export.

## Architecture (Practical View)

- CLI/API surface: `cmd/memory-cli`
- Core modules: `internal/index`, `internal/retrieval`, `internal/governance`, `internal/snapshot`, `internal/episode`, `internal/telemetry`
- Storage: local memory root (index, metadata, entries, telemetry)

## Quality Model

- Retrieval modes: `classic`, `hybrid`
- Backends: `sqlite` (default), optional `qdrant`, optional `neo4j`
- Deterministic fallback when semantic confidence is insufficient
- Evaluation harness with useful-rate + determinism gates

## Governance Model

- Mutation writes require reviewer evidence fields.
- Constraint checks enforce cost/traceability/reliability rules.
- Latency degradation forces deterministic fallback behavior.

## What To Tune

- Retrieval mode, top-k, backend
- Embedding endpoint choice
- Constraint env vars
- OTel/OTLP export configuration

## Start Here

- [Installation](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/getting-started/installation.md)
- [Quickstart](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/getting-started/quickstart.md)
- [CLI Commands](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/cli/commands.md)
- [Configuration Reference](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/references/configuration.md)
