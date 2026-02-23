#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

hardening_run="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22-HARDENING.md"
baseline_run="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22.md"

doc_test_init

doc_assert_exists "$hardening_run" "Semantic hardening run artifact exists"
doc_assert_exists "$baseline_run" "Baseline run artifact exists for before/after comparison"

doc_assert_contains "$hardening_run" "SCN-SEM-01" "Hardening run includes semantic scenario evidence"
doc_assert_contains "$hardening_run" "precision_at_3" "Hardening run records improved precision at 3"
doc_assert_contains "$hardening_run" "3/3 = 100%" "Hardening run captures 100% precision at 3"
doc_assert_contains "$hardening_run" "trace_completeness_rate" "Hardening run records trace completeness"
doc_assert_contains "$hardening_run" "4/4 = 100%" "Hardening run records trace completeness at or above 95%"
doc_assert_contains "$hardening_run" "wrong_memory_recall_rate" "Hardening run records wrong-memory metric"
doc_assert_contains "$hardening_run" "0/4 = 0%" "Hardening run records no materially incorrect recall"
doc_assert_contains "$hardening_run" "Before/After Delta vs Baseline Run" "Hardening run includes before/after KPI delta section"
doc_assert_contains "$hardening_run" "66.7% -> 100%" "Hardening run preserves baseline-to-current precision delta"
doc_assert_contains "$hardening_run" "75% -> 100%" "Hardening run preserves baseline-to-current trace delta"

doc_test_finish
