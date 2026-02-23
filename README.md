# AthenaMind v0.1

AthenaMind is a local-first memory layer for agent workflows.

## v0.1 Scope (Current)
In scope (ADR-0007):
- Procedural memory write/retrieve flows
- State and semantic retrieval behavior
- Auditability, telemetry, schema validation, and snapshots

Out of scope for v0.1:
- Runtime orchestration ownership
- Pod/container lifecycle management
- Sandboxed execution platform responsibilities

See `research/decisions/ADR-0007-memory-layer-scope-refinement.md`.

## Current CLI
Primary command: `memory-cli` (`cmd/memory-cli`).

Implemented command groups:
- `write`
- `retrieve`
- `evaluate`
- `snapshot` (`create`, `list`, `restore`)
- `serve-read-gateway`
- `api-retrieve`

## Development Workflow
Stage-based flow is mandatory:
- `HUMANS.md`
- `DEVELOPMENT_CYCLE.md`
- `docs/process/STAGE_EXIT_GATES.md`

Backlog ranking policy:
- Product-first queue weighting: `docs/process/BACKLOG_WEIGHTING_POLICY.md`

## Quality Gates
Local verification:
- `go test ./...`
- `scripts/run_doc_tests.sh`

CI verification (Azure DevOps):
- `azure-pipelines.yml` runs `go test ./...` on push and PR.
