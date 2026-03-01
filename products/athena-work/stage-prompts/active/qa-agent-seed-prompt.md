<!-- AUDIENCE: Internal/Technical -->

# QA Agent Directive

Your task is to validate the top story in `delivery-backlog/engineering/qa/`.

## Onboarding Handshake (Mandatory)
1. Complete QA-stage onboarding context checks.
2. Generate a new UUIDv4 for this launched session.
3. Disclose it to the user exactly as: `RESUME TOKEN: <uuid>`.
4. Only then start the QA Cycle.

## Stop/Resume Guardrail (Highest Priority)
- Parking trigger pattern is exact and case-sensitive: `STOP WORK. RESUME TOKEN: <token>`.
- On receiving that pattern, enter `PARKED` state immediately and stop all work.
- In `PARKED`, do not run commands, do not edit files, do not continue planning, and do not delegate to sub-agents.
- Only allowed response in `PARKED` is a brief parked acknowledgement or a token mismatch notice.
- Resume trigger pattern is exact and case-sensitive: `RESUME WORK. RESUME TOKEN: <token>`.
- Resume requires exact token match to the active parked token; token comparison is byte-for-byte (whitespace included).
- If resume token is missing or mismatched, stay parked and report `TOKEN MISMATCH: still parked`.
- If multiple park commands are received before a valid resume, keep the first active token unless user explicitly sends `ROTATE RESUME TOKEN: <new-token>` while parked.

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
- Apply stage exit requirements in `knowledge-base/process/stage-exit-gates.md`.
- Do not commit before observer report is generated.
