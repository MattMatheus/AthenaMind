<!-- AUDIENCE: Internal/Technical -->

# Architect Agent Directive

Your task is to execute the top architecture story in `delivery-backlog/architecture/active/`.

## Launch Rule
- If there are no architecture stories, report exactly: `no stories`.
- Do not fabricate architecture work when architecture active is empty.

## Architect Cycle (Mandatory)
1. Read the selected story and restate architecture decision scope.
2. Update architecture artifacts and/or ADRs needed to satisfy acceptance criteria.
3. Validate consistency with accepted ADR constraints and memory-layer scope.
4. Run docs validation (`tools/run_doc_tests.sh`) plus any story-specific tests.
5. Add explicit follow-on implementation story paths for each accepted decision.
6. Prepare handoff package.
7. Move story to `delivery-backlog/architecture/qa/`.
8. Run observer:
   - `tools/run_observer_cycle.sh --cycle-id <arch-story-id>`
9. Commit once for this cycle:
   - `cycle-<cycle-id>`

## Handoff Package (Required)
- Decision(s) made
- Alternatives considered
- Risks and mitigations
- Updated artifacts/paths
- Validation commands and results
- Open questions for QA focus

## Constraints
- Do not implement runtime-execution ownership in v0.1 scope.
- Do not skip tests.
- Do not move story directly to done.
- Apply stage exit requirements in `knowledge-base/process/STAGE_EXIT_GATES.md`.
- Do not commit before observer report is generated.
