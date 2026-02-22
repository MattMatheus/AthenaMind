#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

doc_test_init

doc_assert_exists "$root_dir/docs/process/STAGE_EXIT_GATES.md" "Stage exit gates doc exists"
doc_assert_exists "$root_dir/docs/process/PROGRAM_OPERATING_SYSTEM.md" "Program operating system doc exists"
doc_assert_exists "$root_dir/research/roadmap/PROGRAM_STATE_BOARD.md" "Program state board exists"
doc_assert_exists "$root_dir/work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md" "Release bundle template exists"

doc_assert_contains "$root_dir/prompts/active/planning-seed-prompt.md" "traceability metadata" "Planning prompt enforces traceability metadata"
doc_assert_contains "$root_dir/prompts/active/architect-agent-seed-prompt.md" "follow-on implementation story" "Architect prompt enforces follow-on mapping"
doc_assert_contains "$root_dir/prompts/active/pm-refinement-seed-prompt.md" "PROGRAM_STATE_BOARD.md" "PM prompt enforces program board updates"
doc_assert_contains "$root_dir/prompts/active/qa-agent-seed-prompt.md" "release-checkpoint readiness note" "QA prompt enforces release-readiness note"

doc_assert_contains "$root_dir/backlog/engineering/intake/STORY_TEMPLATE.md" '`idea_id`' "Engineering story template includes idea traceability"
doc_assert_contains "$root_dir/backlog/engineering/intake/STORY_TEMPLATE.md" '`adr_refs`' "Engineering story template includes ADR references"
doc_assert_contains "$root_dir/backlog/architecture/intake/ARCH_STORY_TEMPLATE.md" '`phase`' "Architecture story template includes phase"
doc_assert_contains "$root_dir/backlog/engineering/intake/BUG_TEMPLATE.md" '`impact_metric`' "Bug template includes impact metric"

doc_test_finish
