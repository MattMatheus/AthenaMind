#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
target="$root_dir/tools/self_check.sh"
active_root="$root_dir/delivery-backlog/engineering/active"
active_readme="$active_root/README.md"
tmp_story="$active_root/STORY-TEST-SELF-CHECK.md"
tmp_backup="$(mktemp -d)"

cleanup() {
  rm -f "$tmp_story"
  if [[ -f "$tmp_backup/README.md.orig" ]]; then
    cp "$tmp_backup/README.md.orig" "$active_readme"
  fi
  rm -rf "$tmp_backup"
  rm -rf "${stub_dir:-}"
}
trap cleanup EXIT

if [[ ! -x "$target" ]]; then
  echo "FAIL: self_check script must exist and be executable"
  exit 1
fi

cp "$active_readme" "$tmp_backup/README.md.orig"
cat >"$tmp_story" <<'EOF'
# Story: Temporary self-check validation

## Metadata
- `id`: STORY-TEST-SELF-CHECK
EOF

cat >"$active_readme" <<'EOF'
# Engineering Active Queue

Ordered execution queue for engineering implementation.

## Active Sequence
1. `STORY-TEST-SELF-CHECK.md`
EOF

output="$("$target")"

if grep -Fq "branch:" <<<"$output"; then
  echo "PASS: self_check prints branch"
else
  echo "FAIL: self_check prints branch"
  exit 1
fi

if grep -Fq "git: changed=" <<<"$output"; then
  echo "PASS: self_check prints git summary"
else
  echo "FAIL: self_check prints git summary"
  exit 1
fi

if grep -Fq "Active Queue (engineering)" <<<"$output"; then
  echo "PASS: self_check prints active queue section"
else
  echo "FAIL: self_check prints active queue section"
  exit 1
fi

if grep -Fq "STORY-TEST-SELF-CHECK.md" <<<"$output" && grep -Fq "Story: Temporary self-check validation" <<<"$output"; then
  echo "PASS: self_check prints active story id and title"
else
  echo "FAIL: self_check prints active story id and title"
  exit 1
fi

warn_output="$(ATHENA_REQUIRED_BRANCH=definitely-not-current "$target")"
if grep -Fq "WARN: expected branch 'definitely-not-current'" <<<"$warn_output"; then
  echo "PASS: self_check warns for required branch mismatch"
else
  echo "FAIL: self_check warns for required branch mismatch"
  exit 1
fi

stub_dir="$(mktemp -d)"

cat >"$stub_dir/git" <<'EOF'
#!/usr/bin/env bash
exit 1
EOF
chmod +x "$stub_dir/git"

if PATH="$stub_dir:$PATH" "$target" >/dev/null 2>&1; then
  echo "FAIL: self_check should fail when git state checks cannot run"
  exit 1
else
  echo "PASS: self_check fails when git state checks cannot run"
fi

echo "Result: PASS"
