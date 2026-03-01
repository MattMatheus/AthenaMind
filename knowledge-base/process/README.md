# Process Docs (AthenaWork)

## Summary

AthenaWork process assets are in-repo and are part of the supported operator workflow.

## Core Process Entrypoints

- `HUMANS.md`
- `DEVELOPMENT_CYCLE.md`
- `tools/launch_stage.sh`
- `tools/run_observer_cycle.sh`
- `tools/run_stage_tests.sh`

## Onboarding Token Contract

- Every newly launched stage agent must finish onboarding by generating a new UUIDv4.
- The agent must disclose it to the operator exactly as: `RESUME TOKEN: <uuid>`.
- The disclosed token is used with the park/resume guardrail in `AGENTS.md`:
  - `STOP WORK. RESUME TOKEN: <token>`
  - `RESUME WORK. RESUME TOKEN: <token>`
  - if token mismatches, remain parked and report `TOKEN MISMATCH: still parked`

## Work-System Assets

- Stage prompts: `stage-prompts/active/`
- Specialist roles: `staff-personas/`
- Delivery backlog lanes: `delivery-backlog/`
- Operating-system evidence and handoff: `operating-system/`

## Start Here

- [AthenaWork Product Guide](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/product/athenawork.md)
- [AthenaWork Operator Reference](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/product/athenawork-operator-reference.md)
- [AthenaWork Quickstart](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/getting-started/athenawork-quickstart.md)
