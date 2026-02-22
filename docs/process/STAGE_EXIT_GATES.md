# Stage Exit Gates

Deterministic exit gates for the idea -> architecture -> PM -> engineering -> QA -> shipped pipeline.

## Planning Exit Gate
All must pass:
1. Planning session artifact exists in `research/planning/sessions/`.
2. Session `status` is `finalized`.
3. Intake artifacts are created in correct lanes.
4. Next-stage recommendation is explicit (`architect` or `pm`).
5. Planning commit uses `plan-<plan-id>`.

## Architect Exit Gate
All must pass:
1. Architecture story has explicit scope and accepted output package.
2. ADR/artifact updates are committed.
3. Follow-on implementation story paths are listed in handoff.
4. Validation commands are recorded and passing.
5. Story transitions `active -> qa` only (never directly to `done`).
6. Commit uses `arch-<story-id>`.

## PM Refinement Exit Gate
All must pass:
1. Intake validation passes (`scripts/validate_intake_items.sh`).
2. Active queue is ranked and explicit in `backlog/engineering/active/README.md`.
3. Active stories include traceability metadata (`idea_id`, `phase`, `adr_refs`, metric).
4. Program board is updated (`research/roadmap/PROGRAM_STATE_BOARD.md`).
5. PM TODO `Now` contains at least one actionable item.

## Engineering Exit Gate
All must pass:
1. Story acceptance criteria are implemented.
2. Tests updated for touched behavior.
3. `scripts/run_doc_tests.sh` and story-specific tests pass.
4. Handoff package is complete and includes risks/questions.
5. New gaps are recorded as intake artifacts before handoff.
6. Story transitions `active -> qa`.
7. Commit message includes story id.

## QA Exit Gate
All must pass for `PASS`:
1. Acceptance criteria evidence is explicit.
2. Test gate passes.
3. Regression gate passes.
4. Artifact gate passes (handoff present).
5. QA result artifact exists in `backlog/engineering/done/QA-RESULT-<story>.md`.
6. State transition is explicit (`qa -> done` or `qa -> active` with bugs).
7. Commit uses `qa-<story-id>`.

## Shipped Gate (Release Checkpoint)
`done` is not automatically `shipped`.

A change is `shipped` only when all pass:
1. Release bundle exists at `work-system/handoff/RELEASE_BUNDLE_<release-id>.md`.
2. Bundle lists included stories/bugs and QA result evidence links.
3. Bundle records operational risks and rollback direction.
4. Bundle records outcome metric baseline or expected trend.
5. Bundle decision is explicit: `ship` or `hold` with rationale.

Reference template: `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`.
