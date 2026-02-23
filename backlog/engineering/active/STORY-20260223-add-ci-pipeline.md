# Story: Add GitHub Actions CI pipeline

## Metadata
- `id`: STORY-20260223-add-ci-pipeline
- `owner_persona`: personas/Product Manager - Maya.md
- `status`: active
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004]
- `success_metric`: every push to dev branch runs go test and doc tests with pass/fail status visible on the repository
- `release_checkpoint`: deferred

## Problem Statement
- All verification is currently manual shell scripts. The stage exit gates and quality checks are honor-system enforced. A CI pipeline makes enforcement automatic and prevents regressions from reaching the dev branch undetected. This is especially important in an agent-operated codebase where sessions are stateless.

## Scope
- In:
  - Create `.github/workflows/ci.yml`
  - Run `go test ./...` on push to `dev` branch and on pull requests
  - Run `scripts/run_doc_tests.sh` as a separate CI job
  - Fail the pipeline if either step fails
- Out:
  - Deployment pipelines
  - Docker image builds
  - Performance or load testing
  - Changes to existing test scripts

## Acceptance Criteria
1. `.github/workflows/ci.yml` exists and is valid GitHub Actions syntax.
2. Pipeline triggers on push to `dev` and on pull requests targeting `dev`.
3. Pipeline runs `go test ./...` and reports pass/fail.
4. Pipeline runs `scripts/run_doc_tests.sh` and reports pass/fail.
5. Pipeline uses Go 1.22 to match `go.mod`.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- None

## Notes
- Keep the pipeline minimal. Do not add linters, formatters, or other tooling unless explicitly requested in a future story.
- The existing `azure-pipelines.yml` should remain untouched; it may serve a different purpose.
