#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

tool="$root_dir/products/athena-work/tools/ingest_artifact_bundle.sh"

doc_test_init
doc_assert_exists "$tool" "Artifact ingest tool exists"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

repo_root="$tmp_dir/repo"
mkdir -p "$repo_root/products/athena-work/operating-system/state"

cat >"$tmp_dir/backend_read_model_v1.json" <<'JSON'
{
  "schema_version": "1",
  "current_stage": "engineering",
  "next_story": "STORY-TEST-001",
  "engineering": {
    "active": ["STORY-TEST-001"],
    "qa": [],
    "done_count": 10
  },
  "architecture": {
    "active": ["ARCH-TEST-001"],
    "qa_count": 1
  },
  "timeline": [
    {
      "event_type": "transition",
      "result": "accepted",
      "label": "intake_to_active",
      "story_id": "STORY-TEST-001",
      "cycle_id": "STORY-TEST-001",
      "correlation_id": "corr-test-1",
      "timestamp": "2026-02-27T00:00:00Z"
    }
  ]
}
JSON

python3 - "$tmp_dir/backend_read_model_v1.json" "$tmp_dir/bundle.zip" <<'PY'
import sys
import zipfile
from pathlib import Path

source = Path(sys.argv[1])
out_zip = Path(sys.argv[2])
with zipfile.ZipFile(out_zip, "w", compression=zipfile.ZIP_DEFLATED) as zf:
    zf.write(source, arcname="backend_read_model_v1.json")
PY

report_file="$tmp_dir/ingest-report.json"
"$tool" --zip "$tmp_dir/bundle.zip" --root "$repo_root" >"$report_file"

doc_assert_contains "$report_file" "\"code\": \"OK\"" "Ingest returns OK status"
doc_assert_contains "$report_file" "\"migrated_objects\": 1" "Ingest reports v1->v2 migration"

out_state="$(python3 - "$report_file" <<'PY'
import json
import sys
from pathlib import Path

payload = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
print(payload.get("output_state_file", ""))
PY
)"
doc_assert_exists "$out_state" "Ingest writes runtime state output"

python3 - "$out_state" <<'PY'
import json
import sys
from pathlib import Path

payload = json.loads(Path(sys.argv[1]).read_text(encoding="utf-8"))
assert payload["schema_version"] == "2"
assert payload["current_stage"] == "engineering"
assert payload["next_story"] == "STORY-TEST-001"
assert payload["engineering"]["active"] == ["STORY-TEST-001"]
assert isinstance(payload["timeline"], list) and payload["timeline"]
PY
echo "PASS: Ingested state is normalized and preserves required workflow data"

doc_test_finish
