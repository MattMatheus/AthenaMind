# Dogfood Scenario Pack v0.1

## Purpose
Define a repeatable weekly scenario set for evaluating AthenaMind memory behavior across `procedural`, `state`, and `semantic` memory types.

## Version
- `pack_version`: `v0.1`
- `effective_date`: `2026-02-22`
- `owner`: `QA Engineer - Iris.md`

## Scoring Loop (Repeatable)
1. Run each scenario in review-first mode.
2. Capture minimal telemetry contract fields per run:
   - `session_id`, `scenario_id`, `memory_type`, `operation`, `result`, `policy_gate`, `trace_id`, `latency_ms`, `operator_verdict`
3. Score outcomes against KPI intent:
   - `resume_success_rate`, `preference_adherence_rate`, `precision_at_3`, `wrong_memory_recall_rate`, `trace_completeness_rate`, `review_gate_bypass_rate`, `p95_retrieval_latency_ms`
4. Classify weak signals/failures as:
   - `policy`, `retrieval`, `memory-quality`, `operator-ux`
5. Produce one prioritized follow-on action before closing the run.

## Scenario Set

### SCN-PROC-01: Reapply Repository Guardrails
- `memory_type`: `procedural`
- Objective:
  - Validate that previously learned workflow constraints (stage order, branch discipline, no mid-cycle commits) are recalled and followed without re-prompting.
- Setup:
  - Start a new cycle with at least one active engineering story.
- Steps:
  1. Request a standard stage launch and handoff path.
  2. Observe whether stored process rules are applied during transitions.
  3. Record any missed guardrail references.
- Success gate:
  - `operator_verdict` should be `correct` and no policy bypass recorded.

### SCN-STATE-01: Resume Half-Completed Cycle State
- `memory_type`: `state`
- Objective:
  - Validate recovery of cycle context (current story lane, pending QA requirements, observer/commit obligations) after interruption.
- Setup:
  - Pause after engineering implementation and resume later.
- Steps:
  1. Resume session with only cycle id and story id.
  2. Continue from pending stage steps.
  3. Confirm no duplicate or skipped transitions.
- Success gate:
  - Resume proceeds without rebuilding context from scratch.

### SCN-SEM-01: Retrieve Correct Architecture Baseline Context
- `memory_type`: `semantic`
- Objective:
  - Validate top-3 retrieval relevance for architecture/ADR guidance tied to active implementation work.
- Setup:
  - Run a story that requires architecture-to-engineering mapping.
- Steps:
  1. Ask for governing artifacts and current queue context.
  2. Compare retrieved context with current delivery-backlog/ADR truth.
  3. Mark incorrect recall events.
- Success gate:
  - `precision_at_3` meets target band intent and wrong-memory recall is zero.

### SCN-PROC-02: QA Defect Triage to Intake
- `memory_type`: `procedural`
- Objective:
  - Validate deterministic bug severity mapping (`P0-P3`) and lane transitions when QA finds defects.
- Setup:
  - Execute QA review against at least one story with known weak signals.
- Steps:
  1. Apply QA rubric severity mapping.
  2. File intake defects with complete evidence.
  3. Move story to correct lane.
- Success gate:
  - Priority mapping and transitions are consistent with rubric.
