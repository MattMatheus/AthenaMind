#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
default_root="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

bundle_zip=""
root_dir="$default_root"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --zip)
      bundle_zip="${2:-}"
      shift 2
      ;;
    --root)
      root_dir="${2:-}"
      shift 2
      ;;
    *)
      echo "usage: $0 --zip <artifact-bundle.zip> [--root <repo-root>]" >&2
      exit 2
      ;;
  esac
done

if [[ -z "$bundle_zip" ]]; then
  echo "error: --zip is required" >&2
  exit 2
fi

if [[ ! -f "$bundle_zip" ]]; then
  echo "error: bundle not found: $bundle_zip" >&2
  exit 2
fi

state_file="${ATHENA_RUNTIME_STATE_FILE:-$root_dir/products/athena-work/operating-system/state/runtime/backend_read_model_v1.local.json}"
mkdir -p "$(dirname "$state_file")"

python3 - "$bundle_zip" "$state_file" <<'PY'
from __future__ import annotations

import json
import sys
import zipfile
from datetime import datetime, timezone
from pathlib import Path


bundle_path = Path(sys.argv[1])
state_path = Path(sys.argv[2])


def fail(reason: str) -> None:
    payload = {
        "code": "ERR_ARTIFACT_BUNDLE_INVALID",
        "reason": reason,
        "accepted_objects": 0,
        "rejected_objects": 1,
        "migrated_objects": 0,
    }
    print(json.dumps(payload, indent=2))
    raise SystemExit(1)


def now_utc() -> str:
    return datetime.now(timezone.utc).strftime("%Y-%m-%dT%H:%M:%SZ")


def as_list(value, fallback):
    if isinstance(value, list):
        return value
    return list(fallback)


with zipfile.ZipFile(bundle_path, "r") as zf:
    candidates = [
        name
        for name in zf.namelist()
        if name.endswith("backend_read_model_v1.json") or name.endswith("backend_read_model_v2.json")
    ]
    if not candidates:
        fail("bundle must include backend_read_model_v1.json or backend_read_model_v2.json")

    source_name = sorted(candidates)[0]
    try:
        raw = zf.read(source_name).decode("utf-8")
        payload = json.loads(raw)
    except Exception as exc:  # noqa: BLE001
        fail(f"unable to parse read-model JSON: {exc}")

if not isinstance(payload, dict):
    fail("read-model JSON must be an object")

required = ["current_stage", "next_story", "engineering", "architecture", "timeline"]
missing = [key for key in required if key not in payload]
if missing:
    fail(f"missing required fields: {', '.join(missing)}")

migrated = 0
schema_version = str(payload.get("schema_version", "1"))
if schema_version in {"", "1"}:
    migrated = 1
schema_version = "2"

engineering = payload.get("engineering", {})
if not isinstance(engineering, dict):
    fail("engineering field must be an object")
architecture = payload.get("architecture", {})
if not isinstance(architecture, dict):
    fail("architecture field must be an object")

normalized = {
    "schema_version": schema_version,
    "state_version": datetime.now(timezone.utc).strftime("%Y-%m-%d.%H%M%S"),
    "generated_at": now_utc(),
    "current_stage": str(payload.get("current_stage", "engineering")),
    "next_story": str(payload.get("next_story", "none")),
    "blockers": as_list(payload.get("blockers"), ["none"]),
    "required_confirmation": payload.get(
        "required_confirmation",
        {
            "required": True,
            "status": "unconfirmed",
            "reason": "Direction-changing actions require valid human confirmation.",
        },
    ),
    "engineering": {
        "active": as_list(engineering.get("active"), []),
        "qa": as_list(engineering.get("qa"), []),
        "done_count": int(engineering.get("done_count", 0)),
    },
    "architecture": {
        "active": as_list(architecture.get("active"), []),
        "qa_count": int(architecture.get("qa_count", 0)),
    },
    "drift_alerts": as_list(payload.get("drift_alerts"), ["none"]),
    "research_comm_exception": payload.get(
        "research_comm_exception",
        {
            "active": False,
            "note": "Undocumented agent communication is blocked outside research mode.",
        },
    ),
    "timeline": as_list(payload.get("timeline"), []),
}

state_path.write_text(json.dumps(normalized, indent=2) + "\n", encoding="utf-8")

report = {
    "code": "OK",
    "reason": "artifact bundle ingested",
    "bundle_path": str(bundle_path),
    "source_object": source_name,
    "output_state_file": str(state_path),
    "accepted_objects": 1,
    "rejected_objects": 0,
    "migrated_objects": migrated,
}
print(json.dumps(report, indent=2))
PY
