# ADR-0011: Memory Mutation Review Workflow Contract

## Status
Accepted

## Context
Governance requires memory mutations to occur outside autonomous runs and under human review, but the workflow contract, actor responsibilities, evidence requirements, and rejection path were not explicit.

## Decision
Adopt a formal mutation review workflow for v0.1/v1.0:

- Allowed initiation stages:
  - `planning`
  - `architect`
  - `pm`
- Required actors:
  - `proposer`: drafts mutation request and artifacts.
  - `reviewer`: human operator with approval authority.
  - `qa`: verifies workflow evidence and transition compliance.

Approval evidence minimum:
1. Proposed diff or artifact delta list.
2. Reason-for-change statement linked to architecture/roadmap context.
3. Risk assessment with mitigation and rollback note.
4. Review decision metadata (`approved` or `rejected`, reviewer id, timestamp).

Outcome transitions:
- Approved:
  - mutation artifacts can be committed
  - memory entry status may transition to `approved`
  - story proceeds through normal stage transitions
- Rejected:
  - mutation is not applied
  - item returns to backlog refinement or active rework path
  - rejection reason becomes mandatory audit evidence

Audit contract:
- Every mutation decision must be traceable in committed artifacts or handoff docs.
- Missing evidence is a QA-fail condition.

## Consequences
- Positive:
  - Makes governance enforceable and testable.
  - Reduces ambiguity in review/rework decisions.
- Negative:
  - Adds process overhead for high-frequency edits.
- Neutral:
  - Workflow may later be automated, but evidence requirements remain.

## Alternatives Considered
- Option A: reviewer discretion without required evidence checklist
  - Rejected due to inconsistent audit quality.
- Option B: auto-approve low-risk edits
  - Rejected for v0.1 due to trust/safety priority.

## Validation Plan
- Architecture artifact must define sequence and required artifacts.
- QA checks must fail when approval/rejection evidence is incomplete.
