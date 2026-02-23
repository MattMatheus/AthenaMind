# Engineering Handoff: STORY-20260223-readme-v01-alignment

## What Changed
- Rewrote `README.md` to describe only implemented v0.1 memory-layer capabilities.
- Added concrete quick-start command examples for:
  - `go run ./cmd/memory-cli write ...`
  - `go run ./cmd/memory-cli retrieve ...`
- Added explicit roadmap pointer to `product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`.
- Added preserved long-term vision doc at `knowledge-base/product/VISION.md`.
- Updated `knowledge-base/product/README.md` to include the new vision artifact.
- Added docs regression test `tools/test_readme_v01_alignment.sh` and wired it into `tools/run_doc_tests.sh`.
- Updated `product-research/roadmap/PROGRAM_STATE_BOARD.md` queue counts to match filesystem truth during this cycle’s test gate run.

## Why It Changed
- Story requires README alignment to current v0.1 behavior so agents operate from accurate implementation context.
- Long-term vision needed to be preserved outside README to avoid mixing future direction with current delivery scope.
- Added test coverage to prevent README drift back toward unimplemented claims.

## Test Updates Made
- New: `tools/test_readme_v01_alignment.sh`
  - validates README contains v0.1 scope + write/retrieve examples
  - validates README links phased plan and preserved vision doc
  - validates README excludes unimplemented-term claims
- Updated runner: `tools/run_doc_tests.sh` includes the new test script.

## Test Run Results
- `tools/run_doc_tests.sh` passed.
- `go test ./...` passed.

## Open Risks / Questions
- `knowledge-base/product/VISION.md` is preserved from current product-vision sources and summarized for long-term direction; if exact historical README wording is required verbatim, a follow-up archival artifact may be needed.
- `tools/run_doc_tests.sh` still emits non-blocking warning during launcher tests when `memory-cli` is not installed on PATH (`memory bootstrap skipped; 'memory-cli' not found`).

## Recommended QA Focus Areas
- Verify README acceptance criteria directly:
  - implemented v0.1 scope only
  - no unimplemented-term claims in README
  - quick-start write/retrieve examples present
  - phased-plan link present
  - vision content preserved in `knowledge-base/product/VISION.md`
- Verify new docs test fails appropriately if forbidden terms are reintroduced into README.

## New Gaps Discovered
- None.
