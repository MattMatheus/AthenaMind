# AthenaWork Quickstart

## Summary

AthenaWork is available directly in this repository for operator-driven workflows.

## Start Here

1. Read `HUMANS.md` for operator expectations.
2. Launch planning or engineering stage.
3. Follow launcher checklist output.

## Minimal Cycle

From repo root:

```bash
./tools/launch_stage.sh engineering
./tools/launch_stage.sh qa
./tools/run_observer_cycle.sh --cycle-id <story-id>
```

Then commit once for the cycle:

```bash
git commit -m "cycle-<cycle-id>"
```

## Non-Technical Operator Pattern

1. Pick stage goal in plain language.
2. Select or confirm specialist role from `staff-personas/STAFF_DIRECTORY.md`.
3. Ask agent to execute only the stage checklist.
4. Require observer report before accepting completion.

## Related Docs

- [AthenaWork Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork.md)
- [AthenaWork Operator Reference](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork-operator-reference.md)
