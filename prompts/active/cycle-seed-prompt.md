<!-- AUDIENCE: Internal/Technical -->

# Cycle Seed Directive

Your task is to drain the engineering active backlog by alternating engineering and QA cycles.

## Cycle Loop (Mandatory)
1. Run `scripts/launch_stage.sh engineering`.
2. If output is exactly `no stories`, stop and report completion.
3. Execute the engineering cycle for the selected story.
4. Run `scripts/launch_stage.sh qa`.
5. Execute the QA cycle for the story in `backlog/engineering/qa/`.
6. Repeat from step 1 until `backlog/engineering/active/` is drained.

## Commit Discipline
- Engineering commits must include story id.
- QA commits must use: `qa-<story-id>` only.

## Constraints
- Do not skip tests.
- Do not bypass backlog states.
- Do not continue if branch is not `dev`.
