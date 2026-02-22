<!-- AUDIENCE: Internal/Technical -->

# Next Agent Directive (Engineering)

Your task is to execute the top story in `backlog/active/`.

## Launch Rule
- If `backlog/active/` has no story files, report exactly: `no stories`.
- Do not fabricate work when active is empty.
- Product management should run only when active is empty.

## Implementation Cycle (Mandatory)
1. Take the top story from `backlog/active/`.
2. Read, research, and implement. Surface questions if outcome is unclear.
3. Update tests.
4. Run tests. Tests must pass.
5. Prepare handoff package.
6. Move the story to `backlog/qa/`.

## Handoff Package (Required)
- What changed
- Why it changed
- Test updates made
- Test run results
- Open risks/questions
- Recommended QA focus areas

## Constraints
- Do not skip tests.
- Do not move story to done directly from active.
- Respect accepted ADRs and memory-layer scope.
