# Research Brief: Governance and Human-in-the-Loop Controls

## Source Inputs
- `README.md` (vision prompt)
- `research/ingest/pillars.md`
- `research/architecture/MODULE_BOUNDARIES_V01.md`
- `research/architecture/MVP_CONSTRAINTS_V01.md`
- `research/references/external-sources-2026-02-22-governance.md`

## Problem Statement
AthenaMind needs governance controls for memory promotion, conflict resolution, and autonomy safety. Current source material defines intent but not sufficient risk thresholds.

## Key Insights (80/20)
- A governance framework can be decided now, including decision points and required audit events.
- Numeric thresholds and policy strictness for autonomous actions are not specified in current inputs.
- Human review pathways should be tied to confidence, impact scope, and policy risk class.

## Governance Framework (Decidable Now)
- Decision points requiring policy evaluation:
  - Memory promotion (repo-local -> global)
  - Conflict resolution (global vs local directives)
  - Autonomous writeback after runtime tasks
  - External side-effecting actions (repo writes, API mutations)

- Required policy artifacts:
  - Risk-class matrix (low/medium/high)
  - Review trigger matrix (confidence, scope, side effects)
  - Escalation rules and override roles

- Required telemetry:
  - Policy decision id
  - Inputs considered (confidence/provenance/risk-class)
  - Outcome (allow, deny, require review)
  - Actor (human, system, role)

## Decisions Blocked by Missing Founder Inputs
- Exact confidence thresholds for automatic promotion
- Scope boundary for "high risk" actions requiring mandatory human approval
- Whether default mode should be conservative (review-first) or permissive (auto-first)
- Role model for overrides (solo founder only vs delegated operators)

## Recommended Next Action
- Collect founder decisions via a compact policy questionnaire, then finalize governance ADR.
