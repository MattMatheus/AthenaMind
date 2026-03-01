#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
guard="$root_dir/tools/check_markdown_drift.sh"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

write_state() {
  local state_path="$1"
  local version="$2"
  local story="$3"
  cat >"$state_path" <<EOF
{"state_version":"$version","engineering":{"active":["$story"]}}
EOF
}

write_readme() {
  local readme_path="$1"
  local version="$2"
  local story="$3"
  cat >"$readme_path" <<EOF
# Engineering Active Queue

Ordered execution queue for engineering stories.
<!-- projection_version: $version -->

## Active Sequence
1. \`$story.md\`
EOF
}

# Scenario 1: ordering + stale revision
state_a="$tmp_dir/state-a.json"
readme_a="$tmp_dir/README-a.md"
stories_a="$tmp_dir/stories-a"
mkdir -p "$stories_a"
write_state "$state_a" "test-a" "STORY-A"
echo "# Story A" >"$stories_a/STORY-A.md"
write_readme "$readme_a" "old-version" "STORY-OTHER"

out_a="$(ATHENA_MD_STATE_FILE="$state_a" ATHENA_MD_ACTIVE_README="$readme_a" ATHENA_MD_ACTIVE_STORY_DIR="$stories_a" "$guard" --dry-run 2>&1 || true)"
if grep -Fq "drift_class: ordering_conflict" <<<"$out_a" && \
   grep -Fq "drift_class: stale_revision" <<<"$out_a" && \
   grep -Fq "remediation_id: run_markdown_sync_worker" <<<"$out_a" && \
   grep -Fq "remediation_id: refresh_projection_revision" <<<"$out_a"; then
  echo "PASS: drift guard reports deterministic ordering/stale classes with remediation ids"
else
  echo "FAIL: drift guard ordering/stale regression output mismatch"
  echo "$out_a"
  exit 1
fi

# Scenario 2: missing artifact
state_b="$tmp_dir/state-b.json"
readme_b="$tmp_dir/README-b.md"
stories_b="$tmp_dir/stories-b"
mkdir -p "$stories_b"
write_state "$state_b" "test-b" "STORY-MISSING"
write_readme "$readme_b" "test-b" "STORY-MISSING"

out_b="$(ATHENA_MD_STATE_FILE="$state_b" ATHENA_MD_ACTIVE_README="$readme_b" ATHENA_MD_ACTIVE_STORY_DIR="$stories_b" "$guard" --dry-run 2>&1 || true)"
if grep -Fq "drift_class: missing_artifact" <<<"$out_b" && \
   grep -Fq "remediation_id: restore_or_update_state" <<<"$out_b"; then
  echo "PASS: drift guard reports deterministic missing-artifact class with remediation id"
else
  echo "FAIL: drift guard missing-artifact regression output mismatch"
  echo "$out_b"
  exit 1
fi

echo "Result: PASS"
