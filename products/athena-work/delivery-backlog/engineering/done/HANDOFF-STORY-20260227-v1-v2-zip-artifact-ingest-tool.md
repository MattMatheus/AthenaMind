# Engineering Handoff: STORY-20260227-v1-v2-zip-artifact-ingest-tool

## What Changed
- Added `products/athena-work/tools/ingest_artifact_bundle.sh`:
  - accepts `.zip` bundle input
  - validates required state object presence and structure
  - normalizes/migrates V1 payloads to schema `2`
  - writes deterministic canonical state output to workspace state file
  - emits machine-readable ingest report with accepted/rejected/migrated counts
- Added `products/athena-work/tools/test_v1_v2_zip_artifact_ingest.sh`.
- Wired test into canonical docs test runner (`tools/run_doc_tests.sh`).
- Documented operator command in `products/athena-work/HUMANS.md`.

## Why It Changed
- Enable user-provided artifact bundle ingestion for realistic board seeding and migration flows.
- Replace manual state-file edits with an auditable ingest command.

## Test Updates Made
- Added end-to-end ingest test that creates V1 bundle zip, runs tool, and verifies normalized output schema/data.

## Test Run Results
- `products/athena-work/tools/test_v1_v2_zip_artifact_ingest.sh` -> PASS
- `products/athena-work/tools/test_program_state_consistency.sh` -> PASS
- `products/athena-work/tools/validate_intake_items.sh` -> PASS

## Open Risks/Questions
- Current ingest scope targets canonical state object; multi-object bundle expansion can be added in later iterations if needed.

## Recommended QA Focus Areas
- Verify ingest behavior with malformed zip payloads and missing required fields.
- Verify board/API reflects ingested state immediately in running local environment.

## New Gaps Discovered
- none
