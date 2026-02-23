# QA Regression Rubric and Defect Triage Rules

## Purpose
Standardize QA pass/fail decisions and bug severity mapping for stories reviewed in `delivery-backlog/qa/`.

## Pass/Fail Gates (Deterministic)
A story is `PASS` only if all gates below pass.

1. Acceptance criteria gate:
   - All listed acceptance criteria are satisfied with explicit evidence.
2. Test gate:
   - Required test commands run.
   - All required tests pass.
3. Regression gate:
   - No known regressions in touched scope.
4. Artifact gate:
   - Required handoff package exists for the story.
5. Cycle-closure readiness gate:
   - QA output is ready for observer reporting and single cycle commit.

If any gate fails, result is `FAIL`.

## Defect Priority Mapping (Aligned to BUG_TEMPLATE)
- `P0`: release-blocking, data loss/corruption, security-critical
- `P1`: major functional regression, acceptance criteria blocked
- `P2`: moderate defect with workaround
- `P3`: minor defect, polish, or low-impact inconsistency

Priority assignment rule:
1. Select highest applicable impact class.
2. If uncertain between two levels, choose higher severity and note uncertainty in evidence.

## Minimum Evidence Requirements for QA Bug Filing
Every filed bug must include:
- `source_story` path under review
- expected behavior
- actual behavior
- reproduction steps (minimum 3 steps)
- evidence references (test output, logs, screenshots, or doc diffs)
- assigned priority and rationale

## QA Handoff Examples
### Example PASS
- Story: `delivery-backlog/qa/STORY-YYYYMMDD-example-pass.md`
- Result: `PASS`
- Why:
  - acceptance criteria met
  - tests passed
  - no regression evidence
  - handoff/QA artifacts present
- Transition:
  - move to `delivery-backlog/done/`
  - run observer and close cycle with `cycle-<cycle-id>` commit

### Example FAIL
- Story: `delivery-backlog/qa/STORY-YYYYMMDD-example-fail.md`
- Result: `FAIL`
- Why:
  - acceptance criterion blocked by defect
  - regression found in touched scope
- Defect filing:
  - create `delivery-backlog/intake/BUG-YYYYMMDD-<slug>.md` with `P0-P3` priority and required evidence
- Transition:
  - move story back to `delivery-backlog/active/` with linked bug path(s)
  - run observer and close cycle with `cycle-<cycle-id>` commit
