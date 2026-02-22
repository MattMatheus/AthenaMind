# Bug: launch_stage engineering ignores Active README ordering on BSD sed

## Metadata
- `id`: BUG-20260222-launch-stage-readme-parse-fallback
- `priority`: P2
- `reported_by`: QA Engineer - Iris
- `source_story`: backlog/active/STORY-20260222-doc-test-harness-standardization.md
- `status`: intake

## Priority Definitions
- `P0`: release-blocking, data loss/corruption, security-critical
- `P1`: major functional regression, acceptance criteria blocked
- `P2`: moderate defect with workaround
- `P3`: minor defect, polish, or low-impact inconsistency

## Summary
`scripts/launch_stage.sh engineering` uses a sed expression with `\+` that does not parse as expected on BSD sed, so ordered entries in `backlog/active/README.md` are not read. The launcher silently falls back to alphabetical file ordering, which can execute stories out of ranked sequence.

## Expected Behavior
- Top story should be selected from `backlog/active/README.md` sequence when present and valid.

## Actual Behavior
- README parsing returns empty and launcher selects first filename alphabetically.

## Reproduction Steps
1. Ensure `backlog/active/README.md` lists a top story not alphabetically first.
2. Run `./scripts/launch_stage.sh engineering`.
3. Observe reported `story:` follows alphabetical order instead of queue order.

## Evidence
- Manual check of current sed extraction command returns no output on local environment.
- Launcher reported `STORY-20260222-doc-test-harness-standardization.md` while README rank listed `STORY-20260222-state-transition-checklist.md` first.

## Suggested Fix Direction (Optional)
- Replace BSD-incompatible sed quantifier usage with portable expression (`-E` + `+`, or POSIX-compliant pattern) and add a unit-style script test for README queue parsing.
