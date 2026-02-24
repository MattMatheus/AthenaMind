#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
runner="$root_dir/tools/run_stage_tests.sh"

tmp_dir="$(mktemp -d)"
trap 'rm -rf "$tmp_dir"' EXIT

stub_bin="$tmp_dir/bin"
mkdir -p "$stub_bin"
log_file="$tmp_dir/commands.log"

cat >"$stub_bin/go" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf 'go %s\n' "$*" >>"${RUN_STAGE_TESTS_LOG:?}"
exit 0
EOF
chmod +x "$stub_bin/go"

doc_stub="$tmp_dir/doc_stub.sh"
cat >"$doc_stub" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
printf 'docs\n' >>"${RUN_STAGE_TESTS_LOG:?}"
exit 0
EOF
chmod +x "$doc_stub"

run_case() {
  local label="$1"
  shift
  : >"$log_file"
  if ! env PATH="$stub_bin:$PATH" RUN_STAGE_TESTS_LOG="$log_file" STAGE_DOC_TEST_CMD="$doc_stub" "$@" >/dev/null 2>&1; then
    echo "FAIL: $label"
    exit 1
  fi
}

run_case "push mode runs targeted tests" "$runner" --mode push
if grep -Fq "docs" "$log_file" && grep -Fq "go test ./internal/governance ./internal/episode ./cmd/memory-cli" "$log_file"; then
  echo "PASS: push mode runs docs + targeted tests"
else
  echo "FAIL: push mode did not run docs + targeted tests"
  cat "$log_file"
  exit 1
fi

run_case "pr mode runs full suite" "$runner" --mode pr
if grep -Fq "docs" "$log_file" && grep -Fq "go test ./..." "$log_file"; then
  echo "PASS: pr mode runs docs + full suite"
else
  echo "FAIL: pr mode did not run docs + full suite"
  cat "$log_file"
  exit 1
fi

run_case "auto mode detects PR context" BUILD_REASON=PullRequest "$runner" --mode auto
if grep -Fq "go test ./..." "$log_file"; then
  echo "PASS: auto mode detects PR context"
else
  echo "FAIL: auto mode did not detect PR context"
  cat "$log_file"
  exit 1
fi

echo "Result: PASS"

