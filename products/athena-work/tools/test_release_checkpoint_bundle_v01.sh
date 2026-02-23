#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/tools/lib/doc_test_harness.sh"

bundle="$root_dir/operating-system/handoff/RELEASE_BUNDLE_v0.1-initial-2026-02-22.md"
board="$root_dir/product-research/roadmap/PROGRAM_STATE_BOARD.md"

doc_test_init

doc_assert_exists "$bundle" "Release checkpoint bundle exists"
doc_assert_contains "$bundle" "## Decision" "Release bundle includes decision section"
if rg -q '^[[:space:]]*-[[:space:]]*`(ship|hold)`' "$bundle"; then
  echo "PASS: Release bundle records explicit hold/ship decision"
else
  echo "FAIL: Release bundle records explicit hold/ship decision"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi
doc_assert_contains "$bundle" "Included stories" "Release bundle includes scope stories"
doc_assert_contains "$bundle" "QA result artifacts" "Release bundle includes QA evidence"
doc_assert_contains "$bundle" "Validation commands/results" "Release bundle includes validation evidence"
doc_assert_contains "$bundle" "Rollback direction" "Release bundle includes rollback direction"

doc_assert_contains "$board" "RELEASE_BUNDLE_v0.1-initial-2026-02-22.md" "Program board references release checkpoint bundle"

doc_test_finish
