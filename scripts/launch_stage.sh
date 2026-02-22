#!/usr/bin/env bash
set -euo pipefail

stage="${1:-engineering}"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
active_dir="$root_dir/backlog/engineering/active"
active_readme="$active_dir/README.md"
arch_active_dir="$root_dir/backlog/architecture/active"
arch_active_readme="$arch_active_dir/README.md"
required_branch="dev"

if ! git -C "$root_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "abort: not a git repository at $root_dir" >&2
  exit 1
fi

current_branch="$(git -C "$root_dir" branch --show-current)"
if [ "$current_branch" != "$required_branch" ]; then
  echo "abort: active branch is '$current_branch'; expected '$required_branch'" >&2
  exit 1
fi

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

    rel_story="${top_story#"$root_dir"/}"
    cat <<EOF
launch: prompts/active/next-agent-seed-prompt.md
cycle: engineering
story: $rel_story
checklist:
  1) read story and clarify open questions
  2) implement required changes
  3) update tests
  4) run tests (must pass)
  5) commit changes with story id in message
  6) prepare handoff package
  7) move story to backlog/engineering/qa
EOF
    ;;
  qa)
    cat <<EOF
launch: prompts/active/qa-agent-seed-prompt.md
cycle: qa
checklist:
  1) review story in backlog/engineering/qa against acceptance criteria
  2) validate tests/regression risk
  3) file defects in backlog/engineering/intake with P0-P3 if found
  4) move story to backlog/engineering/done or backlog/engineering/active
  5) commit QA artifacts and state changes as: qa-<story-id>
EOF
    ;;
  pm)
    cat <<EOF
launch: prompts/active/pm-refinement-seed-prompt.md
cycle: pm
checklist:
  1) review/refine items from backlog/engineering/intake
  2) rank and move selected items to backlog/engineering/active
  3) update backlog/engineering/active/README.md sequence
  4) update engineering directive only if needed
  5) commit refinement outputs and state changes
EOF
    ;;
  architect)
    top_arch_story="$(select_top_story_from_lane "$arch_active_dir" "$arch_active_readme" || true)"
    if [ -z "${top_arch_story:-}" ]; then
      echo "no stories"
      exit 0
    fi

    rel_arch_story="${top_arch_story#"$root_dir"/}"
    cat <<EOF
launch: prompts/active/architect-agent-seed-prompt.md
cycle: architect
story: $rel_arch_story
checklist:
  1) read story and clarify architecture decision scope
  2) update ADRs/architecture artifacts
  3) run docs validation tests
  4) commit changes as: arch-<story-id>
  5) prepare handoff package
  6) move story to backlog/architecture/qa
EOF
    ;;
  cycle)
    cat <<EOF
launch: prompts/active/cycle-seed-prompt.md
cycle: engineering+qa loop
loop:
  - run: scripts/launch_stage.sh engineering
  - if output is "no stories": stop
  - execute engineering story cycle
  - run: scripts/launch_stage.sh qa
  - execute QA cycle
  - repeat until active backlog is drained
EOF
    ;;
  *)
    echo "usage: scripts/launch_stage.sh [engineering|qa|pm|architect|cycle]" >&2
    exit 1
    ;;
esac
