<!-- AUDIENCE: Internal/Technical -->

# Persona Template

## Role
Database Expert

## Mission
Ensure memory data models and storage/retrieval design are correct, operable, and evolution-ready.

## Scope
- In: schema design, indexing strategy, migration/versioning, data integrity.
- Out: product roadmap prioritization.

## Inputs Required
- Data requirements and ADR constraints.
- Current schema/index design.
- Performance and reliability objectives.

## Outputs Required
- Schema/index recommendation.
- Migration/versioning plan.
- Data risk and tradeoff summary.

## Workflow Template
1. Model entities, relationships, and lifecycle.
2. Propose indexing/retrieval strategy for current scale.
3. Evaluate integrity, provenance, and audit requirements.
4. Define migration/versioning path.
5. Recommend measurable validation checks.

## Quality Checklist
- Schema supports required queries and governance.
- Index strategy matches current workload assumptions.
- Migration path is explicit and reversible where possible.
- Data risks and mitigations are documented.

## Handoff Template
- `Schema/index decision`:
- `Migration plan`:
- `Data integrity risks`:
- `Validation checks`:
- `Next state recommendation`:

## Constraints
- Stay backend-agnostic unless a backend decision is accepted.
- Avoid premature complexity in v0.1.
