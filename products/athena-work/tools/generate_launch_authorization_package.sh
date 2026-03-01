#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"

timestamp="$(date -u +"%Y%m%dT%H%M%SZ")"
created_at="$(date -u +"%Y-%m-%dT%H:%M:%SZ")"
short_sha="$(git -C "$root_dir" rev-parse --short HEAD)"
branch="$(git -C "$root_dir" rev-parse --abbrev-ref HEAD)"

program_board="$root_dir/product-research/roadmap/PROGRAM_STATE_BOARD.md"
release_bundle="$root_dir/products/athena-work/operating-system/handoff/RELEASE_BUNDLE_v0.1-initial-2026-02-22.md"
runtime_state_file="${ATHENA_RUNTIME_STATE_FILE:-$root_dir/products/athena-work/operating-system/state/runtime/backend_read_model_v1.local.json}"
state_file="$root_dir/products/athena-work/operating-system/state/backend_read_model_v1.json"

extract_count() {
  local key="$1"
  sed -n "s/^[[:space:]]*-[[:space:]]*\`$key\`:[[:space:]]*\([0-9][0-9]*\).*/\1/p" "$program_board"
}

engineering_active_count="$(extract_count engineering_active_count)"
engineering_intake_count="$(extract_count engineering_intake_count)"
engineering_qa_count="$(extract_count engineering_qa_count)"

if [[ -z "$engineering_active_count" ]]; then engineering_active_count=0; fi
if [[ -z "$engineering_intake_count" ]]; then engineering_intake_count=0; fi
if [[ -z "$engineering_qa_count" ]]; then engineering_qa_count=0; fi

test_doc_status="unknown"
if "$root_dir/products/athena-work/tools/test_workspace_ui_read_only_board_v1.sh" >/dev/null 2>&1; then
  test_doc_status="pass"
else
  test_doc_status="fail"
fi

security_gate_status="not-required"
if "$root_dir/products/athena-work/tools/check_security_nonce_gate.sh" >/dev/null 2>&1; then
  security_gate_status="pass"
else
  security_gate_status="blocked"
fi

confirmation_status="$(
  python3 - "$runtime_state_file" "$state_file" <<'PY'
import json
import sys
from pathlib import Path

runtime_path = Path(sys.argv[1])
default_path = Path(sys.argv[2])
path = runtime_path if runtime_path.exists() else default_path
if not path.exists():
    print("unconfirmed")
    raise SystemExit(0)
payload = json.loads(path.read_text(encoding="utf-8"))
print((payload.get("required_confirmation", {}) or {}).get("status", "unconfirmed"))
PY
)"

queue_readiness="pass"
if [[ "$engineering_active_count" -ne 0 || "$engineering_intake_count" -ne 0 || "$engineering_qa_count" -ne 0 ]]; then
  queue_readiness="blocked"
fi

confirmation_readiness="pass"
if [[ "$confirmation_status" != "confirmed" ]]; then
  confirmation_readiness="blocked"
fi

security_readiness="pass"
if [[ "$security_gate_status" == "blocked" ]]; then
  security_readiness="blocked"
fi

blockers=()
if [[ "$engineering_active_count" -ne 0 ]]; then
  blockers+=("engineering_active_count must be 0 before launch")
fi
if [[ "$engineering_intake_count" -ne 0 ]]; then
  blockers+=("engineering_intake_count must be 0 before launch")
fi
if [[ "$engineering_qa_count" -ne 0 ]]; then
  blockers+=("engineering_qa_count must be 0 before launch")
fi
if [[ "$test_doc_status" != "pass" ]]; then
  blockers+=("workspace UI doc test failed")
fi
if [[ "$security_gate_status" == "blocked" ]]; then
  blockers+=("security nonce gate check failed for protected changes")
fi
if [[ "$confirmation_status" != "confirmed" ]]; then
  blockers+=("direction confirmation status is not confirmed")
fi
if [[ ! -f "$release_bundle" ]]; then
  blockers+=("release bundle file missing")
fi

status="blocked"
if [[ "${#blockers[@]}" -eq 0 ]]; then
  status="ready"
fi

out_dir="$root_dir/products/athena-work/operating-system/handoff"
mkdir -p "$out_dir"
out_file="$out_dir/LAUNCH_AUTHZ_${timestamp}.json"

blockers_json="[]"
if [[ "${#blockers[@]}" -gt 0 ]]; then
  blockers_json="$(printf '%s\n' "${blockers[@]}" | python3 -c 'import json,sys; print(json.dumps([l.rstrip("\n") for l in sys.stdin if l.rstrip("\n")]))')"
fi

python3 - "$out_file" "$created_at" "$short_sha" "$branch" "$status" \
  "$engineering_active_count" "$engineering_intake_count" "$engineering_qa_count" \
  "$test_doc_status" "$security_gate_status" "$release_bundle" "$blockers_json" \
  "$queue_readiness" "$confirmation_readiness" "$security_readiness" <<'PY'
import json
import sys
from pathlib import Path

(
    out_file,
    created_at,
    short_sha,
    branch,
    status,
    eng_active,
    eng_intake,
    eng_qa,
    doc_status,
    sec_status,
    release_bundle,
    blockers_json,
    queue_readiness,
    confirmation_readiness,
    security_readiness,
) = sys.argv[1:]

payload = {
    "schema_version": "1",
    "created_at_utc": created_at,
    "status": status,
    "commit": {"short_sha": short_sha, "branch": branch},
    "queue_summary": {
        "engineering_active_count": int(eng_active),
        "engineering_intake_count": int(eng_intake),
        "engineering_qa_count": int(eng_qa),
    },
    "gates": {
        "workspace_ui_doc_test": doc_status,
        "security_nonce_gate": sec_status,
    },
    "readiness_signals": {
        "queue_readiness": queue_readiness,
        "confirmation_readiness": confirmation_readiness,
        "security_gate_readiness": security_readiness,
    },
    "required_confirmation_markers": {
        "confirmed_by": "required",
        "authorization_phrase": "LAUNCH AUTHORIZED",
        "dev_to_prod_requires_human_approval": True,
    },
    "release_bundle_path": str(Path(release_bundle)),
    "blockers": json.loads(blockers_json),
}
payload["blocker_count"] = len(payload["blockers"])

Path(out_file).write_text(json.dumps(payload, indent=2) + "\n", encoding="utf-8")
print(
    json.dumps(
        {
            "code": "OK",
            "path": str(Path(out_file)),
            "status": status,
            "blockers": payload["blockers"],
            "blocker_count": payload["blocker_count"],
            "readiness_signals": payload["readiness_signals"],
        },
        indent=2,
    )
)
PY
