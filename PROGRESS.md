# Progress Log

## 2026-02-23T09:07:44Z
- Started overnight `experiment` branch work focused on memory system + work system reliability.
- Baseline: `go test ./...` passed in this clone before edits.

## 2026-02-23T09:07:44Z — Decision: Retrieval should degrade, not fail, on partial corruption
- Implemented retrieval hardening in `/internal/retrieval/retrieval.go`:
  - Candidate loading now skips invalid/corrupt entries instead of aborting the full query.
  - Retrieval warnings now aggregate causes (invalid entries skipped, embedding fallback reasons) with de-duplication.
  - If all candidates are invalid, retrieval now reports that skipped count in the error.
- Added test coverage in `/internal/retrieval/retrieval_test.go` for corrupt metadata skip behavior.

## 2026-02-23T09:07:44Z — Decision: Work tooling should remain safe by default but testable on non-`dev` branches
- Updated `/products/athena-work/tools/launch_stage.sh`:
  - Branch guard now uses `ATHENA_REQUIRED_BRANCH` with default `dev`.
- Updated launch-stage integration tests to pin required branch explicitly so they work on isolated local branches:
  - `/products/athena-work/tools/test_launch_stage_readme_queue.sh`
  - `/products/athena-work/tools/test_launch_stage_memory_integration.sh`

## 2026-02-23T09:07:44Z — Decision: Observer episode write-back should infer stage context
- Updated `/products/athena-work/tools/run_observer_cycle.sh`:
  - Added policy-stage inference (`planning`, `architect`, fallback `pm`) from cycle/story context.
  - Episode write-back now uses inferred stage instead of hardcoded `pm`.
- Extended `/products/athena-work/tools/test_observer_cycle_memory_integration.sh` to assert stage inference behavior.

## 2026-02-23T09:07:44Z — Decision: Intake IDs must be globally unique across active intake lanes
- Updated `/products/athena-work/tools/validate_intake_items.sh`:
  - Added duplicate `id` detection across engineering + architecture intake files.
- Added regression test `/products/athena-work/tools/test_validate_intake_duplicate_ids.sh`.
- Wired new test into `/products/athena-work/tools/run_doc_tests.sh`.

## 2026-02-23T09:07:44Z — Operator input captured
- Podman is available for optional containerized storage experiments.
- Decision for now: defer Podman-based storage changes until current no-install hardening set is fully validated.

## 2026-02-23T09:07:44Z — Decision: Doc harness tests must be lane-absence tolerant
- During verification, `tools/run_doc_tests.sh` surfaced a robustness bug:
  - `test_program_state_consistency.sh` assumed `delivery-backlog/engineering/qa` always existed and hard-failed with `find` errors when absent.
- Fixed by adding directory-aware counting helper:
  - Missing lane directories now count as `0` instead of failing the entire test run.
- File updated: `/products/athena-work/tools/test_program_state_consistency.sh`.

## 2026-02-23T09:11:53Z — Decision: Corrupt embedding rows must not break global retrieval
- Updated `/internal/index/sqlite_store.go`:
  - `GetEmbeddingRecords` now skips malformed `vector_json` rows instead of returning an error for the full read.
- Added regression test `/internal/index/index_test.go`:
  - `TestGetEmbeddingRecordsSkipsCorruptVectors` seeds good+bad rows and verifies only valid vectors are returned.
- Validation:
  - `go test ./internal/index` passed.
  - `go test ./...` passed after the change.
