#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

doc_test_init

program_board="$root_dir/product-research/roadmap/PROGRAM_STATE_BOARD.md"
readiness="$root_dir/product-research/roadmap/CODING_READINESS_DECISION_2026-02-22.md"
plan_session="$root_dir/product-research/planning/sessions/PLAN-20260222-idea-generation-session.md"
research_backlog="$root_dir/product-research/roadmap/RESEARCH_BACKLOG.md"

doc_assert_exists "$program_board" "Program state board exists"
doc_assert_exists "$readiness" "Coding readiness decision exists"
doc_assert_exists "$plan_session" "Planning session exists"

extract_count() {
  local key="$1"
  local value
  value="$(sed -n "s/^[[:space:]]*-[[:space:]]*\`$key\`:[[:space:]]*\([0-9][0-9]*\).*/\1/p" "$program_board")"
  printf '%s' "$value"
}

assert_equal() {
  local expected="$1"
  local actual="$2"
  local label="$3"
  if [[ "$expected" == "$actual" ]]; then
    echo "PASS: $label"
  else
    echo "FAIL: $label (expected=$expected actual=$actual)"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
}

count_matching_files() {
  local dir="$1"
  local include_pattern="$2"
  local exclude_pattern="${3:-}"
  if [[ ! -d "$dir" ]]; then
    printf '0'
    return 0
  fi
  if [[ -n "$exclude_pattern" ]]; then
    find "$dir" -maxdepth 1 -type f -name "$include_pattern" ! -name "$exclude_pattern" | wc -l | tr -d ' '
    return 0
  fi
  find "$dir" -maxdepth 1 -type f -name "$include_pattern" | wc -l | tr -d ' '
}

actual_eng_intake="$(count_matching_files "$root_dir/delivery-backlog/engineering/intake" '*.md' '*TEMPLATE*')"
actual_eng_active="$(count_matching_files "$root_dir/delivery-backlog/engineering/active" '*.md' 'README.md')"
actual_eng_qa="$(count_matching_files "$root_dir/delivery-backlog/engineering/qa" '*.md')"
actual_eng_done="$(count_matching_files "$root_dir/delivery-backlog/engineering/done" 'STORY-*.md')"
actual_arch_intake="$(count_matching_files "$root_dir/delivery-backlog/architecture/intake" '*.md' '*TEMPLATE*')"
actual_arch_active="$(count_matching_files "$root_dir/delivery-backlog/architecture/active" '*.md' 'README.md')"
actual_arch_qa="$(count_matching_files "$root_dir/delivery-backlog/architecture/qa" '*.md')"
actual_arch_done="$(count_matching_files "$root_dir/delivery-backlog/architecture/done" 'ARCH-*.md')"

assert_equal "$(extract_count engineering_intake_count)" "$actual_eng_intake" "Program board engineering intake count matches"
assert_equal "$(extract_count engineering_active_count)" "$actual_eng_active" "Program board engineering active count matches"
assert_equal "$(extract_count engineering_qa_count)" "$actual_eng_qa" "Program board engineering QA count matches"
assert_equal "$(extract_count engineering_done_story_count)" "$actual_eng_done" "Program board engineering done count matches"
assert_equal "$(extract_count architecture_intake_count)" "$actual_arch_intake" "Program board architecture intake count matches"
assert_equal "$(extract_count architecture_active_count)" "$actual_arch_active" "Program board architecture active count matches"
assert_equal "$(extract_count architecture_qa_count)" "$actual_arch_qa" "Program board architecture QA count matches"
assert_equal "$(extract_count architecture_done_story_count)" "$actual_arch_done" "Program board architecture done count matches"

if rg -q '^-[[:space:]]*`GO`' "$readiness"; then
  if rg -q 'delivery-backlog/engineering/active/' "$readiness"; then
    echo "FAIL: GO readiness decision should not reference active blocker paths"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  else
    echo "PASS: GO readiness decision references resolved blocker paths"
  fi
fi

if rg -q '^-[[:space:]]*`status`:[[:space:]]*draft' "$plan_session"; then
  has_downstream_done=0
  while IFS= read -r line; do
    path="$(printf '%s' "$line" | sed -E 's/.*`([^`]+)`.*/\1/')"
    base="$(basename "$path")"
    if [[ -f "$root_dir/delivery-backlog/engineering/done/$base" || -f "$root_dir/delivery-backlog/architecture/done/$base" ]]; then
      has_downstream_done=1
      break
    fi
  done < <(rg -n 'delivery-backlog/(engineering|architecture)/(intake|active)/[^`]+\.md' "$plan_session" | sed 's/^[0-9]*://')

  if [[ "$has_downstream_done" -eq 1 ]]; then
    echo "FAIL: planning session is draft while linked artifacts are already done"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  else
    echo "PASS: planning draft has no downstream done artifacts"
  fi
else
  echo "PASS: planning session status is finalized"
fi

if [[ "$actual_eng_active" -eq 0 ]]; then
  if rg -qi '^## Now' "$program_board" && rg -qi 'PM refinement' "$program_board"; then
    echo "PASS: roadmap now section reflects empty active queue behavior"
  else
    echo "FAIL: roadmap now section missing PM refinement guidance for empty active queue"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
fi

if rg -q '^## Next$' "$research_backlog" && rg -q 'KPI snapshot' "$research_backlog"; then
  echo "PASS: research backlog next section aligned to operating-system control work"
else
  echo "FAIL: research backlog next section missing control-plane follow-on"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi

doc_test_finish
