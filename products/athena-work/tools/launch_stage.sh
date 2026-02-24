#!/usr/bin/env bash
set -euo pipefail

stage="${1:-engineering}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
active_dir="$root_dir/delivery-backlog/engineering/active"
active_readme="$active_dir/README.md"
arch_active_dir="$root_dir/delivery-backlog/architecture/active"
arch_active_readme="$arch_active_dir/README.md"
required_branch="${ATHENA_REQUIRED_BRANCH:-dev}"
repo_id="${ATHENA_REPO_ID:-$(basename "$root_dir")}"
memory_cli_bin="${MEMORY_CLI_BIN:-memory-cli}"
preflight_mode="${ATHENA_PREFLIGHT_MODE:-enforce}"
resume_context_file="$root_dir/operating-system/observer/RESUME_CONTEXT.md"

default_memory_root() {
  if [ -n "${ATHENA_MEMORY_ROOT:-}" ]; then
    printf '%s\n' "$ATHENA_MEMORY_ROOT"
    return 0
  fi
  if [ "${ATHENA_MEMORY_IN_REPO:-0}" = "1" ]; then
    printf '%s\n' "$root_dir/memory"
    return 0
  fi
  printf '%s\n' "${HOME}/.athena/memory/$repo_id"
}

memory_root="$(default_memory_root)"

abort_or_warn() {
  local message="$1"
  if [ "$preflight_mode" = "warn" ]; then
    echo "warning: $message" >&2
    return 0
  fi
  echo "abort: $message" >&2
  exit 1
}

ensure_resume_context_if_needed() {
  if [ "$stage" = "planning" ]; then
    return 0
  fi

  if ! compgen -G "$root_dir/operating-system/observer/OBSERVER-REPORT-*.md" >/dev/null; then
    return 0
  fi

  if [ ! -f "$resume_context_file" ]; then
    abort_or_warn "resume context missing at operating-system/observer/RESUME_CONTEXT.md; run tools/run_observer_cycle.sh once to seed operator context"
  fi
}

validate_story_readiness() {
  local story_path="$1"
  local lane_name="$2"

  if ! grep -Eq '^## Metadata' "$story_path"; then
    abort_or_warn "$lane_name story '$story_path' is missing '## Metadata' section"
  fi
  if ! grep -Eq '`idea_id`:' "$story_path"; then
    abort_or_warn "$lane_name story '$story_path' is missing metadata field 'idea_id'"
  fi
  if ! grep -Eq '`phase`:' "$story_path"; then
    abort_or_warn "$lane_name story '$story_path' is missing metadata field 'phase'"
  fi
  if ! grep -Eq '`adr_refs`:' "$story_path"; then
    abort_or_warn "$lane_name story '$story_path' is missing metadata field 'adr_refs'"
  fi
  if ! grep -Eq '^## Acceptance Criteria' "$story_path"; then
    abort_or_warn "$lane_name story '$story_path' is missing '## Acceptance Criteria' section"
  fi
}

if ! git -C "$root_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "abort: not a git repository at $root_dir" >&2
  exit 1
fi

current_branch="$(git -C "$root_dir" branch --show-current)"
if [ "$current_branch" != "$required_branch" ]; then
  echo "abort: active branch is '$current_branch'; expected '$required_branch'" >&2
  exit 1
fi

ensure_resume_context_if_needed

