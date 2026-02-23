# CLI Commands

## Summary
Command reference for all currently supported `memory-cli` operations.

## Root Command
```bash
memory-cli <write|retrieve|snapshot|serve-read-gateway|api-retrieve|evaluate> [flags]
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
- `--approved` (legacy alias for approved decision)
- `--rework-notes` and `--re-reviewed-by` (required for rejected decision)

## `retrieve`
Runs semantic retrieval with deterministic fallback.

Required:
- `--query`

Optional:
- `--root` (default `memory`)
- `--domain`
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

## References
- `cmd/memory-cli/main.go`
