# Common Errors

## Summary
Quick lookup for frequent `memory-cli` failures and direct fixes.

## Command/Input Errors
- `--query is required`
  - Fix: provide `--query` for `retrieve` or `api-retrieve`.
- `--session-id is required`
  - Fix: provide `--session-id` for `api-retrieve`.
- `--id --title --type --domain are required`
  - Fix: include all required write fields.

## Governance Errors
- `ERR_MUTATION_NOT_ALLOWED_DURING_AUTONOMOUS_RUN`
  - Fix: run write/restore outside autonomous cycle.
- `ERR_MUTATION_STAGE_INVALID`
  - Fix: use `--stage planning|architect|pm`.
- `ERR_MUTATION_EVIDENCE_REQUIRED`
  - Fix: provide `--reason --risk --notes`.

## Snapshot Errors
- `ERR_SNAPSHOT_MANIFEST_INVALID`
  - Fix: ensure required create/manifest fields and `scope=full`.
- `ERR_SNAPSHOT_COMPATIBILITY_BLOCKED`
  - Fix: restore from compatible schema major version.
- `ERR_SNAPSHOT_INTEGRITY_CHECK_FAILED`
  - Fix: regenerate snapshot; verify payload integrity.

## Quality/Constraint Errors
- `ERR_CONSTRAINT_COST_BUDGET_EXCEEDED`
  - Fix: adjust operation/budget env or split workflow.
- `ERR_CONSTRAINT_TRACEABILITY_INCOMPLETE`
  - Fix: ensure trace/session fields are present.
- `ERR_API_CLI_PARITY_MISMATCH`
  - Fix: investigate gateway divergence from local retrieve contract.

## Escalation
- Create a bug using `backlog/engineering/intake/BUG_TEMPLATE.md`.
- Include exact command, output, and environment details.
