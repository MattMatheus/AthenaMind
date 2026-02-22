# Story: Publish Initial v0.1 Release Checkpoint Bundle

## Metadata
- `id`: STORY-20260222-release-checkpoint-bundle-v01
- `owner_persona`: Product Manager - Maya.md
- `status`: intake
- `idea_id`: PLAN-20260222-architecture-gap-to-execution
- `phase`: v0.1
- `adr_refs`: [ADR-0002, ADR-0008]
- `success_metric`: one complete release bundle exists with explicit ship/hold decision and evidence links
- `release_checkpoint`: required

## Problem Statement
- v0.1 has significant done volume but no release bundle artifact yet, which blocks a reliable `done` vs `shipped` boundary.

## Scope
- In:
  - Create first release bundle using canonical template.
  - Link QA artifacts and validation results for included scope.
  - Record explicit ship/hold decision and rollback direction.
- Out:
  - New feature implementation.
  - Release automation tooling changes.

## Acceptance Criteria
1. Release bundle artifact is created and references included stories/bugs and QA evidence.
2. Decision field is explicit (`ship` or `hold`) with rationale.
3. Program board links to the created release checkpoint evidence.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- Existing QA and handoff artifacts in `backlog/engineering/done/`

## Notes
- This is a control-plane quality gate story, not a production runtime change.
