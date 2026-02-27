#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
cd "$root_dir"

strict=0

usage() {
  cat <<'USAGE'
usage: tools/workspace_status.sh [--strict]

Human-first workspace status summary for AthenaWork.

Outputs:
- branch safety
- git delta summary
- lane counts
- active queue view
- workspace drift warnings

Flags:
  --strict   return non-zero when drift warnings are present
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --strict)
      strict=1
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown arg '$1'" >&2
      usage
      exit 1
      ;;
  esac
done

if ! git rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "error: repository context required"
  exit 1
fi

required_branch="${ATHENA_REQUIRED_BRANCH:-dev}"
current_branch="$(git branch --show-current)"
branch_status="ok"
if [[ "$current_branch" != "$required_branch" ]]; then
  branch_status="warn"
fi

git_status="$(git status --porcelain)"
changed_count="$(printf '%s\n' "$git_status" | awk 'NF' | wc -l | tr -d ' ')"
untracked_count="$(printf '%s\n' "$git_status" | awk '/^\?\?/{c++} END{print c+0}')"
staged_count="$(printf '%s\n' "$git_status" | awk 'substr($0,1,1) ~ /[A-Z]/{c++} END{print c+0}')"
unstaged_count="$(printf '%s\n' "$git_status" | awk 'substr($0,2,1) ~ /[A-Z]/{c++} END{print c+0}')"

count_lane_files() {
  local lane_dir="$1"
  if [[ ! -d "$lane_dir" ]]; then
    printf 'missing'
    return 0
  fi
  find "$lane_dir" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' ! -name '*TEMPLATE*' | wc -l | tr -d ' '
}

engineering_root="$root_dir/delivery-backlog/engineering"
architecture_root="$root_dir/delivery-backlog/architecture"
active_readme="$engineering_root/active/README.md"

drift_count=0
warn() {
  drift_count=$((drift_count + 1))
  echo "WARN: $1"
}

echo "Workspace Status"
echo "root: $root_dir"
echo "branch: $current_branch (required: $required_branch, status: $branch_status)"
echo "git: changed=$changed_count staged=$staged_count unstaged=$unstaged_count untracked=$untracked_count"
echo

echo "Lane Counts"
for state in intake ready active qa blocked done archive; do
  value="$(count_lane_files "$engineering_root/$state")"
  echo "engineering/$state: $value"
done
for state in intake ready active qa blocked done archive; do
  value="$(count_lane_files "$architecture_root/$state")"
  echo "architecture/$state: $value"
done
echo

echo "Active Queue (engineering)"
active_story_count=0
if [[ -d "$engineering_root/active" ]]; then
  active_story_count="$(find "$engineering_root/active" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | wc -l | tr -d ' ')"
fi
echo "story_files: $active_story_count"

if [[ -f "$active_readme" ]]; then
  if grep -Fq '\n' "$active_readme"; then
    warn "active README appears to contain literal '\\n' sequences"
  fi
  if grep -qi 'No active stories' "$active_readme" && [[ "$active_story_count" -gt 0 ]]; then
    warn "active README says 'No active stories' but story files exist"
  fi

  queue_refs=()
  while IFS= read -r row; do
    [[ -n "$row" ]] && queue_refs+=("$row")
  done < <(sed -En 's/^[[:space:]]*[0-9]+\.[[:space:]]*`([^`]+)`.*/\1/p' "$active_readme")

  if [[ "${#queue_refs[@]}" -eq 0 && "$active_story_count" -gt 0 ]]; then
    warn "active README has no explicit ordered queue while active stories exist"
  fi

  if [[ "${#queue_refs[@]}" -gt 0 ]]; then
    for i in "${!queue_refs[@]}"; do
      ref="${queue_refs[$i]}"
      resolved="$ref"
      if [[ "${resolved#/}" == "$resolved" ]]; then
        if [[ "$resolved" == */* ]]; then
          resolved="$root_dir/$resolved"
        else
          resolved="$engineering_root/active/$resolved"
        fi
      fi
      marker="ok"
      if [[ ! -f "$resolved" ]]; then
        marker="missing"
        warn "queue entry references missing file: $ref"
      fi
      if [[ "$i" -lt 5 ]]; then
        echo "$((i + 1)). $ref [$marker]"
      fi
    done
  fi
else
  warn "missing active queue readme at delivery-backlog/engineering/active/README.md"
fi
echo

observer_dir="$root_dir/operating-system/observer"
resume_context="$observer_dir/RESUME_CONTEXT.md"
if compgen -G "$observer_dir/OBSERVER-REPORT-*.md" >/dev/null && [[ ! -f "$resume_context" ]]; then
  warn "observer reports exist but resume context is missing"
fi

if [[ "$branch_status" != "ok" ]]; then
  warn "branch safety violation: launcher default expects '$required_branch'"
fi

if [[ "$drift_count" -eq 0 ]]; then
  echo "status: OK"
else
  echo "status: WARN ($drift_count drift warnings)"
fi

if [[ "$strict" -eq 1 && "$drift_count" -gt 0 ]]; then
  exit 1
fi
