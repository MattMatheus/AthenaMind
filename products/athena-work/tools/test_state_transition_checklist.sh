#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

checklist="$root_dir/delivery-backlog/STATE_TRANSITION_CHECKLIST.md"
eng_prompt="$root_dir/stage-prompts/active/next-agent-seed-prompt.md"
qa_prompt="$root_dir/stage-prompts/active/qa-agent-seed-prompt.md"
pm_prompt="$root_dir/stage-prompts/active/pm-refinement-seed-prompt.md"

doc_test_init

doc_assert_exists "$checklist" "State transition checklist exists"
doc_assert_contains "$checklist" "active -> qa" "Checklist defines active to qa transition"
doc_assert_contains "$checklist" "qa -> done" "Checklist defines qa to done transition"
doc_assert_contains "$checklist" "qa -> active" "Checklist defines qa to active failure path"
doc_assert_contains "$checklist" "Cycle Closure (Observer + Commit)" "Checklist defines cycle closure transition"
doc_assert_contains "$checklist" "Required artifacts" "Checklist defines artifact requirements"
doc_assert_contains "$checklist" "Approval" "Checklist defines approval ownership"

doc_assert_contains "$eng_prompt" "STATE_TRANSITION_CHECKLIST.md" "Engineering prompt references checklist"
doc_assert_contains "$qa_prompt" "STATE_TRANSITION_CHECKLIST.md" "QA prompt references checklist"
doc_assert_contains "$pm_prompt" "STATE_TRANSITION_CHECKLIST.md" "PM prompt references checklist"
doc_assert_contains "$qa_prompt" "run_observer_cycle.sh" "QA prompt includes observer command"

doc_test_finish
