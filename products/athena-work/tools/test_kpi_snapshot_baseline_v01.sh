#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/tools/lib/doc_test_harness.sh"

template="$root_dir/operating-system/metrics/KPI_SNAPSHOT_TEMPLATE.md"
snapshot="$root_dir/operating-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md"
program_board="$root_dir/product-research/roadmap/PROGRAM_STATE_BOARD.md"
founder_snapshot="$root_dir/product-research/roadmap/FOUNDER_SNAPSHOT.md"

doc_test_init

doc_assert_exists "$template" "KPI template exists"
doc_assert_exists "$snapshot" "KPI baseline snapshot exists"

doc_assert_contains "$snapshot" "## Summary" "KPI baseline includes summary"
doc_assert_contains "$snapshot" "Lead time:" "KPI baseline includes lead time"
doc_assert_contains "$snapshot" "QA turnaround:" "KPI baseline includes QA turnaround"
doc_assert_contains "$snapshot" "Defect escape:" "KPI baseline includes defect escape"
doc_assert_contains "$snapshot" "Reopen rate:" "KPI baseline includes reopen rate"
doc_assert_contains "$snapshot" "Handoff completeness:" "KPI baseline includes handoff completeness"
doc_assert_contains "$snapshot" "Retrieval quality gate pass rate:" "KPI baseline includes retrieval quality gate pass rate"
doc_assert_contains "$snapshot" "Traceability completeness:" "KPI baseline includes traceability completeness"
doc_assert_contains "$snapshot" "ADR-0008" "KPI baseline includes ADR-0008 interpretation"
doc_assert_contains "$snapshot" "## Actions" "KPI baseline includes actions"

doc_assert_contains "$program_board" "operating-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md" "Program board references KPI baseline snapshot"
doc_assert_contains "$founder_snapshot" "operating-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md" "Founder snapshot references KPI baseline snapshot"

doc_test_finish
