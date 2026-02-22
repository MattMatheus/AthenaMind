# Dogfooding Loop Design and Telemetry KPI Set (v0.1)

## Purpose
Translate `research/ingest/memory-best-practices.md` into a product + engineering operating loop that is compatible with AthenaMind v0.1 memory-layer scope.

## Inputs
- `research/ingest/memory-best-practices.md`
- `research/decisions/ADR-0002-north-star-and-success-metrics.md`
- `research/decisions/ADR-0006-governance-and-hitl-policy.md`
- `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- `research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`

## Scope Filter Applied
- Included:
  - Memory write quality (ingestion, consolidation, conflict handling)
  - Memory retrieval quality (precision, recall usefulness, explainability)
  - Memory maintenance quality (staleness control, link integrity, pruning)
- Excluded:
  - Runtime execution orchestration ownership
  - Pod/container lifecycle ownership

## Product + Engineering Meeting Outcome
This cycle establishes an evidence loop where each dogfooding session creates measurable memory outcomes and feeds weekly policy/implementation adjustments.

## Dogfooding Loop (Weekly Cadence)
1. Define scenarios
   - Select 3 to 5 recurring operator tasks requiring cross-session continuity.
   - Tag each scenario by memory type: `procedural`, `state`, `semantic`.
2. Run instrumented sessions
   - Execute scenarios in review-first mode.
   - Capture retrieval traces, memory writes, policy decisions, and operator interventions.
3. Score against KPI set
   - Compute KPI values from session telemetry.
   - Compare against ADR-0008 target bands.
4. Triage and decide
   - Label failures as `policy`, `retrieval`, `memory-quality`, or `operator-ux`.
   - Approve one reversible change batch for the next cycle.
5. Apply and verify
   - Ship small bounded changes.
   - Re-run at least one prior failing scenario to verify improvement.

## KPI Set (v0.1)
### Continuity KPIs
- `resume_success_rate`
  - Definition: percent of resumed sessions that complete without re-collecting already known preferences/context.
  - Why: direct proxy for continuity value.
- `preference_adherence_rate`
  - Definition: percent of actions aligned to stored operator/repo preferences.
  - Why: validates long-term behavior adherence.

### Retrieval Quality KPIs
- `precision_at_3`
  - Definition: proportion of top-3 retrieved memories judged relevant by operator.
  - Why: favors precision over volume per vision constraints.
- `wrong_memory_recall_rate`
  - Definition: percent of sessions with at least one materially incorrect recall used in decision flow.
  - Why: primary trust risk indicator.

### Governance and Explainability KPIs
- `trace_completeness_rate`
  - Definition: percent of memory-influenced actions with full trace fields: loaded, recalled, changed, why.
  - Why: aligns with auditability requirement.
- `review_gate_bypass_rate`
  - Definition: percent of medium/high-risk actions that bypass required review.
  - Why: governance integrity metric; expected to stay at 0.

### Efficiency/Cost KPIs
- `memory_ops_cost_per_session`
  - Definition: normalized memory operation count/cost per completed session.
  - Why: confirms strict budget policy behavior.
- `p95_retrieval_latency_ms`
  - Definition: p95 latency for retrieval requests used in active sessions.
  - Why: keeps continuity gains from regressing usability.

## Telemetry Event Contract (Minimal)
Each memory-relevant session event should include:
- `session_id`
- `scenario_id`
- `memory_type` (`procedural|state|semantic`)
- `operation` (`retrieve|write|promote|reject|review`)
- `result` (`success|fail`)
- `policy_gate` (`none|low|medium|high`)
- `trace_id`
- `latency_ms`
- `operator_verdict` (`correct|partially_correct|incorrect`)

## 80/20 Implementation Guidance
- Start with manual scoring for `operator_verdict` before automation.
- Instrument only KPI-required fields first; avoid broad telemetry expansion.
- Treat retrieval precision and trace completeness as first failure gates.
- Keep all changes reversible and tied to one KPI delta hypothesis.

## Assumptions and Verification Needs
- Assumption: memory-best-practices claims of large token/latency gains are directional, not yet validated for AthenaMind.
- Required verification in next cycles:
  - Benchmark continuity and precision changes on AthenaMind scenarios.
  - Validate any performance claims with local measurements before roadmap commitment.

## Immediate Next Actions
1. Apply baseline priorities from `research/roadmap/MEMORY_BEST_PRACTICES_RUBRIC_V01.md`.
2. Define the initial 3 to 5 scenarios and owners.
3. Add telemetry fields to current memory CLI/session traces.
4. Publish first KPI snapshot using `work-system/metrics/KPI_SNAPSHOT_TEMPLATE.md`.
