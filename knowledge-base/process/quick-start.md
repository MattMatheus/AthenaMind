# AthenaWork Quick Start

## Summary

Run AthenaWork stage loops directly from this repository.

## Minimal Stage Loop

From repo root:

```bash
./tools/launch_stage.sh engineering
./tools/launch_stage.sh qa
./tools/run_observer_cycle.sh --cycle-id <cycle-id>
```

## Practical Pairing

- AthenaWork constrains the execution workflow.
- AthenaMind provides memory context, retrieval, and evaluation.

## Related

- [AthenaWork Product Guide](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/product/athenawork.md)
- [AthenaWork Operator Reference](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/product/athenawork-operator-reference.md)
