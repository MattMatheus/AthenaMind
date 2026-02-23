# AthenaMind

AthenaMind is a Go-based memory toolchain for agent workflows. It supports governed writes, retrieval, bootstrap payload generation, snapshots, episode logging, and embedding-backed semantic search over a local memory root.

## Current Operating Model
- Toolchain root: `/Users/foundry/Source/orchestrator/AthenaMind`
- Active memory root: `/Users/foundry/Source/orchestrator/AthenaMind-Memory/core`
- Optional work root: `/Users/foundry/Source/orchestrator/AthenaMind-Memory/work`
- Embeddings: Azure OpenAI when `AZURE_OPENAI_ENDPOINT` and credentials are present
- Latency fallback policy: configurable with `MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS` (`0` disables latency fallback)

## Command Surface
Implemented command families in `cmd/memory-cli`:
- `write`
- `retrieve`
- `evaluate`
- `bootstrap`
- `verify` (`embeddings`)
- `reindex-all`
- `crawl`
- `snapshot` (`create`, `list`, `restore`)
- `episode` (`write`, `list`)
- `serve-read-gateway`
- `api-retrieve`

## Practical Workflow
1. Load environment:
```bash
set -a; source /Users/foundry/Source/orchestrator/AthenaMind/.env; set +a
```

2. Ensure memory roots exist:
```bash
mkdir -p /Users/foundry/Source/orchestrator/AthenaMind-Memory/core /Users/foundry/Source/orchestrator/AthenaMind-Memory/work
```

3. Crawl docs into memory (collision-safe IDs are path-based and deterministic):
```bash
go run ./cmd/memory-cli crawl \
  --root /Users/foundry/Source/orchestrator/AthenaMind-Memory/core \
  --dir /Users/foundry/Source/orchestrator/AthenaMind/docs \
  --domain docs-crawl \
  --reviewer system
```

4. Build missing embeddings:
```bash
go run ./cmd/memory-cli reindex-all \
  --root /Users/foundry/Source/orchestrator/AthenaMind-Memory/core
```

5. Verify embedding coverage:
```bash
go run ./cmd/memory-cli verify embeddings \
  --root /Users/foundry/Source/orchestrator/AthenaMind-Memory/core
```

6. Retrieve:
```bash
go run ./cmd/memory-cli retrieve \
  --root /Users/foundry/Source/orchestrator/AthenaMind-Memory/core \
  --query "memory lifecycle" \
  --domain docs-crawl
```

## Test Status
| Gate | Status | Notes |
|---|---|---|
| Targeted governance and memory-cli tests | ✔ PASS | `go test ./internal/governance` and targeted `./cmd/memory-cli` tests are passing |
| Full memory-cli package | ✖ FAIL | Known failing test: `TestSemanticConfidenceGate` |
| Full repository suite (`go test ./...`) | ✖ FAIL | Not currently green due to the memory-cli package failure above |

## References
- Scope boundary: `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- CLI docs: `docs/cli/commands.md`, `docs/cli/examples.md`
- Process docs: `docs/process/STAGE_EXIT_GATES.md`, `docs/process/BACKLOG_WEIGHTING_POLICY.md`
