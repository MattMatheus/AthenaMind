# Argo Multi-Agent Control Plane Research

## Objective
- Define a deployment and delegation pattern that moves AthenaWork execution off laptops into a visible, governed environment.

## Evaluation Criteria
- Governance fit: stage gates, backlog semantics, and auditability.
- Operational control: kill switch, fan-out bounds, and failure containment.
- Delivery speed: effort to get first production-grade cycle running.
- Drift risk: runtime dependency and toolchain consistency.

## Option Matrix
| Option | Governance Fit | Operational Control | Delivery Speed | Drift Risk | Notes |
|---|---|---|---|---|---|
| Parent -> Child Argo Workflows | High | High | Medium | Medium | Strong visibility and policy control; preferred baseline. |
| Direct Pod Launch by Agents | Low | Medium | High | High | Fast to start but weak lineage/audit and harder to govern. |
| Hybrid (Workflows + Restricted Direct Pods) | Medium-High | High | Medium | Medium-High | Works if direct pods are confined to explicit experimental lanes. |

## Recommended Baseline
- Use parent -> child Argo Workflows as the default delegation model.
- Keep direct pod launch as an explicit experiment-only escape hatch with strict policy boundaries.
- Preserve existing AthenaWork scripts initially and wrap them in workflow steps.

## Minimum Guardrails
- Namespace isolation:
  - governed lane namespace for production workflow runs
  - experiment namespace for temporary direct pod runs
- RBAC:
  - agent service accounts can submit specific workflow templates
  - no create pod permission in governed namespace
- Runtime limits:
  - workflow TTL
  - semaphore/mutex fan-out cap
  - max retry count and timeout defaults
- Observability:
  - mandatory workflow labels (`cycle_id`, `story_id`, `parent_run_id`, `worker_id`)
  - artifact/log retention policy

## Linked Artifacts
- Architecture story: `delivery-backlog/architecture/intake/ARCH-20260225-argo-multi-agent-control-plane-research.md`
- ADR draft: `operating-system/decisions/WSI-ADR-0002-argo-multi-agent-delegation-model.md`
- Experiment plan: `operating-system/experiments/EXP-20260225-argo-control-plane-pilot.md`
- Transfer contract: `knowledge-base/architecture/AGENT_STATE_TRANSFER_AND_OBSERVABILITY.md`
