# Story: Establish User-Facing Release Notes and Known Gaps Process

## Metadata
- `id`: STORY-20260223-user-docs-release-notes-known-gaps-v01
- `owner_persona`: Technical Writer - Clara.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0002, ADR-0008]
- `success_metric`: every release checkpoint has user-facing release notes with explicit known limitations aligned to release bundle decision
- `release_checkpoint`: required

## Problem Statement
- Release communication currently exists as internal bundle artifacts but lacks a dedicated user-facing release notes and known gaps output.

## Scope
- In:
  - Define release notes structure under `knowledge-base/release-notes/`.
  - Include known limitations and mitigation guidance.
  - Align release notes with release bundle ship/hold decision and QA evidence.
- Out:
  - Release decision authority changes.
  - Product scope expansion beyond v0.1.

## Acceptance Criteria
1. Release notes template includes highlights, fixes, limitations, and evidence links.
2. Known limitations are explicit and consistent with QA and bundle artifacts.
3. Release notes are referenced from docs index/navigation.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- `knowledge-base/process/STAGE_EXIT_GATES.md`
- `product-research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`

## Notes
- Preserve clear separation between internal release gate evidence and external user communication.
