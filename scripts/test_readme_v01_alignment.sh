#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

readme="$root_dir/README.md"
vision="$root_dir/docs/product/VISION.md"

doc_test_init

doc_assert_exists "$readme" "README exists"
doc_assert_exists "$vision" "Vision document exists"
doc_assert_contains "$readme" "What v0.1 Delivers Today" "README includes current v0.1 scope section"
doc_assert_contains "$readme" "go run ./cmd/memory-cli write" "README includes write example"
doc_assert_contains "$readme" "go run ./cmd/memory-cli retrieve" "README includes retrieve example"
doc_assert_contains "$readme" "research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md" "README links phased plan"
doc_assert_contains "$readme" "docs/product/VISION.md" "README links preserved vision doc"
doc_assert_contains "$vision" "AthenaMind Product Vision (Long-Term)" "Vision document title present"

for forbidden in "FAISS" "Podman" "SQLite" "embeddings" "cloud"; do
  if grep -qi "$forbidden" "$readme"; then
    echo "FAIL: README should avoid unimplemented term '$forbidden'"
    DOC_TEST_FAILURES=$((DOC_TEST_FAILURES + 1))
  else
    echo "PASS: README avoids unimplemented term '$forbidden'"
  fi
done

doc_test_finish
