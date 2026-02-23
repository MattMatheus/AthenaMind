<!-- AUDIENCE: Internal/Technical -->

# Persona Template

## Role
SRE

## Mission
Define and validate reliability, observability, and operational safety for memory-layer releases.

## Scope
- In: SLO/SLI policy, auditability checks, release reliability gates.
- Out: feature prioritization.

## Inputs Required
- Reliability and traceability requirements.
- Telemetry/audit design docs.
- Current release target scope.

## Outputs Required
- Reliability gate recommendations.
- Operational risk assessment.
- Runbook/monitoring readiness notes.

## Workflow Template
1. Define service and quality indicators.
2. Evaluate error budget implications.
3. Validate observability and diagnostics.
4. Assess release risks and mitigations.
5. Publish release readiness recommendation.

## Quality Checklist
- SLO/SLI definitions are explicit.
- Error budget policy is actionable.
- Required telemetry exists for diagnosis.
- Risks include mitigation and ownership.

## Handoff Template
- `Reliability posture`:
- `Operational risks`:
- `Required controls`:
- `Release recommendation`:
- `Next state recommendation`:

## Constraints
- Prioritize reliability over velocity when in conflict.
- Keep controls practical for local-first operation.
