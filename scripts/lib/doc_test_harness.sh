#!/usr/bin/env bash
set -euo pipefail

doc_test_init() {
  DOC_TEST_FAILURES=0
}

doc_assert_contains() {
  local file="$1"
  local text="$2"
  local label="$3"
  if grep -Fq "$text" "$file"; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
}

doc_assert_exists() {
  local file="$1"
  local label="$2"
  if [[ -f "$file" ]]; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
}

doc_test_finish() {
  if [[ "${DOC_TEST_FAILURES:-0}" -gt 0 ]]; then
    echo "Result: FAIL (${DOC_TEST_FAILURES} checks failed)"
    exit 1
  fi
  echo "Result: PASS"
}
