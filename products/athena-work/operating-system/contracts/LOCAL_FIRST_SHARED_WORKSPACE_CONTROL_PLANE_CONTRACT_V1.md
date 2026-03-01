# Local-First Shared Workspace Control-Plane Contract v1

## Scope
Defines authoritative state entities, transition API semantics, and immutable observer-linked event model for AthenaWork local-first shared workspace workflows.

## Canonical State Authority
- Authoritative source: control-plane backend state store.
- Derived source: markdown lane artifacts (`intake`, `active`, `qa`, `done`, `blocked`, `archive`).
- Reconciliation rule: markdown updates are projections from authoritative state transitions; markdown edits alone are non-authoritative.
- Exception path: undocumented agent-to-agent communication is permitted only in research mode and must emit auditable exception events.

## Entity Model

### `story`
- Identity: `story_id`
- Core fields: `lane`, `status`, `owner_persona`, `phase`, `idea_id`, `adr_refs`
- Concurrency guard: `version`

### `cycle`
- Identity: `cycle_id`
- Core fields: `stage`, `started_at`, `closed_at`, `operator`, `observer_report_path`

### `transition`
- Identity: `transition_id`
- Core fields: `transition_type`, `requested_by`, `requested_at`, `precondition_snapshot`, `result`

### `observer_event`
- Identity: `event_id`
- Core fields: `event_type`, `timestamp`, `cycle_id`, `story_id`, `session_id`, `trace_id`, `payload_hash`

## Transition API Contract

## Request
`POST /api/v1/workflow/transitions`

Required fields:
- `transition_type`
- `story_id` (required for story-scoped transitions)
- `cycle_id`
- `session_id`
- `trace_id`
- `requested_by`
- `expected_story_version` (for optimistic concurrency)
- `lane_context` (`research` | `standard`)
- `direction_confirmation_id` (required for direction-changing transitions)

## Response
- `200 OK`: transition accepted and committed; returns updated entity snapshot and emitted event IDs.
- `409 Conflict`: stale version or violated lane precondition.
- `422 Unprocessable Entity`: missing required fields or invalid state transition.
- `423 Locked`: transition blocked pending human direction confirmation.
- `403 Forbidden`: undocumented agent communication attempted outside research mode.

## Core Transition Types and Preconditions
1. `ARCH_PROMOTE_INTAKE_TO_ACTIVE`
- Preconditions: story exists in `architecture/intake`; no active duplicate.

2. `ARCH_MOVE_ACTIVE_TO_QA`
- Preconditions: story is `architecture/active`; architecture outputs + handoff present.

3. `PM_PROMOTE_ENG_INTAKE_TO_ACTIVE`
- Preconditions: story exists in `engineering/intake`; active queue ordering update included.

4. `ENG_MOVE_ACTIVE_TO_QA`
- Preconditions: implementation + required tests completed.

5. `QA_ACCEPT_TO_DONE`
- Preconditions: QA result attached; acceptance criteria passed.

6. `QA_RETURN_TO_ACTIVE`
- Preconditions: defect or criteria gap documented with severity.

7. `CYCLE_OBSERVER_RECORD`
- Preconditions: observer report generated and linked to `cycle_id`.

8. `DIRECTION_CHANGE_REQUEST`
- Preconditions: explicit `direction_confirmation_id` pointing to human-approved workflow artifact.

## Failure Code Set
- `ERR_PRECONDITION_LANE_MISMATCH`
- `ERR_PRECONDITION_STATUS_MISMATCH`
- `ERR_CONCURRENCY_VERSION_STALE`
- `ERR_HUMAN_CONFIRMATION_REQUIRED`
- `ERR_RESEARCH_ONLY_EXCEPTION`
- `ERR_ARTIFACT_MISSING`
- `ERR_TRANSITION_UNKNOWN`

## Immutable Event Contract
Each event must include:
- `event_id`
- `event_type`
- `cycle_id`
- `story_id` (nullable for cycle-wide events)
- `session_id`
- `trace_id`
- `timestamp`
- `actor`
- `transition_id`
- `result` (`accepted` | `rejected`)
- `error_code` (nullable)
- `payload_hash`

Research-only exception event:
- `event_type`: `research_comm_exception`
- Required payload: `justification`, `auditable_artifact_path`, `approver` (optional)

## Policy Mapping
- Undocumented agent communication:
  - `lane_context=research`: allowed only with exception event.
  - `lane_context=standard`: rejected with `ERR_RESEARCH_ONLY_EXCEPTION`.
- Direction-changing actions:
  - blocked with `ERR_HUMAN_CONFIRMATION_REQUIRED` until `direction_confirmation_id` is provided.
- Human and UI optimization:
  - APIs return concise machine-readable errors for agents.
  - Planning/UI artifacts must preserve low-vision-friendly defaults.

## Risks and Tradeoffs
- Consistency vs complexity:
  - More reliable concurrency and audits, but requires strict projection/reconciliation.
- Local portability vs deployability:
  - Local-first defaults improve operator startup; optional remote deployment compatibility requires adapter boundaries.
