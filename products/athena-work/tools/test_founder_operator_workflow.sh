#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

workflow_doc="$root_dir/knowledge-base/process/operator-daily-workflow.md"

doc_test_init

doc_assert_exists "$workflow_doc" "Operator workflow doc exists"
doc_assert_contains "$workflow_doc" "## Startup Routine" "Workflow includes startup routine"
doc_assert_contains "$workflow_doc" "## Engineering + QA Cycle Loop" "Workflow includes engineering and QA cycle loop"
doc_assert_contains "$workflow_doc" "## Shutdown Routine" "Workflow includes shutdown routine"
doc_assert_contains "$workflow_doc" "tools/launch_stage.sh engineering" "Workflow references engineering launcher"
doc_assert_contains "$workflow_doc" "tools/launch_stage.sh qa" "Workflow references QA launcher"
doc_assert_contains "$workflow_doc" "stage-prompts/active/next-agent-seed-prompt.md" "Workflow references engineering prompt"
doc_assert_contains "$workflow_doc" "stage-prompts/active/qa-agent-seed-prompt.md" "Workflow references QA prompt"
doc_assert_contains "$workflow_doc" "tools/run_observer_cycle.sh" "Workflow references observer script"
doc_assert_contains "$workflow_doc" "cycle-<cycle-id>" "Workflow references cycle commit format"
doc_assert_contains "$workflow_doc" "If engineering launch returns" "Workflow includes empty backlog instruction"
doc_assert_contains "$workflow_doc" "no stories" "Workflow includes explicit no-stories token"
doc_assert_contains "$workflow_doc" "If QA finds blocking defects" "Workflow includes QA failure instruction"
doc_assert_contains "$workflow_doc" "ATHENA_REQUIRED_BRANCH" "Workflow includes branch discipline"
doc_assert_contains "$workflow_doc" "default dev" "Workflow includes explicit default branch token"
doc_assert_contains "$workflow_doc" "command escalation" "Workflow includes escalation rules"

doc_test_finish
