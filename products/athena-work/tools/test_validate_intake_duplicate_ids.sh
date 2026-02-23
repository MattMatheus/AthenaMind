#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
validator="$root_dir/tools/validate_intake_items.sh"
eng_intake="$root_dir/delivery-backlog/engineering/intake"

tmp_one="$eng_intake/STORY-TEST-DUPLICATE-ID-ONE.md"
tmp_two="$eng_intake/STORY-TEST-DUPLICATE-ID-TWO.md"

cleanup() {
  rm -f "$tmp_one" "$tmp_two"
}
trap cleanup EXIT

cat >"$tmp_one" <<'EOF'
# Story: Duplicate ID A

## Metadata
- `id`: STORY-TEST-DUPLICATE-ID
- `owner_persona`: Product Manager - Maya
- `status`: intake
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-TEST]
- `success_metric`: duplicate-id validation catches collisions
- `release_checkpoint`: deferred
EOF

cat >"$tmp_two" <<'EOF'
# Story: Duplicate ID B

## Metadata
- `id`: STORY-TEST-DUPLICATE-ID
- `owner_persona`: Software Architect - Ada
- `status`: intake
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-TEST]
- `success_metric`: duplicate-id validation catches collisions
- `release_checkpoint`: deferred
EOF

if "$validator" >/dev/null 2>"$root_dir/.tmp.intake-duplicate.err"; then
  echo "FAIL: intake validator passed despite duplicate IDs"
  rm -f "$root_dir/.tmp.intake-duplicate.err"
  exit 1
fi

if grep -Fq "duplicate intake id 'STORY-TEST-DUPLICATE-ID'" "$root_dir/.tmp.intake-duplicate.err"; then
  echo "PASS: intake validator rejects duplicate IDs"
else
  echo "FAIL: intake validator did not report duplicate ID failure"
  cat "$root_dir/.tmp.intake-duplicate.err"
  rm -f "$root_dir/.tmp.intake-duplicate.err"
  exit 1
fi

rm -f "$root_dir/.tmp.intake-duplicate.err"
echo "Result: PASS"
