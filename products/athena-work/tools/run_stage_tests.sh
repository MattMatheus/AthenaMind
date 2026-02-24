#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
cd "$root_dir"

mode="auto"

usage() {
  cat <<'USAGE'
usage: tools/run_stage_tests.sh [--mode auto|pr|push]

Runs docs tests plus the correct Go test scope:
  - push/non-PR: targeted tests
  - PR: full suite

Overrides:
  STAGE_DOC_TEST_CMD       (default: ./tools/run_doc_tests.sh)
  STAGE_TARGETED_TEST_CMD  (default: go test ./internal/governance ./internal/episode ./cmd/memory-cli)
  STAGE_FULL_TEST_CMD      (default: go test ./...)
USAGE
}

while [[ $# -gt 0 ]]; do
  case "$1" in
    --mode)
      mode="${2:-}"
      shift 2
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "error: unknown arg '$1'" >&2
      usage
      exit 1
      ;;
  esac
done

case "$mode" in
  auto|pr|push) ;;
  *)
    echo "error: --mode must be one of auto|pr|push" >&2
    exit 1
    ;;
esac

is_pr_context() {
  if [[ "${mode}" == "pr" ]]; then
    return 0
  fi
  if [[ "${mode}" == "push" ]]; then
    return 1
  fi

  if [[ "${BUILD_REASON:-}" == "PullRequest" ]]; then
    return 0
  fi
  if [[ -n "${SYSTEM_PULLREQUEST_PULLREQUESTID:-}" ]]; then
    return 0
  fi
  if [[ "${GITHUB_EVENT_NAME:-}" == "pull_request" || "${GITHUB_EVENT_NAME:-}" == "pull_request_target" ]]; then
    return 0
  fi
  if [[ -n "${CI_MERGE_REQUEST_IID:-}" || "${CI_PIPELINE_SOURCE:-}" == "merge_request_event" ]]; then
    return 0
  fi
  return 1
}

doc_cmd="${STAGE_DOC_TEST_CMD:-./tools/run_doc_tests.sh}"
targeted_cmd="${STAGE_TARGETED_TEST_CMD:-go test ./internal/governance ./internal/episode ./cmd/memory-cli}"
full_cmd="${STAGE_FULL_TEST_CMD:-go test ./...}"

echo "running docs tests: $doc_cmd"
bash -lc "$doc_cmd"

if is_pr_context; then
  echo "test scope: full (PR context)"
  echo "running: $full_cmd"
  bash -lc "$full_cmd"
else
  echo "test scope: targeted (push/non-PR context)"
  echo "running: $targeted_cmd"
  bash -lc "$targeted_cmd"
fi

