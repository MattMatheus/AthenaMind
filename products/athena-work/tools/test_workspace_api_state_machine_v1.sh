#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

contract="$root_dir/operating-system/contracts/WORKSPACE_API_STATE_MACHINE_V1.md"

doc_test_init

doc_assert_exists "$contract" "Workspace API state machine contract exists"
doc_assert_contains "$contract" "POST /api/v1/workflow/promote" "Contract defines promote endpoint"
doc_assert_contains "$contract" "POST /api/v1/workflow/fail" "Contract defines fail endpoint"
doc_assert_contains "$contract" "POST /api/v1/workflow/close_cycle" "Contract defines close_cycle endpoint"
doc_assert_contains "$contract" "POST /api/v1/workflow/confirm_direction" "Contract defines confirm_direction endpoint"
doc_assert_contains "$contract" "ERR_POLICY_RESEARCH_ONLY_EXCEPTION" "Contract defines research-mode rejection code"
doc_assert_contains "$contract" "ERR_CONFIRM_DIRECTION_REQUIRED" "Contract defines direction confirmation rejection code"
doc_assert_contains "$contract" "ERR_CONFIRM_DIRECTION_EXPIRED" "Contract defines direction confirmation expiry rejection code"
doc_assert_contains "$contract" "ERR_CONFIRM_DIRECTION_SUPERSEDED" "Contract defines direction confirmation superseded rejection code"
doc_assert_contains "$contract" "lane_context" "Contract defines research lane context requirement"
doc_assert_contains "$contract" "source_agent" "Contract defines research exception source agent audit field"
doc_assert_contains "$contract" "target_agent" "Contract defines research exception target agent audit field"
doc_assert_contains "$contract" "reason" "Contract defines research exception reason audit field"
doc_assert_contains "$contract" "confirmed_by" "Contract defines direction confirmation model confirmed_by field"
doc_assert_contains "$contract" "confirmed_at" "Contract defines direction confirmation model confirmed_at field"
doc_assert_contains "$contract" "scope" "Contract defines direction confirmation model scope field"
doc_assert_contains "$contract" "expiry" "Contract defines direction confirmation model expiry field"
doc_assert_contains "$contract" "immutable event" "Contract defines immutable transition event write"
doc_assert_contains "$contract" "code" "Contract includes machine-readable response code field"
doc_assert_contains "$contract" "reason" "Contract includes machine-readable response reason field"
doc_assert_contains "$contract" "next_action" "Contract includes machine-readable response next_action field"
doc_assert_contains "$contract" "correlation_id" "Contract includes machine-readable response correlation_id field"
doc_assert_contains "$contract" "<= 300ms" "Contract defines local response p95 target"

doc_test_finish
