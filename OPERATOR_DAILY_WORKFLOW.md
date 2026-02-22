# Founder-Operator Daily Workflow (Codex App)

Single-operator daily script for running AthenaMind delivery cycles in v0.1.

## Scope
- Local single-operator workflow only.
- Team/multi-operator workflows are deferred to v2.0.

## Startup Routine
1. Open workspace at repo root and confirm branch discipline:
   - Run `git branch --show-current`
   - If branch is not `dev`, switch/fix before running stages.
2. Check current cycle context:
   - Read `CYCLE_INDEX.md`
   - Read `DEVELOPMENT_CYCLE.md`
3. Run baseline docs checks:
   - `scripts/run_doc_tests.sh`
4. Run Go toolchain preflight for memory CLI work:
   - `scripts/check_go_toolchain.sh`
5. Confirm active queue:
   - `backlog/engineering/active/README.md`
6. Confirm program board context:
   - `research/roadmap/PROGRAM_STATE_BOARD.md`

## Planning Stage Loop (As Needed)
1. Launch:
   - `scripts/launch_stage.sh planning`
2. Follow returned seed prompt:
   - `prompts/active/planning-seed-prompt.md`
3. Execute interactive planning session:
   - capture notes in `research/planning/sessions/` using `research/planning/PLANNING_SESSION_TEMPLATE.md`
   - convert ideas into engineering and/or architecture intake stories
   - recommend next stage (`architect` and/or `pm`) based on decision needs
   - mark planning session `status: finalized` after intake artifacts are linked
4. Close planning cycle:
   - run `scripts/run_observer_cycle.sh --cycle-id <plan-id>`
   - commit once: `cycle-<cycle-id>`

## Engineering + QA Cycle Loop
1. Launch engineering:
   - `scripts/launch_stage.sh engineering`
2. Follow returned seed prompt:
   - `prompts/active/next-agent-seed-prompt.md`
3. Execute top active story:
   - implement required artifacts
   - update tests
   - run `scripts/run_doc_tests.sh` plus any story-specific tests
   - move story from `backlog/engineering/active/` to `backlog/engineering/qa/` with handoff package
   - do not commit yet
4. Launch QA:
   - `scripts/launch_stage.sh qa`
5. Follow returned seed prompt:
   - `prompts/active/qa-agent-seed-prompt.md`
6. For top `backlog/engineering/qa/` story:
   - validate acceptance criteria and regression risk
   - if defects exist, file bug(s) in `backlog/engineering/intake/` and return story to `backlog/engineering/active/`
   - if quality bar is met, move story to `backlog/engineering/done/`
7. Close cycle:
   - run `scripts/run_observer_cycle.sh --cycle-id <story-id>`
   - commit once: `cycle-<cycle-id>`

## Architect Stage Loop (As Needed)
1. Launch:
   - `scripts/launch_stage.sh architect`
2. Follow returned seed prompt:
   - `prompts/active/architect-agent-seed-prompt.md`
3. Execute top architecture story:
   - update ADRs/architecture docs
   - run docs validation tests
   - move story from `backlog/architecture/active/` to `backlog/architecture/qa/`
4. Close architecture cycle:
   - run `scripts/run_observer_cycle.sh --cycle-id <arch-story-id>`
   - commit once: `cycle-<cycle-id>`

## PM Refinement Loop
1. Launch:
   - `scripts/launch_stage.sh pm`
2. Follow returned seed prompt:
   - `prompts/active/pm-refinement-seed-prompt.md`
3. Refine intake and re-rank active queue in `backlog/engineering/active/README.md`.
4. Update `research/roadmap/PROGRAM_STATE_BOARD.md` counts and Now/Next priorities.
5. Ensure intake/active stories have traceability metadata (`idea_id`, `phase`, `adr_refs`, metric).
6. Close PM cycle:
   - run `scripts/run_observer_cycle.sh --cycle-id PM-<date>-<slug>`
   - commit once: `cycle-<cycle-id>`

## If X Then Y Rules
- If engineering launch returns `no stories`:
  - Do not fabricate work.
  - Run PM refinement to move/refine intake items into active.
- If QA finds blocking defects:
  - File `P0-P3` bugs via `backlog/engineering/intake/BUG_TEMPLATE.md`.
  - Move story back to `backlog/engineering/active/`.
  - Prioritize bug resolution before further promotion.
- If tests fail:
  - Do not move state forward.
  - Fix failures first, then re-run test commands.
- If a process gap is discovered mid-story:
  - Add intake story via `backlog/engineering/intake/STORY_TEMPLATE.md` before handoff.

## Escalation Rules
- Use command escalation when required by execution environment policy.
- Do not bypass sandbox/security controls; request approval with clear command purpose.
- Do not use destructive git/file operations unless explicitly requested.

## Shutdown Routine
1. Confirm no partial state transitions remain.
2. Confirm observer report exists for each completed cycle in this session.
3. Ensure latest cycle changes are committed with `cycle-<cycle-id>` messages.
4. Leave next actionable queue visible in `backlog/engineering/active/README.md`.
5. Capture any new ideas/defects in `backlog/engineering/intake/` before ending session.
6. If items should be considered shipped, produce a release bundle from `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`.
