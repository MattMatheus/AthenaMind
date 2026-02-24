# Athena Platform

Athena is split into two products:

- `AthenaMind`: memory layer for agentic coding workflows (primary runtime in this repo).
- `AthenaWork`: staged operating system for planning/engineering/QA/PM loops (fully available in this repo).

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

1. Install prerequisites: [Installation](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/getting-started/installation.md)
2. Download precompiled binaries: [Precompiled Binaries](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/getting-started/binaries.md)
3. Run first end-to-end flow: [Quickstart](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/getting-started/quickstart.md)
4. Learn command surface: [CLI Commands](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/cli/commands.md)
5. Configure observability: [OTel/OTLP Setup](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/how-to/memory-cli-otel-setup.md)

## Documentation Map

- Full docs index: [knowledge-base/INDEX.md](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/INDEX.md)
- AthenaMind product docs: [knowledge-base/product/athenamind.md](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/product/athenamind.md)
- AthenaWork product docs: [knowledge-base/product/athenawork.md](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/product/athenawork.md)

## AthenaWork Operator Paths

- Human operator guide: [HUMANS.md](https://github.com/MattMatheus/AthenaMind/blob/main/HUMANS.md)
- Agent rules: [AGENTS.md](https://github.com/MattMatheus/AthenaMind/blob/main/AGENTS.md)
- Stage launcher: [tools/launch_stage.sh](https://github.com/MattMatheus/AthenaMind/blob/main/tools/launch_stage.sh)
- Specialist directory: [staff-personas/STAFF_DIRECTORY.md](https://github.com/MattMatheus/AthenaMind/blob/main/staff-personas/STAFF_DIRECTORY.md)
- Queue and work-system artifacts: [delivery-backlog](https://github.com/MattMatheus/AthenaMind/blob/main/delivery-backlog), [operating-system](https://github.com/MattMatheus/AthenaMind/blob/main/operating-system)


## Public Testing

- Program: [PUBLIC_TESTING.md](https://github.com/MattMatheus/AthenaMind/blob/main/PUBLIC_TESTING.md)
- Contribution guide: [CONTRIBUTING.md](https://github.com/MattMatheus/AthenaMind/blob/main/CONTRIBUTING.md)
- Testing protocol: [TESTING.md](https://github.com/MattMatheus/AthenaMind/blob/main/TESTING.md)
- Security and support: [SECURITY.md](https://github.com/MattMatheus/AthenaMind/blob/main/SECURITY.md), [SUPPORT.md](https://github.com/MattMatheus/AthenaMind/blob/main/SUPPORT.md)

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

See [LICENSE](https://github.com/MattMatheus/AthenaMind/blob/main/LICENSE).
