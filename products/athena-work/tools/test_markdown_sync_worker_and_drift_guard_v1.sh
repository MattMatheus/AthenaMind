#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

worker="$root_dir/tools/markdown_sync_worker.sh"
guard="$root_dir/tools/check_markdown_drift.sh"
state_file="$root_dir/products/athena-work/operating-system/state/backend_read_model_v1.json"
projected_active="$root_dir/products/athena-work/delivery-backlog/engineering/active/README.md"
projected_board="$root_dir/products/athena-work/operating-system/observer/LATEST_BOARD_READ_MODEL.md"
projected_timeline="$root_dir/products/athena-work/operating-system/observer/LATEST_TIMELINE_READ_MODEL.md"
projected_docs="$root_dir/products/athena-work/operating-system/observer/LATEST_DOCS_INDEX_READ_MODEL.md"

doc_test_init

doc_assert_exists "$worker" "Markdown sync worker exists"
doc_assert_exists "$guard" "Markdown drift guard exists"
doc_assert_exists "$state_file" "Canonical backend read-model state file exists"
doc_assert_contains "$worker" "dry-run" "Sync worker supports dry-run mode"
doc_assert_contains "$guard" "ordering_conflict" "Drift guard checks ordering conflict class"
doc_assert_contains "$guard" "missing_artifact" "Drift guard checks missing artifact class"
doc_assert_contains "$guard" "stale_revision" "Drift guard checks stale revision class"
doc_assert_contains "$guard" "remediation_id:" "Drift guard emits deterministic remediation ids"

doc_assert_exists "$projected_active" "Projected engineering active README exists"
doc_assert_exists "$projected_board" "Projected board markdown exists"
doc_assert_exists "$projected_timeline" "Projected timeline markdown exists"
doc_assert_exists "$projected_docs" "Projected docs index markdown exists"

worker_output="$("$worker" --dry-run)"
if grep -Eq '^status: (clean|drift)$' <<<"$worker_output"; then
  echo "PASS: sync worker dry-run emits deterministic status"
else
  echo "FAIL: sync worker dry-run did not emit expected status"
  echo "$worker_output"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi

guard_output="$("$guard" --dry-run)"
if grep -Fq "status: pass" <<<"$guard_output"; then
  echo "PASS: drift guard reports pass after projection sync"
else
  echo "FAIL: drift guard did not report pass"
  echo "$guard_output"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi

doc_test_finish
