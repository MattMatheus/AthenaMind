# Experiment: Argo Control Plane Pilot for AthenaWork

## Hypothesis
- Running AthenaWork stage execution through Argo Workflows (while keeping existing shell scripts and backlog semantics) will improve team visibility and operational control without degrading cycle throughput.

## Scope
- In:
  - Containerize current stage commands (`tools/launch_stage.sh`, `tools/run_stage_tests.sh`, `tools/run_observer_cycle.sh`).
  - Execute one full engineering -> QA cycle in Argo.
  - Capture logs/artifacts and validate kill-switch behavior.
- Out:
  - Replacing backlog storage model.
  - Rewriting AthenaWork scripts into a new orchestration DSL.
  - Org-wide build-tool standardization across all repositories.

## Success Criteria
- At least one full stage cycle is executed in Argo and visible in shared UI.
- Parent workflow can submit and monitor at least one child worker workflow with correlation id.
- Operator can terminate runaway workflow execution within 60 seconds.
- No stage-gate regressions compared to baseline local flow.
- At least one non-canonical state transfer attempt is denied and captured in audit/telemetry.

## Duration
- 5 working days from pilot start.

## Rollback Trigger
- More than one P1 incident related to workflow orchestration reliability or security policy violations.
- Any successful non-canonical state transfer in governed namespace without explicit experiment approval.

## Result
- status: pending
- notes:
  - Linked architecture story: `delivery-backlog/architecture/intake/ARCH-20260225-argo-multi-agent-control-plane-research.md`
  - Linked ADR draft: `operating-system/decisions/WSI-ADR-0002-argo-multi-agent-delegation-model.md`
