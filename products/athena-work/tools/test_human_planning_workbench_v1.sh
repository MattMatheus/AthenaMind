#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

ui_file="$root_dir/products/athena-work/ui/index.html"
api_file="$root_dir/products/athena-work/ui/local_control_plane_api.py"

doc_test_init

doc_assert_exists "$ui_file" "Planning workbench UI exists"
doc_assert_exists "$api_file" "Planning workbench API exists"

doc_assert_contains "$ui_file" "Planning Workbench" "UI includes planning workbench section"
doc_assert_contains "$ui_file" "Direction" "UI includes direction field"
doc_assert_contains "$ui_file" "Constraints" "UI includes constraints field"
doc_assert_contains "$ui_file" "Risks" "UI includes risks field"
doc_assert_contains "$ui_file" "Next stage" "UI includes next stage selector"
doc_assert_contains "$ui_file" "Confirm direction" "UI includes explicit human confirmation action"
doc_assert_contains "$ui_file" "Confirmation status: unconfirmed" "UI shows immediate confirmation status"
doc_assert_contains "$ui_file" "Export to workflow artifact" "UI includes workflow artifact export action"
doc_assert_contains "$ui_file" "What happens next" "UI includes explicit next-step summary"

doc_assert_contains "$api_file" "/api/v1/planning/export" "API exposes planning export endpoint"
doc_assert_contains "$api_file" "PLAN-WORKBENCH-" "Export writes planning artifacts"
doc_assert_contains "$api_file" "direction, constraints, risks, and next_stage are required" "Export endpoint validates required fields"

doc_test_finish