select_top_story_from_lane() {
  local lane_dir="$1"
  local lane_readme="$2"
  local candidate

  if [ -f "$lane_readme" ]; then
    while IFS= read -r candidate; do
      [ -z "$candidate" ] && continue
      if [ "${candidate#/}" = "$candidate" ]; then
        if [[ "$candidate" == */* ]]; then
          candidate="$root_dir/$candidate"
        else
          candidate="$lane_dir/$candidate"
        fi
      fi
      if [ -f "$candidate" ]; then
        echo "$candidate"
        return 0
      fi
    done < <(sed -En 's/^[[:space:]]*[0-9]+\.[[:space:]]*`([^`]+)`.*/\1/p' "$lane_readme")
  fi

  find "$lane_dir" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | sort | head -n1
}

emit_memory_bootstrap_context() {
  local bootstrap_output
  local session_id
  local scenario_id

  if ! command -v "$memory_cli_bin" >/dev/null 2>&1; then
    echo "warning: memory bootstrap skipped; '$memory_cli_bin' not found" >&2
    return 0
  fi

  session_id="launch-$stage-$(date -u +"%Y%m%dT%H%M%SZ")"
  scenario_id="$stage"
  if ! bootstrap_output="$("$memory_cli_bin" bootstrap --root "$memory_root" --repo "$repo_id" --session-id "$session_id" --scenario "$scenario_id" 2>&1)"; then
    echo "warning: memory bootstrap skipped; bootstrap command failed: $bootstrap_output" >&2
    return 0
  fi

  echo "memory_bootstrap_context:"
  while IFS= read -r line; do
    echo "  $line"
  done <<< "$bootstrap_output"
}

case "$stage" in
  engineering)
    top_story="$(select_top_story_from_lane "$active_dir" "$active_readme" || true)"
    if [ -z "${top_story:-}" ]; then
      echo "no stories"
      exit 0
    fi

    if [ "${top_story#/}" = "$top_story" ]; then
      if [[ "$top_story" == */* ]]; then
        top_story="$root_dir/$top_story"
      else
        top_story="$active_dir/$top_story"
      fi
    fi

    if [ ! -f "$top_story" ]; then
      echo "abort: top active story not found at '$top_story'" >&2
      exit 1
    fi
    validate_story_readiness "$top_story" "engineering"

    rel_story="${top_story#"$root_dir"/}"
    cat <<EOF
launch: stage-prompts/active/next-agent-seed-prompt.md
cycle: engineering
story: $rel_story
checklist:
  1) read story and clarify open questions
  2) implement required changes
  3) update tests
  4) run tools/run_stage_tests.sh (must pass); add story-specific tests as needed
  5) prepare handoff package
  6) move story to delivery-backlog/engineering/qa
  7) do not commit yet (cycle-level commit after observer step)
EOF
    emit_memory_bootstrap_context
    ;;
  qa)
    cat <<EOF
launch: stage-prompts/active/qa-agent-seed-prompt.md
cycle: qa
checklist:
  1) review story in delivery-backlog/engineering/qa against acceptance criteria
  2) validate tests/regression risk
  3) file defects in delivery-backlog/engineering/intake with P0-P3 if found
  4) move story to delivery-backlog/engineering/done or delivery-backlog/engineering/active
  5) run observer: tools/run_observer_cycle.sh --cycle-id <story-id>
  6) commit once for the full cycle with message: cycle-<cycle-id>
EOF
    emit_memory_bootstrap_context
    ;;
  pm)
    cat <<EOF
launch: stage-prompts/active/pm-refinement-seed-prompt.md
cycle: pm
checklist:
  1) review/refine items from delivery-backlog/engineering/intake
  2) rank and move selected items to delivery-backlog/engineering/active
  3) update delivery-backlog/engineering/active/README.md sequence
  4) update engineering directive only if needed
  5) run observer: tools/run_observer_cycle.sh --cycle-id PM-<date>-<slug>
  6) commit once for the full cycle with message: cycle-<cycle-id>
EOF
    emit_memory_bootstrap_context
    ;;
  planning)
    cat <<EOF
launch: stage-prompts/active/planning-seed-prompt.md
cycle: planning
checklist:
  1) run an interactive idea-generation session with the human operator
  2) capture structured notes in product-research/planning/sessions using the planning template
  3) convert session output into intake stories (engineering and/or architecture) using canonical templates
  4) recommend next stage: architect (for decisions) and/or pm (for prioritization)
  5) run observer: tools/run_observer_cycle.sh --cycle-id <plan-id>
  6) commit once for the full cycle with message: cycle-<cycle-id>
EOF
    emit_memory_bootstrap_context
    ;;
  architect)
    top_arch_story="$(select_top_story_from_lane "$arch_active_dir" "$arch_active_readme" || true)"
    if [ -z "${top_arch_story:-}" ]; then
      echo "no stories"
      exit 0
    fi

    validate_story_readiness "$top_arch_story" "architecture"
    rel_arch_story="${top_arch_story#"$root_dir"/}"
    cat <<EOF
launch: stage-prompts/active/architect-agent-seed-prompt.md
cycle: architect
story: $rel_arch_story
checklist:
  1) read story and clarify architecture decision scope
  2) update ADRs/architecture artifacts
  3) run docs validation tests
  4) prepare handoff package
  5) move story to delivery-backlog/architecture/qa
  6) do not commit yet (cycle-level commit after observer step)
EOF
    emit_memory_bootstrap_context
    ;;
  cycle)
    cat <<EOF
launch: stage-prompts/active/cycle-seed-prompt.md
cycle: engineering+qa loop
loop:
  - run: tools/launch_stage.sh engineering
  - if output is "no stories": stop
  - execute engineering story cycle
  - run: tools/launch_stage.sh qa
  - execute QA cycle
  - run observer: tools/run_observer_cycle.sh --cycle-id <story-id>
  - commit once: cycle-<cycle-id>
  - repeat until active backlog is drained
EOF
    emit_memory_bootstrap_context
    ;;
  *)
    echo "usage: tools/launch_stage.sh [engineering|qa|pm|planning|architect|cycle]" >&2
    exit 1
    ;;
esac
