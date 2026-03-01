#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
observer_script="$root_dir/tools/run_observer_cycle.sh"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
direction_log="$tmp_dir/direction-confirmations.jsonl"

report_path="$tmp_dir/observer-report.md"
output="$(WORKSPACE_API_ENABLED=true WORKSPACE_API_CURL_BIN=missing-curl "$observer_script" --cycle-id STORY-TEST-OBSERVER-ADAPTER --output "$report_path")"
if grep -Fq "workspace_api_adapter:" <<<"$output" && grep -Fq "status: fallback" <<<"$output"; then
  echo "PASS: observer emits fallback adapter status when API is unavailable"
else
  echo "FAIL: observer did not emit expected fallback adapter status"
  echo "$output"
  exit 1
fi

err_file="$tmp_dir/observer.err"
set +e
ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" "$observer_script" --cycle-id STORY-TEST-OBSERVER-DIRECTION --output "$tmp_dir/observer-direction.md" >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_CONFIRM_DIRECTION_REQUIRED" "$err_file"; then
  echo "PASS: observer blocks direction change without explicit confirmation id"
else
  echo "FAIL: observer did not enforce direction confirmation id gate"
  cat "$err_file"
  exit 1
fi

if ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-2 ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='observer' ATHENA_DIRECTION_EXPIRY='2099-12-31T23:59:59Z' ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" "$observer_script" --cycle-id STORY-TEST-OBSERVER-DIRECTION-OK --output "$tmp_dir/observer-direction-ok.md" >/dev/null; then
  echo "PASS: observer proceeds when direction confirmation id is provided"
else
  echo "FAIL: observer did not continue after explicit direction confirmation"
  exit 1
fi

set +e
ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-OLD ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='observer' ATHENA_DIRECTION_EXPIRY='2000-01-01T00:00:00Z' ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" "$observer_script" --cycle-id STORY-TEST-OBSERVER-DIRECTION-EXPIRED --output "$tmp_dir/observer-direction-expired.md" >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_CONFIRM_DIRECTION_EXPIRED" "$err_file"; then
  echo "PASS: observer rejects expired direction confirmations with deterministic reason"
else
  echo "FAIL: observer did not reject expired direction confirmation"
  cat "$err_file"
  exit 1
fi

set +e
ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-SUPERSEDED ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='observer' ATHENA_DIRECTION_EXPIRY='2099-12-31T23:59:59Z' ATHENA_DIRECTION_SUPERSEDED=true ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_log" "$observer_script" --cycle-id STORY-TEST-OBSERVER-DIRECTION-SUPERSEDED --output "$tmp_dir/observer-direction-superseded.md" >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_CONFIRM_DIRECTION_SUPERSEDED" "$err_file"; then
  echo "PASS: observer rejects superseded direction confirmations with deterministic reason"
else
  echo "FAIL: observer did not reject superseded direction confirmation"
  cat "$err_file"
  exit 1
fi

set +e
ATHENA_UNDOCUMENTED_AGENT_COMM=true "$observer_script" --cycle-id STORY-TEST-OBSERVER-COMM-BLOCK --output "$tmp_dir/observer-comm-block.md" >/dev/null 2>"$err_file"
code=$?
set -e
if [[ $code -ne 0 ]] && grep -Fq "ERR_POLICY_RESEARCH_ONLY_EXCEPTION" "$err_file"; then
  echo "PASS: observer rejects undocumented agent communication outside research lane"
else
  echo "FAIL: observer did not reject non-research undocumented communication"
  cat "$err_file"
  exit 1
fi

audit_log="$tmp_dir/research-comm-exceptions.jsonl"
cat >"$tmp_dir/story.md" <<'EOF'
# Story
EOF
research_output="$(ATHENA_UNDOCUMENTED_AGENT_COMM=true \
  ATHENA_LANE_CONTEXT=research \
  ATHENA_RESEARCH_COMM_EXCEPTION=true \
  ATHENA_SOURCE_AGENT=agent-gamma \
  ATHENA_TARGET_AGENT=agent-delta \
  ATHENA_REASON='compare alternate synthesis paths' \
  ATHENA_AUDIT_LOG_PATH="$audit_log" \
  "$observer_script" --cycle-id STORY-TEST-OBSERVER-COMM-ALLOW --story "$tmp_dir/story.md" --output "$tmp_dir/observer-comm-allow.md")"
if grep -Fq "research_comm_exception:" <<<"$research_output" && grep -Fq '"cycle_id":"STORY-TEST-OBSERVER-COMM-ALLOW"' "$audit_log" && grep -Fq '"source_agent":"agent-gamma"' "$audit_log" && grep -Fq '"target_agent":"agent-delta"' "$audit_log" && grep -Fq '"reason":"compare alternate synthesis paths"' "$audit_log"; then
  echo "PASS: observer logs audited research communication exception with required fields"
else
  echo "FAIL: observer did not log audited research communication exception"
  echo "$research_output"
  cat "$audit_log"
  exit 1
fi

direction_audit="$tmp_dir/direction-confirmations.jsonl"
if ATHENA_DIRECTION_CHANGE=true ATHENA_DIRECTION_CONFIRMATION_ID=CONFIRM-TEST-REPORT ATHENA_DIRECTION_CONFIRMED_BY='Matt' ATHENA_DIRECTION_CONFIRMED_AT='2026-02-27T00:00:00Z' ATHENA_DIRECTION_SCOPE='observer-report' ATHENA_DIRECTION_EXPIRY='2099-12-31T23:59:59Z' ATHENA_DIRECTION_AUDIT_LOG_PATH="$direction_audit" "$observer_script" --cycle-id STORY-TEST-OBSERVER-DIRECTION-REPORT --output "$tmp_dir/observer-direction-report.md" >/dev/null; then
  :
else
  echo "FAIL: observer did not accept valid direction confirmation for report evidence test"
  exit 1
fi
if grep -Fq "## Direction Confirmation Evidence" "$tmp_dir/observer-direction-report.md" && grep -Fq "confirmation_status: valid" "$tmp_dir/observer-direction-report.md" && grep -Fq '"confirmation_id":"CONFIRM-TEST-REPORT"' "$direction_audit"; then
  echo "PASS: observer report and audit log include direction confirmation evidence"
else
  echo "FAIL: observer did not emit direction confirmation evidence in report/audit"
  cat "$tmp_dir/observer-direction-report.md"
  cat "$direction_audit"
  exit 1
fi

echo "Result: PASS"
