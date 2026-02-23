#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
cd "$root_dir"
memory_cli_bin="${MEMORY_CLI_BIN:-memory-cli}"
memory_root="${ATHENA_MEMORY_ROOT:-$root_dir/memory}"
repo_id="${ATHENA_REPO_ID:-$(basename "$root_dir")}"

cycle_id=""
story_path=""
out_path=""

infer_policy_stage() {
  local cycle_value="$1"
  local story_value="$2"
  local lower_cycle
  local lower_story

  lower_cycle="$(printf '%s' "$cycle_value" | tr '[:upper:]' '[:lower:]')"
  lower_story="$(printf '%s' "$story_value" | tr '[:upper:]' '[:lower:]')"

  if [[ "$lower_story" == *"/architecture/"* ]] || [[ "$lower_cycle" == arch-* ]] || [[ "$lower_cycle" == *"architect"* ]]; then
    echo "architect"
    return 0
  fi
  if [[ "$lower_cycle" == plan-* ]] || [[ "$lower_cycle" == *"planning"* ]]; then
    echo "planning"
    return 0
  fi
  echo "pm"
}

emit_episode_writeback() {
  local files_csv="$1"
  local story_id="$2"
  local summary_file="$3"
  local decisions_file="$4"
  local policy_stage="$5"
  local write_output
  local session_id

  if ! command -v "$memory_cli_bin" >/dev/null 2>&1; then
    echo "warning: episode write-back skipped; '$memory_cli_bin' not found" >&2
    return 0
  fi

  session_id="observer-$cycle_slug"
  if ! write_output="$("$memory_cli_bin" episode write \
    --root "$memory_root" \
    --repo "$repo_id" \
    --session-id "$session_id" \
    --cycle-id "$cycle_id" \
    --story-id "$story_id" \
    --outcome "success" \
    --summary-file "$summary_file" \
    --files-changed "$files_csv" \
    --decisions-file "$decisions_file" \
    --stage "$policy_stage" \
    --reviewer "cycle-observer" \
    --decision "approved" \
    --notes "observer cycle closure write-back" \
    --reason "record cycle outcome for stage bootstrap context" \
    --risk "low risk append-only memory update" 2>&1)"; then
    echo "warning: episode write-back skipped; memory-cli command failed: $write_output" >&2
    return 0
  fi

  echo "memory episode write-back: recorded for cycle '$cycle_id'"
}

usage() {
  cat <<USAGE
usage: tools/run_observer_cycle.sh --cycle-id <id> [--story <path>] [--output <path>]

Generates a deterministic observer report from current git diff.
Default output path:
  operating-system/observer/OBSERVER-REPORT-<cycle-id>.md
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --cycle-id)
      cycle_id="${2:-}"
      shift 2
      ;;
    --story)
      story_path="${2:-}"
      shift 2
      ;;
    --output)
      out_path="${2:-}"
      shift 2
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

if [[ -z "$cycle_id" ]]; then
  echo "error: --cycle-id is required" >&2
  usage
  exit 1
fi

sanitize() {
  printf '%s' "$1" | tr '[:space:]/' '--' | tr -cd '[:alnum:]_.-'
}

cycle_slug="$(sanitize "$cycle_id")"
if [[ -z "$cycle_slug" ]]; then
  echo "error: cycle id produced empty slug" >&2
  exit 1
fi

observer_dir="$root_dir/operating-system/observer"
mkdir -p "$observer_dir"

if [[ -z "$out_path" ]]; then
  out_path="$observer_dir/OBSERVER-REPORT-$cycle_slug.md"
fi

if [[ -n "$story_path" && ! -f "$story_path" ]]; then
  echo "error: --story file not found: $story_path" >&2
  exit 1
fi

branch="$(git branch --show-current)"
generated_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"

tmp_staged="$(mktemp)"
tmp_unstaged="$(mktemp)"
tmp_untracked="$(mktemp)"
tmp_summary="$(mktemp)"
tmp_decisions="$(mktemp)"
trap 'rm -f "$tmp_staged" "$tmp_unstaged" "$tmp_untracked" "$tmp_summary" "$tmp_decisions"' EXIT

