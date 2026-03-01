#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

state_file="${ATHENA_MD_STATE_FILE:-$root_dir/products/athena-work/operating-system/state/backend_read_model_v1.json}"
active_readme="${ATHENA_MD_ACTIVE_README:-$root_dir/products/athena-work/delivery-backlog/engineering/active/README.md}"
active_story_dir="${ATHENA_MD_ACTIVE_STORY_DIR:-$root_dir/products/athena-work/delivery-backlog/engineering/active}"
dry_run=false

usage() {
  cat <<USAGE
usage: tools/check_markdown_drift.sh [--dry-run]

Checks critical markdown drift classes:
  - ordering_conflict
  - missing_artifact
  - stale_revision
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --dry-run)
      dry_run=true
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

if [[ ! -f "$state_file" ]]; then
  echo "drift_class: missing_artifact"
  echo "drift: missing_artifact state file missing: $state_file"
  echo "remediation_id: restore_canonical_state"
  echo "remediation: restore canonical state file then rerun drift guard"
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
expected_active="$tmp_dir/engineering-active-expected.md"

python3 - "$state_file" "$expected_active" <<'PY'
import json
import sys

state_path, out_active = sys.argv[1:]
state = json.load(open(state_path, "r", encoding="utf-8"))
active = state.get("engineering", {}).get("active", [])
state_version = state.get("state_version", "unknown")

with open(out_active, "w", encoding="utf-8") as f:
    f.write("# Engineering Active Queue\n\n")
    f.write("Ordered execution queue for engineering stories.\n")
    f.write(f"<!-- projection_version: {state_version} -->\n\n")
    f.write("## Active Sequence\n")
    for i, story in enumerate(active, start=1):
        f.write(f"{i}. `{story}.md`\n")
PY

failures=0

if [[ ! -f "$active_readme" ]] || ! cmp -s "$active_readme" "$expected_active"; then
  echo "drift_class: ordering_conflict"
  echo "drift: ordering_conflict engineering active queue differs from canonical projection"
  echo "remediation_id: run_markdown_sync_worker"
  echo "remediation: run tools/markdown_sync_worker.sh to re-project queue order"
  failures=$((failures + 1))
fi

projection_version="$(sed -n 's/^<!-- projection_version: \(.*\) -->$/\1/p' "$active_readme" | head -n1)"
state_version="$(python3 -c "import json;print(json.load(open('$state_file', 'r', encoding='utf-8')).get('state_version','unknown'))")"
if [[ -z "$projection_version" || "$projection_version" != "$state_version" ]]; then
  echo "drift_class: stale_revision"
  echo "drift: stale_revision projection_version='$projection_version' state_version='$state_version'"
  echo "remediation_id: refresh_projection_revision"
  echo "remediation: run tools/markdown_sync_worker.sh to refresh projected revision"
  failures=$((failures + 1))
fi

while IFS= read -r story; do
  [[ -z "$story" ]] && continue
  path="$active_story_dir/${story}.md"
  if [[ ! -f "$path" ]]; then
    echo "drift_class: missing_artifact"
    echo "drift: missing_artifact active story markdown missing: ${path#$root_dir/}"
    echo "remediation_id: restore_or_update_state"
    echo "remediation: restore missing artifact or update canonical state before sync"
    failures=$((failures + 1))
  fi
done < <(python3 -c "import json; s=json.load(open('$state_file','r',encoding='utf-8')); print('\n'.join(s.get('engineering',{}).get('active',[])))")

if [[ "$failures" -gt 0 ]]; then
  if "$dry_run"; then
    echo "status: fail"
    echo "action: remediation_required"
    echo "why: critical_drift_detected"
    echo "next: run tools/markdown_sync_worker.sh and re-check"
  fi
  exit 1
fi

echo "status: pass"
echo "action: none"
echo "why: no_critical_drift"
echo "next: continue"
