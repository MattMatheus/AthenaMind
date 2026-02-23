# Story: Add self-check helper for stage readiness

## Metadata
- `id`: STORY-20260223-self-check-helper
- `owner_persona`: staff-personas/SRE - Nia.md
- `status`: active
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0004, ADR-0013]
- `success_metric`: Self-check finishes under 5s and surfaces current branch, git cleanliness, and active story count in one run.
- `release_checkpoint`: deferred

## Problem Statement
New agents spend time opening multiple docs/paths to learn whether the repo is on the right branch, has pending changes, or has an active story queued. This slows execution and increases risk of skipping required stages.

## Scope
- In: single command that reports branch, git cleanliness summary, and engineering queue state.
- Out: memory CLI operations, automatic queue promotion, CI/test execution.

## Acceptance Criteria
1. Running `./tools/self_check.sh` prints branch name, warns if not `dev`, and shows git status summary.
2. The command reports engineering active queue items (IDs and titles) or says "No active stories" if empty.
3. The script exits 0 when it runs successfully and 1 when required repo state checks fail to run (e.g., git unavailable).

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- None

## Notes
- Keep runtime under 5 seconds on typical laptop.
