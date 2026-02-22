#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

cycle_index="$root_dir/CYCLE_INDEX.md"

doc_test_init

doc_assert_exists "$cycle_index" "Cycle index exists at repo root"
doc_assert_contains "$cycle_index" "## First 5 Minutes" "Cycle index includes first 5 minutes section"
doc_assert_contains "$cycle_index" "scripts/launch_stage.sh" "Cycle index links stage launcher script"
doc_assert_contains "$cycle_index" "prompts/active/next-agent-seed-prompt.md" "Cycle index links engineering prompt"
doc_assert_contains "$cycle_index" "prompts/active/qa-agent-seed-prompt.md" "Cycle index links QA prompt"
doc_assert_contains "$cycle_index" "prompts/active/pm-refinement-seed-prompt.md" "Cycle index links PM prompt"
doc_assert_contains "$cycle_index" "backlog/engineering/active/" "Cycle index links engineering backlog states"
doc_assert_contains "$cycle_index" "backlog/architecture/active/" "Cycle index links architecture backlog states"
doc_assert_contains "$cycle_index" "personas/STAFF_DIRECTORY.md" "Cycle index links staff directory"
doc_assert_contains "$cycle_index" "research/handoff.md" "Cycle index links handoff docs"
doc_assert_contains "$cycle_index" "no stories" "Cycle index includes no-stories behavior"
doc_assert_contains "$cycle_index" "expected 'dev'" "Cycle index includes branch safety rule"

doc_test_finish
