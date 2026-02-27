# Agent State Transfer and Observability Contract

## Objective
- Define the required interface for agent-to-agent state transfer using AthenaMind as the control point.
- Define the observability contract for retrieval/transfer behavior analysis.
- Evaluate whether alternate state-transfer methods should be allowed.

## Primary Interface Decision
- Canonical transfer interface: AthenaMind Go binary in this repository.
- Requirement: all agent state/memory transfer operations pass through this interface.
- Rationale:
  - One policy choke point.
  - Uniform audit records.
  - Uniform OpenTelemetry spans for retrieval quality tuning and incident response.

## Transfer Envelope (Required Fields)
- `request_id`
- `session_id`
- `cycle_id`
- `source_agent_id`
- `target_agent_id`
- `story_id` (or `task_id`)
- `domain`
- `memory_ids` (IDs only in logs/metrics; do not log full memory payload by default)
- `policy_stage`
- `decision` (`allow|deny|redact`)
- `decision_reason`
- `risk_level`

## Enforcement Rules
- Default deny for unknown destination, unknown domain, or missing metadata.
- Transfer must fail closed if policy evaluation or audit write fails.
- Destination allowlist required for cross-agent transfer.
- Redaction required before delivery when policy determines partial disclosure.

## OpenTelemetry Span Contract
- Root span: `memory.transfer`
- Child spans:
  - `memory.retrieve`
  - `memory.policy_eval`
  - `memory.redact`
  - `memory.deliver`
- Minimum span attributes:
  - `request_id`
  - `cycle_id`
  - `source_agent_id`
  - `target_agent_id`
  - `domain`
  - `result_count`
  - `latency_ms`
  - `decision`
  - `risk_level`
  - `fallback_used`
  - `cache_hit`

## Retrieval Pattern Monitoring
- Required dashboard views:
  - p50/p95 latency by domain and agent pair.
  - decision distribution (`allow|deny|redact`) over time.
  - fallback rate and cache hit rate.
  - no-result retrieval rate.
  - policy deny reasons by category.
- Required alert starters:
  - sudden deny-rate spike per domain.
  - sudden latency increase for `memory.retrieve`.
  - high fallback rate regression.

## Alternate State Transfer Methods: Risk Evaluation

### Channels to Treat as Non-Canonical
- Direct pod-to-pod HTTP/gRPC memory transfer bypassing AthenaMind.
- Writing transferable state directly to shared object storage.
- Side-channel transfer via workflow annotations/labels.
- Ad hoc files in shared volumes not mediated by AthenaMind.

### Channel Control Matrix
| Channel | Governed Lane Policy | Enforcement Point | Audit Requirement |
|---|---|---|---|
| AthenaMind Go binary transfer API | Allow | Binary policy layer + namespace RBAC | Full envelope + OTel span chain |
| Direct pod-to-pod HTTP/gRPC | Deny by default | NetworkPolicy + service account RBAC | Deny event with source/target metadata |
| Shared volume file drop | Deny by default | Pod security policy + mount restrictions | Deny event + workload identity |
| Object storage side channel | Deny by default | IAM policy + egress policy | Deny event + bucket/path metadata |
| Workflow annotation/label payload | Deny for state payload | Admission policy + schema checks | Validation failure event |
| Explicit experiment channel | Allow only in experiment namespace | Namespace policy + TTL/quota + kill switch | Parity audit fields + parity OTel attributes |

### Risk
- Policy bypass and incomplete audit trail.
- Inconsistent redaction and data classification handling.
- Fragmented telemetry, reducing retrieval optimization quality.
- Higher incident response and compliance risk due to missing lineage.

### Policy Position
- Governed lanes: disallow non-canonical transfer methods.
- Experimental lanes: allow only if explicitly marked, time-boxed, and fully observed.
- Any experimental channel must emit equivalent audit + OTel fields and must have operator kill control.

## Recommended Rollout
1. Enforce AthenaMind-only transfer in governed Argo namespaces.
2. Add admission/policy checks blocking disallowed transfer patterns.
3. Validate telemetry completeness in stage QA before promotion.
4. Run time-boxed experiments for any alternate channel and either:
   - promote after controls parity is proven, or
   - deprecate and block permanently.

## Related Artifacts
- `knowledge-base/architecture/ARGO_MULTI_AGENT_CONTROL_PLANE_RESEARCH.md`
- `knowledge-base/architecture/ATHENAMIND_TRANSFER_PROTOCOL_V1.md`
- `operating-system/decisions/WSI-ADR-0002-argo-multi-agent-delegation-model.md`
- `operating-system/experiments/EXP-20260225-argo-control-plane-pilot.md`
