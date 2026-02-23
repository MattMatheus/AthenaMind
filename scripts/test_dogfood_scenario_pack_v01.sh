#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"
pack="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_PACK_V01.md"
run="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22.md"
follow_on_story_id="STORY-20260222-dogfood-semantic-retrieval-hardening-v01"

doc_test_init

doc_assert_exists "$pack" "Dogfood scenario pack exists"
doc_assert_exists "$run" "Dogfood first run artifact exists"
follow_on_path="$(rg --files "$root_dir/backlog/engineering" | rg "/${follow_on_story_id}\\.md$" || true)"
if [[ -z "$follow_on_path" ]]; then
  doc_assert_contains "$pack" "$follow_on_story_id" "Prioritized follow-on story is referenced in canonical artifacts"
else
  doc_assert_contains "$follow_on_path" "precision_at_3" "Follow-on story ties to KPI delta"
fi

doc_assert_contains "$pack" "pack_version" "Scenario pack includes explicit version metadata"
doc_assert_contains "$pack" "SCN-PROC-01" "Scenario pack includes procedural scenario"
doc_assert_contains "$pack" "SCN-STATE-01" "Scenario pack includes state scenario"
doc_assert_contains "$pack" "SCN-SEM-01" "Scenario pack includes semantic scenario"
doc_assert_contains "$pack" "Scoring Loop (Repeatable)" "Scenario pack defines repeatable scoring loop"

doc_assert_contains "$run" "KPI-Relevant Snapshot Annotations" "First run includes KPI annotations"
doc_assert_contains "$run" "Failure Classification Summary" "First run classifies failures"
doc_assert_contains "$run" "Prioritized Follow-On Action" "First run records prioritized follow-on action"
doc_assert_contains "$run" "$follow_on_story_id" "First run references prioritized follow-on story id"

doc_test_finish
