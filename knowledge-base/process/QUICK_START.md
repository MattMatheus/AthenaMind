# Quick Start Guide

This page exists to help agents that need a concise, deterministic next step rather than wading through the full workflow manuals. Follow the numbered checklist below before launching any stage.

## 1. Branch and workspace check
- Confirm branch policy before running any `launch_stage` scripts.
  - Default required branch is `dev`.
  - Override with `ATHENA_REQUIRED_BRANCH=<branch>` for isolated workflows.
- If you need to explore a side effort, create a feature branch prefixed with `codex/` (e.g., `codex/experiment2`).

## 2. Command palette
- `./tools/launch_stage.sh engineering` (after verifying there are active stories) and `./tools/launch_stage.sh qa|pm|architect|cycle` for their respective stages.
- `./tools/run_observer_cycle.sh --cycle-id <cycle-id>` once per cycle, then commit with `cycle-<cycle-id>` and include `operating-system/observer/OBSERVER-REPORT-<cycle-id>.md`.
- `go test ./...` keeps the Azure DevOps gate green; run locally before pushing.

## 3. Queue signal interpretation
1. Read `product-research/roadmap/PROGRAM_STATE_BOARD.md`. If any queue count is non-zero, follow the `Now` section before launching another stage.
2. If the engineering active queue is empty, trigger PM refinement (`./tools/launch_stage.sh pm`) before re-running engineering.
3. If a launch script replies `no stories`, stop and feed PM/architect to replenish the intake queue. Do not invent work.

## 4. Memory layer reminder (AthenaMind skill)
- Whenever context matters (intake decisions, QA evidence, release rationale), run a memory CLI command to fetch the latest snapshot:
  ```bash
  ATHENA_MEMORY_ROOT="${ATHENA_MEMORY_ROOT:-./memory/core}"
  go run ./products/athena-mind/cmd/memory-cli retrieve \
    --root "$ATHENA_MEMORY_ROOT" \
    --query "current cycle context"
  ```
- Verify embeddings or health when unsure about retrieval quality:
  ```bash
  go run ./products/athena-mind/cmd/memory-cli verify embeddings --root ...
  ```
- Use the AthenaMind skill whenever the task hinges on memory indexing, retrieval, or CLI behavior. It collapses the context so less capable agents do not lose track.

## 5. Work-system sync points
- Stage exit gates: `knowledge-base/process/STAGE_EXIT_GATES.md` (pass each gate before transitioning). Failure means repeat the current stage.
- Backlog weighting policy: `knowledge-base/process/BACKLOG_WEIGHTING_POLICY.md` (product-first ranking without estimates).
- Documents that change workflow or tooling rules (e.g., commit format, exit gates, stage commands) must be updated in `HUMANS.md`, `AGENTS.md`, and any affected `stage-prompts/active/` files.

## 6. Simplified run-through (for cases when you just need to proceed)
1. Verify `PROGRAM_STATE_BOARD.md` for open queues.
2. If there is work, run `launch_stage` for that domain (engineering → QA). If there is none, run PM refinement to refill active stories.
3. After QA finishes a cycle, run the observer, build the release bundle evidence if needed, and commit `cycle-<cycle-id>`.
4. Always include memory CLI context captures when the next steps rely on aggregated knowledge.

Use this quick guide whenever you need a deterministic action and are unsure which long-form doc to read next.
