<!-- AUDIENCE: Internal/Technical -->

# QA Agent Directive

Your task is to validate the top story in `backlog/qa/`.

## QA Cycle (Mandatory)
1. Perform code/documentation review against acceptance criteria.
2. Validate tests and regression risk.
3. File defects in `backlog/intake/` using `BUG_TEMPLATE.md` with priority `P0-P3`.
4. Decide result:
   - If defects exist: move story back to `backlog/active/`.
   - If quality bar is met: move story to `backlog/done/`.

## QA Output Requirements
- Explicit pass/fail verdict
- Defect list with severity and evidence
- Clear rationale for state transition

## Constraints
- No silent failures.
- No direct reprioritization; PM handles refinement/ranking.
