# Athena Platform (Slim)

Athena is split into two products:

- `AthenaMind`: memory layer for agentic coding workflows (primary runtime in this repo).
- `AthenaWork`: staged operating system for planning/engineering/QA/PM loops (documented here, operator pack archived externally).

This repository is a **slim product distribution** focused on user-facing runtime and docs.

## Who This Is For

- Engineers and data scientists operating memory-assisted coding agents.
- Technical operators tuning retrieval quality, governance controls, and telemetry export.
- Teams evaluating local-first memory workflows before broader deployment.

## What You Can Do Today

- Write and retrieve governed memory entries.
- Evaluate retrieval quality with deterministic fallback checks.
- Use snapshots for rollback and recovery.
- Capture episodes for cross-session continuity.
- Run local read gateway and API retrieval fallback.
- Export telemetry events and OpenTelemetry traces (including OTLP collectors).

## Quick Start

1. Install prerequisites: [Installation](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/getting-started/installation.md)
2. Run first end-to-end flow: [Quickstart](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/getting-started/quickstart.md)
3. Learn command surface: [CLI Commands](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/cli/commands.md)
4. Configure observability: [OTel/OTLP Setup](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/how-to/MEMORY_CLI_OTEL_SETUP.md)

## Documentation Map

- Full docs index: [knowledge-base/INDEX.md](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/INDEX.md)
- AthenaMind product docs: [knowledge-base/product/athenamind.md](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenamind.md)
- AthenaWork product docs: [knowledge-base/product/athenawork.md](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork.md)

## Slim Boundary

Non user-facing research/process content was moved out of this repo to keep distribution focused.

Archived internal pack location:
- `/Users/foundry/Experiments/Archived/AthenaMind-internal-2026-02-24`


## Public Testing

- Program: [PUBLIC_TESTING.md](/Users/foundry/Experiments/Current/AthenaMind/PUBLIC_TESTING.md)
- Contribution guide: [CONTRIBUTING.md](/Users/foundry/Experiments/Current/AthenaMind/CONTRIBUTING.md)
- Testing protocol: [TESTING.md](/Users/foundry/Experiments/Current/AthenaMind/TESTING.md)
- Security and support: [SECURITY.md](/Users/foundry/Experiments/Current/AthenaMind/SECURITY.md), [SUPPORT.md](/Users/foundry/Experiments/Current/AthenaMind/SUPPORT.md)

## Runtime Commands

`cmd/memory-cli` supports:

- `write`
- `retrieve`
- `evaluate`
- `bootstrap`
- `verify`
- `snapshot` (`create|list|restore`)
- `episode` (`write|list`)
- `crawl`
- `reindex-all`
- `reembed-changed`
- `sync-qdrant`
- `serve-read-gateway`
- `api-retrieve`
- `telemetry tail`

## Current Validation Status

- `go test ./...` passes in this repo.

## License

See [LICENSE](/Users/foundry/Experiments/Current/AthenaMind/LICENSE).
