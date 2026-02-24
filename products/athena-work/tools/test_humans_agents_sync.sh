#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
cd "$root_dir"

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "FAIL: repository context required"
  exit 1
fi

# Files that define stage behavior and process rules.
watched=(
  "DEVELOPMENT_CYCLE.md"
  "tools/launch_stage.sh"
  "tools/run_observer_cycle.sh"
  "operating-system/README.md"
  "operating-system/observer/README.md"
  "knowledge-base/process/stage-exit-gates.md"
  "knowledge-base/process/program-operating-system.md"
  "knowledge-base/process/README.md"
  "delivery-backlog/engineering/intake/STORY_TEMPLATE.md"
  "delivery-backlog/engineering/intake/BUG_TEMPLATE.md"
  "delivery-backlog/architecture/intake/ARCH_STORY_TEMPLATE.md"
  "stage-prompts/active/planning-seed-prompt.md"
  "stage-prompts/active/next-agent-seed-prompt.md"
  "stage-prompts/active/architect-agent-seed-prompt.md"
  "stage-prompts/active/qa-agent-seed-prompt.md"
  "stage-prompts/active/pm-refinement-seed-prompt.md"
  "stage-prompts/active/cycle-seed-prompt.md"
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
