#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

api_file="$root_dir/products/athena-work/ui/local_control_plane_api.py"
ui_file="$root_dir/products/athena-work/ui/index.html"
gen_tool="$root_dir/products/athena-work/tools/generate_launch_authorization_package.sh"
val_tool="$root_dir/products/athena-work/tools/validate_launch_authorization_package.sh"

doc_test_init

doc_assert_exists "$gen_tool" "Launch authorization package generator exists"
doc_assert_exists "$val_tool" "Launch authorization package validator exists"
doc_assert_contains "$api_file" "/api/v1/launch/package" "API exposes launch package generation endpoint"
doc_assert_contains "$api_file" "/api/v1/launch/validate" "API exposes launch package validation endpoint"
doc_assert_contains "$api_file" "/api/v1/launch/latest" "API exposes latest launch package endpoint"
doc_assert_contains "$ui_file" "Launch Authorization" "UI includes launch authorization panel"
doc_assert_contains "$ui_file" "Generate launch package" "UI includes launch package action"
doc_assert_contains "$ui_file" "Validate latest package" "UI includes launch validation action"
doc_assert_contains "$ui_file" "Flight Director" "UI launch panel includes operator role cue"
doc_assert_contains "$ui_file" "Launch readiness:" "UI includes launch readiness summary strip"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
package_out="$("$gen_tool" --root "$root_dir" | python3 -c 'import json,sys; print(json.load(sys.stdin)["path"])')"
doc_assert_exists "$package_out" "Launch authorization package file generated"

if "$val_tool" --package "$package_out" --root "$root_dir" >/dev/null 2>&1; then
  echo "PASS: Generated launch package validates"
else
  echo "FAIL: Generated launch package validates"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi

doc_test_finish
