# Evaluation Gate Workflow

## Goal

Evaluate retrieval usefulness and safety gates on a fixed query set.

## Command

```bash
go run ./cmd/memory-cli evaluate --root memory
```

## Key Fields To Review

- `status`
- `top1_useful_rate`
- `fallback_determinism`
- `selection_mode_reporting`
- `source_trace_completeness`
- `latency_p95_ms`

## Typical Gate Expectations

- query set size >= 50
- useful rate >= 0.80
- fallback determinism == 1.0
- selection/source completeness == 1.0

## Related Docs

- [Configuration Reference](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/references/configuration.md)
- [CLI Commands](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/cli/commands.md)
