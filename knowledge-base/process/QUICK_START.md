# AthenaWork Quick Start (Slim)

## Summary

Run AthenaWork stage loops from the archived operator pack while using AthenaMind in this repo for runtime memory behavior.

## Operator Pack Location

- `/Users/foundry/Experiments/Archived/AthenaMind-internal-2026-02-24/products/athena-work`

## Minimal Stage Loop

From the archived AthenaWork root:

```bash
./tools/launch_stage.sh engineering
./tools/launch_stage.sh qa
./tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

## Practical Pairing

- AthenaWork: governs the delivery loop.
- AthenaMind: powers memory operations and quality checks.

## Related

- [AthenaWork Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork.md)
- [AthenaWork Operator Reference](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork-operator-reference.md)
