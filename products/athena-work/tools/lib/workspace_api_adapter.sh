#!/usr/bin/env bash

workspace_api_enabled() {
  [[ "${WORKSPACE_API_ENABLED:-false}" == "true" ]]
}

workspace_api_endpoint() {
  printf '%s' "${WORKSPACE_API_ENDPOINT:-http://127.0.0.1:8787}"
}

workspace_api_curl_bin() {
  printf '%s' "${WORKSPACE_API_CURL_BIN:-curl}"
}

workspace_api_health_ok() {
  local curl_bin endpoint
  curl_bin="$(workspace_api_curl_bin)"
  endpoint="$(workspace_api_endpoint)"

  if ! command -v "$curl_bin" >/dev/null 2>&1; then
    return 1
  fi

  "$curl_bin" -fsS "${endpoint%/}/health" >/dev/null 2>&1
}

workspace_api_emit_status() {
  local context="$1"

  if ! workspace_api_enabled; then
    return 0
  fi

  if workspace_api_health_ok; then
    cat <<EOF_STATUS
workspace_api_adapter:
  status: connected
  action: use_api
  why: endpoint_healthy
  next: continue
EOF_STATUS
    return 0
  fi

  cat <<EOF_STATUS
workspace_api_adapter:
  status: fallback
  action: use_scripts
  why: endpoint_unavailable
  next: start_api
EOF_STATUS
}

workspace_api_handle_direction_confirmation() {
  local context="$1"
  local cycle_id="${2:-${ATHENA_CYCLE_ID:-$context}}"
  local story_id="${3:-${ATHENA_STORY_ID:-unscoped}}"
  local session_id="${4:-${ATHENA_SESSION_ID:-$context}}"
  local curl_bin endpoint payload confirmation_id confirmed_by confirmed_at scope expiry
  local validation_result audit_log audit_log_dir timestamp
  local escaped_context escaped_cycle escaped_story escaped_session
  local escaped_confirmation escaped_confirmed_by escaped_confirmed_at escaped_scope escaped_expiry

  if [[ "${ATHENA_DIRECTION_CHANGE:-false}" != "true" ]]; then
    return 0
  fi

  validation_result="$(workspace_api_direction_confirmation_state)"
  if [[ "$validation_result" == "missing" ]]; then
    echo "abort: ERR_CONFIRM_DIRECTION_REQUIRED direction change requires confirmation_id, confirmed_by, confirmed_at, scope, and expiry" >&2
    return 1
  fi
  if [[ "$validation_result" == "expired" ]]; then
    echo "abort: ERR_CONFIRM_DIRECTION_EXPIRED direction confirmation has expired" >&2
    return 1
  fi
  if [[ "$validation_result" == "superseded" ]]; then
    echo "abort: ERR_CONFIRM_DIRECTION_SUPERSEDED direction confirmation is superseded" >&2
    return 1
  fi

  confirmation_id="${ATHENA_DIRECTION_CONFIRMATION_ID:-}"
  confirmed_by="${ATHENA_DIRECTION_CONFIRMED_BY:-}"
  confirmed_at="${ATHENA_DIRECTION_CONFIRMED_AT:-}"
  scope="${ATHENA_DIRECTION_SCOPE:-}"
  expiry="${ATHENA_DIRECTION_EXPIRY:-}"
  workspace_api_emit_direction_confirmation_status "$context"

  audit_log="${ATHENA_DIRECTION_AUDIT_LOG_PATH:-operating-system/observer/DIRECTION_CONFIRMATIONS.jsonl}"
  audit_log_dir="$(dirname "$audit_log")"
  mkdir -p "$audit_log_dir"
  timestamp="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
  escaped_context="$(workspace_api_json_escape "$context")"
  escaped_cycle="$(workspace_api_json_escape "$cycle_id")"
  escaped_story="$(workspace_api_json_escape "$story_id")"
  escaped_session="$(workspace_api_json_escape "$session_id")"
  escaped_confirmation="$(workspace_api_json_escape "$confirmation_id")"
  escaped_confirmed_by="$(workspace_api_json_escape "$confirmed_by")"
  escaped_confirmed_at="$(workspace_api_json_escape "$confirmed_at")"
  escaped_scope="$(workspace_api_json_escape "$scope")"
  escaped_expiry="$(workspace_api_json_escape "$expiry")"
  printf '{"event_type":"direction_confirmation","timestamp":"%s","context":"%s","cycle_id":"%s","story_id":"%s","session_id":"%s","confirmation_id":"%s","confirmed_by":"%s","confirmed_at":"%s","scope":"%s","expiry":"%s","status":"accepted"}\n' \
    "$timestamp" "$escaped_context" "$escaped_cycle" "$escaped_story" "$escaped_session" "$escaped_confirmation" "$escaped_confirmed_by" "$escaped_confirmed_at" "$escaped_scope" "$escaped_expiry" >>"$audit_log"

  if ! workspace_api_enabled; then
    return 0
  fi

  if ! workspace_api_health_ok; then
    echo "warning: direction confirmation API call skipped; workspace API unavailable" >&2
    return 0
  fi

  curl_bin="$(workspace_api_curl_bin)"
  endpoint="$(workspace_api_endpoint)"
  payload="{\"confirmation_id\":\"${confirmation_id}\",\"context\":\"${context}\",\"confirmed_by\":\"${confirmed_by}\",\"confirmed_at\":\"${confirmed_at}\",\"scope\":\"${scope}\",\"expiry\":\"${expiry}\"}"

  if ! "$curl_bin" -fsS -X POST "${endpoint%/}/api/v1/workflow/confirm_direction" -H 'content-type: application/json' -d "$payload" >/dev/null 2>&1; then
    echo "warning: direction confirmation API call failed; continuing with explicit confirmation id" >&2
  fi

  return 0
}

