# Story: Split main.go into internal packages

## Metadata
- `id`: STORY-20260223-split-main-go-into-packages
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: qa
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0003, ADR-0009]
- `success_metric`: cmd/memory-cli/main.go contains only CLI wiring (flag parsing and subcommand dispatch); all domain logic lives in internal packages with independent unit tests
- `release_checkpoint`: deferred

## Problem Statement
- The memory CLI is a single 1,800-line main.go file. ADR-0003 defines module boundaries (memory-governance, audit-telemetry, retrieval, etc.) but the code does not reflect them. Agent sessions working on the CLI will produce better results against focused 200-300 line packages than a monolith. This is a prerequisite for rapid feature development in upcoming sprints.

## Scope
- In:
  - Extract `internal/index` (index loading, schema compat, entry CRUD, index file I/O)
  - Extract `internal/governance` (write policy enforcement, audit records, constraint checks)
  - Extract `internal/retrieval` (scoring, confidence gating, fallback chain, candidate ranking)
  - Extract `internal/telemetry` (event emission, normalization helpers)
  - Extract `internal/snapshot` (create/list/restore, manifest handling, checksums)
  - Extract `internal/gateway` (HTTP read gateway handler, API retrieve with fallback)
  - Reduce `cmd/memory-cli/main.go` to thin CLI wiring: flag parsing, subcommand dispatch, and calls into internal packages
  - Move shared types to `internal/types` or co-locate with owning package
  - All existing tests must continue to pass
- Out:
  - New features or behavior changes
  - External API changes
  - Dependency additions beyond stdlib

## Acceptance Criteria
1. `cmd/memory-cli/main.go` contains only flag parsing, subcommand dispatch, and calls to internal packages.
2. Each extracted package has at least one unit test that exercises its core function independently.
3. `go test ./...` passes with no regressions against existing test suite.
4. No behavior changes: CLI input/output contracts remain identical.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- None (foundational refactor)

## Notes
- This is a structural refactor only. Do not add features, change behavior, or add dependencies.
- Preserve all existing test scenarios; move tests to appropriate package test files where they test internal logic directly.
- Integration-level tests (full CLI flow) should remain in `cmd/memory-cli/` or a dedicated integration test file.
