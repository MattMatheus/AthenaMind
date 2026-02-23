# Release Notes V0 1 0a089890

# Release Notes: v0.1

## Summary
v0.1 delivers the AthenaMind memory-layer baseline: governed writes, deterministic retrieval fallback, and quality-evaluation workflows.

## User-Visible Additions
- `memory-cli write` and `memory-cli retrieve` baseline.
- Retrieval response contract with `selected_id`, `selection_mode`, `source_path`.
- `memory-cli evaluate` quality report workflow.
- Snapshot create/list/restore lifecycle.
- Read gateway and `api-retrieve` parity/fallback behavior.

## Fixes and Hardening
- Deterministic fallback behavior for low-confidence retrieval cases.
- Mutation policy enforcement with mandatory review evidence.
- Schema/version checks and compatibility-mode behavior for minor upgrades.

## Known Limitations
- Scope is memory-layer only (no runtime orchestration ownership).
- Snapshot scope is `full` only in current implementation.
- Shipped decision remains controlled by release checkpoint bundle policy.

## Evidence
- `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-memory-cli-v1-write-retrieve-baseline.md`
- `delivery-backlog/engineering/done/QA-RESULT-STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
- `operating-system/handoff/RELEASE_BUNDLE_v0.1-initial-2026-02-22.md`

## Support
- Troubleshooting: `knowledge-base/troubleshooting/common-errors.md`
- FAQ: `knowledge-base/faq/common-questions.md`
