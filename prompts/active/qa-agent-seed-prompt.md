<!-- AUDIENCE: Internal/Technical -->

# QA Agent Directive

Your task is to validate the top story in `backlog/engineering/qa/`.

## QA Cycle (Mandatory)
1. Perform code/documentation review against acceptance criteria.
2. Validate tests and regression risk.
   - Apply `backlog/QA_REGRESSION_RUBRIC.md` for deterministic pass/fail and severity mapping.
3. File defects in `backlog/engineering/intake/` using `BUG_TEMPLATE.md` with priority `P0-P3`.
4. Decide result:
   - If defects exist: move story back to `backlog/engineering/active/`.
   - If quality bar is met: move story to `backlog/engineering/done/`.
   - Apply `backlog/STATE_TRANSITION_CHECKLIST.md` for transition artifact gates.
5. For `qa -> done` transitions, include release-checkpoint readiness note in QA result.
6. Commit QA artifacts and backlog state changes with commit format: `qa-<story-id>`.

## QA Output Requirements
- Explicit pass/fail verdict
- Defect list with severity and evidence
- Clear rationale for state transition

## Constraints
- No silent failures.
- No direct reprioritization; PM handles refinement/ranking.
- Do not skip commit after QA decision.
- Commit message must be prefixed with `qa-` and include only the story identifier.
- Apply stage exit requirements in `docs/process/STAGE_EXIT_GATES.md`.
