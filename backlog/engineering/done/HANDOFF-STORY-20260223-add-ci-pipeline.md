# Engineering Handoff: STORY-20260223-add-ci-pipeline

## What Changed
- Added GitHub Actions workflow at `.github/workflows/ci.yml`.
- Configured triggers for:
  - push to `dev`
  - pull requests targeting `dev`
- Added separate CI jobs:
  - `go-test` running `go test ./...`
  - `doc-tests` running `scripts/run_doc_tests.sh`
- Set Go toolchain in both jobs to Go `1.22.x` via `actions/setup-go@v5`.
- Left existing `azure-pipelines.yml` unchanged per story notes.

## Why It Changed
- Story requires automated enforcement of quality gates so pushes and PRs to `dev` always run the core test suite and doc validation.

## Test Updates Made
- No repository test code changes were required for this story.

## Test Run Results
- `go test ./...` passed.
- `scripts/run_doc_tests.sh` passed.

## Open Risks / Questions
- CI now runs in both Azure Pipelines and GitHub Actions. If both remain enabled, duplicate checks may appear; this is not blocking for this story.

## Recommended QA Focus Areas
- Confirm workflow syntax and trigger behavior in GitHub Actions UI on:
  - push to `dev`
  - PR target `dev`
- Confirm both jobs fail independently when their respective command fails.

## New Gaps Discovered
- None.
