#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
default_root="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

package_path=""
root_dir="$default_root"

while [[ $# -gt 0 ]]; do
  case "$1" in
    --package)
      package_path="${2:-}"
      shift 2
      ;;
    --root)
      root_dir="${2:-}"
      shift 2
      ;;
    *)
      echo "usage: $0 --package <launch-authz.json> [--root <repo-root>]" >&2
      exit 2
      ;;
  esac
done

if [[ -z "$package_path" ]]; then
  echo "error: --package is required" >&2
  exit 2
fi

if [[ ! -f "$package_path" ]]; then
  echo "error: package file not found: $package_path" >&2
  exit 2
fi

python3 - "$package_path" <<'PY'
import json
import sys
from pathlib import Path

path = Path(sys.argv[1])
payload = json.loads(path.read_text(encoding="utf-8"))

required_top = [
    "schema_version",
    "created_at_utc",
    "status",
    "commit",
    "queue_summary",
    "gates",
    "required_confirmation_markers",
    "release_bundle_path",
    "blockers",
]
for key in required_top:
    if key not in payload:
        print(f"FAIL: missing required field: {key}")
        raise SystemExit(1)

if payload["required_confirmation_markers"].get("authorization_phrase") != "LAUNCH AUTHORIZED":
    print("FAIL: authorization phrase marker mismatch")
    raise SystemExit(1)

if not isinstance(payload.get("blockers"), list):
    print("FAIL: blockers must be a list")
    raise SystemExit(1)

status = payload.get("status")
if status not in {"ready", "blocked"}:
    print("FAIL: status must be ready|blocked")
    raise SystemExit(1)

print("PASS: launch authorization package structure is valid")
PY
