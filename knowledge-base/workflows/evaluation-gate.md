# Evaluation Gate Workflow

## Summary
Run `evaluate` to verify semantic usefulness and deterministic safety gates.

## Preconditions
- Memory corpus populated with representative entries.
- Query set available (default: `cmd/memory-cli/testdata/eval-query-set-v1.json`).

## Steps
1. Run:
```bash
go run ./cmd/memory-cli evaluate --root memory
```
2. Inspect JSON output:
   - `status` (`PASS`, `WATCH`, or `FAIL`)
   - `top1_useful_rate`
   - `fallback_determinism`
   - `selection_mode_reporting`
   - `source_trace_completeness`
3. Investigate `failing_queries` if status is not `PASS`.

## Pass Criteria
- At least 50 evaluated queries.
- Useful rate >= 80%.
- Determinism and trace metrics at 100%.

## References
- `product-research/decisions/ADR-0012-semantic-retrieval-quality-gates-v1.md`
- `cmd/memory-cli/main.go`
