# Story: Bootstrap local control-plane runtime with Docker Compose

## Metadata
- `id`: STORY-20260227-docker-compose-local-control-plane-bootstrap
- `owner_persona`: staff-personas/SRE - Nia.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0004, WSI-ADR-0002]
- `success_metric`: `docker compose -f docker-compose.local.yml up --build` starts core services healthy in <=120s on baseline laptop.
- `release_checkpoint`: required

## Problem Statement
A shared workspace backend is not usable if local startup is fragile or slow. We need a reproducible local runtime as the baseline developer/operator path.

## Scope
- In:
  - `docker-compose.local.yml` for API, DB, UI, and optional worker services
  - health checks and dependency ordering
  - `.env.example` and quickstart runbook
- Out:
  - full transition logic
  - full UI features

## Acceptance Criteria
1. Compose stack starts API + DB + UI with health checks passing.
2. Startup, teardown, and reset commands are documented and tested.
3. Local data persists across restarts unless explicit reset command is used.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- ARCH-20260227-shared-workspace-control-plane-contract

## Notes
- Keep service names aligned with future deployment manifests.
