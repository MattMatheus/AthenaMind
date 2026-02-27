# WSI-ADR-0002: Argo Multi-Agent Delegation Model

## Status
Proposed

## Context
- AthenaWork currently executes from operator laptops, which limits shared visibility and standardized runtime behavior.
- We need multi-agent delegation for research, implementation, and validation work while preserving stage controls and kill-switch authority.
- Three delegation options are under consideration:
  - Parent agent submits child Argo Workflows.
  - Parent agent launches Kubernetes pods directly.
  - Hybrid: workflows for governed lanes, direct pods for ad hoc tasks.

## Decision
- Primary model: parent workflow submits child Argo Workflows using `WorkflowTemplateRef`.
- Guardrail policy:
  - Direct pod creation by agent runtimes is disallowed in governed lanes.
  - Optional direct pods are allowed only in explicit experimental namespaces with TTL, low quotas, and operator kill control.
- Phase-1 execution keeps current AthenaWork scripts unchanged and wraps them as workflow steps.
- State-transfer policy:
  - Canonical agent state transfer path is AthenaMind interface only.
  - Non-canonical transfer channels are denied by default in governed lanes.
  - Any exception requires explicit experiment scope, equivalent telemetry/audit fields, and operator kill control.

## Consequences
- Positive:
  - End-to-end execution lineage, retries, logs, and artifacts are visible in a single control plane.
  - Standard policy controls (RBAC, quota, timeout, retry) apply uniformly to delegated work.
  - Parent-child correlation ids make multi-agent traceability practical for QA and audits.
- Negative:
  - Additional operational complexity in workflow template lifecycle and cluster policy management.
  - Requires image and dependency governance to avoid drift between local and cluster runs.

## Validation Plan
- Run `operating-system/experiments/EXP-20260225-argo-control-plane-pilot.md`.
- Produce architecture output package with:
  - control-plane component diagram
  - worker execution contract (input/output/status/correlation id)
  - security profile (RBAC, namespace model, quotas, timeout/retry defaults)
  - non-canonical state-transfer channel matrix with allow/deny policy and enforcement point
- Acceptance checks:
  - at least one full engineering -> QA cycle executes in Argo
  - parent workflow delegates at least one worker task as child workflow
  - operator kill path can terminate runaway execution quickly
  - at least one attempted non-canonical transfer path is blocked and logged
