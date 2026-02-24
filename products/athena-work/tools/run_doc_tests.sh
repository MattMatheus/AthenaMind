#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

profile="${ATHENA_DOC_TEST_PROFILE:-auto}"

if [[ "$profile" == "auto" ]]; then
  if [[ -d "$root_dir/product-research" ]]; then
    profile="full"
  else
    profile="slim"
  fi
fi

if [[ "$profile" != "full" && "$profile" != "slim" ]]; then
  echo "error: ATHENA_DOC_TEST_PROFILE must be auto|full|slim" >&2
  exit 1
fi

always_tests=(
  "test_qa_regression_rubric.sh"
  "test_state_transition_checklist.sh"
  "test_launch_stage_readme_queue.sh"
  "test_launch_stage_memory_integration.sh"
  "test_dogfood_scenario_pack_v01.sh"
  "test_dogfood_semantic_hardening_v01.sh"
  "test_kpi_snapshot_delta_post_hardening_v01.sh"
  "test_humans_agents_sync.sh"
  "test_intake_validation.sh"
  "test_validate_intake_duplicate_ids.sh"
  "test_go_toolchain_readiness.sh"
  "test_observer_cycle_policy.sh"
  "test_observer_cycle_memory_integration.sh"
  "test_run_stage_tests_scope.sh"
)

legacy_research_tests=(
  "test_goals_scorecard_v01.sh"
  "test_phased_plan_v01_v03.sh"
  "test_coding_readiness_gate.sh"
  "test_doc_test_harness_standardization.sh"
  "test_readme_v01_alignment.sh"
  "test_docs_navigation_hardening.sh"
  "test_founder_operator_workflow.sh"
  "test_memory_snapshot_plan.sh"
  "test_kpi_snapshot_baseline_v01.sh"
  "test_release_checkpoint_bundle_v01.sh"
  "test_program_state_consistency.sh"
  "test_stage_exit_pipeline.sh"
)

echo "doc test profile: $profile"

for test_script in "${always_tests[@]}"; do
  "$root_dir/tools/$test_script"
done

if [[ "$profile" == "full" ]]; then
  for test_script in "${legacy_research_tests[@]}"; do
    "$root_dir/tools/$test_script"
  done
else
  echo "note: skipping legacy research-coupled tests in slim profile"
  for test_script in "${legacy_research_tests[@]}"; do
    echo "skip: $test_script"
  done
fi
