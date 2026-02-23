#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

rubric="$root_dir/delivery-backlog/QA_REGRESSION_RUBRIC.md"
bug_template="$root_dir/delivery-backlog/engineering/intake/BUG_TEMPLATE.md"
qa_prompt="$root_dir/stage-prompts/active/qa-agent-seed-prompt.md"

doc_test_init

doc_assert_exists "$rubric" "QA regression rubric doc exists"
doc_assert_contains "$rubric" "Pass/Fail Gates (Deterministic)" "Rubric defines deterministic gates"
doc_assert_contains "$rubric" "If any gate fails, result is" "Rubric defines fail condition"
doc_assert_contains "$rubric" "FAIL" "Rubric includes explicit FAIL token"
doc_assert_contains "$rubric" "P0" "Rubric includes P0 mapping"
doc_assert_contains "$rubric" "P1" "Rubric includes P1 mapping"
doc_assert_contains "$rubric" "P2" "Rubric includes P2 mapping"
doc_assert_contains "$rubric" "P3" "Rubric includes P3 mapping"
doc_assert_contains "$rubric" "Minimum Evidence Requirements for QA Bug Filing" "Rubric defines evidence minimums"
doc_assert_contains "$rubric" "Example PASS" "Rubric includes pass handoff example"
doc_assert_contains "$rubric" "Example FAIL" "Rubric includes fail handoff example"
doc_assert_contains "$qa_prompt" "QA_REGRESSION_RUBRIC.md" "QA prompt references rubric"
doc_assert_contains "$bug_template" "P0" "Bug template includes priority definitions"

doc_test_finish
