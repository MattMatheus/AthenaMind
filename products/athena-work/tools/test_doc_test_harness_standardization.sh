#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

standard_doc="$root_dir/product-research/roadmap/DOC_TEST_HARNESS_STANDARD.md"
engineering_prompt="$root_dir/stage-prompts/active/next-agent-seed-prompt.md"
runner="$root_dir/tools/run_doc_tests.sh"
shared_lib="$root_dir/tools/lib/doc_test_harness.sh"
migrated_test="$root_dir/tools/test_goals_scorecard_v01.sh"

doc_test_init

doc_assert_exists "$standard_doc" "Harness standard doc exists"
doc_assert_exists "$runner" "Canonical runner exists"
doc_assert_exists "$shared_lib" "Shared harness library exists"

doc_assert_contains "$standard_doc" "tools/run_doc_tests.sh" "Standard doc defines canonical command path"
doc_assert_contains "$standard_doc" "tools/test_<story-scope>.sh" "Standard doc defines test location pattern"
doc_assert_contains "$engineering_prompt" "tools/run_doc_tests.sh" "Engineering prompt references canonical command"
doc_assert_contains "$migrated_test" "tools/lib/doc_test_harness.sh" "Existing story test migrated to shared harness"

doc_test_finish
