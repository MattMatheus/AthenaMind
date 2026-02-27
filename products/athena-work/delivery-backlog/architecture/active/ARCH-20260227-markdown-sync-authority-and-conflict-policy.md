# Architecture Story: Define markdown sync authority and conflict policy

## Metadata
- `id`: ARCH-20260227-markdown-sync-authority-and-conflict-policy
- `owner_persona`: Software Architect - Ada.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0010, ADR-0013, WSI-ADR-0002]
- `decision_owner`: staff-personas/Software Architect - Ada.md
- `success_metric`: Sync policy is approved with deterministic conflict handling and auditable reconciliation flow.

## Decision Scope
Define how backend state and markdown artifacts stay consistent, including conflict detection, resolution priority, and operator override rules.

## Problem Statement
The current markdown-first model is mutable by many actors. Migrating to a backend authoritative model without a strict sync policy will create silent divergence and operator mistrust.

## Inputs
- ADRs: ADR-0010, ADR-0013, WSI-ADR-0002
- Architecture docs: stage-exit gates, observer cycle policy, intake templates
- Constraints: retain repo-readable artifacts; no silent overwrite of human edits

## Outputs Required
- ADR updates: state authority and sync rules
- Architecture artifacts:
  - sync direction matrix (backend -> markdown, markdown -> backend exceptions)
  - conflict taxonomy (`content_conflict`, `ordering_conflict`, `missing_artifact`, `stale_revision`)
  - reconciliation workflow with operator decision hooks
- Risk/tradeoff notes:
  - strict authority vs human flexibility
  - sync frequency vs performance

## Acceptance Criteria
1. Single-writer authority model is explicit with permitted manual-edit exceptions.
2. Conflict detection and resolution actions are deterministic and testable.
3. Drift alarms and blocking conditions are defined for stage-critical artifacts.
4. Sync/read-model design supports simple, high-clarity human planning views and low-latency agent consumption paths.

## QA Focus
Validate that sync policy preserves observer/report auditability and does not permit hidden state mutation.

## Intake Promotion Checklist (intake -> ready)
- [ ] Decision scope is explicit and bounded.
- [ ] Problem statement describes urgency and impact.
- [ ] Required inputs are listed (ADRs, architecture docs, constraints).
- [ ] Separation rule verified: architecture output, not implementation output.
- [ ] Required outputs are concrete and reviewable in QA handoff.
- [ ] Risks/tradeoffs include mitigation and owner.
