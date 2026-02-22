# Observer Reports

Observer reports are generated once per completed cycle, before the cycle commit.

## Command
- `scripts/run_observer_cycle.sh --cycle-id <cycle-id> [--story <path>]`

## Policy
- Stage operators should not commit during intermediate stage transitions.
- At cycle boundary, run observer to capture deterministic diff inventory and process-memory deltas.
- Commit once with all cycle artifacts (engineering/qa artifacts, observer report, queue/program updates).

## Required Artifact
- `work-system/observer/OBSERVER-REPORT-<cycle-id>.md`

Use `OBSERVER_REPORT_TEMPLATE.md` as the structural contract for manual edits or review.
