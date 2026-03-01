# WSI-ADR-0002: Local-First Shared Workspace Control-Plane Authority and Transition Policy

## Status
Accepted

## Context
AthenaWork currently coordinates stage progress through markdown lane files and shell tools. Concurrent human and agent operations require a deterministic control-plane contract so lane transitions, state authority, and audit behavior are consistent.

## Decision
Define a local-first control-plane contract with these rules:

1. Authority model
- Backend control plane is authoritative for workflow state.
- Markdown lane files are derived projections from authoritative state.
- Exception: research-mode undocumented agent-to-agent communication is allowed only when recorded as auditable `research_comm_exception` events.

2. Transition policy
- All core stage transitions execute through a transition API with deterministic preconditions and machine-readable failure codes.
- Direction-changing transitions are hard-blocked until a human confirmation record is supplied and linked.
- Non-research lanes must reject undocumented agent-to-agent communication attempts.

3. Audit/event policy
- Every accepted or rejected transition emits immutable events with required correlation fields: `cycle_id`, `story_id`, `session_id`, and `trace_id`.
- Observer linkage is required for cycle-close visibility.

4. Optimization defaults
- Agent workflow: minimal transition friction, deterministic preconditions, concise machine-readable errors.
- Human workflow and UI: planning-first clarity and low-vision-friendly defaults.

## Consequences
- Positive:
  - Eliminates authority ambiguity between markdown and runtime state.
  - Enables deterministic automation and safer concurrency.
  - Creates explicit guardrails for research-only communication exceptions.
- Negative:
  - Adds schema and lifecycle governance overhead.
  - Requires synchronization logic from authoritative state to markdown views.

## Validation Plan
- Verify launcher/observer lifecycle transitions map to declared transition IDs.
- Verify non-research undocumented communication attempts return hard-block failures.
- Verify direction-changing transitions fail without human confirmation evidence.
- Verify emitted events always include required correlation fields.
