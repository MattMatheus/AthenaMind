#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
launch_script="$root_dir/tools/launch_stage.sh"
current_branch="$(git -C "$root_dir" branch --show-current)"

output="$(WORKSPACE_API_ENABLED=true WORKSPACE_API_CURL_BIN=missing-curl ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" pm)"
if grep -Fq "workspace_api_adapter:" <<<"$output" && grep -Fq "status: fallback" <<<"$output"; then
  echo "PASS: launch_stage emits fallback adapter status when API is unavailable"
else
  echo "FAIL: launch_stage did not emit expected fallback adapter status"
  echo "$output"
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
err_file="$tmp_dir/launch.err"
direction_log="$tmp_dir/direction-confirmations.jsonl"
set +e
ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" pm >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_CONFIRM_DIRECTION_REQUIRED" "$err_file"; then
  echo "PASS: launch_stage blocks direction change without explicit confirmation id"
else
  echo "FAIL: launch_stage did not enforce direction confirmation id gate"
  cat "$err_file"
  exit 1
fi

confirm_output="$(ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-1 ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='qa-stage' ATHENA_DIRECTION_EXPIRY='2099-12-31T23:59:59Z' ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" qa)"
if grep -Fq "launch: stage-prompts/active/qa-agent-seed-prompt.md" <<<"$confirm_output" && grep -Fq "direction_confirmation:" <<<"$confirm_output" && grep -Fq "status: valid" <<<"$confirm_output"; then
  echo "PASS: launch_stage proceeds when direction confirmation id is provided"
else
  echo "FAIL: launch_stage did not continue after explicit direction confirmation"
  echo "$confirm_output"
  exit 1
fi

set +e
ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-OLD ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='pm-stage' ATHENA_DIRECTION_EXPIRY='2000-01-01T00:00:00Z' ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" pm >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_CONFIRM_DIRECTION_EXPIRED" "$err_file"; then
  echo "PASS: launch_stage rejects expired direction confirmations with deterministic reason"
else
  echo "FAIL: launch_stage did not reject expired direction confirmation"
  cat "$err_file"
  exit 1
fi

set +e
ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-SUPERSEDED ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='pm-stage' ATHENA_DIRECTION_EXPIRY='2099-12-31T23:59:59Z' ATHENA_DIRECTION_SUPERSEDED=true ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" pm >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_CONFIRM_DIRECTION_SUPERSEDED" "$err_file"; then
  echo "PASS: launch_stage rejects superseded direction confirmations with deterministic reason"
else
  echo "FAIL: launch_stage did not reject superseded direction confirmation"
  cat "$err_file"
  exit 1
fi

set +e
ATHENA_UNDOCUMENTED_AGENT_COMM=true ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" pm >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_POLICY_RESEARCH_ONLY_EXCEPTION" "$err_file"; then
  echo "PASS: launch_stage rejects undocumented agent communication outside research lane"
else
  echo "FAIL: launch_stage did not reject non-research undocumented communication"
  cat "$err_file"
  exit 1
fi

audit_log="$tmp_dir/research-comm-exceptions.jsonl"
research_output="$(ATHENA_REQUIRED_BRANCH="$current_branch" \
  ATHENA_UNDOCUMENTED_AGENT_COMM=true \
  ATHENA_LANE_CONTEXT=research \
  ATHENA_RESEARCH_COMM_EXCEPTION=true \
  ATHENA_SOURCE_AGENT=agent-alpha \
  ATHENA_TARGET_AGENT=agent-beta \
  ATHENA_REASON='investigate parser behavior' \
  ATHENA_CYCLE_ID=PLAN-RESEARCH-1 \
  ATHENA_STORY_ID=STORY-RESEARCH-1 \
  ATHENA_SESSION_ID=SESSION-RESEARCH-1 \
  ATHENA_AUDIT_LOG_PATH="$audit_log" \
  "$launch_script" pm)"
if grep -Fq "research_comm_exception:" <<<"$research_output" && grep -Fq '"source_agent":"agent-alpha"' "$audit_log" && grep -Fq '"target_agent":"agent-beta"' "$audit_log" && grep -Fq '"reason":"investigate parser behavior"' "$audit_log"; then
  echo "PASS: launch_stage logs audited research communication exception with required fields"
else
  echo "FAIL: launch_stage did not log audited research communication exception"
  echo "$research_output"
  cat "$audit_log"
  exit 1
fi

planning_output="$(ATHENA_REQUIRED_BRANCH="$current_branch" ATHENA_DIRECTION_TEXT='Prioritize control plane hardening' ATHENA_DIRECTION_CONSTRAINTS='No human coding; low-vision-first' ATHENA_DIRECTION_NEXT_STAGE='architect' ATHENA_DIRECTION_CONFIRMED_BY='Matt' "$launch_script" planning)"
if grep -Fq "planning_direction_summary:" <<<"$planning_output" && grep -Fq "direction: Prioritize control plane hardening" <<<"$planning_output" && grep -Fq "next_stage: architect" <<<"$planning_output" && grep -Fq "confirmed_by: Matt" <<<"$planning_output"; then
  echo "PASS: launch_stage planning emits one-screen direction confirmation summary"
else
  echo "FAIL: launch_stage planning did not emit expected direction confirmation summary"
  echo "$planning_output"
  exit 1
fi

echo "Result: PASS"
