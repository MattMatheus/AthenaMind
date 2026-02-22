#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$root_dir"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "FAIL: repository context required"
  exit 1
fi

# Files that define stage behavior and process rules.
watched=(
  "DEVELOPMENT_CYCLE.md"
  "scripts/launch_stage.sh"
  "work-system/README.md"
  "prompts/active/next-agent-seed-prompt.md"
  "prompts/active/qa-agent-seed-prompt.md"
  "prompts/active/pm-refinement-seed-prompt.md"
  "prompts/active/cycle-seed-prompt.md"
)

# Files that must be touched when behavior changes.
sync_targets=(
  "HUMANS.md"
  "AGENTS.md"
)

changed_files="$(git status --porcelain | sed -E 's/^.. //')"
if [[ -z "$changed_files" ]]; then
  echo "PASS: no working-tree changes"
  echo "Result: PASS"
  exit 0
fi

watched_changed=0
for f in "${watched[@]}"; do
  if printf '%s\n' "$changed_files" | grep -Fxq "$f"; then
    watched_changed=1
    break
  fi
done

if [[ "$watched_changed" -eq 0 ]]; then
  echo "PASS: no workflow-rule files changed"
  echo "Result: PASS"
  exit 0
fi

sync_touched=0
for f in "${sync_targets[@]}"; do
  if printf '%s\n' "$changed_files" | grep -Fxq "$f"; then
    sync_touched=1
    break
  fi
done

if [[ "$sync_touched" -eq 1 ]]; then
  echo "PASS: workflow-rule change includes HUMANS.md/AGENTS.md sync"
  echo "Result: PASS"
  exit 0
fi

echo "FAIL: workflow-rule files changed without updating HUMANS.md or AGENTS.md"
echo "Changed files:"
printf '%s\n' "$changed_files"
echo "Required sync targets: HUMANS.md or AGENTS.md"
echo "Result: FAIL"
exit 1
