<!-- AUDIENCE: Internal/Technical -->

# QA Agent Directive

Your task is to validate the top story in `delivery-backlog/engineering/qa/`.

## QA Cycle (Mandatory)
1. Perform code/documentation review against acceptance criteria.
2. Validate tests and regression risk.
   - Apply `delivery-backlog/QA_REGRESSION_RUBRIC.md` for deterministic pass/fail and severity mapping.
3. File defects in `delivery-backlog/engineering/intake/` using `BUG_TEMPLATE.md` with priority `P0-P3`.
4. Decide result:
   - If defects exist: move story back to `delivery-backlog/engineering/active/`.
   - If quality bar is met: move story to `delivery-backlog/engineering/done/`.
   - Apply `delivery-backlog/STATE_TRANSITION_CHECKLIST.md` for transition artifact gates.
5. For `qa -> done` transitions, include release-checkpoint readiness note in QA result.
6. Run observer:
   - `tools/run_observer_cycle.sh --cycle-id <story-id> --story <path-to-story>`
7. Commit once for the full cycle:
   - `cycle-<cycle-id>`

## QA Output Requirements
- Explicit pass/fail verdict
- Defect list with severity and evidence
- Clear rationale for state transition

## Constraints
- No silent failures.
- No direct reprioritization; PM handles refinement/ranking.
- Apply stage exit requirements in `knowledge-base/process/STAGE_EXIT_GATES.md`.
- Do not commit before observer report is generated.
