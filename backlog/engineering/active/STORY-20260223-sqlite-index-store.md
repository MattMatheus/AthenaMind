# Story: Migrate index store from JSON files to SQLite

## Metadata
- `id`: STORY-20260223-sqlite-index-store
- `owner_persona`: personas/Product Manager - Maya.md
- `status`: active
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0005, ADR-0009, ADR-0010]
- `success_metric`: index and metadata queries use SQLite WAL mode with identical CLI behavior and no regression on retrieval quality benchmarks
- `release_checkpoint`: required

## Problem Statement
- The current index is a JSON file read entirely into memory on every operation. ADR-0005 specifies SQLite WAL as the v0.1 storage default. As the memory store grows, full-file reads become a scaling bottleneck. SQLite provides queryable, concurrent-safe storage that also serves as the foundation for embedding storage in the retrieval upgrade story.

## Scope
- In:
  - Add SQLite WAL-mode database as the index and metadata store
  - Implement adapter interface so retrieval/write logic is storage-agnostic
  - Migrate index.yaml read/write to SQLite adapter
  - Keep markdown content files on disk (SQLite is for structured index/metadata only)
  - Migration path: auto-import existing index.yaml into SQLite on first run
  - Add `go-sqlite3` or `modernc.org/sqlite` dependency (pure Go preferred for portability)
- Out:
  - Embedding storage (separate story)
  - Removal of file-based fallback
  - Schema changes to memory entry model

## Acceptance Criteria
1. Index and metadata stored in SQLite with WAL mode enabled.
2. Existing CLI commands (write, retrieve, evaluate, snapshot) work identically against SQLite backend.
3. Auto-migration imports existing index.yaml into SQLite on first run.
4. `go test ./...` passes with no regressions.
5. Retrieval evaluation benchmark produces identical or better results.
6. Adapter interface allows future backend implementations.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260223-split-main-go-into-packages (strongly recommended: index logic should be in internal/index)

## Notes
- Pure Go SQLite driver (`modernc.org/sqlite`) is preferred over CGO-based `go-sqlite3` for cross-platform portability and simpler CI.
- Markdown content files stay on disk — SQLite holds only the index, metadata, and future embeddings.
