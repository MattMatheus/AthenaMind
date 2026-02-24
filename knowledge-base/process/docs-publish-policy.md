# Docs Publish Policy

## Policy
- Canonical docs source-of-truth is markdown in this repository.
- Published website docs are generated artifacts and must not be edited directly in the host.
- Documentation changes are complete only when:
  1. Repo markdown updates are reviewed and merged.
  2. `tools/run_doc_tests.sh` passes.
  3. Docs publish workflow completes successfully.

## Publish Target
- Site: `athena.teamorchestrator.com`
- Docs path: `/docs/`

## Build And Publish
- Local build command:
  - `tools/build_docs_site.sh`
- CI workflow:
  - `.github/workflows/docs-publish.yml`

## Required Ownership
- Clara owns doc content quality and coverage.
- Engineering/QA/PM must route behavior and API changes through Clara docs updates before cycle closure.

