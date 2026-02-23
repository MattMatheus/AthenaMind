# AthenaMind Operator Skill

## Purpose
Operate AthenaMind as a memory system for indexing, retrieval, and bootstrap support.

Repository split note:
- AthenaMind product root: `/Users/foundry/AthenaMind/products/athena-mind`
- AthenaWork product root: `/Users/foundry/AthenaMind/products/athena-work`
- Dedicated skills:
  - `/Users/foundry/AthenaMind/skills/athena-mind/SKILL.md`
  - `/Users/foundry/AthenaMind/skills/athena-work/SKILL.md`

## Working Directories
- Toolchain root: `/Users/foundry/AthenaMind/products/athena-mind`
- Memory root (recommended): `/Users/foundry/AthenaMind/products/athena-mind/memory/core`
- Work root (optional): `/Users/foundry/AthenaMind/products/athena-mind/memory/work`

## Environment
1. Load environment variables:
   ```bash
   set -a; source /Users/foundry/AthenaMind/.env; set +a
   ```
2. For context-first retrieval (no latency fallback), set:
   ```bash
   export MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS=0
   ```

## Initialize Memory
```bash
mkdir -p /Users/foundry/AthenaMind/products/athena-mind/memory/core /Users/foundry/AthenaMind/products/athena-mind/memory/work
```

## Core Commands
Run from `/Users/foundry/AthenaMind/products/athena-mind`.

1. Bootstrap
```bash
go run ./cmd/memory-cli bootstrap \
  --root /Users/foundry/AthenaMind/products/athena-mind/memory/core \
  --repo AthenaMind \
  --session-id setup-session \
  --scenario setup
```

2. Crawl markdown content
```bash
go run ./cmd/memory-cli crawl \
  --root /Users/foundry/AthenaMind/products/athena-mind/memory/core \
  --dir /Users/foundry/AthenaMind/knowledge-base \
  --domain docs-crawl \
  --reviewer system
```

3. Reindex embeddings (Azure-backed when env is loaded)
```bash
go run ./cmd/memory-cli reindex-all \
  --root /Users/foundry/AthenaMind/products/athena-mind/memory/core
```

4. Verify embedding coverage
```bash
go run ./cmd/memory-cli verify embeddings \
  --root /Users/foundry/AthenaMind/products/athena-mind/memory/core
```

5. Retrieve
```bash
go run ./cmd/memory-cli retrieve \
  --root /Users/foundry/AthenaMind/products/athena-mind/memory/core \
  --query "memory lifecycle" \
  --domain docs-crawl
```

## Verification
- Confirm semantic mode appears in retrieve output:
  - `selection_mode: embedding_semantic`
- Confirm embeddings coverage:
  ```bash
  sqlite3 /Users/foundry/AthenaMind/products/athena-mind/memory/core/index.db \
  "select (select count(*) from entries),(select count(*) from embeddings),(select count(*) from entries e left join embeddings em on em.entry_id=e.id where em.entry_id is null);"
  ```

## Operational Notes
- If external network is sandbox-restricted, run embedding-related commands with unrestricted network.
- `crawl` now generates deterministic path-based IDs, so duplicate basenames do not collide.
- Do not commit `.env` files.
