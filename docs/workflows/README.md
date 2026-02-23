# Workflows

## Summary
Task-oriented product workflows that combine concepts and CLI usage into reproducible user outcomes.

## Intended Audience
- Users completing common memory workflows end to end.
- Internal teams validating workflow reliability and quality gates.

## Preconditions
- Getting Started and CLI prerequisites are complete.
- Relevant ADR and acceptance criteria are available for traceability.

## Main Flow
1. Start with `write-retrieve-baseline.md`.
2. Run `evaluation-gate.md` to validate retrieval quality.
3. Run `snapshot-recovery.md` for create/list/restore.
4. Capture outputs and link to QA evidence.

## Failure Modes
- Workflow steps omit required preconditions.
- Expected outputs are not explicitly defined.
- Regression behavior is undocumented after story updates.

## References
- `docs/getting-started/README.md`
- `docs/cli/README.md`
- `backlog/engineering/done/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
- `backlog/engineering/done/STORY-20260222-dogfood-semantic-retrieval-hardening-v01.md`

## Pages
- `write-retrieve-baseline.md`
- `evaluation-gate.md`
- `snapshot-recovery.md`
