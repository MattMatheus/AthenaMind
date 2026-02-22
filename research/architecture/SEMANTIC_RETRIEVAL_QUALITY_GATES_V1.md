# Semantic Retrieval Quality Gates (v1)

## Purpose
Define measurable pass/fail gates for semantic retrieval behavior in the memory CLI v1 architecture.

## Metric Definitions
- Top-1 Useful Rate:
  - percentage of benchmark queries where top selected semantic result is judged useful.
- Fallback Determinism:
  - percentage of repeated runs where fallback outcome is identical for same input/corpus/config.
- Selection Mode Reporting:
  - percentage of responses that include `selection_mode`.
- Source Trace Completeness:
  - percentage of responses that include both `selected_id` and `source_path`.

## v1 Thresholds
- Top-1 Useful Rate: `>= 80%` (minimum 50 benchmark queries).
- Fallback Determinism: `100%`.
- Selection Mode Reporting: `100%`.
- Source Trace Completeness: `100%`.

## Pass/Fail Rules
- PASS:
  - all thresholds meet or exceed targets.
- FAIL:
  - any threshold misses target.
- WATCH:
  - optional internal signal for borderline tuning when top-1 useful rate is `80%` to `82%`.
  - WATCH is still PASS for v1 unless another metric fails.

## Benchmark and Harness Contract
Required harness inputs:
- fixed benchmark corpus snapshot id
- fixed query set id (minimum 50 queries)
- retrieval config snapshot (confidence/ambiguity thresholds)

Required harness outputs:
- metric summary table with numerators/denominators
- failing query list for top-1 useful misses
- deterministic replay proof for fallback runs
- artifact version references (`corpus_id`, `query_set_id`, `config_id`)

## Fallback Interaction Rule
If semantic confidence/ambiguity gate fails, engine must use deterministic fallback:
1. exact-key lookup
2. path-priority lookup

Quality evaluation must separately report semantic wins vs fallback-triggered results.

## QA Checklist
1. Verify benchmark corpus/query/config are pinned and recorded.
2. Verify threshold calculations match contract definitions.
3. Verify fallback ordering and determinism evidence.
4. Verify every response artifact includes `selection_mode`, `selected_id`, and `source_path`.

## References
- `research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`
- `research/decisions/ADR-0009-file-memory-cli-v1-architecture.md`
- `research/decisions/ADR-0012-semantic-retrieval-quality-gates-v1.md`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
