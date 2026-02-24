# AthenaWork Operator Reference

## Purpose

Operational reference for running AthenaWork directly in this repository.

## Operator Root Paths

- `HUMANS.md`
- `DEVELOPMENT_CYCLE.md`
- `tools/`
- `stage-prompts/active/`
- `delivery-backlog/`
- `operating-system/`
- `staff-personas/`

## Stage Launch Commands

```bash
./tools/launch_stage.sh planning
./tools/launch_stage.sh architect
./tools/launch_stage.sh engineering
./tools/launch_stage.sh qa
./tools/launch_stage.sh pm
./tools/launch_stage.sh cycle
```

Observer:

```bash
./tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

## Required Execution Rules

- Respect lane boundaries (`engineering` vs `architecture`).
- Run `./tools/run_stage_tests.sh` before handoff.
- Run observer after completed cycles.
- Commit exactly once per cycle (`cycle-<cycle-id>`).
- Treat `done` as QA-complete, not auto-shipped.

## Human/Non-Technical Steering Model

- Non-technical operators choose stage + goal.
- Stage prompt and specialist persona constrain agent behavior.
- Queue artifacts provide explicit acceptance criteria.
- Observer and release-handoff artifacts provide auditability.

## Related

- [AthenaWork Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork.md)
- [AthenaWork Cycle Overview](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/workflows/athenawork-cycle.md)
- [AthenaMind Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenamind.md)
