#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
checklist="$root_dir/research/roadmap/CODING_READINESS_GATE_CHECKLIST.md"
decision="$root_dir/research/roadmap/CODING_READINESS_DECISION_2026-02-22.md"
path_doc="$root_dir/PRE_CODING_PATH.md"
active_readme="$root_dir/backlog/active/README.md"

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

assert_exists() {
  local file="$1"
  local label="$2"
  if [[ -f "$file" ]]; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    failures=$((failures + 1))
  fi
}

assert_exists "$checklist" "Checklist artifact exists"
assert_exists "$decision" "Decision artifact exists"

assert_contains "$checklist" "Applied Run (2026-02-22)" "Checklist applied run recorded"
assert_contains "$checklist" "Result summary:" "Checklist result summary recorded"
assert_contains "$decision" "NO-GO" "Decision includes explicit no-go outcome"
assert_contains "$decision" "Blockers (Converted and Ranked in Backlog)" "Decision includes blockers section"

assert_contains "$path_doc" "CODING_READINESS_GATE_CHECKLIST.md" "Pre-coding path references checklist artifact"
assert_contains "$path_doc" "CODING_READINESS_DECISION_2026-02-22.md" "Pre-coding path references decision artifact"

assert_contains "$active_readme" "STORY-20260222-state-transition-checklist.md" "Active queue includes blocker story 1"
assert_contains "$active_readme" "STORY-20260222-qa-regression-rubric.md" "Active queue includes blocker story 2"
assert_contains "$active_readme" "STORY-20260222-doc-test-harness-standardization.md" "Active queue includes blocker story 3"
assert_contains "$active_readme" "STORY-20260222-founder-operator-workflow.md" "Active queue includes blocker story 4"
assert_contains "$active_readme" "STORY-20260222-docs-navigation-hardening.md" "Active queue includes blocker story 5"

if [[ "$failures" -gt 0 ]]; then
  echo "Result: FAIL ($failures checks failed)"
  exit 1
fi

echo "Result: PASS"
