#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"
checklist="$root_dir/research/roadmap/CODING_READINESS_GATE_CHECKLIST.md"
decision="$root_dir/research/roadmap/CODING_READINESS_DECISION_2026-02-22.md"
path_doc="$root_dir/PRE_CODING_PATH.md"

assert_story_tracked() {
  local story="$1"
  local label="$2"
  if [[ -f "$root_dir/backlog/active/$story" || -f "$root_dir/backlog/qa/$story" || -f "$root_dir/backlog/done/$story" ]]; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
}

doc_test_init

doc_assert_exists "$checklist" "Checklist artifact exists"
doc_assert_exists "$decision" "Decision artifact exists"

doc_assert_contains "$checklist" "Applied Run (2026-02-22)" "Checklist applied run recorded"
doc_assert_contains "$checklist" "Result summary:" "Checklist result summary recorded"
doc_assert_contains "$decision" "NO-GO" "Decision includes explicit no-go outcome"
doc_assert_contains "$decision" "Blockers (Converted and Ranked in Backlog)" "Decision includes blockers section"

doc_assert_contains "$path_doc" "CODING_READINESS_GATE_CHECKLIST.md" "Pre-coding path references checklist artifact"
doc_assert_contains "$path_doc" "CODING_READINESS_DECISION_2026-02-22.md" "Pre-coding path references decision artifact"

assert_story_tracked "STORY-20260222-state-transition-checklist.md" "Blocker story 1 is tracked in backlog"
assert_story_tracked "STORY-20260222-qa-regression-rubric.md" "Blocker story 2 is tracked in backlog"
assert_story_tracked "STORY-20260222-doc-test-harness-standardization.md" "Blocker story 3 is tracked in backlog"
assert_story_tracked "STORY-20260222-founder-operator-workflow.md" "Blocker story 4 is tracked in backlog"
assert_story_tracked "STORY-20260222-docs-navigation-hardening.md" "Blocker story 5 is tracked in backlog"

doc_test_finish
