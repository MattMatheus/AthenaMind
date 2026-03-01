<!-- AUDIENCE: Internal/Technical -->

# Cycle Seed Directive

Your task is to drain the engineering active backlog by alternating engineering and QA cycles.

## Accessibility Profile (Mandatory)
Before step 1, select and record active profile(s) for this cycle:
- `low-vision-default`
- `high-variability-attention`
- both (if needed)
Maintain single-task visibility and concise pass/fail markers throughout the loop.

## Cycle Loop (Mandatory)
1. Run `tools/launch_stage.sh engineering`.
2. If output is exactly `no stories`, stop and report completion.
3. Execute the engineering cycle for the selected story.
4. Run `tools/launch_stage.sh qa`.
5. Execute the QA cycle for the story in `delivery-backlog/engineering/qa/`.
6. Run observer at cycle boundary:
   - `tools/run_observer_cycle.sh --cycle-id <story-id> --story <path-to-story>`
7. Commit once for the full cycle:
   - `cycle-<cycle-id>`
8. Repeat from step 1 until `delivery-backlog/engineering/active/` is drained.

## Commit Discipline
- Do not commit during intermediate stage transitions.
- Use exactly one commit per completed cycle.
- Commit format: `cycle-<cycle-id>`.

## Constraints
- Do not skip tests.
- Do not bypass backlog states.
- Do not continue if branch is not `dev`.
