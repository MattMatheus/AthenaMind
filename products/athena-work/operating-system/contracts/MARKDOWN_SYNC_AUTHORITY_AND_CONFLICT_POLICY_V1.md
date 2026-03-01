# Markdown Sync Authority and Conflict Policy v1

## Scope
Defines authority boundaries, sync direction matrix, conflict taxonomy, and reconciliation workflow for backend-authoritative state projected into markdown artifacts.

## Single-Writer Authority Model
- Canonical writer: backend control plane.
- Derived artifacts: markdown lane files and queue readmes.
- Permitted manual-edit exceptions:
  - `operator_override_note` (human-authored annotation only)
  - `research_comm_exception` documentation in research lane artifacts
- Non-permitted behavior:
  - direct markdown mutation that attempts to alter canonical state without transition API.

## Sync Direction Matrix
| Source | Target | Allowed | Conditions | Result |
|---|---|---|---|---|
| backend state | markdown lanes | yes | successful transition commit | projected artifact update |
| backend state | queue README ordering | yes | promotion/reorder transition accepted | deterministic ordered list update |
| markdown lanes | backend state | no (default) | none | reject with policy violation |
| markdown lanes | backend state | conditional | explicit operator override transition | reconciliation review workflow |

## Conflict Taxonomy
1. `content_conflict`
- Same artifact version receives divergent content across competing write intents.
- Resolution: authoritative backend snapshot wins; emit conflict event and preserve losing content in audit payload.

2. `ordering_conflict`
- Queue ordering differs between derived markdown and authoritative ranking.
- Resolution: authoritative order projection; emit reorder event.

3. `missing_artifact`
- Required derived artifact is absent after transition commit.
- Resolution: regenerate artifact; if regeneration fails, block stage transition completion.

4. `stale_revision`
- Transition request references outdated artifact/state version.
- Resolution: reject request; require caller refresh and retry.

## Reconciliation Workflow
1. Detect drift during transition post-commit projection or periodic drift check.
2. Classify conflict using taxonomy.
3. Apply deterministic auto-resolution where policy allows.
4. If exception path required, open operator decision hook:
- capture `operator_id`
- capture `decision_reason`
- capture `approval_artifact_path`
5. Emit immutable reconciliation events with required correlation fields.
6. Reproject markdown artifacts and verify zero unresolved critical drift.

## Drift Alarms and Blocking Conditions
Stage-critical artifacts:
- active/qa lane story files
- active queue README ordering
- observer report references tied to cycle closure

Alarm policy:
- warn for non-critical drift, continue with reconciliation queue.
- block for stage-critical drift until resolved.

Blocking error codes:
- `ERR_SYNC_DRIFT_CRITICAL`
- `ERR_SYNC_OVERRIDE_APPROVAL_REQUIRED`
- `ERR_SYNC_RECONCILE_FAILED`

## Human and Agent Experience Targets
- Human planning views:
  - high-clarity markdown projections with stable ordering.
  - low-vision-friendly defaults preserved in UI-facing acceptance criteria.
- Agent consumption paths:
  - low-latency authoritative reads.
  - deterministic errors and retry-safe semantics.

## Audit and Observer Linkage
Every sync/reconciliation event must include:
- `cycle_id`
- `story_id`
- `session_id`
- `trace_id`
- `conflict_type` (when applicable)
- `resolution_action`
- `operator_override` (boolean)
