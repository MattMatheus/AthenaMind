# Architecture Handoff: ARCH-20260222-semantic-retrieval-quality-gates

## Decision(s) Made
- Accepted quantitative semantic retrieval gates for v1, including top-1 usefulness threshold and deterministic fallback requirements.
- Defined a reproducible benchmark/harness contract with minimum query volume and pinned corpus/config references.
- Made fallback interaction explicit in quality evaluation and QA pass/fail logic.

## Alternatives Considered
- Keep qualitative "usually useful" quality language only.
  - Rejected because QA pass/fail remains subjective.
- Set initial threshold at `>= 90%` top-1 usefulness.
  - Rejected due to high delivery risk for v1 baseline.

## Risks and Mitigations
- Risk: benchmark curation quality may bias results.
  - Mitigation: require pinned benchmark ids and explicit failing-query outputs.
- Risk: stricter gates delay implementation completion.
  - Mitigation: thresholds chosen to balance measurable quality and v1 schedule feasibility.

## Updated Artifacts
- `product-research/decisions/ADR-0012-semantic-retrieval-quality-gates-v1.md`
- `product-research/architecture/SEMANTIC_RETRIEVAL_QUALITY_GATES_V1.md`
- `product-research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `delivery-backlog/architecture/qa/ARCH-20260222-semantic-retrieval-quality-gates.md`

## Validation Commands and Results
- `tools/run_doc_tests.sh` -> PASS

## Open Questions for QA Focus
- Are threshold definitions and formulas precise enough for deterministic scoring across runs?
- Does the fallback-determinism evidence contract provide enough proof for pass/fail adjudication?
