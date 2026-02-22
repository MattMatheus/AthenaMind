#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

doc_test_init

doc_assert_exists "$root_dir/scripts/run_observer_cycle.sh" "Observer script exists"
doc_assert_exists "$root_dir/work-system/observer/README.md" "Observer README exists"
doc_assert_exists "$root_dir/work-system/observer/OBSERVER_REPORT_TEMPLATE.md" "Observer report template exists"

doc_assert_contains "$root_dir/DEVELOPMENT_CYCLE.md" "cycle-<cycle-id>" "Development cycle defines cycle commit format"
doc_assert_contains "$root_dir/DEVELOPMENT_CYCLE.md" "run_observer_cycle.sh" "Development cycle references observer command"
doc_assert_contains "$root_dir/docs/process/STAGE_EXIT_GATES.md" "Cycle Closure Gate (Observer + Commit)" "Stage exits include cycle closure gate"
doc_assert_contains "$root_dir/prompts/active/qa-agent-seed-prompt.md" "run_observer_cycle.sh" "QA prompt requires observer"
doc_assert_contains "$root_dir/prompts/active/cycle-seed-prompt.md" "one commit per completed cycle" "Cycle prompt enforces single cycle commit"

doc_test_finish
