# AthenaWork Product Guide

## What AthenaWork Is

AthenaWork is the operator-facing workflow system that constrains how agents execute work around AthenaMind.

It is designed so non-technical users can steer agents through explicit stages, queue policies, and evidence checkpoints.

## Why It Matters

- Provides stage guardrails (`planning`, `architect`, `engineering`, `qa`, `pm`, `cycle`).
- Enforces queue discipline and handoff structure.
- Captures observer evidence and cycle continuity.
- Reduces agent drift by routing execution through canonical prompts and specialist roles.

## Canonical Operator Assets

- Root operator guide: `HUMANS.md`
- Agent rules: `AGENTS.md`
- Stage prompts: `stage-prompts/active/`
- Specialist directory: `staff-personas/STAFF_DIRECTORY.md`
- Queue system: `delivery-backlog/`
- Work OS artifacts: `operating-system/`
- Stage launchers and checks: `tools/`

## Typical Use Pattern

1. Start a stage with `./tools/launch_stage.sh <stage>`.
2. Execute queue item under that stage prompt.
3. Validate with `./tools/run_stage_tests.sh`.
4. Run observer: `./tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
5. Commit once per cycle.

## Related Docs

- [AthenaWork Operator Reference](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork-operator-reference.md)
- [AthenaWork Quickstart](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/getting-started/athenawork-quickstart.md)
- [AthenaMind Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenamind.md)
