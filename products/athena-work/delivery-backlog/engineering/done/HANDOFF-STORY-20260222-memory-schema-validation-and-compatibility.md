# Engineering Handoff: STORY-20260222-memory-schema-validation-and-compatibility

## What Changed
- Added strict index and metadata schema validation in `<repo>/cmd/memory-cli/main.go`.
- Enforced required field/type/status/path constraints for index entries and metadata.
- Added schema compatibility behavior:
  - reject unsupported major versions (`ERR_SCHEMA_MAJOR_UNSUPPORTED`),
  - reject missing/invalid versions (`ERR_SCHEMA_VERSION_INVALID`),
  - allow newer minor versions with compatibility warning (`WARN_SCHEMA_MINOR_NEWER_COMPAT`).
- Retrieval now validates metadata contracts during candidate loading and fails deterministically on violations.
- Expanded tests in `<repo>/cmd/memory-cli/main_test.go` for unsupported major, newer-minor compatibility, and invalid metadata enforcement.

## Why It Changed
- Implements ADR-0010 and architecture policy requirements for deterministic schema/version contract enforcement before write/retrieve operations proceed.

## Test Updates Made
- Updated Go tests in `<repo>/cmd/memory-cli/main_test.go` to cover schema validation success/failure paths and compatibility mode behavior.

## Test Run Results
- `go test ./...` -> PASS
- `tools/run_doc_tests.sh` -> PASS

## Open Risks/Questions
- Compatibility warning currently emits to stderr; downstream callers should preserve stderr for diagnostics.

## Recommended QA Focus Areas
- Validate index/metadata field enforcement with malformed fixtures.
- Verify deterministic errors for invalid/missing schema versions.
- Confirm newer minor versions are accepted while unsupported major versions fail.

## New Gaps Discovered During Implementation
- None.
