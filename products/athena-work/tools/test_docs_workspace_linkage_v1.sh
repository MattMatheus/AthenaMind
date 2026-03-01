#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
api_file="$root_dir/products/athena-work/ui/local_control_plane_api.py"

if python3 - "$api_file" <<'PY'
import importlib.util
import pathlib
import sys

api_path = pathlib.Path(sys.argv[1])
spec = importlib.util.spec_from_file_location("local_control_plane_api", api_path)
module = importlib.util.module_from_spec(spec)
spec.loader.exec_module(module)

state = {
    "state_version": "docs-link-test",
    "next_story": "none",
    "engineering": {
        "active": [],
        "qa": [
            "STORY-20260227-docs-linkage-and-sync-projection-hardening-v1",
            "STORY-20260227-left-anchored-landscape-layout-v1",
            "STORY-20260227-kanban-entity-relationship-implementation-v1",
        ],
        "done_count": 0,
    },
    "architecture": {"active": [], "qa_count": 0},
    "blockers": ["none"],
    "timeline": [],
}

docs_payload = module.build_docs_index(state)
docs = docs_payload.get("docs", [])
paths = {item.get("path") for item in docs}
ids = {item.get("id") for item in docs}

assert "products/athena-work/delivery-backlog/engineering/active/README.md" in paths
assert "products/athena-work/delivery-backlog/engineering/done/STORY-20260227-docs-linkage-and-sync-projection-hardening-v1.md" in paths
assert "products/athena-work/delivery-backlog/engineering/done/STORY-20260227-left-anchored-landscape-layout-v1.md" in paths
assert "products/athena-work/delivery-backlog/engineering/done/STORY-20260227-kanban-entity-relationship-implementation-v1.md" in paths
assert "engineering-active-queue" in ids
assert "story-story-20260227-docs-linkage-and-sync-projection-hardening-v1" in ids
PY
then
  echo "PASS: docs workspace links are generated from canonical queue/story state"
  echo "Result: PASS"
else
  echo "FAIL: docs workspace canonical linkage regression"
  echo "Result: FAIL"
  exit 1
fi
