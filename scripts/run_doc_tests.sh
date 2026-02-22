#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"

"$root_dir/scripts/test_goals_scorecard_v01.sh"
"$root_dir/scripts/test_phased_plan_v01_v03.sh"
"$root_dir/scripts/test_coding_readiness_gate.sh"
"$root_dir/scripts/test_doc_test_harness_standardization.sh"
"$root_dir/scripts/test_docs_navigation_hardening.sh"
"$root_dir/scripts/test_founder_operator_workflow.sh"
