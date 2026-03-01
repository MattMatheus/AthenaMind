#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
contracts_file="$root_dir/products/athena-work/operating-system/contracts/SECURITY_BACKING_CONTRACT_PATHS.txt"

if [[ ! -f "$contracts_file" ]]; then
  echo "FAIL: missing contracts path list: $contracts_file"
  exit 1
fi

base_ref="${ATHENA_SECURITY_GATE_BASE_REF:-origin/main}"
if git -C "$root_dir" rev-parse --verify "$base_ref" >/dev/null 2>&1; then
  diff_range="$base_ref...HEAD"
else
  diff_range="HEAD~1..HEAD"
fi

changed_files="$(git -C "$root_dir" diff --name-only "$diff_range" || true)"
if [[ -z "$changed_files" ]]; then
  echo "PASS: no changed files detected for security gate"
  exit 0
fi

security_change=0
while IFS= read -r protected; do
  [[ -z "$protected" ]] && continue
  if grep -Fxq "$protected" <<<"$changed_files"; then
    security_change=1
    break
  fi
done <"$contracts_file"

if [[ "$security_change" -eq 0 ]]; then
  echo "PASS: no security-backing contract changes detected"
  exit 0
fi

required_vars=(
  ATHENA_SECURITY_CHANGE_NONCE
  ATHENA_SECURITY_NONCE_ISSUED_AT_UTC
  ATHENA_SECURITY_CHANGE_OTP_PRIMARY
  ATHENA_SECURITY_CHANGE_OTP_SECURITY
  ATHENA_SECURITY_CHANGE_ACK
)

for key in "${required_vars[@]}"; do
  if [[ -z "${!key:-}" ]]; then
    echo "FAIL: required gate input missing: $key"
    exit 1
  fi
done

if [[ ! "${ATHENA_SECURITY_CHANGE_NONCE}" =~ ^[A-Z0-9]{12,64}$ ]]; then
  echo "FAIL: ATHENA_SECURITY_CHANGE_NONCE must match [A-Z0-9]{12,64}"
  exit 1
fi

if [[ "${ATHENA_SECURITY_CHANGE_ACK}" != "LAUNCH AUTHORIZED" ]]; then
  echo "FAIL: ATHENA_SECURITY_CHANGE_ACK must equal 'LAUNCH AUTHORIZED'"
  exit 1
fi

window_minutes="${ATHENA_SECURITY_GATE_WINDOW_MINUTES:-15}"
if [[ ! "$window_minutes" =~ ^[0-9]+$ || "$window_minutes" -le 0 ]]; then
  echo "FAIL: ATHENA_SECURITY_GATE_WINDOW_MINUTES must be a positive integer"
  exit 1
fi

age_minutes="$(
  python3 - "$ATHENA_SECURITY_NONCE_ISSUED_AT_UTC" <<'PY'
from datetime import datetime, timezone
import sys

value = sys.argv[1]
if value.endswith("Z"):
    value = value[:-1] + "+00:00"
issued = datetime.fromisoformat(value)
if issued.tzinfo is None:
    issued = issued.replace(tzinfo=timezone.utc)
now = datetime.now(timezone.utc)
delta = now - issued.astimezone(timezone.utc)
print(int(delta.total_seconds() // 60))
PY
)"

if [[ ! "$age_minutes" =~ ^-?[0-9]+$ ]]; then
  echo "FAIL: unable to parse ATHENA_SECURITY_NONCE_ISSUED_AT_UTC"
  exit 1
fi

if [[ "$age_minutes" -lt 0 ]]; then
  echo "FAIL: nonce issued time is in the future"
  exit 1
fi

if (( age_minutes > window_minutes )); then
  echo "FAIL: nonce expired (${age_minutes}m > window ${window_minutes}m)"
  exit 1
fi

echo "PASS: security nonce gate validated for protected contract changes"
