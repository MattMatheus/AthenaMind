#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"
howto_doc="$root_dir/knowledge-base/how-to/GO_TOOLCHAIN_SETUP.md"
workflow_doc="$root_dir/knowledge-base/process/OPERATOR_DAILY_WORKFLOW.md"

doc_test_init

doc_assert_exists "$root_dir/tools/check_go_toolchain.sh" "Go toolchain preflight script exists"
doc_assert_exists "$howto_doc" "Go setup how-to doc exists"
doc_assert_contains "$howto_doc" 'Source of truth: `go.mod`' "How-to references go.mod as source of truth"
doc_assert_contains "$howto_doc" "tools/check_go_toolchain.sh" "How-to references toolchain preflight command"
doc_assert_contains "$howto_doc" "go test ./..." "How-to references canonical Go test command"
doc_assert_contains "$workflow_doc" "tools/check_go_toolchain.sh" "Daily workflow includes Go toolchain preflight"

required_go="$(awk '/^go[[:space:]]+[0-9]+\.[0-9]+/{print $2; exit}' "$root_dir/go.mod")"
if [[ -z "$required_go" ]]; then
  echo "FAIL: go.mod includes required Go version"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
else
  if grep -Fq "Current minimum: \`go $required_go\`" "$howto_doc"; then
    echo "PASS: How-to minimum Go version matches go.mod ($required_go)"
  else
    echo "FAIL: How-to minimum Go version matches go.mod ($required_go)"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
fi

doc_test_finish
