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

## Cross-Artifact Dependency Chain
1. Foundation and scope ADRs define what can exist in v0.1.
2. Module-boundary ADRs constrain ownership and enforcement points.
3. Constraints architecture binds design to measurable controls.
4. Memory CLI v1 architecture operationalizes file-memory interfaces for near-term implementation.

## Unresolved Gaps (Captured as Intake)
1. Retrieval quality evaluation and acceptance gates for semantic lookup:
   - `backlog/architecture/intake/ARCH-20260222-semantic-retrieval-quality-gates.md`
2. `/memory` schema and index evolution/versioning policy:
   - `backlog/architecture/intake/ARCH-20260222-memory-schema-versioning-policy.md`
3. Human-review workflow contract for memory mutation approvals:
   - `backlog/architecture/intake/ARCH-20260222-memory-mutation-review-workflow.md`

## Usage Guidance
- PM and Engineering should treat this map as the architecture entrypoint for v0.1 scoping.
- Any new architecture artifact should add itself to this map and link governing ADR(s).
