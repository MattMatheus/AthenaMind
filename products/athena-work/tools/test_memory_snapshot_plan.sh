#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/tools/lib/doc_test_harness.sh"

brief="$root_dir/product-research/roadmap/MEMORY_SNAPSHOT_DESIGN_BRIEF_POST_V01.md"

doc_test_init

doc_assert_exists "$brief" "Memory snapshot design brief exists"
doc_assert_contains "$brief" "Snapshot Use Cases" "Brief defines snapshot use cases"
doc_assert_contains "$brief" "Restore Semantics" "Brief defines restore semantics"
doc_assert_contains "$brief" "Data Model and Versioning Implications" "Brief defines data model/versioning"
doc_assert_contains "$brief" "Integration Points (Current Modules)" "Brief identifies module integration points"
doc_assert_contains "$brief" "post-v0.1" "Brief marks rollout as post-v0.1"
doc_assert_contains "$brief" "ADR-0007" "Brief aligns scope with ADR-0007"

doc_test_finish
