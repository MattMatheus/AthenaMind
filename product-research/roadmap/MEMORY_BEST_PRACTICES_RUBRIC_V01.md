# Memory Best Practices Rubric (v0.1)

## Purpose
Provide a living, evidence-driven rubric for deciding which memory techniques are baseline now vs experiment later for AthenaMind v0.1.

## Decision Model
Each technique is scored 1 to 5 on:
- `evidence_quality`: primary-source and/or local measured support
- `scope_fit_v01`: alignment with memory-layer-only scope (ADR-0007)
- `reversibility`: ease of rollback if results regress
- `instrumentation_cost`: effort to measure impact reliably

Weighted priority score:
- `0.40 * evidence_quality`
- `0.30 * scope_fit_v01`
- `0.20 * reversibility`
- `0.10 * instrumentation_cost`

Status labels:
- `Adopted`: baseline behavior for v0.1
- `Trial`: bounded experiment with success/fail gates
- `Deferred`: not for current phase
- `Rejected`: conflicts with v0.1 strategy

## Baseline First (Highest-Quality Techniques)
1. `Adopted` - Atomic memory extraction (no transcript stuffing)
   - Priority score: `4.7`
   - Why first: highest impact on precision, traceability, and continuity.
   - KPI linkage: memory precision, memory error rate, trace completeness.
   - Evidence required to keep: weekly precision/error trend stays green/yellow in ADR-0008 bands.

2. `Adopted` - Precision-first retrieval with small candidate sets and abstain path
   - Priority score: `4.6`
   - Why first: directly reduces wrong-memory recalls and trust failures.
   - KPI linkage: memory precision, wrong memory recall rate, trusted autonomy rate.
   - Evidence required to keep: no sustained increase in false recalls.

3. `Adopted` - Mandatory memory trace contract per memory-influenced action
   - Priority score: `4.6`
   - Why first: core trust and governance requirement, enables reliable diagnosis.
   - KPI linkage: trace completeness, trusted autonomy rate.
   - Evidence required to keep: trace completeness remains >= 99% green band target.

4. `Adopted` - Policy-gated promotion and conflict handling (local -> global)
   - Priority score: `4.5`
   - Why first: controls cross-session drift while preserving operator intent.
   - KPI linkage: preference adherence rate, review gate bypass rate, memory error rate.
   - Evidence required to keep: bypass rate remains 0 for medium/high risk flows.

5. `Adopted` - Consolidation/dedup on write path
   - Priority score: `4.3`
   - Why first: reduces contradictory memory and retrieval noise.
   - KPI linkage: rework reduction, memory precision, continuity success.
   - Evidence required to keep: reduced duplicate/conflict incidents in weekly review.

## Bounded Experiment Queue
1. `Trial` - Salience scoring with controlled decay
   - Priority score: `3.8`
   - Hypothesis: improves retrieval quality over time by demoting stale items.
   - Guardrail: do not drop below continuity and precision floors.
   - Exit gate: promote to `Adopted` only if precision/continuity improve for 2 consecutive snapshots.

2. `Trial` - Memory graph/topology-assisted retrieval
   - Priority score: `3.7`
   - Hypothesis: better multi-hop recall for complex workflows.
   - Guardrail: retrieval latency p95 and precision must not regress.
   - Exit gate: promote only with measured benefit on complex scenario subset.

3. `Trial` - Scheduled hygiene pass (staleness/pruning/link checks)
   - Priority score: `3.7`
   - Hypothesis: lowers long-run drift and context bloat.
   - Guardrail: no accidental loss of high-salience memories.
   - Exit gate: measurable reduction in stale/conflict findings with stable continuity.

## Not For v0.1 Baseline
- `Deferred` - Meta-agent skill evolution / agentic crossover loops
  - Reason: promising but high complexity and weak local evidence.
- `Rejected` - Runtime/pod/container lifecycle ownership as memory-core behavior
  - Reason: out of scope by ADR-0007.
- `Rejected` - Unverified headline performance claims as planning inputs
  - Reason: must be validated on AthenaMind workloads before adoption.

## Agent-Centric Suggestions (Practical)
- Keep retrieval outputs concise and ranked, with explicit abstain reason when confidence is low.
- Prefer stable memory schemas over rapid structure churn; drift hurts reuse quality.
- Preserve provenance on every recalled item; confidence without source is not trustworthy.
- Optimize for predictable behavior under long-running sessions, not only short demos.

## Operating Cadence
- Weekly: score techniques against KPIs and update status (`Adopted|Trial|Deferred|Rejected`).
- Monthly: re-rank trial items by measured impact and rollback cost.
- Promotion rule: no technique moves to `Adopted` without two consecutive positive snapshots and no governance regression.

## Evidence Links
- `product-research/roadmap/DOGFOODING_LOOP_AND_TELEMETRY_KPI_SET_V01.md`
- `product-research/decisions/ADR-0002-north-star-and-success-metrics.md`
- `product-research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- `product-research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`
