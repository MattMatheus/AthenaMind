# Matt Action Queue: 2026-02-27

Ordered list of actions that require Matt (human operator/approver) before agents can run at full speed.

## 1) Approve architecture stories currently in QA
- [ ] Review and accept:
  - `delivery-backlog/architecture/qa/ARCH-20260227-shared-workspace-control-plane-contract.md`
  - `delivery-backlog/architecture/qa/ARCH-20260227-markdown-sync-authority-and-conflict-policy.md`
- [ ] Validate handoffs:
  - `delivery-backlog/architecture/qa/HANDOFF-ARCH-20260227-shared-workspace-control-plane-contract.md`
  - `delivery-backlog/architecture/qa/HANDOFF-ARCH-20260227-markdown-sync-authority-and-conflict-policy.md`
- [ ] Decision needed from Matt:
  - `PASS`: move both ARCH stories (and handoffs) to `delivery-backlog/architecture/done/`.
  - `FAIL`: return story to `delivery-backlog/architecture/active/` with gap notes.

## 2) Confirm architecture prioritization for next wave
- [ ] Decide whether to activate:
  - `delivery-backlog/architecture/intake/ARCH-20260227-kanban-workspace-information-architecture-v1.md`
- [ ] Decision needed from Matt:
  - `YES`: promote to `delivery-backlog/architecture/active/` and place first in active README.
  - `NO`: keep in intake and prioritize direct engineering follow-ons instead.

## 3) Provide PM direction for engineering queue refill
- [ ] Confirm Friday push objective:
  - prioritize AthenaWork v0.2 workspace throughput over additional v0.1 cleanup.
- [ ] Approve creation/promotion of a fresh engineering active sequence from current architecture outputs.
- [ ] Decision needed from Matt:
  - rank order for first 3 engineering stories (agent execution order).

## 4) Direction-change confirmation artifact (human-only gate)
- [ ] Confirm direction-change scope and publish an explicit confirmation record.
- [ ] Required fields (per control-plane contract):
  - `confirmation_id`
  - `confirmed_by` (Matt)
  - `confirmed_at` (UTC timestamp)
  - `scope`
  - `expiry`
- [ ] Use existing artifact as baseline reference:
  - `products/athena-work/operating-system/handoff/HUMAN_APPROVAL-2026-02-27-matt-local-first-workspace.md`

## 5) Security launch gate inputs (human secrets/approval)
- [ ] If protected contracts are part of this push, provide:
  - fresh nonce (`./tools/generate_security_nonce.sh`)
  - primary OTP
  - security OTP
  - authorization phrase: `LAUNCH AUTHORIZED`
- [ ] Reason:
  - current launch auth package is blocked by security gate:
  - `products/athena-work/operating-system/handoff/LAUNCH_AUTHZ_20260227T132941Z.json`

## 6) Release posture decision for this session
- [ ] Choose explicit posture before heavy execution:
  - `delivery-only` (build/QA backlog, no release promotion), or
  - `release-candidate` (prepare updated release bundle and signoff evidence).
- [ ] If release-candidate, Matt must approve `dev -> main` promotion only after bundle review.

## 7) Go/No-Go call for agent sprint
- [ ] Final Matt call:
  - `GO`: proceed with architecture closure + PM refinement + engineering/QA cycle execution.
  - `NO-GO`: hold with listed blockers and revisit after missing approvals.

## Fast Start Recommendation
1. Mark Item 1 as PASS/FAIL for both ARCH stories.
2. Set Item 2 (activate Kanban architecture story: YES/NO).
3. Set Item 3 rank order for top 3 engineering stories.
4. Issue Item 4 direction confirmation artifact.
5. Set Item 7 to GO.
