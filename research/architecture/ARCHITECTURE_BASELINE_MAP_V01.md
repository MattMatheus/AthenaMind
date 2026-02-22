# AthenaMind v0.1 Architecture Baseline Map

## Purpose
Provide one cross-reference map from architecture domains to governing ADRs and current architecture artifacts so implementation and planning can operate from a single baseline.

## Baseline Domains

### 1) Product Scope and Principles
- Governing ADRs:
  - `research/decisions/ADR-0001-three-pillar-foundation.md`
  - `research/decisions/ADR-0004-v01-mvp-constraints-framework.md`
  - `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- Architecture artifacts:
  - `research/architecture/FOUNDATION_PILLARS.md`
  - `research/architecture/MVP_CONSTRAINTS_V01.md`
- Notes:
  - v0.1 scope is memory-layer-first with explicit runtime exclusions.

### 2) Success Metrics and Delivery Targets
- Governing ADRs:
  - `research/decisions/ADR-0002-north-star-and-success-metrics.md`
  - `research/decisions/ADR-0008-v01-goals-and-scorecard-target-bands.md`
- Architecture artifacts:
  - `research/architecture/MVP_CONSTRAINTS_V01.md`
- Notes:
  - Metrics define acceptable quality and rollout gates; architecture must remain measurable.

### 3) Module Boundaries and Ownership
- Governing ADRs:
  - `research/decisions/ADR-0003-v01-module-boundaries.md`
  - `research/decisions/ADR-0006-governance-and-hitl-policy.md`
  - `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
- Architecture artifacts:
  - `research/architecture/MODULE_BOUNDARIES_V01.md`
- Notes:
  - Boundaries constrain cross-module behavior and prevent hidden side effects.

### 4) Memory Storage and Retrieval Stack
- Governing ADRs:
  - `research/decisions/ADR-0005-v01-local-storage-and-embedding-stack.md`
  - `research/decisions/ADR-0007-memory-layer-scope-refinement.md`
  - `research/decisions/ADR-0009-file-memory-cli-v1-architecture.md`
- Architecture artifacts:
  - `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
  - `research/architecture/FOUNDATION_PILLARS.md`
- Notes:
  - Current baseline is repo-local file memory rooted at `/memory` with Go CLI runtime and semantic-first retrieval + deterministic fallback.

### 5) Governance, HITL, and Mutation Policy
- Governing ADRs:
  - `research/decisions/ADR-0006-governance-and-hitl-policy.md`
  - `research/decisions/ADR-0009-file-memory-cli-v1-architecture.md`
- Architecture artifacts:
  - `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
  - `research/architecture/MVP_CONSTRAINTS_V01.md`
- Notes:
  - Pre-MVP memory mutation requires reviewed workflows outside autonomous runs.

## Planned vs Implemented Status Matrix

| Baseline domain | Planned baseline | Implemented now | Open execution links |
| --- | --- | --- | --- |
| Product Scope and Principles | Memory-layer-only v0.1 scope and exclusions | Implemented via ADR set and baseline artifacts | None |
| Success Metrics and Delivery Targets | KPI snapshot + release checkpoint evidence | Partially implemented (targets defined, runtime reporting in progress) | `backlog/engineering/active/STORY-20260222-kpi-snapshot-baseline-v01.md` |
| Module Boundaries and Ownership | Phase boundary contract and enforceable constraints | Partially implemented (contract finalized, enforcement in progress) | `backlog/engineering/active/STORY-20260222-mvp-constraint-enforcement-v01.md` |
| Memory Storage and Retrieval Stack | File-memory CLI v1 baseline with telemetry, v0.2 snapshot, v0.3 read gateway | Partially implemented (CLI baseline complete, telemetry and phase-next execution in progress) | `backlog/engineering/active/STORY-20260222-memory-cli-telemetry-contract-v01.md`, `backlog/engineering/active/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`, `backlog/engineering/active/STORY-20260222-memory-api-read-gateway-v03.md` |
| Governance, HITL, and Mutation Policy | Human-review mutation workflow + release checkpoint bundle | Partially implemented (workflow contract complete, checkpoint packaging in progress) | `backlog/engineering/active/STORY-20260222-release-checkpoint-bundle-v01.md` |

## Cross-Artifact Dependency Chain
1. Foundation and scope ADRs define what can exist in v0.1.
2. Module-boundary ADRs constrain ownership and enforcement points.
3. Constraints architecture binds design to measurable controls.
4. Memory CLI v1 architecture operationalizes file-memory interfaces for near-term implementation.

## Phase Boundary Contract (Current)
- Canonical allocation and ownership contract:
  - `research/decisions/ADR-0013-phase-boundary-contract-v01-v03.md`
  - `research/architecture/MODULE_BOUNDARIES_V01.md` section `Phase Boundary Contract (v0.1-v0.3)`
- Summary:
  - `v0.1`: memory-layer reliability + governance in local-first mode.
  - `v0.2`: semantic quality maturity + telemetry depth.
  - `v0.3`: optional API-wrapper/cloud-ready extensions with CLI fallback.
  - `workspace-runtime`: external integration surface, not v0.1 core ownership.

## Resolved v0.1 Architecture Gaps
1. Retrieval quality evaluation and acceptance gates:
   - `research/decisions/ADR-0012-semantic-retrieval-quality-gates-v1.md`
2. `/memory` schema and version compatibility policy:
   - `research/decisions/ADR-0010-memory-schema-versioning-policy.md`
3. Human-review mutation workflow contract:
   - `research/decisions/ADR-0011-memory-mutation-review-workflow-contract.md`

## Open Architecture Work (Current Queue)
1. Architecture intake currently has no open architecture stories:
   - `backlog/architecture/active/README.md` (active sequence is empty)
   - `backlog/architecture/intake/ARCH_STORY_TEMPLATE.md` (template only)
2. Architecture definitions are complete for current v0.1-v0.3 planning set and now tracked through implementation stories:
   - `backlog/architecture/done/ARCH-20260222-mvp-constraint-contract-hardening.md`
   - `backlog/architecture/done/ARCH-20260222-memory-telemetry-event-contract-v01.md`
   - `backlog/architecture/done/ARCH-20260222-memory-snapshot-architecture-v02.md`
   - `backlog/architecture/done/ARCH-20260222-memory-api-wrapper-contract-v03.md`
3. Remaining open work is implementation-side in engineering active queue:
   - `backlog/engineering/active/STORY-20260222-memory-cli-telemetry-contract-v01.md`
   - `backlog/engineering/active/STORY-20260222-mvp-constraint-enforcement-v01.md`
   - `backlog/engineering/active/STORY-20260222-kpi-snapshot-baseline-v01.md`
   - `backlog/engineering/active/STORY-20260222-release-checkpoint-bundle-v01.md`
   - `backlog/engineering/active/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
   - `backlog/engineering/active/STORY-20260222-memory-api-read-gateway-v03.md`

## Usage Guidance
- PM and Engineering should treat this map as the architecture entrypoint for v0.1 scoping.
- Any new architecture artifact should add itself to this map and link governing ADR(s).
