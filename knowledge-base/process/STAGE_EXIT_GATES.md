# Stage Exit Gates

Deterministic exit gates for the idea -> architecture -> PM -> engineering -> QA -> shipped pipeline.

## Planning Exit Gate
All must pass:
1. Planning session artifact exists in `product-research/planning/sessions/`.
2. Session `status` is `finalized`.
3. Intake artifacts are created in correct lanes.
4. Next-stage recommendation is explicit (`architect` or `pm`).
5. Planning outputs are ready for observer capture in cycle closure.

## Architect Exit Gate
All must pass:
1. Architecture story has explicit scope and accepted output package.
2. ADR/artifact updates are complete.
3. Follow-on implementation story paths are listed in handoff.
4. Validation commands are recorded and passing.
5. Story transitions `active -> qa` only (never directly to `done`).
6. Outputs are ready for observer capture in cycle closure.

## PM Refinement Exit Gate
All must pass:
1. Intake validation passes (`tools/validate_intake_items.sh`).
2. Active queue is ranked and explicit in `delivery-backlog/engineering/active/README.md`.
3. Product-first backlog weighting is applied per `knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md`.
4. Active stories include traceability metadata (`idea_id`, `phase`, `adr_refs`, metric).
5. Program board is updated (`product-research/roadmap/PROGRAM_STATE_BOARD.md`).
6. PM TODO `Now` contains at least one actionable item.

## Engineering Exit Gate
All must pass:
1. Story acceptance criteria are implemented.
2. Tests updated for touched behavior.
3. `tools/run_doc_tests.sh` and story-specific tests pass.
4. `go test ./...` passes locally and is eligible for CI enforcement in Azure DevOps.
5. Handoff package is complete and includes risks/questions.
6. New gaps are recorded as intake artifacts before handoff.
7. Story transitions `active -> qa`.
8. No stage-level commit is made before observer step.

## QA Exit Gate
All must pass for `PASS`:
1. Acceptance criteria evidence is explicit.
2. Test gate passes.
3. Regression gate passes.
4. Artifact gate passes (handoff present).
5. QA result artifact exists in `delivery-backlog/engineering/done/QA-RESULT-<story>.md`.
6. State transition is explicit (`qa -> done` or `qa -> active` with bugs).
7. No stage-level commit is made before observer step.

## Cycle Closure Gate (Observer + Commit)
A cycle is closed only when all pass:
1. Observer report exists at `operating-system/observer/OBSERVER-REPORT-<cycle-id>.md`.
2. Observer report includes a diff inventory and workflow-sync checklist.
3. Required sync targets were updated when workflow behavior changed (`HUMANS.md`, `AGENTS.md`, `DEVELOPMENT_CYCLE.md`, stage-prompts/gates as needed).
4. Exactly one cycle commit is created with subject `cycle-<cycle-id>`.
5. Cycle commit includes observer report and all cycle artifacts (handoff/QA/program-state updates as applicable).

## Shipped Gate (Release Checkpoint)
`done` is not automatically `shipped`.

A change is `shipped` only when all pass:
1. Release bundle exists at `operating-system/handoff/RELEASE_BUNDLE_<release-id>.md`.
2. Bundle lists included stories/bugs and QA result evidence links.
3. Bundle records operational risks and rollback direction.
4. Bundle records outcome metric baseline or expected trend.
5. Bundle decision is explicit: `ship` or `hold` with rationale.

Reference template: `operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`.
