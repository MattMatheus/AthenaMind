# AthenaWork Cycle Overview

## Summary

AthenaWork runs a governed delivery cycle to keep agent behavior constrained and auditable.

## Cycle Steps

1. Planning/PM intake defines and ranks work.
2. Engineering executes top active story.
3. QA validates acceptance criteria and regressions.
4. Observer records cycle report and continuity state.
5. PM reprioritizes queue for next cycle.

## Commands

```bash
./tools/launch_stage.sh engineering
./tools/launch_stage.sh qa
./tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

For continuous loop:

```bash
./tools/launch_stage.sh cycle
```

## Artifacts You Should Expect

- story/bug docs in `delivery-backlog/engineering/`
- architecture artifacts in `delivery-backlog/architecture/`
- observer report in `operating-system/observer/`
- release decision artifacts in `operating-system/handoff/`
