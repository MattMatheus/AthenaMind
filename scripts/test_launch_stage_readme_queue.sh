#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
active_readme="$root_dir/backlog/engineering/active/README.md"
launch_script="$root_dir/scripts/launch_stage.sh"
story_a="$root_dir/backlog/engineering/active/STORY-TEST-README-ORDER-A.md"
story_b="$root_dir/backlog/engineering/active/STORY-TEST-README-ORDER-B.md"
expected_story="backlog/engineering/active/STORY-TEST-README-ORDER-B.md"

tmp_dir="$(mktemp -d)"
restore_readme() {
  if [[ -f "$tmp_dir/README.md.orig" ]]; then
    cp "$tmp_dir/README.md.orig" "$active_readme"
  fi
  rm -f "$story_a" "$story_b"
  rm -rf "$tmp_dir"
}
trap restore_readme EXIT

cp "$active_readme" "$tmp_dir/README.md.orig"
cat >"$story_a" <<'EOF'
# Temporary test story A
EOF
cat >"$story_b" <<'EOF'
# Temporary test story B
EOF

cat >"$active_readme" <<'EOF'
# Active Queue

Ordered execution queue for engineering implementation.

## Active Sequence
1. `STORY-TEST-README-ORDER-B.md`
2. `STORY-TEST-README-ORDER-A.md`
EOF

output="$("$launch_script" engineering)"
if grep -Fq "story: $expected_story" <<<"$output"; then
  echo "PASS: launch_stage engineering honors Active README ordering"
else
  echo "FAIL: launch_stage engineering ignored Active README ordering"
  echo "Expected story: $expected_story"
  echo "Actual output:"
  echo "$output"
  exit 1
fi

echo "Result: PASS"
