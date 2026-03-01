# Story: Launch readiness summary strip v1

## Metadata
- `id`: STORY-20260227-launch-readiness-summary-strip-v1
- `owner_persona`: staff-personas/Product Designer - Clara.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: Operator can assess launch readiness state in <=30 seconds from the first screen without opening additional panels.
- `release_checkpoint`: required

## Problem Statement
- Launch authorization tooling exists, but the top-of-screen signal for immediate go/no-go state is not yet explicit.

## Scope
- In:
  - add a first-screen launch-readiness summary strip with status and blocker count
  - include explicit signals for queue readiness, confirmation readiness, and security-gate readiness
  - wire summary strip to launch package output so status is deterministic
  - keep summary readable in low-vision mode and keyboard-first workflows
- Out:
  - automatic release execution
  - external notification integrations

## Acceptance Criteria
1. Board first screen shows launch-readiness status and blocker count without opening launch panel.
2. Summary strip status is derived from latest generated launch package and updates after package generation.
3. If readiness is blocked, strip shows concise, actionable blocker text.
4. No regressions to existing docs/workbench/ingest interactions.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260227-release-launch-authorization-workbench-v1

## Notes
- Use direct language suitable for high-variability-attention operation.
