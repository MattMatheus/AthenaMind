#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

if ! command -v git >/dev/null 2>&1; then
  echo "error: git is required" >&2
  exit 1
fi

if ! git -C "$root_dir" rev-parse --is-inside-work-tree >/dev/null 2>&1; then
  echo "error: unable to read repository state" >&2
  exit 1
fi

required_branch="${ATHENA_REQUIRED_BRANCH:-dev}"
current_branch="$(git -C "$root_dir" branch --show-current 2>/dev/null || true)"
if [[ -z "$current_branch" ]]; then
  echo "error: unable to determine current branch" >&2
  exit 1
fi

if ! git_status="$(git -C "$root_dir" status --porcelain 2>/dev/null)"; then
  echo "error: unable to collect git status" >&2
  exit 1
fi

changed_count="$(printf '%s\n' "$git_status" | awk 'NF' | wc -l | tr -d ' ')"
staged_count="$(printf '%s\n' "$git_status" | awk 'substr($0,1,1) ~ /[A-Z]/{c++} END{print c+0}')"
unstaged_count="$(printf '%s\n' "$git_status" | awk 'substr($0,2,1) ~ /[A-Z]/{c++} END{print c+0}')"
untracked_count="$(printf '%s\n' "$git_status" | awk '/^\?\?/{c++} END{print c+0}')"

engineering_root="$root_dir/delivery-backlog/engineering"
active_root="$engineering_root/active"
active_readme="$active_root/README.md"

echo "Self Check"
echo "root: $root_dir"
echo "branch: $current_branch"
if [[ "$current_branch" != "$required_branch" ]]; then
  echo "WARN: expected branch '$required_branch'"
fi
echo "git: changed=$changed_count staged=$staged_count unstaged=$unstaged_count untracked=$untracked_count"
echo
echo "Active Queue (engineering)"

queue_refs=()
if [[ -f "$active_readme" ]]; then
  while IFS= read -r row; do
    [[ -n "$row" ]] && queue_refs+=("$row")
  done < <(sed -En 's/^[[:space:]]*[0-9]+\.[[:space:]]*`([^`]+)`.*/\1/p' "$active_readme")
fi

resolved_paths=()
if [[ "${#queue_refs[@]}" -gt 0 ]]; then
  for ref in "${queue_refs[@]}"; do
    resolved="$ref"
    if [[ "${resolved#/}" == "$resolved" ]]; then
      if [[ "$resolved" == */* ]]; then
        resolved="$root_dir/$resolved"
      else
        resolved="$active_root/$resolved"
      fi
    fi
    [[ -f "$resolved" ]] && resolved_paths+=("$resolved")
  done
elif [[ -d "$active_root" ]]; then
  while IFS= read -r f; do
    resolved_paths+=("$f")
  done < <(find "$active_root" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | sort)
fi

story_count="${#resolved_paths[@]}"
echo "active_story_count: $story_count"
if [[ "$story_count" -eq 0 ]]; then
  echo "No active stories"
  exit 0
fi

for path in "${resolved_paths[@]}"; do
  story_id="$(basename "$path")"
  story_title="$(sed -n 's/^# \(.*\)$/\1/p' "$path" | head -n 1)"
  if [[ -z "$story_title" ]]; then
    story_title="(missing title)"
  fi
  echo "- $story_id: $story_title"
done
