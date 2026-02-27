#!/usr/bin/env bash
set -euo pipefail

script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
root_dir="$(git -C "$script_dir" rev-parse --show-toplevel 2>/dev/null || (cd "$script_dir/.." && pwd))"
source "$root_dir/tools/lib/doc_test_harness.sh"

doc_test_init

assert_dir_exists() {
  local path="$1"
  local label="$2"
  if [[ -d "$path" ]]; then
    echo "PASS: $label"
  else
    echo "FAIL: $label"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  fi
}

engineering_root="$root_dir/delivery-backlog/engineering"
architecture_root="$root_dir/delivery-backlog/architecture"
active_readme="$engineering_root/active/README.md"

for state in intake ready active qa blocked done archive; do
  assert_dir_exists "$engineering_root/$state" "engineering lane state exists: $state"
done

for state in intake ready active qa blocked done archive; do
  assert_dir_exists "$architecture_root/$state" "architecture lane state exists: $state"
done

if [[ -f "$active_readme" ]]; then
  echo "PASS: engineering active README exists"
else
  echo "FAIL: engineering active README exists"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
fi

if [[ -f "$active_readme" ]] && grep -Fq '\n' "$active_readme"; then
  echo "FAIL: engineering active README contains literal \\n sequences"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
else
  echo "PASS: engineering active README newline formatting is valid"
fi

active_story_count=0
if [[ -d "$engineering_root/active" ]]; then
  active_story_count="$(find "$engineering_root/active" -maxdepth 1 -type f -name '*.md' ! -name 'README.md' | wc -l | tr -d ' ')"
fi

if [[ -f "$active_readme" ]] && grep -qi 'No active stories' "$active_readme" && [[ "$active_story_count" -gt 0 ]]; then
  echo "FAIL: active README says 'No active stories' while active story files exist"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
else
  echo "PASS: active README empty-queue claim is consistent"
fi

queue_refs=()
if [[ -f "$active_readme" ]]; then
  while IFS= read -r row; do
    [[ -n "$row" ]] && queue_refs+=("$row")
  done < <(sed -En 's/^[[:space:]]*[0-9]+\.[[:space:]]*`([^`]+)`.*/\1/p' "$active_readme")
fi

if [[ "${#queue_refs[@]}" -eq 0 && "$active_story_count" -gt 0 ]]; then
  echo "FAIL: active stories exist but README queue ordering is missing"
  DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
else
  echo "PASS: active queue ordering presence is consistent"
fi

dupes="$(
  {
    for ref in "${queue_refs[@]:-}"; do
      [[ -z "$ref" ]] && continue
      printf '%s\n' "$ref"
    done
  } | sort | uniq -d
)"

if [[ -n "$dupes" ]]; then
  while IFS= read -r dup; do
    [[ -z "${dup:-}" ]] && continue
    echo "FAIL: duplicate queue reference in README: ${dup}"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  done <<< "$dupes"
else
  echo "PASS: queue references are unique"
fi

for ref in "${queue_refs[@]:-}"; do
  [[ -z "$ref" ]] && continue
  resolved="$ref"
  if [[ "${resolved#/}" == "$resolved" ]]; then
    if [[ "$resolved" == */* ]]; then
      resolved="$root_dir/$resolved"
    else
      resolved="$engineering_root/active/$resolved"
    fi
  fi

  if [[ ! -f "$resolved" ]]; then
    echo "FAIL: queue reference does not exist: $ref"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  else
    echo "PASS: queue reference exists: $ref"
  fi
done

doc_test_finish
