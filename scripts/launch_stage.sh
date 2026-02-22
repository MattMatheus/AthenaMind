#!/usr/bin/env bash
set -euo pipefail

stage="${1:-engineering}"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
active_dir="$root_dir/backlog/active"
active_readme="$active_dir/README.md"
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

select_top_active_story() {
  local from_readme
  if [ -f "$active_readme" ]; then
    from_readme="$(sed -n 's/^[[:space:]]*[0-9]\+\.[[:space:]]*`\([^`]\+\)`.*/\1/p' "$active_readme" | head -n1 || true)"
    if [ -n "$from_readme" ]; then
      echo "$from_readme"
      return 0
    fi
  fi

  find "$active_dir" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | sort | head -n1
}

case "$stage" in
  engineering)
    count=$(find "$active_dir" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | wc -l | tr -d ' ')
    if [ "$count" -eq 0 ]; then
      echo "no stories"
      exit 0
    fi

    top_story="$(select_top_active_story)"
    if [ -z "${top_story:-}" ]; then
      echo "abort: unable to determine top active story" >&2
      exit 1
    fi

    if [ "${top_story#/}" = "$top_story" ]; then
      top_story="$root_dir/$top_story"
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
  7) move story to backlog/qa
EOF
    ;;
  qa)
    echo "launch: prompts/active/qa-agent-seed-prompt.md"
    ;;
  pm)
    echo "launch: prompts/active/pm-refinement-seed-prompt.md"
    ;;
  *)
    echo "usage: scripts/launch_stage.sh [engineering|qa|pm]" >&2
    exit 1
    ;;
esac
