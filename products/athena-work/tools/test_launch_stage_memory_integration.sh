#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
launch_script="$root_dir/tools/launch_stage.sh"
current_branch="$(git -C "$root_dir" branch --show-current)"
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
if [[ "${1:-}" == "bootstrap" ]]; then
  echo '{"memory_entries":[{"id":"stub-entry"}]}'
  exit 0
fi
echo "unexpected command: $*" >&2
exit 1
EOF
chmod +x "$stub_bin/memory-cli"

output="$(PATH="$stub_bin:$PATH" MEMORY_CLI_TEST_LOG="$args_log" ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" pm)"
if grep -Fq "memory_bootstrap_context:" <<<"$output" && grep -Fq "stub-entry" <<<"$output"; then
  echo "PASS: launch_stage appends memory bootstrap context when memory-cli is available"
else
  echo "FAIL: launch_stage did not append memory bootstrap context"
  echo "$output"
  exit 1
fi

if grep -Fq "bootstrap --root $default_memory_root --repo $repo_id" "$args_log" && grep -Fq " --scenario pm" "$args_log"; then
  echo "PASS: launch_stage invokes memory-cli bootstrap with stage scenario"
else
  echo "FAIL: launch_stage bootstrap invocation missing expected args"
  cat "$args_log"
  exit 1
fi

warn_file="$tmp_dir/warn.log"
MEMORY_CLI_BIN="memory-cli-missing-for-test" ATHENA_REQUIRED_BRANCH="$current_branch" "$launch_script" qa > /dev/null 2>"$warn_file"
if grep -Fq "warning: memory bootstrap skipped" "$warn_file"; then
  echo "PASS: launch_stage degrades gracefully when memory-cli is unavailable"
else
  echo "FAIL: launch_stage did not warn when memory-cli was unavailable"
  cat "$warn_file"
  exit 1
fi

echo "Result: PASS"
