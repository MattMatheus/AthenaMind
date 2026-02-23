# Architecture Intake Refinement Guide

Use this checklist before promoting an item from `delivery-backlog/architecture/intake/` to `delivery-backlog/architecture/ready/`.

## Intake Quality Gate
1. Decision scope is explicit and bounded (what decision is in scope and what is out of scope).
2. Problem statement describes why the decision is needed now.
3. Inputs list all required ADRs, architecture docs, and constraints.
4. Outputs required are concrete and verifiable:
   - ADR updates
   - architecture artifacts
   - risk/tradeoff notes
5. Acceptance criteria are testable and reference concrete artifacts.
6. QA focus describes what review should verify.

## Separation Rule: Architecture vs Engineering
- Keep work in architecture lane when the primary output is a decision, ADR, policy, contract, or design standard.
- Move work to engineering lane when the primary output is implementation code, runtime behavior changes, tests, or deployment changes.
- If an architecture item identifies implementation tasks, keep those tasks out of the architecture story and file them as engineering intake stories.

## Validation Failure Handling (PM Refinement)
Before promoting items to active queues, run:
- `tools/validate_intake_items.sh`

If validation fails:
1. Fix missing metadata fields or invalid status values in the affected intake file.
2. Move misfiled `ARCH-*` stories from engineering intake to architecture intake.
3. Move misfiled `STORY-*` or `BUG-*` items from architecture intake to engineering intake.
4. Re-run validation until it passes, then proceed with ranking/promotion.

## Required Output Package for Architect Stage
1. Decision record updates (new ADR or ADR revision).
2. Architecture artifact updates (maps, standards, contracts, or design docs).
3. Risk register entry with:
   - risk description
   - impact
   - mitigation
   - owner and follow-up trigger
4. Handoff package in `delivery-backlog/architecture/qa/` including:
   - decisions made
   - alternatives considered
   - risks and mitigations
   - updated artifact paths
   - validation commands and results
   - open QA questions
