# Architecture Handoff: ARCH-20260227-markdown-sync-authority-and-conflict-policy

## Decision(s) Made
- Accepted markdown sync authority ADR:
  - `operating-system/decisions/WSI-ADR-0003-markdown-sync-authority-and-conflict-policy.md`
- Published sync/conflict policy contract:
  - `operating-system/contracts/MARKDOWN_SYNC_AUTHORITY_AND_CONFLICT_POLICY_V1.md`

## Alternatives Considered
- Fully flexible markdown-first edits with eventual reconciliation.
  - Rejected due to hidden mutation and drift risk.
- Strict no-exceptions model for markdown edits.
  - Rejected because audited operator exceptions are required for practical operations.

## Risks and Mitigations
- Risk: strict authority reduces human flexibility.
  - Mitigation: explicit operator override hooks with audit records.
- Risk: sync checks add latency.
  - Mitigation: deterministic taxonomy and auto-resolution for common cases.
- Risk: unresolved drift compromises stage integrity.
  - Mitigation: stage-critical blocking rules and observer-linked alarms.

## Acceptance Criteria Mapping
1. Single-writer authority explicit with exceptions: satisfied in ADR and sync matrix.
2. Deterministic conflict handling: satisfied via taxonomy + deterministic resolutions.
3. Drift alarms and blocking conditions: satisfied with critical artifact list and blocking codes.
4. Human clarity + agent latency goals: satisfied with projection and consumption targets.

## Updated Artifacts
- `operating-system/decisions/WSI-ADR-0003-markdown-sync-authority-and-conflict-policy.md`
- `operating-system/contracts/MARKDOWN_SYNC_AUTHORITY_AND_CONFLICT_POLICY_V1.md`
- `delivery-backlog/architecture/qa/ARCH-20260227-markdown-sync-authority-and-conflict-policy.md`

## Validation Commands and Results
- `tools/run_doc_tests.sh` -> PASS

## Open Questions for QA Focus
- Are conflict classes mutually exclusive and complete for current workflow operations?
- Are stage-critical blocking rules strict enough to prevent hidden state divergence?
- Do override hooks provide enough governance evidence without operator overload?
