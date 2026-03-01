# Story: Release launch authorization workbench v1

## Metadata
- `id`: STORY-20260227-release-launch-authorization-workbench-v1
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Operator can generate and review a launch authorization package in <=2 minutes before any `dev -> prod` action.
- `release_checkpoint`: required

## Problem Statement
- Launch controls exist, but operator experience for final authorization is split across docs, env variables, and CI output.

## Scope
- In:
  - workspace UI panel to summarize launch prerequisites (tests, queue state, security-gate signals)
  - machine-readable launch manifest artifact generated from current repo state
  - explicit authorization package preview for operator signoff workflow
  - local validation command for launch package completeness
- Out:
  - automatic `dev -> prod` execution
  - cloud secret management integration

## Acceptance Criteria
1. Operator can generate a launch authorization package from the workspace without manual file stitching.
2. Package includes commit digest, test gate summary, and required confirmation markers.
3. Missing prerequisites are shown as explicit blockers with actionable remediation text.
4. Output aligns with low-vision and high-variability-attention interaction profile requirements.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-ui-zip-artifact-ingest-workbench-v1
- STORY-20260227-theme-system-planning-and-implementation-v1

## Notes
- Keep operator control explicit; do not bypass human launch approval gates.
