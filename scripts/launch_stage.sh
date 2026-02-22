#!/usr/bin/env bash
set -euo pipefail

stage="${1:-engineering}"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
active_dir="$root_dir/backlog/active"
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

case "$stage" in
  engineering)
    count=$(find "$active_dir" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | wc -l | tr -d ' ')
    if [ "$count" -eq 0 ]; then
      echo "no stories"
      exit 0
    fi
    echo "launch: prompts/active/next-agent-seed-prompt.md"
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
