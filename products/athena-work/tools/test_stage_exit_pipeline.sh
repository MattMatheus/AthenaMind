#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/tools/lib/doc_test_harness.sh"

doc_test_init

doc_assert_exists "$root_dir/knowledge-base/process/STAGE_EXIT_GATES.md" "Stage exit gates doc exists"
doc_assert_exists "$root_dir/knowledge-base/process/PROGRAM_OPERATING_SYSTEM.md" "Program operating system doc exists"
doc_assert_exists "$root_dir/product-research/roadmap/PROGRAM_STATE_BOARD.md" "Program state board exists"
doc_assert_exists "$root_dir/operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md" "Release bundle template exists"
doc_assert_exists "$root_dir/operating-system/observer/OBSERVER_REPORT_TEMPLATE.md" "Observer report template exists"

doc_assert_contains "$root_dir/stage-prompts/active/planning-seed-prompt.md" "traceability metadata" "Planning prompt enforces traceability metadata"
doc_assert_contains "$root_dir/stage-prompts/active/architect-agent-seed-prompt.md" "follow-on implementation story" "Architect prompt enforces follow-on mapping"
doc_assert_contains "$root_dir/stage-prompts/active/pm-refinement-seed-prompt.md" "PROGRAM_STATE_BOARD.md" "PM prompt enforces program board updates"
doc_assert_contains "$root_dir/stage-prompts/active/qa-agent-seed-prompt.md" "release-checkpoint readiness note" "QA prompt enforces release-readiness note"
doc_assert_contains "$root_dir/stage-prompts/active/qa-agent-seed-prompt.md" "run_observer_cycle.sh" "QA prompt enforces observer step"
doc_assert_contains "$root_dir/knowledge-base/process/STAGE_EXIT_GATES.md" "Cycle Closure Gate (Observer + Commit)" "Stage exits include observer closure gate"

doc_assert_contains "$root_dir/delivery-backlog/engineering/intake/STORY_TEMPLATE.md" '`idea_id`' "Engineering story template includes idea traceability"
doc_assert_contains "$root_dir/delivery-backlog/engineering/intake/STORY_TEMPLATE.md" '`adr_refs`' "Engineering story template includes ADR references"
doc_assert_contains "$root_dir/delivery-backlog/architecture/intake/ARCH_STORY_TEMPLATE.md" '`phase`' "Architecture story template includes phase"
doc_assert_contains "$root_dir/delivery-backlog/engineering/intake/BUG_TEMPLATE.md" '`impact_metric`' "Bug template includes impact metric"

doc_test_finish
