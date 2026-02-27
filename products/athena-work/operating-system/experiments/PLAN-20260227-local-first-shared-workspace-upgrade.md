# Planning Session: Local-First Shared Workspace Upgrade

## Metadata
- `id`: PLAN-20260227-local-first-shared-workspace-upgrade
- `status`: finalized
- `owner`: Founder Operator
- `date_utc`: 2026-02-27
- `phase_target`: v0.2

## Problem Framing
Human operators are evolving AthenaWork manually across markdown lanes and scripts. This causes state drift, onboarding friction, and weak shared visibility across human+agent collaboration. The current model works for local single-operator execution but is not robust for a shared workspace UX.

## Target Outcomes
1. Preserve local-first operation with Docker Compose as the default runtime.
2. Introduce a durable state backend for concurrent human+agent collaboration.
3. Keep markdown artifacts as auditable outputs, not the primary mutable state.
4. Provide a visual board and timeline so operators can navigate workflow state without reading many docs.
5. Provide an optimal work environment for both humans and agents, with equal priority on operator clarity and agent execution reliability.

## Constraints
- Local-first must remain first-class.
- Existing AthenaWork stage semantics and gates must remain enforceable.
- Migration must be incremental and reversible.
- Existing scripts should continue to work during transition.
- Undocumented agent-to-agent communication is permitted only during research-mode execution and must be auditable.
- Direction-changing workflow actions must include explicit human confirmation before execution.
- Implementation workflow must be optimized for agents (low-friction execution, deterministic gates, minimal ambiguity).
- Planning workflow and UI must be optimized for humans (clear intent, short paths, explicit next actions).
- UI must be low-vision-friendly by default: high contrast, large typography, uncluttered layout, and direct language.

## Options Considered
1. Keep file-only state machine.
- Rejected for shared use due to drift and lack of transactional guarantees.
2. Replace everything with cloud-native SaaS backend.
- Rejected due to local-first requirement.
3. Hybrid local control plane (Docker Compose) with backend state and markdown sync.
- Selected baseline.

## Selected Direction
Adopt a local-first control plane:
- `workspace-api` enforces state transitions and gates.
- Document database stores canonical workflow documents.
- Event log captures immutable transition history.
- `workspace-ui` provides visual board/timeline.
- Sync worker exports canonical markdown artifacts back into repo paths.

## Delivery Plan (AthenaWork)

### Stage 1: Architect (decision and contract hardening)
1. Define backend state contract and transition API.
2. Define markdown synchronization policy and conflict semantics.
3. Define Compose topology and local ops contract.
4. Define research-mode communication exception policy boundaries.
5. Define human direction-confirmation contract and artifact model.

### Stage 2: PM (intake refinement and sequencing)
1. Rank implementation stories in dependency order.
2. Set release checkpoint criteria for v0.2 pilot.

### Stage 3: Engineering + QA cycles
1. Compose bootstrap and dev/runtime scripts.
2. State machine API and transition gate engine.
3. Read-only workspace UI (board + cycle timeline).
4. Markdown sync worker + drift guard enforcement.
5. CLI adapter for launcher/observer compatibility.
6. End-to-end validation and rollout playbook.

## Intake Artifacts Created
Architecture intake:
- `products/athena-work/delivery-backlog/architecture/intake/ARCH-20260227-shared-workspace-control-plane-contract.md`
- `products/athena-work/delivery-backlog/architecture/intake/ARCH-20260227-markdown-sync-authority-and-conflict-policy.md`

Engineering intake:
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-docker-compose-local-control-plane-bootstrap.md`
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-workspace-api-state-machine-v1.md`
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-workspace-ui-read-only-board-v1.md`
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-markdown-sync-worker-and-drift-guard-v1.md`
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-cli-adapter-for-launcher-observer.md`
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-human-direction-confirmation-gate.md`
- `products/athena-work/delivery-backlog/engineering/intake/STORY-20260227-research-communication-exception-audit.md`

## Risks and Mitigations
- Drift between DB and markdown.
  - Mitigation: single-writer sync policy + deterministic export + drift tests.
- Operator confusion during migration.
  - Mitigation: read-only UI first, then guarded write paths.
- Over-complex first release.
  - Mitigation: v0.2 scope cap with read-only UI and essential transitions only.

## Success Metrics
- `operator_navigation_time_p50`: <= 60s to identify next actionable item.
- `state_drift_incidents_per_cycle`: 0 in pilot cycles.
- `workspace_transition_success_rate`: >= 99% for valid transitions.
- `local_startup_time_compose`: <= 120s cold start on baseline laptop.
- `agent_stage_execution_overhead_seconds_p50`: <= 15s from launch to actionable task context.
- `human_planning_completion_time_p50`: <= 10 minutes to capture/confirm direction and produce next-stage decision.
- `ui_accessibility_low_vision_pass_rate`: 100% for defined baseline checks (contrast, font size, spacing, focus clarity).

## Recommended Next Stage
`architect`

Rationale: shared workspace requires contract-level decisions (state authority, sync policy, transition semantics) before PM can safely rank implementation stories.
