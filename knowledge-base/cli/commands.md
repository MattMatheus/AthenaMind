# CLI Commands

## Summary
Command reference for all currently supported `memory-cli` operations.

## Root Command
```bash
memory-cli <write|retrieve|snapshot|serve-read-gateway|api-retrieve|evaluate|bootstrap|reindex-all|crawl|reembed-changed|sync-qdrant|verify|episode|telemetry> [flags]
```

## `write`
Creates or updates an entry.

Required:
- `--id`
- `--title`
- `--type prompt|instruction`
- `--domain`
- `--body` or `--body-file`
- `--stage planning|architect|pm`
- `--reviewer`
- `--decision approved|rejected`
- `--reason`
- `--risk`
- `--notes`

Optional:
- `--root` (default `memory`)
- `--session-id`
- `--scenario-id`
- `--memory-type`
- `--operator-verdict`
- `--telemetry-file`
- `--embedding-endpoint` (default `http://localhost:11434`)
- `--approved` (legacy alias for approved decision)
- `--rework-notes` and `--re-reviewed-by` (required for rejected decision)

## `retrieve`
Runs semantic retrieval with deterministic fallback.

Required:
- `--query`

Optional:
- `--root` (default `memory`)
- `--domain`
- `--embedding-endpoint` (default `http://localhost:11434`)
- `--mode classic|hybrid` (default `classic`)
- `--top-k` (default `5`)
- `--retrieval-backend sqlite|qdrant|neo4j` (default `sqlite`)
- `--session-id`
- `--scenario-id`
- `--memory-type`
- `--operator-verdict`
- `--telemetry-file`

## `evaluate`
Runs retrieval quality evaluation and prints a JSON report.

Optional:
- `--root` (default `memory`)
- `--query-file` (default `cmd/memory-cli/testdata/eval-query-set-v1.json`)
- `--corpus-id`
- `--query-set-id`
- `--config-id`
- `--embedding-endpoint` (default `http://localhost:11434`)
- `--mode classic|hybrid` (default `classic`)
- `--top-k` (default `5`)
- `--retrieval-backend sqlite|qdrant|neo4j` (default `sqlite`)
- telemetry flags (`--session-id`, `--scenario-id`, `--memory-type`, `--operator-verdict`, `--telemetry-file`)

## `snapshot`
Snapshot subcommands:
- `snapshot create`
  - required: `--created-by`, `--reason`
  - optional: `--scope` (`full` only), `--root`, `--session-id`
- `snapshot list`
  - optional: `--root`
- `snapshot restore`
  - required: `--snapshot-id`, restore review evidence (`--reviewer --decision --reason --risk --notes`)
  - optional: `--stage` (default `pm`), `--root`, `--session-id`, rejection fields

## `serve-read-gateway`
Starts local read gateway.

Optional:
- `--root` (default `memory`)
- `--addr` (default `127.0.0.1:8788`)

## `api-retrieve`
Calls read gateway and enforces parity with local CLI contract.

Required:
- `--query`
- `--session-id`

Optional:
- `--root` (default `memory`)
- `--domain`
- `--gateway-url`
- `--mode classic|hybrid` (default `classic`)
- `--top-k` (default `5`)
- `--retrieval-backend sqlite|qdrant|neo4j` (default `sqlite`)

## `bootstrap`
Builds a memory bootstrap payload for agent startup.

Required:
- `--repo`
- `--session-id`
- `--scenario`

Optional:
- `--root` (default `memory`)
- telemetry flags (`--memory-type`, `--operator-verdict`, `--telemetry-file`)

Bootstrap payload schema:
- top-level: `repo`, `session_id`, `scenario`, `generated_at`, `memory_entries`, optional `episode`
- `memory_entries[]`: `id`, `selection_mode`, `source_path`, `confidence`, `reason`, `type`, `domain`, `title`
- `episode` (when available from episode store): `repo`, `scenario`, `cycle_id`, `story_id`, `outcome`, `summary`, `timestamp_utc`

## `episode`
Episode subcommands:
- `episode write`
  - required: `--repo`, `--session-id`, `--cycle-id`, `--story-id`, `--outcome`, `--summary` or `--summary-file`, `--decisions` or `--decisions-file`
  - required governance review: `--reviewer --decision --reason --risk --notes`
  - optional: `--files-changed`, `--stage` (default `pm`), `--root`, `--telemetry-file`, rejection evidence fields
- `episode list`
  - required: `--repo`
  - optional: `--root`

## `reindex-all`
Rebuilds missing embeddings for currently indexed entries.

Optional:
- `--root` (default `memory`)
- `--embedding-endpoint` (default `http://localhost:11434`)

## `crawl`
Crawls markdown docs and indexes them as instructions with deterministic path-based IDs.

Required:
- `--dir`

Optional:
- `--root` (default `memory`)
- `--domain` (default `auto-crawled`)
- `--reviewer` (default `system`)
- `--embedding-endpoint` (default `http://localhost:11434`)

## `reembed-changed`
Re-indexes and re-embeds changed markdown files for incremental consistency after cycle changes.

Required:
- `--files-changed` (comma-separated file paths)

Optional:
- `--root` (default `memory`)
- `--repo-root` (default `.`)
- `--domain` (default `auto-crawled`, used for newly discovered files)
- `--reviewer` (default `system`)
- `--session-id`
- `--embedding-endpoint` (default `http://localhost:11434`)

## `sync-qdrant`
Pushes local embedding records into a Qdrant collection for backend experiments.

Optional:
- `--root` (default `memory`)
- `--qdrant-url` (default env `ATHENA_QDRANT_URL` or `http://localhost:6333`)
- `--collection` (default env `ATHENA_QDRANT_COLLECTION` or `athena_memories`)
- `--batch-size` (default `128`)

## `verify`
Verification subcommands:
- `verify embeddings`
  - optional: `--root`, `--show-missing`
  - reports embedding coverage for indexed entries
- `verify health`
  - optional: `--root`, `--query`, `--domain`, `--session-id`, `--embedding-endpoint`
  - runs semantic retrieval health check and reports pass/fail

## `telemetry`
Telemetry subcommands:
- `telemetry tail`
  - optional: `--root` (default `memory`)
  - optional: `--lines` (default `20`, recent records per source)
  - optional: `--source events|retrieval|both` (default `events`)
  - optional filters: `--operation`, `--result`, `--session-id`
  - optional live mode: `--follow` (stream new records), `--follow-poll-ms` (default `500`), `--follow-seconds` (default `0`, run until interrupted)
  - optional: `--telemetry-file` (default `<root>/telemetry/events.jsonl`)
  - optional: `--retrieval-metrics-file` (default `<root>/telemetry/retrieval-metrics.jsonl`)
  - prints a JSON payload with recent records from local telemetry files; when `--follow` is set, streams additional JSON lines as new records are written

## References
- `cmd/memory-cli/main.go`
- `cmd/memory-cli/commands.go`