git diff --cached --name-status > "$tmp_staged"
git diff --name-status > "$tmp_unstaged"
git ls-files --others --exclude-standard | sed 's/^/A\t/' > "$tmp_untracked"

combined_status="$(cat "$tmp_staged" "$tmp_unstaged" "$tmp_untracked" | awk 'NF' | sort -u || true)"
files_changed_csv="$(printf '%s\n' "$combined_status" | awk -F'\t' 'NF >= 2 {for (i = 2; i <= NF; i++) if ($i != "") print $i}' | sort -u | paste -sd, -)"

idea_id="unknown"
adr_refs="unknown"
if [[ -n "$story_path" ]]; then
  extracted_idea="$(sed -n 's/^[[:space:]]*-[[:space:]]*`idea_id`:[[:space:]]*\(.*\)$/\1/p' "$story_path" | head -n1)"
  extracted_adr="$(sed -n 's/^[[:space:]]*-[[:space:]]*`adr_refs`:[[:space:]]*\(.*\)$/\1/p' "$story_path" | head -n1)"
  [[ -n "$extracted_idea" ]] && idea_id="$extracted_idea"
  [[ -n "$extracted_adr" ]] && adr_refs="$extracted_adr"
fi

{
  echo "# Observer Report: $cycle_id"
  echo
  echo "## Metadata"
  echo "- cycle_id: $cycle_id"
  echo "- generated_at_utc: $generated_at"
  echo "- branch: $branch"
  if [[ -n "$story_path" ]]; then
    rel_story="${story_path#"$root_dir/"}"
    echo "- story_path: $rel_story"
  fi
  echo "- idea_id: $idea_id"
  echo "- adr_refs: $adr_refs"
  echo
  echo "## Diff Inventory"
  if [[ -n "$combined_status" ]]; then
    while IFS= read -r row; do
      [[ -z "$row" ]] && continue
      status="$(printf '%s' "$row" | cut -f1)"
      path="$(printf '%s' "$row" | cut -f2-)"
      echo "- $status $path"
    done <<< "$combined_status"
  else
    echo "- No tracked or untracked file deltas detected."
  fi
  echo
  echo "## Workflow-Sync Checks"
  echo "- [ ] If workflow behavior changed, confirm HUMANS.md, AGENTS.md, and DEVELOPMENT_CYCLE.md were updated."
  echo "- [ ] If prompts changed, confirm corresponding stage docs and gates were updated."
  echo "- [ ] If backlog state changed, confirm queue order and status fields are synchronized."
  echo
  echo "## Memory Promotions"
  echo "- Durable decisions to promote:"
  echo "- New risks/tradeoffs to promote:"
  echo "- Reusable implementation patterns to promote:"
  echo
  echo "## Release Impact"
  echo "- [ ] release_checkpoint impact evaluated for stories touched in this cycle."
  echo "- [ ] If release-bound scope changed, update release bundle inputs."
} > "$out_path"

story_id="$cycle_slug"
story_label="none"
if [[ -n "$story_path" ]]; then
  story_file="$(basename "$story_path")"
  story_id="${story_file%.md}"
  story_label="${story_path#"$root_dir/"}"
fi

cat > "$tmp_summary" <<EOF
Cycle $cycle_id observer report generated on branch $branch at $generated_at.
Story: $story_label
Report: ${out_path#"$root_dir/"}
EOF

cat > "$tmp_decisions" <<EOF
Generated deterministic observer report and wrote cycle outcome for bootstrap reuse.
EOF

policy_stage="$(infer_policy_stage "$cycle_id" "$story_label")"
emit_episode_writeback "$files_changed_csv" "$story_id" "$tmp_summary" "$tmp_decisions" "$policy_stage"

printf '%s\n' "wrote: ${out_path#"$root_dir/"}"
