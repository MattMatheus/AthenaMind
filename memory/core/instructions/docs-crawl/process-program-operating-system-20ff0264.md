# Process Program Operating System 20ff0264

# Program Operating System

This is the control-plane contract for strategic alignment and execution traceability.

## Required Artifacts
- Program state board: `product-research/roadmap/PROGRAM_STATE_BOARD.md`
- Phase plan: `product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`
- PM backlog notes: `knowledge-base/process/PM-TODO.md`
- Stage exit gates: `knowledge-base/process/stage-exit-gates.md`
- Observer report stream: `operating-system/observer/`

## Control-Plane Rules
1. Every new story/bug includes phase + ADR + metric traceability metadata.
2. PM refinement updates program board and active queue in the same cycle.
3. Planning sessions are finalized when downstream artifacts are created.
4. Readiness and roadmap artifacts must reflect current backlog state.
5. `done` is not `shipped` until release checkpoint bundle is approved.
6. Stage-level commits are disallowed; commit once per cycle after observer report generation.

## Observer Rule
- At cycle boundary run `tools/run_observer_cycle.sh --cycle-id <cycle-id>`.
- Observer reports must be committed with the cycle commit (`cycle-<cycle-id>`).
- Observer is responsible for durable metadata deltas (workflow sync checks, memory promotions, release-impact notes).

## Sync Rules
When process behavior changes, update:
- `HUMANS.md`
- `AGENTS.md`
- `DEVELOPMENT_CYCLE.md`
- affected stage prompts under `stage-prompts/active/`
- this process section under `knowledge-base/process/`

## No-Time-Estimate Rule
The operating system does not require delivery date or duration estimates.
Prioritization is value/risk/sequence based, not schedule based.
