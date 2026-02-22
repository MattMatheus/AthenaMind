# AthenaMind Foundation: Three-Pillar Memory Architecture

## Purpose
This document defines the baseline architecture for AthenaMind and is the default reference for product and engineering decisions.

## Core Thesis
AthenaMind should treat memory and execution as a unified system made of three pillars:
1. Procedural Memory (versioned Markdown)
2. State & Semantic Memory (Document DB + episodic state)
3. Working Environment Memory (isolated task pods)

This model is designed to minimize context rot, preserve user intent across sessions, and keep execution auditable and bounded.

## Pillar 1: Procedural Memory
Definition: Human-readable, versioned instructions and SOPs stored as Markdown.

Primary contents:
- Global operating rules
- Reusable skills/procedures
- Examples and policy constraints

Scope split:
- Global: cross-repo directives and operating standards
- Local: repository-specific directives (for example, `.athena/preferences.md`)

Why it exists:
- Makes agent behavior explicit and reviewable
- Keeps critical rules in Git history
- Provides deterministic startup context

## Pillar 2: State & Semantic Memory
Definition: Persistent system memory for workflow state, episodic logs, and retrievable preferences.

Primary contents:
- Run/session episodes
- Workflow state and checkpoints
- User/team preferences with confidence and provenance

Retrieval modes:
- Exact lookup for structured preferences/state
- Semantic retrieval for related prior episodes and patterns

Why it exists:
- Prevents drift across sessions
- Enables compounding improvements from prior runs
- Supports traceability of what was recalled and why

## Pillar 3: Working Environment Memory
Definition: Ephemeral, isolated runtime where code/data is cloned and operated on safely.

Primary contents:
- Task-specific cloned repository/data
- Runtime-local context and temporary artifacts
- Semantic navigation/indexing for that workspace (for example via LSP)

Why it exists:
- Safe execution boundaries (host protection)
- High-fidelity code understanding without full-file dumps
- Disposable short-term memory with explicit persistence back to Pillar 2

## Cross-Pillar Contracts
- Bootstrap order:
  1. Load global procedural memory
  2. Load local repo procedural memory (if present)
  3. Retrieve relevant state/semantic memory
  4. Execute in isolated working environment
  5. Persist validated learnings back to state memory

- Memory write policy:
  - Local repo preferences remain local by default
  - Cross-repo preferences promoted to global profile only with confidence/provenance

- Auditability:
  - Every recall and write should capture source, reason, and run/session id

- Cost and safety bounds:
  - Isolation required for untrusted code execution
  - Resource/time limits required per run

## Design Principles (Foundation)
- Local-first, cloud-optional
- Deterministic over implicit behavior
- Traceability over hidden adaptation
- Precision retrieval over context bulk
- Reversible decisions and incremental rollout

## Open Design Questions
- Promotion rules from local preference to global profile
- Conflict resolution between global and local procedural rules
- Minimal schema for episodic memory that is useful without over-modeling
- Human review gates for low-confidence memory writes

## Foundation Decision
AthenaMind architecture decisions should default to this three-pillar model unless a formal ADR supersedes it.
