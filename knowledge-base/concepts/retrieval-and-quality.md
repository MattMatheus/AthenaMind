# Retrieval And Quality

## Summary
Retrieval uses semantic scoring first, then deterministic fallbacks when confidence is insufficient.

## Selection Pipeline
1. Load approved candidates from index + metadata.
2. Compute semantic score.
   - embedding score is used only when query/doc vectors are model+dimension compatible
   - fresher/session-matching embeddings receive a small deterministic freshness bonus
3. If confidence gate passes, return semantic result.
4. Otherwise fallback in fixed order:
   - exact key match on `id`
   - path-priority deterministic order

## Output Contract
Every retrieve response includes:
- `selected_id`
- `selection_mode`
- `source_path`
- `confidence`
- `reason`
- `fallback_used`
- `semantic_hit`
- `precision_hint`

## Observability
Retrieve telemetry also tracks rolling quality proxies:
- `semantic_hit_rate`
- `fallback_rate`
- `precision_proxy`

## Quality Gates (Evaluate)
- Minimum query set size: 50
- `top1_useful_rate >= 0.80`
- `fallback_determinism == 1.0`
- `selection_mode_reporting == 1.0`
- `source_trace_completeness == 1.0`

## References
- `product-research/decisions/ADR-0012-semantic-retrieval-quality-gates-v1.md`
- `cmd/memory-cli/main.go`
