#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/tools/lib/doc_test_harness.sh"
plan="$root_dir/product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md"
roadmap="$root_dir/product-research/roadmap/RESEARCH_BACKLOG.md"
handoff="$root_dir/product-research/handoff.md"

doc_test_init

if [[ ! -f "$plan" ]]; then
  echo "FAIL: phased plan missing at $plan"
  exit 1
fi

doc_assert_contains "$plan" "## Phase v0.1" "Plan covers phase v0.1"
doc_assert_contains "$plan" "## Phase v0.2" "Plan covers phase v0.2"
doc_assert_contains "$plan" "## Phase v0.3" "Plan covers phase v0.3"
doc_assert_contains "$plan" "### Exit Criteria (Success Gates)" "Each phase defines success gates"
doc_assert_contains "$plan" "### Major Risks" "Each phase defines major risks"
doc_assert_contains "$plan" "ADR-0001" "Plan maps ADR-0001 constraint"
doc_assert_contains "$plan" "ADR-0002" "Plan maps ADR-0002 constraint"
doc_assert_contains "$plan" "ADR-0003" "Plan maps ADR-0003 constraint"
doc_assert_contains "$plan" "ADR-0004" "Plan maps ADR-0004 constraint"
doc_assert_contains "$plan" "ADR-0005" "Plan maps ADR-0005 constraint"
doc_assert_contains "$plan" "ADR-0006" "Plan maps ADR-0006 constraint"
doc_assert_contains "$plan" "ADR-0007" "Plan maps ADR-0007 constraint"

doc_assert_contains "$roadmap" "PHASED_IMPLEMENTATION_PLAN_V01_V03.md" "Roadmap reflects phased plan"
doc_assert_contains "$handoff" "PHASED_IMPLEMENTATION_PLAN_V01_V03.md" "Handoff reflects phased plan"

doc_test_finish