workspace_api_direction_confirmation_state() {
  local confirmation_id confirmed_by confirmed_at scope expiry now_utc

  if [[ "${ATHENA_DIRECTION_CHANGE:-false}" != "true" ]]; then
    printf 'not_required'
    return 0
  fi

  confirmation_id="${ATHENA_DIRECTION_CONFIRMATION_ID:-}"
  confirmed_by="${ATHENA_DIRECTION_CONFIRMED_BY:-}"
  confirmed_at="${ATHENA_DIRECTION_CONFIRMED_AT:-}"
  scope="${ATHENA_DIRECTION_SCOPE:-}"
  expiry="${ATHENA_DIRECTION_EXPIRY:-}"
  if [[ -z "$confirmation_id" || -z "$confirmed_by" || -z "$confirmed_at" || -z "$scope" || -z "$expiry" ]]; then
    printf 'missing'
    return 0
  fi

  if [[ "${ATHENA_DIRECTION_SUPERSEDED:-false}" == "true" || -n "${ATHENA_DIRECTION_SUPERSEDED_BY:-}" ]]; then
    printf 'superseded'
    return 0
  fi

  now_utc="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
  if [[ "$expiry" < "$now_utc" ]]; then
    printf 'expired'
    return 0
  fi

  printf 'valid'
}

workspace_api_emit_direction_confirmation_status() {
  local context="$1"
  local state confirmation_id confirmed_by confirmed_at scope expiry

  state="$(workspace_api_direction_confirmation_state)"
  if [[ "$state" == "not_required" ]]; then
    return 0
  fi

  confirmation_id="${ATHENA_DIRECTION_CONFIRMATION_ID:-unconfirmed}"
  confirmed_by="${ATHENA_DIRECTION_CONFIRMED_BY:-unconfirmed}"
  confirmed_at="${ATHENA_DIRECTION_CONFIRMED_AT:-unconfirmed}"
  scope="${ATHENA_DIRECTION_SCOPE:-unconfirmed}"
  expiry="${ATHENA_DIRECTION_EXPIRY:-unconfirmed}"

  cat <<EOF_DIRECTION
direction_confirmation:
  status: $state
  context: $context
  confirmation_id: $confirmation_id
  confirmed_by: $confirmed_by
  confirmed_at: $confirmed_at
  scope: $scope
  expiry: $expiry
EOF_DIRECTION
}

