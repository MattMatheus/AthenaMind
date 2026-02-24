# AthenaWork Operator Reference

## Purpose

Provide a compact, technical reference for operators running AthenaWork from the archived pack.

## Operator Pack Root

- `/Users/foundry/Experiments/Archived/AthenaMind-internal-2026-02-24/products/athena-work`

## Key Files

- `HUMANS.md`
- `DEVELOPMENT_CYCLE.md`
- `tools/launch_stage.sh`
- `tools/run_observer_cycle.sh`
- `tools/run_stage_tests.sh`

## Stage Commands

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

## Operating Rules (Condensed)

- Stage order: planning/architect as needed, then engineering -> qa -> pm.
- Respect queue model; do not invent work when active queue is empty.
- Run observer at cycle boundary.
- Commit once per completed cycle with cycle-aligned artifact bundle.

## Recommended Pairing With AthenaMind

- Use AthenaMind CLI retrieval before stage transitions for current memory context.
- Use AthenaMind episode/snapshot features for continuity and rollback evidence.

## Related

- [AthenaWork Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork.md)
- [AthenaMind Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenamind.md)
