<!-- AUDIENCE: Internal/Technical -->

# Next Agent Directive (Engineering)

Your task is to execute the top story in `backlog/engineering/active/`.

## Launch Rule
- If `backlog/engineering/active/` has no story files, report exactly: `no stories`.
- Do not fabricate work when active is empty.
- Product management should run only when active is empty.

## Implementation Cycle (Mandatory)
1. Take the top story from `backlog/engineering/active/`.
2. Read, research, and implement. Surface questions if outcome is unclear.
   - Apply `backlog/STATE_TRANSITION_CHECKLIST.md` before moving story state.
3. Update tests.
4. Run tests with the canonical docs command (`scripts/run_doc_tests.sh`) plus any story-specific test commands. Tests must pass.
5. Commit changes with an appropriate message that includes the story id.
6. Prepare handoff package.
7. Move the story to `backlog/engineering/qa/`.

## Handoff Package (Required)
- What changed
- Why it changed
- Test updates made
- Test run results
- Open risks/questions
- Recommended QA focus areas
- New gaps discovered during implementation (as intake story paths in `backlog/engineering/intake/`)

## Constraints
- Do not skip tests.
- Use `scripts/run_doc_tests.sh` as the default docs validation entrypoint.
- Do not skip commit after passing tests.
- Do not move story to done directly from active.
- Respect accepted ADRs and memory-layer scope.
- If a gap is discovered, log a new intake story before handoff.
