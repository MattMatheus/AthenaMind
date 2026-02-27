# Architecture Story: Define shared workspace control-plane contract (local-first)

## Metadata
- `id`: ARCH-20260227-shared-workspace-control-plane-contract
- `owner_persona`: Software Architect - Ada.md
- `status`: intake
- `idea_id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `phase`: v0.2
- `adr_refs`: [ADR-0004, ADR-0013, WSI-ADR-0002]
- `decision_owner`: staff-personas/Software Architect - Ada.md
- `success_metric`: A reviewed contract defines state machine entities, transition API, and event model with no unresolved authority ambiguity.

## Decision Scope
Define the canonical state authority and API contract for a local-first shared workspace backend serving human and agent operators.

## Problem Statement
Current workflow state is spread across markdown lanes and scripts, creating ambiguity under concurrent usage. A visual shared UI requires a consistent backend contract with deterministic transitions.

## Inputs
- ADRs: ADR-0004, ADR-0013, WSI-ADR-0002
- Architecture docs: AthenaWork workflow docs and stage exit gates
- Constraints: local-first default via Docker Compose; backward compatibility with launcher/observer scripts

## Outputs Required
- ADR updates: control-plane authority and transition policy ADR
- Architecture artifacts:
  - state entity model (`story`, `cycle`, `transition`, `observer_event`)
  - transition API contract (`request`, `preconditions`, `failure codes`)
  - event contract for immutable audit log
- Risk/tradeoff notes:
  - consistency vs complexity
  - local portability vs deployability

## Acceptance Criteria
1. Canonical state authority is explicitly defined (backend authoritative, markdown derived) with documented exceptions.
2. Transition API and failure semantics are specified for all core transitions used by AthenaWork.
3. Event model includes required correlation fields (`cycle_id`, `story_id`, `session_id`, `trace_id`) and supports observer linkage.
4. Contract defines research-only exception path for undocumented agent-to-agent communication and hard-block behavior for non-research lanes.
5. Contract defines required human direction-confirmation checkpoint for direction-changing transitions.
6. Contract includes explicit optimization goals for agent implementation workflow (minimal transition friction, deterministic preconditions, concise machine-readable errors).

## QA Focus
Validate that all existing AthenaWork lane transitions can be represented without semantic loss.

## Intake Promotion Checklist (intake -> ready)
- [ ] Decision scope is explicit and bounded.
- [ ] Problem statement describes urgency and impact.
- [ ] Required inputs are listed (ADRs, architecture docs, constraints).
- [ ] Separation rule verified: architecture output, not implementation output.
- [ ] Required outputs are concrete and reviewable in QA handoff.
- [ ] Risks/tradeoffs include mitigation and owner.
