#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/tools/lib/doc_test_harness.sh"
adr="$root_dir/product-research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md"
pm_todo="$root_dir/knowledge-base/process/PM-TODO.md"
roadmap="$root_dir/product-research/roadmap/RESEARCH_BACKLOG.md"

doc_test_init

if [[ ! -f "$adr" ]]; then
  echo "FAIL: ADR file missing at $adr"
  exit 1
fi

doc_assert_contains "$adr" "Task continuity success rate" "ADR maps task continuity metric"
doc_assert_contains "$adr" "Rework reduction" "ADR maps rework metric"
doc_assert_contains "$adr" "Trusted autonomy rate" "ADR maps trusted autonomy metric"
doc_assert_contains "$adr" "Memory precision" "ADR maps memory precision metric"
doc_assert_contains "$adr" "Memory error rate" "ADR maps memory error rate metric"
doc_assert_contains "$adr" "Cost predictability" "ADR maps cost predictability metric"
doc_assert_contains "$adr" "Trace completeness" "ADR maps trace completeness metric"
doc_assert_contains "$adr" "Weekly active memory-enabled workflows" "ADR maps adoption metric (weekly active workflows)"
doc_assert_contains "$adr" "Multi-session workflow completion rate" "ADR maps adoption metric (multi-session completion)"

doc_assert_contains "$pm_todo" "ADR-0008-v01-goals-and-scorecard-target-bands.md" "PM TODO references accepted decision"
doc_assert_contains "$roadmap" "ADR-0008" "Roadmap references accepted decision"

doc_test_finish
