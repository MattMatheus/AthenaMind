# Clara Documentation Program

## Summary
Execution brief for building and maintaining comprehensive user-facing documentation with Technical Writer Clara.

## Intended Audience
- Product Manager and operator coordinating release readiness.
- Technical Writer Clara executing documentation cycles.

## Preconditions
- Source-of-truth artifacts are available:
  - accepted ADRs in `research/decisions/`
  - implementation stories in `backlog/engineering/done/`
  - QA evidence in `backlog/engineering/done/QA-RESULT-*`
- Documentation backlog stories exist in `backlog/engineering/intake/` or `backlog/engineering/active/`.

## Main Flow
1. PM promotes documentation stories into active queue based on value/risk.
2. Clara executes doc story scope and updates user-facing pages under `docs/`.
3. Clara removes stale/conflicting docs and captures remaining gaps.
4. QA validates documentation acceptance criteria and references.
5. Observer runs, then one cycle commit is made using `cycle-<cycle-id>`.

## Operating Modes
### Continuous Mode (Every Relevant Cycle)
- Trigger when any done story changes user-visible behavior.
- Required outputs:
  - updated docs pages
  - linked evidence references
  - gaps captured as intake items

### Release Mode (On-Demand Before Release Decision)
- Trigger during release checkpoint preparation.
- Required outputs:
  - updated `docs/release-notes/` artifact
  - known limitations section aligned with QA evidence
  - cross-check with release bundle in `work-system/handoff/`

## Failure Modes
- Documentation claims drift from accepted ADR or QA-backed behavior.
- Core user journeys lack complete, executable guidance.
- Release notes diverge from release bundle ship/hold decision.

## References
- `personas/Technical Writer - Clara.md`
- `docs/process/STAGE_EXIT_GATES.md`
- `scripts/run_observer_cycle.sh`
- `work-system/handoff/RELEASE_BUNDLE_TEMPLATE.md`
