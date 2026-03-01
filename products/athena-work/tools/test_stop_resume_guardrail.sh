#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

doc_test_init

doc_assert_contains "$root_dir/AGENTS.md" "## Stop/Resume Guardrail (Highest Priority)" "AGENTS defines stop/resume guardrail section"
doc_assert_contains "$root_dir/AGENTS.md" "STOP WORK. RESUME TOKEN: <token>" "AGENTS defines exact park trigger"
doc_assert_contains "$root_dir/AGENTS.md" "RESUME WORK. RESUME TOKEN: <token>" "AGENTS defines exact resume trigger"
doc_assert_contains "$root_dir/AGENTS.md" "TOKEN MISMATCH: still parked" "AGENTS defines mismatch behavior"

doc_assert_contains "$root_dir/HUMANS.md" "PARKED" "HUMANS defines parked state behavior"
doc_assert_contains "$root_dir/DEVELOPMENT_CYCLE.md" "TOKEN MISMATCH: still parked" "Development cycle defines mismatch behavior"

for prompt in \
  "$root_dir/stage-prompts/active/planning-seed-prompt.md" \
  "$root_dir/stage-prompts/active/next-agent-seed-prompt.md" \
  "$root_dir/stage-prompts/active/architect-agent-seed-prompt.md" \
  "$root_dir/stage-prompts/active/qa-agent-seed-prompt.md" \
  "$root_dir/stage-prompts/active/pm-refinement-seed-prompt.md" \
  "$root_dir/stage-prompts/active/cycle-seed-prompt.md"; do
  doc_assert_contains "$prompt" "## Stop/Resume Guardrail (Highest Priority)" "$(basename "$prompt") includes guardrail section"
  doc_assert_contains "$prompt" "STOP WORK. RESUME TOKEN: <token>" "$(basename "$prompt") includes park trigger"
  doc_assert_contains "$prompt" "RESUME WORK. RESUME TOKEN: <token>" "$(basename "$prompt") includes resume trigger"
  doc_assert_contains "$prompt" "TOKEN MISMATCH: still parked" "$(basename "$prompt") includes mismatch behavior"
done

doc_assert_contains "$root_dir/tools/launch_stage.sh" "park_resume_guardrail:" "Launcher emits guardrail block"
doc_assert_contains "$root_dir/tools/launch_stage.sh" "TOKEN MISMATCH: still parked" "Launcher guardrail includes mismatch behavior"

doc_test_finish
