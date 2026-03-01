#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

state_file="$root_dir/products/athena-work/operating-system/state/backend_read_model_v1.json"
active_readme="$root_dir/products/athena-work/delivery-backlog/engineering/active/README.md"
observer_board="$root_dir/products/athena-work/operating-system/observer/LATEST_BOARD_READ_MODEL.md"
observer_timeline="$root_dir/products/athena-work/operating-system/observer/LATEST_TIMELINE_READ_MODEL.md"
observer_docs="$root_dir/products/athena-work/operating-system/observer/LATEST_DOCS_INDEX_READ_MODEL.md"
dry_run=false

usage() {
  cat <<USAGE
usage: tools/markdown_sync_worker.sh [--dry-run]

Projects canonical backend read model into markdown artifacts.
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
  echo "error: state file not found: $state_file" >&2
  exit 1
fi

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT
next_active="$tmp_dir/engineering-active-readme.md"
next_board="$tmp_dir/latest-board.md"
next_timeline="$tmp_dir/latest-timeline.md"
next_docs="$tmp_dir/latest-docs.md"

python3 - "$state_file" "$next_active" "$next_board" "$next_timeline" "$next_docs" <<'PY'
import json
import sys
from pathlib import Path

state_path, out_active, out_board, out_timeline, out_docs = sys.argv[1:]
state = json.load(open(state_path, "r", encoding="utf-8"))
active = state.get("engineering", {}).get("active", [])
qa = state.get("engineering", {}).get("qa", [])
state_version = state.get("state_version", "unknown")
next_story = state.get("next_story", "none")

with open(out_active, "w", encoding="utf-8") as f:
    f.write("# Engineering Active Queue\n\n")
    f.write("Ordered execution queue for engineering stories.\n")
    f.write(f"<!-- projection_version: {state_version} -->\n\n")
    f.write("## Active Sequence\n")
    for i, story in enumerate(active, start=1):
        f.write(f"{i}. `{story}.md`\n")

with open(out_board, "w", encoding="utf-8") as f:
    f.write("# Latest Board Read Model\n\n")
    f.write(f"- `state_version`: {state_version}\n")
    f.write(f"- `current_stage`: {state.get('current_stage', 'unknown')}\n")
    f.write(f"- `next_story`: {state.get('next_story', 'none')}\n")
    f.write(f"- `engineering_active_count`: {len(active)}\n")
    f.write(f"- `architecture_active_count`: {len(state.get('architecture', {}).get('active', []))}\n")
    f.write(f"- `drift_alerts`: {', '.join(state.get('drift_alerts', ['none']))}\n")

with open(out_timeline, "w", encoding="utf-8") as f:
    f.write("# Latest Timeline Read Model\n\n")
    f.write(f"- `state_version`: {state_version}\n")
    f.write("\n## Events\n")
    for ev in state.get("timeline", []):
        f.write(
            f"- {ev.get('timestamp','n/a')} | {ev.get('label','n/a')} | "
            f"{ev.get('story_id','n/a')} | correlation_id={ev.get('correlation_id','n/a')}\n"
        )

docs = [
    ("humans", "products/athena-work/HUMANS.md"),
    ("development-cycle", "products/athena-work/DEVELOPMENT_CYCLE.md"),
    ("stage-exit-gates", "knowledge-base/process/STAGE_EXIT_GATES.md"),
    ("local-control-quickstart", "knowledge-base/operations/LOCAL_CONTROL_PLANE_QUICKSTART.md"),
    ("engineering-active-queue", "products/athena-work/delivery-backlog/engineering/active/README.md"),
]
story_ids = []
if next_story and next_story != "none":
    story_ids.append(next_story)
story_ids.extend(active)
story_ids.extend(qa)
seen = set()
for story in story_ids:
    if story in seen:
        continue
    seen.add(story)
    for lane in ("active", "qa", "done"):
        path = Path(f"products/athena-work/delivery-backlog/engineering/{lane}/{story}.md")
        if Path("/workspace").joinpath(path).exists():
            docs.append((f"story-{story.lower()}", str(path)))
            break

with open(out_docs, "w", encoding="utf-8") as f:
    f.write("# Latest Docs Index Read Model\n\n")
    f.write(f"- `state_version`: {state_version}\n")
    f.write("\n## Docs\n")
    for doc_id, path in docs:
        f.write(f"- `{doc_id}` -> `{path}`\n")
PY

if "$dry_run"; then
  status="clean"
  for pair in \
    "$active_readme:$next_active" \
    "$observer_board:$next_board" \
    "$observer_timeline:$next_timeline" \
    "$observer_docs:$next_docs"; do
    current="${pair%%:*}"
    expected="${pair##*:}"
    if [[ ! -f "$current" ]] || ! cmp -s "$current" "$expected"; then
      status="drift"
      echo "drift: $current"
      diff -u "$current" "$expected" || true
    fi
  done
  if [[ "$status" == "clean" ]]; then
    echo "status: clean"
    echo "action: none"
    echo "why: markdown_projection_in_sync"
    echo "next: continue"
  else
    echo "status: drift"
    echo "action: apply_sync"
    echo "why: markdown_projection_out_of_sync"
    echo "next: run tools/markdown_sync_worker.sh"
  fi
  exit 0
fi

mkdir -p "$(dirname "$observer_board")" "$(dirname "$observer_timeline")" "$(dirname "$observer_docs")"
cp "$next_active" "$active_readme"
cp "$next_board" "$observer_board"
cp "$next_timeline" "$observer_timeline"
cp "$next_docs" "$observer_docs"

echo "status: synced"
echo "action: projected_markdown"
echo "why: canonical_state_applied"
echo "next: run tools/check_markdown_drift.sh --dry-run"
