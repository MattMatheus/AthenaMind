#!/usr/bin/env bash
set -euo pipefail

stage="${1:-engineering}"
root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
active_dir="$root_dir/backlog/active"

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
