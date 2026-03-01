# Architecture Story: Define Kanban workspace information architecture v1

## Metadata
- `id`: ARCH-20260227-kanban-workspace-information-architecture-v1
- `owner_persona`: staff-personas/Software Architect - Ada.md
- `status`: intake
- `idea_id`: PLAN-20260227-kanban-workspace-ideation
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `decision_owner`: staff-personas/Software Architect - Ada.md
- `success_metric`: Approved architecture contract defines lane/card/docs entities and sync boundaries with zero unresolved ambiguity in v1 Kanban scope.

## Decision Scope
- Canonical information architecture for a true Kanban workspace board integrated with a docs workspace surface.

## Problem Statement
- Kanban and docs features can fragment without a stable architecture model for entities, relationships, and sync authority boundaries.

## Inputs
- ADRs:
  - `operating-system/decisions/WSI-ADR-0002-control-plane-authority-and-transition-policy.md`
  - `operating-system/decisions/WSI-ADR-0003-markdown-sync-authority-and-conflict-policy.md`
- Architecture docs:
  - `operating-system/contracts/LOCAL_FIRST_SHARED_WORKSPACE_CONTROL_PLANE_CONTRACT_V1.md`
  - `operating-system/contracts/MARKDOWN_SYNC_AUTHORITY_AND_CONFLICT_POLICY_V1.md`
- Constraints:
  - Experimental mode allowed; prioritize clarity and implementation utility.

## Outputs Required
- ADR updates:
  - Kanban workspace entity/state model update (lane/card/blocker/cycle/docs linkage).
- Architecture artifacts:
  - Contract for Kanban board read model and docs workspace read model.
  - Sync projection matrix backend -> markdown -> UI read surfaces.
- Risk/tradeoff notes:
  - Complexity vs readability in board density.
  - Sync guarantees vs iteration speed.

## Acceptance Criteria
1. Architecture defines canonical Kanban entities and relationships with deterministic state semantics.
2. Docs workspace model and Kanban linkage are explicit and implementation-ready.
3. Sync authority and drift-guard implications are mapped for all new workspace surfaces.

## QA Focus
- Validate architecture outputs are concrete enough to drive PM refinement and engineering implementation without reinterpretation.

## Intake Promotion Checklist (intake -> ready)
- [ ] Decision scope is explicit and bounded.
- [ ] Problem statement describes urgency and impact.
- [ ] Required inputs are listed (ADRs, architecture docs, constraints).
- [ ] Separation rule verified: architecture output, not implementation output.
- [ ] Required outputs are concrete and reviewable in QA handoff.
- [ ] Risks/tradeoffs include mitigation and owner.
