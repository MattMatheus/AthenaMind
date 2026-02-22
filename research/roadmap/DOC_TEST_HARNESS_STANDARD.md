# Doc Test Harness Standard (v0.1)

## Purpose
Standardize validation of documentation/decision stories so engineering and QA use one repeatable pattern.

## Canonical Command
- `scripts/run_doc_tests.sh`

Engineering stories with doc/ADR output should run the canonical command unless a story requires an additional scoped test.

## Pattern
1. Add a story-specific test script at `scripts/test_<story-scope>.sh`.
2. Source the shared harness helper from `scripts/lib/doc_test_harness.sh`.
3. Use shared assertions:
   - `doc_assert_exists`
   - `doc_assert_contains`
4. End tests with `doc_test_finish`.
5. Add new story-specific test script to `scripts/run_doc_tests.sh`.

## Reusable Template
Use this template as baseline for new doc-driven story tests:

```bash
#!/usr/bin/env bash
set -euo pipefail

root_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
source "$root_dir/scripts/lib/doc_test_harness.sh"

doc_test_init

# Example checks:
doc_assert_exists "$root_dir/path/to/artifact.md" "Artifact exists"
doc_assert_contains "$root_dir/path/to/artifact.md" "Expected text" "Artifact contains expected text"

doc_test_finish
```

## Scope Notes
- This standard is for documentation/ADR/planning stories only.
- Runtime implementation tests and CI pipeline decisions are out of scope for this standard.
