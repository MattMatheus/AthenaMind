# Cycle Index

Single navigation entrypoint for operators and agents running the AthenaMind delivery cycle.

## First 5 Minutes
1. Confirm branch safety:
   - run `git branch --show-current`
   - ensure it matches `ATHENA_REQUIRED_BRANCH` (default: `dev`)
2. Launch the stage you need:
   - `tools/launch_stage.sh planning`
   - `tools/launch_stage.sh architect`
   - `tools/launch_stage.sh engineering`
   - `tools/launch_stage.sh qa`
   - `tools/launch_stage.sh pm`
3. Open the returned seed prompt in `stage-prompts/active/` and follow it as directive.
4. Run docs validation command before handoff/decision points:
   - `tools/run_doc_tests.sh`
5. Apply backlog state movement only through the canonical flow (`engineering/active -> engineering/qa -> engineering/done`, with intake/active loop for defects).
6. Close each cycle with Observer + single commit:
   - `tools/run_observer_cycle.sh --cycle-id <cycle-id>`
   - `git commit -m "cycle-<cycle-id>"`
7. Apply stage and cycle gates:
   - `knowledge-base/process/STAGE_EXIT_GATES.md`

## Branch Rule and Empty Active Behavior
- Branch safety rule: launcher requires branch `ATHENA_REQUIRED_BRANCH` (default `dev`).
- If branch differs, launcher aborts:
  - `abort: active branch is '<branch>'; expected '<required-branch>'`
- If engineering is launched with no active stories, the expected output is:
  - `no stories`

## Canonical References
- Development cycle overview:
  - `DEVELOPMENT_CYCLE.md`
- Stage launch script:
  - `tools/launch_stage.sh`
- Observer script:
  - `tools/run_observer_cycle.sh`
- Stage seed prompts:
  - `stage-prompts/active/planning-seed-prompt.md`
  - `stage-prompts/active/architect-agent-seed-prompt.md`
  - `stage-prompts/active/next-agent-seed-prompt.md`
  - `stage-prompts/active/qa-agent-seed-prompt.md`
  - `stage-prompts/active/pm-refinement-seed-prompt.md`
- Backlog state directories:
  - `delivery-backlog/architecture/intake/`
  - `delivery-backlog/architecture/ready/`
  - `delivery-backlog/architecture/active/`
  - `delivery-backlog/architecture/qa/`
  - `delivery-backlog/architecture/done/`
  - `delivery-backlog/engineering/intake/`
  - `delivery-backlog/engineering/ready/`
  - `delivery-backlog/engineering/active/`
  - `delivery-backlog/engineering/qa/`
  - `delivery-backlog/engineering/done/`
  - `delivery-backlog/engineering/blocked/`
  - `delivery-backlog/engineering/archive/`
- Active queue ordering:
  - `delivery-backlog/engineering/active/README.md`
- Program control plane:
  - `product-research/roadmap/PROGRAM_STATE_BOARD.md`
  - `knowledge-base/process/PROGRAM_OPERATING_SYSTEM.md`
- Observer artifacts:
  - `operating-system/observer/README.md`
  - `operating-system/observer/OBSERVER_REPORT_TEMPLATE.md`
- Release checkpoint template:
  - `operating-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
- Personas directory and role index:
  - `staff-personas/`
  - `staff-personas/STAFF_DIRECTORY.md`
- Handoff docs:
  - `product-research/handoff.md`
  - `delivery-backlog/engineering/qa/HANDOFF-*.md`
  - `delivery-backlog/engineering/done/QA-RESULT-*.md`
- Doc test harness standard:
  - `product-research/roadmap/DOC_TEST_HARNESS_STANDARD.md`
