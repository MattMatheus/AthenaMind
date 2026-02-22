# AthenaMind v0.1 Module Boundaries

## Objective
Translate the three-pillar foundation into explicit module ownership and contracts for implementation planning.

## System Planes
- Memory Plane:
  - `procedural-memory`
  - `state-memory`
  - `semantic-memory`
  - `memory-governance`
- Execution Plane:
  - `workspace-runtime`
  - `semantic-navigation`
- Control Plane:
  - `orchestrator-api`
  - `audit-telemetry`

## Module Responsibilities

### procedural-memory
Owns:
- Load global and repo-local directives
- Validate structure and compute precedence merge

Does not own:
- Episodic state persistence
- Runtime execution

### state-memory
Owns:
- Session, episode, preference, and workflow state persistence
- Schema versioning and migrations

Does not own:
- Embedding/index mechanics
- Policy decisions

### semantic-memory
Owns:
- Embeddings, vector index operations, retrieval ranking inputs
- Retrieval adapters (local-first default)

Does not own:
- Final recall acceptance
- Preference promotion

### workspace-runtime
Owns:
- Ephemeral isolated execution environment lifecycle
- Repo/data clone + run + cleanup lifecycle

Does not own:
- Global memory updates
- Policy exceptions

### semantic-navigation
Owns:
- Symbolic/LSP navigation and targeted snippet extraction

Does not own:
- Long-term memory persistence

### memory-governance
Owns:
- Memory write authorization and promotion policy
- Conflict resolution and human-review pathways

Does not own:
- Raw runtime execution

### audit-telemetry
Owns:
- Trace and event capture for all memory/runtime decisions
- OTEL-aligned span and attribute schema for analysis

Does not own:
- Runtime/policy behavior decisions

### orchestrator-api
Owns:
- Unified API for run requests, memory operations, and orchestration
- Cross-module flow coordination and request-level enforcement

Does not own:
- Internal module data models beyond public contracts

## Core Data Flows
1. Run bootstrap:
- `orchestrator-api` -> `procedural-memory` -> `state-memory`/`semantic-memory` -> `workspace-runtime`

2. Runtime recall:
- `workspace-runtime`/`semantic-navigation` -> `semantic-memory` -> `memory-governance` -> `orchestrator-api`

3. Learning writeback:
- `workspace-runtime` -> `memory-governance` -> `state-memory` (+ `semantic-memory` reindex)

4. Audit stream:
- all modules -> `audit-telemetry`

## Non-Negotiable Boundary Rules
- No direct global-memory writes from runtime modules.
- No retrieval module can bypass governance.
- No module can skip audit emission for recall/write decisions.
- API-level orchestration cannot introduce hidden side effects outside declared contracts.

## v0.1 Implementation Sequencing
1. `procedural-memory`, `state-memory`, `orchestrator-api`
2. `workspace-runtime`, `audit-telemetry`
3. `semantic-navigation`, `semantic-memory`
4. `memory-governance` hardening and policy gates

## Validation Checklist
- Each module has a single owner and test boundary.
- Each cross-module call is interface-defined.
- End-to-end run can produce a complete recall/write audit chain.

## Phase Boundary Contract (v0.1-v0.3)
Use this as the canonical allocation map for backlog sequencing.

| Module | Phase Allocation | Ownership Note |
| --- | --- | --- |
| `procedural-memory` | `v0.1` now | Core memory-layer module; required for directive continuity. |
| `state-memory` | `v0.1` now | Core memory-layer module; owns state persistence and schema checks. |
| `semantic-memory` | `v0.1` now, `v0.2` next | v0.1 baseline retrieval; v0.2 quality/latency maturity work. |
| `memory-governance` | `v0.1` now, `v0.2` next | v0.1 review-first gates; v0.2 policy hardening and throughput tuning. |
| `audit-telemetry` | `v0.1` now, `v0.2` next | v0.1 minimum trace/KPI fields; v0.2 richer scoring and diagnostics. |
| `semantic-navigation` | `v0.2` next | Deferred until v0.1 reliability evidence is stable. |
| `orchestrator-api` | `v0.3` later | Memory API wrapper path; must preserve CLI-local fallback and reversibility. |
| `workspace-runtime` | external integration surface | Out of AthenaMind core ownership for v0.1; do not assign runtime ownership in v0.1 outputs. |

Decision reference:
- `research/decisions/ADR-0013-phase-boundary-contract-v01-v03.md`
