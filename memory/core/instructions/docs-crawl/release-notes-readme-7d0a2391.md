# Release Notes Readme 7d0a2391

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
- `operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- `knowledge-base/process/stage-exit-gates.md`
- `delivery-backlog/engineering/done/QA-RESULT-STORY-*.md`
- `operating-system/handoff/release-bundle-v0.1-initial-2026-02-22.md`

## Pages
- `v0.1-template.md`
- `v0.1.md`
