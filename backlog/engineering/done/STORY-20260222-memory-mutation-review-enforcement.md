# Story: Implement Memory Mutation Review Workflow Enforcement

## Metadata
- `id`: STORY-20260222-memory-mutation-review-enforcement
- `owner_persona`: personas/Product Manager - Maya.md
- `status`: done

## Problem Statement
- Governance policy and workflow contracts are defined, but engineering must enforce required mutation-review evidence and rejection handling in implementation workflows.

## Scope
- In:
- Implement mutation workflow guards so memory writes are only allowed in approved stage contexts.
- Require mutation decision evidence fields for approved/rejected outcomes.
- Enforce rejection behavior: blocked merge/apply until rework and re-review evidence exists.
- Add audit artifact generation/validation hooks required by QA.
- Out:
- Multi-user role systems beyond the current single-operator model.
- Full workflow orchestration UI.

## Acceptance Criteria
1. Memory mutation attempts outside allowed stage contexts are rejected with explicit reason.
2. Approved and rejected paths both produce required audit evidence fields.
3. Tests verify transition behavior for approval, rejection, and missing-evidence failure paths.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `backlog/architecture/done/ARCH-20260222-memory-mutation-review-workflow.md`
- `research/decisions/ADR-0011-memory-mutation-review-workflow-contract.md`
- `research/architecture/MEMORY_MUTATION_REVIEW_WORKFLOW.md`

## Notes
- Keep behavior aligned with existing stage/backlog transition conventions.
