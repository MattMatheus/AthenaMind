#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
observer_script="$root_dir/tools/run_observer_cycle.sh"
repo_id="$(basename "$root_dir")"
default_memory_root="$HOME/.athena/memory/$repo_id"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

stub_bin="$tmp_dir/bin"
mkdir -p "$stub_bin"
args_log="$tmp_dir/memory-cli.args.log"
cat >"$stub_bin/memory-cli" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
log_file="${MEMORY_CLI_TEST_LOG:?}"
printf '%s\n' "$*" >>"$log_file"
if [[ "${1:-}" == "episode" && "${2:-}" == "write" ]]; then
  echo '{"status":"ok"}'
  exit 0
fi
echo "unexpected command: $*" >&2
exit 1
EOF
chmod +x "$stub_bin/memory-cli"

report_path="$tmp_dir/observer-report.md"
story_path="$tmp_dir/story.md"
cat >"$story_path" <<'EOF'
# Test Story

## Metadata
- `idea_id`: direct
- `adr_refs`: [ADR-TEST]
EOF

output="$(PATH="$stub_bin:$PATH" MEMORY_CLI_TEST_LOG="$args_log" "$observer_script" --cycle-id STORY-TEST-MEMORY-INTEGRATION --story "$story_path" --output "$report_path")"
if grep -Fq "memory episode write-back: recorded" <<<"$output"; then
  echo "PASS: observer writes back episode memory when memory-cli is available"
else
  echo "FAIL: observer did not report episode write-back"
  echo "$output"
  exit 1
fi

if [[ -f "$report_path" ]]; then
  echo "PASS: observer report still generated with write-back integration"
else
  echo "FAIL: observer report was not generated"
  exit 1
fi

if grep -Fq "episode write --root $default_memory_root --repo $repo_id" "$args_log" && grep -Fq " --cycle-id STORY-TEST-MEMORY-INTEGRATION " "$args_log"; then
  echo "PASS: observer invokes memory-cli episode write with cycle data"
else
  echo "FAIL: observer episode write args missing expected cycle metadata"
  cat "$args_log"
  exit 1
fi

if grep -Fq " --stage pm " "$args_log"; then
  echo "PASS: observer uses pm policy stage for engineering-style cycle ids"
else
  echo "FAIL: observer did not use expected pm policy stage for engineering-style cycle id"
  cat "$args_log"
  exit 1
fi

plan_log="$tmp_dir/memory-cli.plan.args.log"
plan_report="$tmp_dir/observer-report-plan.md"
PATH="$stub_bin:$PATH" MEMORY_CLI_TEST_LOG="$plan_log" "$observer_script" --cycle-id PLAN-TEST-STAGE-INFERENCE --output "$plan_report" >/dev/null
if grep -Fq " --stage planning " "$plan_log"; then
  echo "PASS: observer infers planning policy stage from planning cycle id"
else
  echo "FAIL: observer did not infer planning policy stage"
  cat "$plan_log"
  exit 1
fi

warn_file="$tmp_dir/warn.log"
MEMORY_CLI_BIN="memory-cli-missing-for-test" "$observer_script" --cycle-id STORY-TEST-MEMORY-MISSING --output "$tmp_dir/observer-report-missing.md" >/dev/null 2>"$warn_file"
if grep -Fq "warning: episode write-back skipped" "$warn_file"; then
  echo "PASS: observer degrades gracefully when memory-cli is unavailable"
else
  echo "FAIL: observer did not warn when memory-cli was unavailable"
  cat "$warn_file"
  exit 1
fi

echo "Result: PASS"
