<!-- AUDIENCE: Internal/Technical -->

# Next Agent Directive (Engineering)

Your task is to execute the top story in `delivery-backlog/engineering/active/`.

## Launch Rule
- If `delivery-backlog/engineering/active/` has no story files, report exactly: `no stories`.
- Do not fabricate work when active is empty.
- Product management should run only when active is empty.

## Implementation Cycle (Mandatory)
1. Take the top story from `delivery-backlog/engineering/active/`.
2. Read, research, and implement. Surface questions if outcome is unclear.
   - Apply `delivery-backlog/STATE_TRANSITION_CHECKLIST.md` before moving story state.
3. Update tests.
4. Run `tools/run_stage_tests.sh` (auto-scoped push vs PR) plus targeted story-specific test commands. Tests must pass.
5. Prepare handoff package.
6. Move the story to `delivery-backlog/engineering/qa/`.
7. Do not commit yet; cycle commit occurs only after QA + observer.

## Handoff Package (Required)
- What changed
- Why it changed
- Test updates made
- Test run results
- Open risks/questions
- Recommended QA focus areas
- New gaps discovered during implementation (as intake story paths in `delivery-backlog/engineering/intake/`)

## Constraints
- Do not skip tests.
- Use `tools/run_doc_tests.sh` as the default docs validation entrypoint.
- Do not move story to done directly from active.
- Respect accepted ADRs and memory-layer scope.
- If a gap is discovered, log a new intake story before handoff.
- Apply stage exit requirements in `knowledge-base/process/STAGE_EXIT_GATES.md`.
