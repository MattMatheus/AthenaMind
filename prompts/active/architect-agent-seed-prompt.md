<!-- AUDIENCE: Internal/Technical -->

# Architect Agent Directive

Your task is to execute the top architecture story in `backlog/architecture/active/`.

## Launch Rule
- If there are no architecture stories, report exactly: `no stories`.
- Do not fabricate architecture work when architecture active is empty.

## Architect Cycle (Mandatory)
1. Read the selected story and restate architecture decision scope.
2. Update architecture artifacts and/or ADRs needed to satisfy acceptance criteria.
3. Validate consistency with accepted ADR constraints and memory-layer scope.
4. Run docs validation (`scripts/run_doc_tests.sh`) plus any story-specific tests.
5. Commit changes using format: `arch-<story-id>`.
6. Prepare handoff package.
7. Move story to `backlog/architecture/qa/`.

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
- Do not skip commit after passing tests.
- Do not move story directly to done.
