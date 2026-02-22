#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
plan="$root_dir/research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md"
roadmap="$root_dir/research/roadmap/RESEARCH_BACKLOG.md"
handoff="$root_dir/research/handoff.md"

failures=0

assert_contains() {
  local file="$1"
  local text="$2"
  local label="$3"
  if grep -Fq "$text" "$file"; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    failures=$((failures + 1))
  fi
}

if [[ ! -f "$plan" ]]; then
  echo "FAIL: phased plan missing at $plan"
  exit 1
fi

assert_contains "$plan" "## Phase v0.1" "Plan covers phase v0.1"
assert_contains "$plan" "## Phase v0.2" "Plan covers phase v0.2"
assert_contains "$plan" "## Phase v0.3" "Plan covers phase v0.3"
assert_contains "$plan" "### Exit Criteria (Success Gates)" "Each phase defines success gates"
assert_contains "$plan" "### Major Risks" "Each phase defines major risks"
assert_contains "$plan" "ADR-0001" "Plan maps ADR-0001 constraint"
assert_contains "$plan" "ADR-0002" "Plan maps ADR-0002 constraint"
assert_contains "$plan" "ADR-0003" "Plan maps ADR-0003 constraint"
assert_contains "$plan" "ADR-0004" "Plan maps ADR-0004 constraint"
assert_contains "$plan" "ADR-0005" "Plan maps ADR-0005 constraint"
assert_contains "$plan" "ADR-0006" "Plan maps ADR-0006 constraint"
assert_contains "$plan" "ADR-0007" "Plan maps ADR-0007 constraint"

assert_contains "$roadmap" "PHASED_IMPLEMENTATION_PLAN_V01_V03.md" "Roadmap reflects phased plan"
assert_contains "$handoff" "PHASED_IMPLEMENTATION_PLAN_V01_V03.md" "Handoff reflects phased plan"

if [[ "$failures" -gt 0 ]]; then
  echo "Result: FAIL ($failures checks failed)"
  exit 1
fi

echo "Result: PASS"
