# Release Notes

## Summary
Versioned release communication for user-visible behavior, known limits, and upgrade impact.

## Intended Audience
- Users deciding whether to adopt a release.
- PM/QA/operators preparing ship or hold decisions.

## Preconditions
- Included stories and QA evidence are finalized.
- Release bundle decision context is available.

## Main Flow
1. Start from `v0.1-template.md`.
2. Publish release-specific note (`v0.1.md`).
3. Link QA evidence and known limitations.
4. Keep aligned with release bundle ship/hold decision.

## Failure Modes
- Notes claim functionality not backed by done + QA evidence.
- Known limitations are omitted, causing expectation mismatch.
- Release notes and release bundle decisions diverge.

## References
- `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- `docs/process/STAGE_EXIT_GATES.md`
- `backlog/engineering/done/QA-RESULT-STORY-*.md`
- `work-system/handoff/RELEASE_BUNDLE_v0.1-initial-2026-02-22.md`

## Pages
- `v0.1-template.md`
- `v0.1.md`
