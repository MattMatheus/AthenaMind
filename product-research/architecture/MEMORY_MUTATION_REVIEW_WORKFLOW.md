# Memory Mutation Review Workflow (v1)

## Purpose
Define the enforceable review contract for memory mutations outside autonomous runs.

## Workflow Sequence
1. Proposer prepares mutation package:
   - target entries/artifacts
   - reason for change
   - risk/rollback notes
2. Reviewer validates package against governance requirements.
3. Reviewer records decision:
   - `approved` with rationale, or
   - `rejected` with required rework notes
4. On approval, mutation artifacts are committed and handed off.
5. QA verifies evidence and transition correctness.

## Actor Responsibilities
- Proposer:
  - prepare complete mutation evidence package
  - ensure no mutation occurs during autonomous run execution
- Reviewer (human operator):
  - approve/reject with explicit rationale
  - ensure risk and rollback notes are present
- QA:
  - validate artifact completeness
  - enforce failure on missing decision evidence

## Required Audit Artifacts
For each mutation:
- list of changed files or entries
- reason-for-change statement
- risk/mitigation note
- decision record:
  - `decision`: `approved` | `rejected`
  - `reviewed_by`
  - `reviewed_at`
  - `decision_notes`

## Approval/Rejection Backlog Mapping
- Approved:
  - story continues normal progression (`active -> qa -> done`)
  - memory artifacts can advance to approved status
- Rejected:
  - return story/item to refinement or active rework
  - preserve rejection reason in handoff/audit docs
  - no mutation merge allowed until re-review

## QA Validation Gates
QA must verify:
1. mutation was initiated in an allowed stage context
2. required audit artifacts exist and are complete
3. approval/rejection outcome matches backlog transition
4. rejection path includes actionable rework notes

If any gate fails, QA result is FAIL and item returns to active/rework.

## References
- `product-research/decisions/ADR-0006-governance-and-hitl-policy.md`
- `product-research/decisions/ADR-0011-memory-mutation-review-workflow-contract.md`
- `product-research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
