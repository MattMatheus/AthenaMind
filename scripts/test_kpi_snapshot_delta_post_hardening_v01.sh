#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

delta_snapshot="$root_dir/work-system/metrics/KPI_SNAPSHOT_2026-02-22_DELTA_POST_HARDENING.md"
baseline_snapshot="$root_dir/work-system/metrics/KPI_SNAPSHOT_2026-02-22_BASELINE.md"
hardening_run="$root_dir/work-system/metrics/DOGFOOD_SCENARIO_RUN_2026-02-22-HARDENING.md"

doc_test_init

doc_assert_exists "$delta_snapshot" "Post-hardening KPI delta snapshot exists"
doc_assert_exists "$baseline_snapshot" "Baseline KPI snapshot exists for before/after comparison"
doc_assert_exists "$hardening_run" "Hardening run evidence exists"

doc_assert_contains "$delta_snapshot" "Before/After Deltas (Impacted Metrics)" "Delta snapshot includes before/after section"
doc_assert_contains "$delta_snapshot" "66.7%" "Delta snapshot preserves baseline precision reference"
doc_assert_contains "$delta_snapshot" "100%" "Delta snapshot records post-hardening metric values"
doc_assert_contains "$delta_snapshot" "+33.3pp" "Delta snapshot records precision delta"
doc_assert_contains "$delta_snapshot" "+25pp" "Delta snapshot records trace completeness delta"
doc_assert_contains "$delta_snapshot" "Updated ADR-0008 Band Interpretation" "Delta snapshot includes updated target-band interpretation"
doc_assert_contains "$delta_snapshot" 'Post-hardening: `Green`' "Delta snapshot records Green-band interpretation after hardening"
doc_assert_contains "$delta_snapshot" "DOGFOOD_SCENARIO_RUN_2026-02-22-HARDENING.md" "Delta snapshot references hardening evidence artifact"

doc_test_finish
