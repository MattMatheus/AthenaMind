# Cycle Index

Single navigation entrypoint for operators and agents running the AthenaMind delivery cycle.

## First 5 Minutes
1. Confirm branch safety: run `git branch --show-current` and verify `dev`.
2. Launch the stage you need:
   - `scripts/launch_stage.sh planning`
   - `scripts/launch_stage.sh architect`
   - `scripts/launch_stage.sh engineering`
   - `scripts/launch_stage.sh qa`
   - `scripts/launch_stage.sh pm`
3. Open the returned seed prompt in `prompts/active/` and follow it as directive.
4. Run docs validation command before handoff/decision points:
   - `scripts/run_doc_tests.sh`
5. Apply backlog state movement only through the canonical flow (`engineering/active -> engineering/qa -> engineering/done`, with intake/active loop for defects).

## Branch Rule and Empty Active Behavior
- Branch safety rule: launcher requires branch `dev`.
- If branch differs, launcher aborts:
  - `abort: active branch is '<branch>'; expected 'dev'`
- If engineering is launched with no active stories, the expected output is:
  - `no stories`

## Canonical References
- Development cycle overview:
  - `DEVELOPMENT_CYCLE.md`
- Stage launch script:
  - `scripts/launch_stage.sh`
- Stage seed prompts:
  - `prompts/active/planning-seed-prompt.md`
  - `prompts/active/architect-agent-seed-prompt.md`
  - `prompts/active/next-agent-seed-prompt.md`
  - `prompts/active/qa-agent-seed-prompt.md`
  - `prompts/active/pm-refinement-seed-prompt.md`
- Backlog state directories:
  - `backlog/architecture/intake/`
  - `backlog/architecture/ready/`
  - `backlog/architecture/active/`
  - `backlog/architecture/qa/`
  - `backlog/architecture/done/`
  - `backlog/engineering/intake/`
  - `backlog/engineering/ready/`
  - `backlog/engineering/active/`
  - `backlog/engineering/qa/`
  - `backlog/engineering/done/`
  - `backlog/engineering/blocked/`
  - `backlog/engineering/archive/`
- Active queue ordering:
  - `backlog/engineering/active/README.md`
- Personas directory and role index:
  - `personas/`
  - `personas/STAFF_DIRECTORY.md`
- Handoff docs:
  - `research/handoff.md`
  - `backlog/engineering/qa/HANDOFF-*.md`
  - `backlog/engineering/done/QA-RESULT-*.md`
- Doc test harness standard:
  - `research/roadmap/DOC_TEST_HARNESS_STANDARD.md`
