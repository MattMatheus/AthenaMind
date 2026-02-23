# Engineering Handoff: STORY-20260223-sqlite-index-store

## What Changed
- Added a storage adapter boundary in `internal/index`:
  - `indexStore` interface (`Load`, `Save`)
  - default implementation `sqliteIndexStore`
- Implemented SQLite-backed index store in `internal/index/sqlite_store.go`:
  - SQLite file at `<root>/index.db`
  - WAL mode enabled via `PRAGMA journal_mode=WAL`
  - `meta` and `entries` tables for index schema/version and entry records
- Updated `internal/index/index.go` to:
  - load index through adapter
  - save index updates to SQLite
  - retain `index.yaml` mirror write for compatibility paths
  - support automatic legacy `index.yaml` import on first SQLite initialization
- Updated snapshot refs in `internal/snapshot/snapshot.go`:
  - include `index.db` in snapshot payload when present
- Expanded index tests in `internal/index/index_test.go`:
  - sqlite initialization on first load
  - legacy `index.yaml` migration to SQLite
  - WAL mode assertion after upsert

## Why It Changed
- Story requires migration from full-file `index.yaml` reads/writes to a SQLite-backed structured index store with WAL mode and a migration path for existing data.
- Adapter seam is required so future storage backends can be introduced without changing retrieval/write call sites.

## Test Updates Made
- Updated: `internal/index/index_test.go`
  - `TestLoadIndexInitializesMissingFile`
  - `TestLoadIndexMigratesLegacyYAMLToSQLite`
  - `TestUpsertEntryEnablesSQLiteWAL`

## Test Run Results
- `go test ./internal/index ./internal/snapshot ./cmd/memory-cli` passed.
- `go test ./...` passed.
- `scripts/run_doc_tests.sh` passed.

## Open Risks / Questions
- Implementation uses the system `sqlite3` binary through `os/exec` rather than an embedded Go SQLite driver. This satisfies SQLite behavior but introduces runtime dependency on `sqlite3` availability.
- Acceptance criterion for retrieval benchmark parity is indirectly covered via existing `go test ./...` retrieval/evaluation tests; no new explicit benchmark fixture was added in this cycle.

## Recommended QA Focus Areas
- Validate SQLite file creation and WAL mode on first write.
- Validate first-run migration path from legacy `index.yaml` into SQLite.
- Validate command parity (`write`, `retrieve`, `evaluate`, `snapshot`) against existing behavior.
- Validate snapshot create/restore includes and restores `index.db` when present.
- Validate behavior when `sqlite3` executable is unavailable or returns command errors.

## New Gaps Discovered
- Potential follow-up: replace CLI-based SQLite execution with embedded Go driver (`modernc.org/sqlite`) to remove runtime binary dependency and improve portability guarantees.
