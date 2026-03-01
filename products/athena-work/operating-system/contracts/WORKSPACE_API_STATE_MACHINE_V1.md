# Workspace API State Machine v1

## Purpose
Define deterministic transition APIs and policy gates for AthenaWork shared workspace operations.

## Response Contract (All Endpoints)
Responses must be machine-readable and concise.

Required fields:
- `code`
- `reason`
- `next_action`
- `correlation_id`

Latency target:
- local-mode p95 transition response: `<= 300ms`

## Canonical Endpoints

### `POST /api/v1/workflow/promote`
Promotes a story between allowed lanes based on transition policy.

### `POST /api/v1/workflow/fail`
Records QA or policy failure and returns story to active/intake as applicable.

### `POST /api/v1/workflow/close_cycle`
Closes a cycle after observer evidence and required artifacts are present.

### `POST /api/v1/workflow/confirm_direction`
Records explicit human direction confirmation for direction-changing transitions.

## Transition Preconditions
1. `promote`
- Story exists and current lane is valid for requested transition.
- Required artifacts for source stage are present.

2. `fail`
- Failure evidence exists and severity is classified.

3. `close_cycle`
- Observer report exists.
- No stage-level commit occurred before observer.

4. `confirm_direction`
- Human identity and decision rationale are captured.
- Confirmation artifact path is recorded.
- Confirmation model fields are required:
  - `confirmed_by`
  - `confirmed_at`
  - `scope`
  - `expiry`

## Deterministic Error Codes
- `ERR_TRANSITION_UNKNOWN`
- `ERR_PRECONDITION_LANE_MISMATCH`
- `ERR_PRECONDITION_ARTIFACT_MISSING`
- `ERR_POLICY_RESEARCH_ONLY_EXCEPTION`
- `ERR_CONFIRM_DIRECTION_REQUIRED`
- `ERR_CONFIRM_DIRECTION_EXPIRED`
- `ERR_CONFIRM_DIRECTION_SUPERSEDED`
- `ERR_CONFIRM_DIRECTION_MISSING_ARTIFACT`
- `ERR_CONCURRENCY_VERSION_STALE`

## Research-Mode Communication Policy
- Undocumented agent-to-agent communication is allowed only when:
  - `lane_context` is `research`
  - request explicitly flags research exception
  - immutable event log entry is written
- If not in research mode, request must be rejected with `ERR_POLICY_RESEARCH_ONLY_EXCEPTION`.
- Research exception audit event must include:
  - `cycle_id`
  - `story_id`
  - `session_id`
  - `source_agent`
  - `target_agent`
  - `reason`

## Direction-Change Policy
- Any direction-changing transition must reference a prior successful `confirm_direction` record.
- Missing confirmation must be rejected with `ERR_CONFIRM_DIRECTION_REQUIRED`.
- Expired confirmation must be rejected with `ERR_CONFIRM_DIRECTION_EXPIRED`.
- Superseded confirmation must be rejected with `ERR_CONFIRM_DIRECTION_SUPERSEDED`.

## Immutable Transition Event Write
Every accepted transition must emit immutable event data with:
- `cycle_id`
- `story_id`
- `session_id`
- `trace_id`
- `correlation_id`
- `transition_type`
- `result`

## Human-Readable Rejection Reasons
- `reason` values must remain readable for UI display while preserving deterministic `code` semantics.

## Notes
- This contract maps to ARCH-20260227 control-plane and markdown sync decisions.
