# Engineering Handoff: STORY-20260222-memory-mutation-review-enforcement

## What Changed
- Implemented mutation workflow enforcement in `<repo>/cmd/memory-cli/main.go`:
  - stage guard (`planning|architect|pm` only),
  - mandatory decision evidence (`--reason`, `--risk`, `--notes`),
  - explicit approval/rejection decisions,
  - autonomous-run write block,
  - rejection gate that blocks apply until rework/re-review evidence exists.
- Added audit artifact generation under `/memory/audits/` for approved and rejected mutation decisions.
- Expanded tests in `<repo>/cmd/memory-cli/main_test.go` for missing evidence failures and rejected-decision blocked apply behavior.

## Why It Changed
- Implements ADR-0011 mutation review contract so writes are allowed only in approved stage contexts with deterministic review evidence enforcement.

## Test Updates Made
- Added and updated tests in `<repo>/cmd/memory-cli/main_test.go` covering approval/rejection and evidence failure paths.

## Test Run Results
- `go test ./...` -> PASS
- `tools/run_doc_tests.sh` -> PASS

## Open Risks/Questions
- Audit artifacts are local file outputs; future orchestration may require centralized aggregation, but v1 contract is satisfied for file-backed workflow.

## Recommended QA Focus Areas
- Verify rejected decisions always block mutation apply.
- Validate required evidence fields for both approval and rejection flows.
- Confirm stage-guard behavior rejects writes outside allowed stages.

## New Gaps Discovered During Implementation
- None.