workspace_api_emit_planning_direction_summary() {
  local direction constraints next_stage confirmed_by
  direction="${ATHENA_DIRECTION_TEXT:-unset}"
  constraints="${ATHENA_DIRECTION_CONSTRAINTS:-unset}"
  next_stage="${ATHENA_DIRECTION_NEXT_STAGE:-unset}"
  confirmed_by="${ATHENA_DIRECTION_CONFIRMED_BY:-unconfirmed}"

  cat <<EOF_PLANNING
planning_direction_summary:
  direction: $direction
  constraints: $constraints
  next_stage: $next_stage
  confirmed_by: $confirmed_by
EOF_PLANNING
}

workspace_api_json_escape() {
  local value="$1"
  value="${value//\\/\\\\}"
  value="${value//\"/\\\"}"
  value="${value//$'\n'/\\n}"
  printf '%s' "$value"
}

workspace_api_handle_research_comm_exception() {
  local context="$1"
  local cycle_id="${2:-${ATHENA_CYCLE_ID:-}}"
  local story_id="${3:-${ATHENA_STORY_ID:-}}"
  local session_id="${4:-${ATHENA_SESSION_ID:-}}"
  local lane_context explicit_exception source_agent target_agent reason
  local audit_log_dir audit_log timestamp
  local escaped_context escaped_cycle escaped_story escaped_session
  local escaped_source escaped_target escaped_reason

  if [[ "${ATHENA_UNDOCUMENTED_AGENT_COMM:-false}" != "true" ]]; then
    return 0
  fi

  lane_context="${ATHENA_LANE_CONTEXT:-standard}"
  explicit_exception="${ATHENA_RESEARCH_COMM_EXCEPTION:-false}"

  if [[ "$lane_context" != "research" ]]; then
    echo "abort: ERR_POLICY_RESEARCH_ONLY_EXCEPTION undocumented agent communication requires lane_context=research" >&2
    return 1
  fi

  if [[ "$explicit_exception" != "true" ]]; then
    echo "abort: ERR_POLICY_RESEARCH_ONLY_EXCEPTION undocumented agent communication requires ATHENA_RESEARCH_COMM_EXCEPTION=true" >&2
    return 1
  fi

  source_agent="${ATHENA_SOURCE_AGENT:-}"
  target_agent="${ATHENA_TARGET_AGENT:-}"
  reason="${ATHENA_REASON:-}"
  if [[ -z "$cycle_id" || -z "$story_id" || -z "$session_id" || -z "$source_agent" || -z "$target_agent" || -z "$reason" ]]; then
    echo "abort: research communication exception requires cycle_id, story_id, session_id, source_agent, target_agent, and reason" >&2
    return 1
  fi

  audit_log="${ATHENA_AUDIT_LOG_PATH:-operating-system/observer/RESEARCH_COMM_EXCEPTIONS.jsonl}"
  audit_log_dir="$(dirname "$audit_log")"
  mkdir -p "$audit_log_dir"
  timestamp="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

  escaped_context="$(workspace_api_json_escape "$context")"
  escaped_cycle="$(workspace_api_json_escape "$cycle_id")"
  escaped_story="$(workspace_api_json_escape "$story_id")"
  escaped_session="$(workspace_api_json_escape "$session_id")"
  escaped_source="$(workspace_api_json_escape "$source_agent")"
  escaped_target="$(workspace_api_json_escape "$target_agent")"
  escaped_reason="$(workspace_api_json_escape "$reason")"

  printf '{"event_type":"research_comm_exception","timestamp":"%s","context":"%s","cycle_id":"%s","story_id":"%s","session_id":"%s","source_agent":"%s","target_agent":"%s","reason":"%s"}\n' \
    "$timestamp" "$escaped_context" "$escaped_cycle" "$escaped_story" "$escaped_session" "$escaped_source" "$escaped_target" "$escaped_reason" >>"$audit_log"

  cat <<EOF_RESEARCH
research_comm_exception:
  status: logged
  action: allow_research_exception
  why: research_lane_with_explicit_flag
  next: continue
EOF_RESEARCH
}
