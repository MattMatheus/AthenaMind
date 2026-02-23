# AthenaMind v0.1

AthenaMind v0.1 is a Go CLI memory layer for agent workflows. It provides file-based memory operations with governance-gated mutations, retrieval evaluation, read access APIs, and deterministic lifecycle artifacts.

## What v0.1 Delivers Today
Implemented command families in `cmd/memory-cli`:
- `write`
- `retrieve`
- `evaluate`
- `snapshot` (`create`, `list`, `restore`)
- `serve-read-gateway`
- `api-retrieve`
- `bootstrap`
- `episode` (`write`, `list`)

Operational capabilities currently implemented:
- File-based memory storage and index operations under a local root (default `memory/`)
- Governance-enforced mutation gates and review evidence requirements for writes/restores
- Constraint checks for operation budgets and autonomy policies
- Telemetry event output for write/retrieve/evaluate/episode flows
- Snapshot create/list/restore with manifest and audit event tracking
- Read-gateway parity checks between API retrieval and local CLI retrieval behavior

Scope boundary reference:
- `research/decisions/ADR-0007-memory-layer-scope-refinement.md`

## Quick Start
Write a memory entry:

```bash
go run ./cmd/memory-cli write \
  --root memory \
  --id handoff-template \
  --title "Handoff Template" \
  --type instruction \
  --domain docs \
  --body "Always include risks, evidence, and next-state recommendation." \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "baseline docs quality" \
  --risk "low; reversible by git revert" \
  --notes "approved for docs baseline"
```

Retrieve a memory entry:

```bash
go run ./cmd/memory-cli retrieve \
  --root memory \
  --query "handoff instruction template"
```

More command examples:
- `docs/cli/examples.md`
- `docs/cli/commands.md`

## Future Direction
Long-term product vision (preserved separately from v0.1 scope):
- `docs/product/VISION.md`

Phased plan (v0.1 -> v0.3):
- `research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`

## Development Workflow
Stage-based flow:
- `HUMANS.md`
- `DEVELOPMENT_CYCLE.md`
- `docs/process/STAGE_EXIT_GATES.md`
- `docs/process/BACKLOG_WEIGHTING_POLICY.md`

## Quality Gates
Run locally before handoff:
- `scripts/run_doc_tests.sh`
- `go test ./...`

CI gate:
- `azure-pipelines.yml` runs `go test ./...` on push and PR.
