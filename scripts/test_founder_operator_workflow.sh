#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

workflow_doc="$root_dir/OPERATOR_DAILY_WORKFLOW.md"

doc_test_init

doc_assert_exists "$workflow_doc" "Operator workflow doc exists"
doc_assert_contains "$workflow_doc" "## Startup Routine" "Workflow includes startup routine"
doc_assert_contains "$workflow_doc" "## Engineering Stage Loop" "Workflow includes engineering loop"
doc_assert_contains "$workflow_doc" "## QA Stage Loop" "Workflow includes QA loop"
doc_assert_contains "$workflow_doc" "## Shutdown Routine" "Workflow includes shutdown routine"
doc_assert_contains "$workflow_doc" "scripts/launch_stage.sh engineering" "Workflow references engineering launcher"
doc_assert_contains "$workflow_doc" "scripts/launch_stage.sh qa" "Workflow references QA launcher"
doc_assert_contains "$workflow_doc" "prompts/active/next-agent-seed-prompt.md" "Workflow references engineering prompt"
doc_assert_contains "$workflow_doc" "prompts/active/qa-agent-seed-prompt.md" "Workflow references QA prompt"
doc_assert_contains "$workflow_doc" "If engineering launch returns" "Workflow includes empty backlog instruction"
doc_assert_contains "$workflow_doc" "no stories" "Workflow includes explicit no-stories token"
doc_assert_contains "$workflow_doc" "If QA finds blocking defects" "Workflow includes QA failure instruction"
doc_assert_contains "$workflow_doc" "branch is not" "Workflow includes branch discipline"
doc_assert_contains "$workflow_doc" "dev" "Workflow includes explicit dev branch token"
doc_assert_contains "$workflow_doc" "command escalation" "Workflow includes escalation rules"

doc_test_finish
