#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"
pack="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_PACK_V01.md"
run="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22.md"
follow_on="$root_dir/backlog/engineering/intake/STORY-20260222-dogfood-semantic-retrieval-hardening-v01.md"

doc_test_init

doc_assert_exists "$pack" "Dogfood scenario pack exists"
doc_assert_exists "$run" "Dogfood first run artifact exists"
doc_assert_exists "$follow_on" "Prioritized follow-on intake story exists"

doc_assert_contains "$pack" "pack_version" "Scenario pack includes explicit version metadata"
doc_assert_contains "$pack" "SCN-PROC-01" "Scenario pack includes procedural scenario"
doc_assert_contains "$pack" "SCN-STATE-01" "Scenario pack includes state scenario"
doc_assert_contains "$pack" "SCN-SEM-01" "Scenario pack includes semantic scenario"
doc_assert_contains "$pack" "Scoring Loop (Repeatable)" "Scenario pack defines repeatable scoring loop"

doc_assert_contains "$run" "KPI-Relevant Snapshot Annotations" "First run includes KPI annotations"
doc_assert_contains "$run" "Failure Classification Summary" "First run classifies failures"
doc_assert_contains "$run" "Prioritized Follow-On Action" "First run records prioritized follow-on action"

doc_assert_contains "$follow_on" "status" "Follow-on story is tracked in intake"
doc_assert_contains "$follow_on" "precision_at_3" "Follow-on story ties to KPI delta"

doc_test_finish
