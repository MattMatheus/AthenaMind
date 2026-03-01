# Agent Branch Strategy

## Purpose
Define a branch model where implementation is performed by agents, while humans operate in management/approval mode.

## Role Mode
- Human role: manager/approver only (no direct code implementation).
- Agent role: all implementation and cycle execution work.

## Branch Topology
1. `main` (release branch)
- Protected, stable, release-only.
- Updated from `dev` after release checkpoint decision.

2. `dev` (integration branch)
- Default execution branch for stage launchers.
- Canonical branch for active agent cycles.

3. `agent/<agent-id>/<cycle-id>` (optional isolation branch)
- Use when parallel agent experiments or risky changes need isolation.
- Rebase or merge back into `dev` after QA and observer evidence are complete.

## Commit and Merge Rules
- One commit per completed cycle: `cycle-<cycle-id>`.
- No stage-level commits before observer report.
- Merge to `main` only after release bundle says ship.
- Human manager may approve/promote, but does not author implementation commits.
- Use date-based release identifiers for ship checkpoints: `<label>-YYYY-MM-DD`.

## Operational Flow
1. Run normal cycles on `dev` by default.
2. If isolation is needed, create `agent/<agent-id>/<cycle-id>` from `dev`.
3. Execute architect/pm/engineering/qa workflow with observer evidence.
4. Integrate back to `dev`.
5. Run release checkpoint and promote `dev -> main` only on explicit ship decision.

## GitHub Settings (Recommended)
- Protect `main`: require PR, CI pass, no force-push, and require at least one human manager approval for `dev -> main`.
- Protect `dev`: require CI pass, restrict force-push.
- Default branch for daily work: `dev`.
- Optional CODEOWNERS: require manager review for `main` promotions.

## Approval and Automation Policy
- `dev -> main`:
  - manual promotion only
  - human approval required
  - agents may prepare PRs and evidence, but must not self-approve
- `dev -> dev`:
  - agent automation allowed (direct push or automation PR)
  - still must pass CI and cycle evidence rules
- `dev -> feature/*`:
  - agent automation allowed for branch creation/update and PR preparation
  - intended for isolated experiments or parallel agent work

## Quick Commands
```bash
# default agent path
git checkout dev

# optional isolated agent branch
git checkout -b agent/<agent-id>/<cycle-id> dev

# after cycle evidence is complete
git checkout dev
git merge --no-ff agent/<agent-id>/<cycle-id>
```

## Audit Notes
- Observer reports and handoff artifacts are required for every cycle commit.
- Direction-changing transitions still require explicit human confirmation artifacts.
