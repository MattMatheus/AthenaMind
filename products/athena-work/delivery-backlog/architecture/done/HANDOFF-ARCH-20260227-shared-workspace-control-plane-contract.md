# Architecture Handoff: ARCH-20260227-shared-workspace-control-plane-contract

## Decision(s) Made
- Accepted control-plane authority and transition policy ADR:
  - `operating-system/decisions/WSI-ADR-0002-control-plane-authority-and-transition-policy.md`
- Published control-plane architecture contract with entity, transition API, and event model:
  - `operating-system/contracts/LOCAL_FIRST_SHARED_WORKSPACE_CONTROL_PLANE_CONTRACT_V1.md`

## Alternatives Considered
- Keep markdown as authoritative state source.
  - Rejected due to concurrency ambiguity.
- Permit undocumented agent communication in all lanes.
  - Rejected due to audit and governance risk.

## Risks and Mitigations
- Risk: projection drift between backend and markdown lanes.
  - Mitigation: markdown is explicitly derived and non-authoritative.
- Risk: additional transition friction for operators.
  - Mitigation: deterministic preconditions and concise machine-readable errors.
- Risk: direction shifts without human review.
  - Mitigation: hard-block on direction-changing transitions until explicit confirmation reference is provided.

## Acceptance Criteria Mapping
1. Canonical state authority explicit: satisfied in ADR + contract authority section.
2. Transition API/failure semantics: satisfied in transition contract and failure code set.
3. Event correlation fields: satisfied in immutable event contract.
4. Research-only exception path: satisfied via `research_comm_exception` policy.
5. Human direction-confirmation checkpoint: satisfied via `DIRECTION_CHANGE_REQUEST` + hard-block policy.
6. Agent workflow optimization goals: satisfied via deterministic preconditions and concise error contract.

## Updated Artifacts
- `operating-system/decisions/WSI-ADR-0002-control-plane-authority-and-transition-policy.md`
- `operating-system/contracts/LOCAL_FIRST_SHARED_WORKSPACE_CONTROL_PLANE_CONTRACT_V1.md`
- `delivery-backlog/architecture/qa/ARCH-20260227-shared-workspace-control-plane-contract.md`

## Validation Commands and Results
- `tools/run_doc_tests.sh` -> PASS

## Open Questions for QA Focus
- Do all live AthenaWork lane transitions map to declared transition types without semantic loss?
- Are failure codes sufficient for deterministic automation and operator diagnosis?
- Is research-only exception audit payload specific enough for governance review?
