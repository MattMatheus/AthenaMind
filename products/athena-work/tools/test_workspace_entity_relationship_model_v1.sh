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

good_state = {
    "state_version": "test-v1",
    "next_story": "STORY-A",
    "engineering": {"active": ["STORY-A"], "qa": [], "done_count": 0},
    "architecture": {"active": ["ARCH-A"], "qa_count": 0},
    "blockers": ["none"],
    "timeline": [
        {
            "event_type": "transition",
            "result": "accepted",
            "label": "active_to_qa",
            "story_id": "STORY-A",
            "cycle_id": "STORY-A",
            "correlation_id": "corr-1",
            "timestamp": "2026-02-27T00:00:00Z",
        }
    ],
}
model = module.build_workspace_entity_model(good_state)
assert model["relationship_errors"] == [], model["relationship_errors"]
assert model["relationships"]["next_story"] == "STORY-A"
assert model["relationships"]["engineering_active_cards"] == ["STORY-A"]
assert model["relationships"]["architecture_active_cards"] == ["ARCH-A"]

bad_state = {
    "state_version": "test-v1-bad",
    "next_story": "STORY-MISSING",
    "engineering": {"active": ["STORY-A"], "qa": ["STORY-A"], "done_count": 0},
    "architecture": {"active": ["STORY-A"], "qa_count": 0},
    "timeline": [
        {
            "event_type": "transition",
            "result": "accepted",
            "label": "active_to_qa",
            "story_id": "STORY-A",
            "cycle_id": "STORY-A",
            "correlation_id": "",
            "timestamp": "",
        }
    ],
}
bad_model = module.build_workspace_entity_model(bad_state)
errors = "\n".join(bad_model["relationship_errors"])
assert "card_in_multiple_engineering_lanes" in errors
assert "card_crosses_engineering_architecture_lanes" in errors
assert "next_story_not_in_active_or_qa" in errors
assert "timeline_missing_correlation_id" in errors
assert "timeline_missing_timestamp" in errors
PY
then
  echo "PASS: canonical workspace entity model resolves consistent relationships"
  echo "PASS: relationship drift checks catch contradictory lane/timeline state"
  echo "Result: PASS"
else
  echo "FAIL: workspace entity model regression checks"
  echo "Result: FAIL"
  exit 1
fi
