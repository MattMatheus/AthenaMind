# Architecture Story: Argo Multi-Agent Control Plane Research and Decision

## Metadata
- `id`: ARCH-20260225-argo-multi-agent-control-plane-research
- `owner_persona`: Software Architect - Ada.md
- `status`: intake
- `idea_id`: direct
- `phase`: v0.3
- `adr_refs`: [WSI-ADR-0002]
- `decision_owner`: staff-personas/Software Architect - Ada.md
- `success_metric`: Approved architecture decision package identifies one execution model and defines guardrails needed to run AthenaWork away from founder laptops.

## Decision Scope
- Decide how AthenaWork agents should delegate sub-work in Kubernetes/Argo.
- Compare and choose between:
  - Child Argo Workflows submitted by a parent workflow.
  - Direct pod launches by agents.
  - Hybrid patterns with strict policy boundaries.
- Define the minimum viable control-plane contract for multi-agent orchestration.

## Problem Statement
Current execution is laptop-centric and not team-visible. We need an execution substrate that preserves AthenaWork process guarantees while enabling shared visibility, auditability, and kill-switch control for agent fan-out.

## Inputs
- ADRs:
  - `operating-system/decisions/WSI-ADR-0001-architect-stage-launch.md`
  - `operating-system/decisions/WSI-ADR-0002-argo-multi-agent-delegation-model.md` (draft)
- Architecture docs:
  - `AGENTS.md`
  - `HUMANS.md`
  - `DEVELOPMENT_CYCLE.md`
  - `knowledge-base/process/stage-exit-gates.md`
  - `knowledge-base/process/backlog-weighting-policy.md`
  - `tools/launch_stage.sh`
  - `tools/run_stage_tests.sh`
- Constraints:
  - Repo is source of truth for future-state tooling.
  - Must support explicit operator kill control for runaway agent behavior.
  - Must preserve stage flow and backlog lane semantics.

## Outputs Required
- ADR updates:
  - New ADR selecting delegation model (parent->child workflow, direct pods, or hybrid) with rationale.
- Architecture artifacts:
  - Control-plane component diagram for parent orchestrator, worker execution units, and observability surfaces.
  - Execution contract for worker task submission (inputs, outputs, correlation IDs, status model).
  - Security and policy profile (RBAC boundaries, namespace model, quotas, retry/timeout defaults).
  - State-transfer channel matrix (canonical vs non-canonical), with per-channel enforcement point and allow/deny policy.
- Risk/tradeoff notes:
  - Operational complexity vs speed tradeoff.
  - Toolchain drift risk (base image vs runtime installs).
  - Runaway fan-out and cost containment controls.
  - Side-channel state transfer risk and mitigation.

## Acceptance Criteria
1. Architecture recommendation selects exactly one primary delegation model and documents why alternatives were rejected.
2. Output package includes concrete guardrails for kill-switch behavior, fan-out limits, and permissions for worker execution.
3. Output package defines a phased rollout path that starts with existing AthenaWork scripts unchanged and adds Argo-native controls incrementally.
4. Output package includes explicit deny-by-default treatment for non-canonical state-transfer channels in governed lanes.

## QA Focus
- Verify separation rule: architecture outputs only (no implementation coupling).
- Verify decision package is actionable by engineering without hidden assumptions.
- Verify risks include owner + mitigation + trigger.

## Intake Promotion Checklist (intake -> ready)
- [x] Decision scope is explicit and bounded.
- [x] Problem statement describes urgency and impact.
- [x] Required inputs are listed (ADRs, architecture docs, constraints).
- [x] Separation rule verified: architecture output, not implementation output.
- [x] Required outputs are concrete and reviewable in QA handoff.
- [x] Risks/tradeoffs include mitigation and owner.
