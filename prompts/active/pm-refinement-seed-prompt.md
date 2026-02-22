<!-- AUDIENCE: Internal/Technical -->

# PM Refinement Directive

Run this cycle when QA has created intake bugs or when `backlog/engineering/active/` is empty.

## PM Refinement Cycle (Mandatory)
1. Review new items in `backlog/engineering/intake/`.
2. Run intake validation before moving any items:
   - `scripts/validate_intake_items.sh`
3. If validation fails:
   - fix missing metadata fields or invalid status values in-place;
   - move misfiled `ARCH-*` items to `backlog/architecture/intake/`;
   - move misfiled `STORY-*`/`BUG-*` items to `backlog/engineering/intake/`;
   - rerun `scripts/validate_intake_items.sh` until it passes.
4. Refine each item into clear, actionable stories/bug cards.
5. Rank and move refined items to `backlog/engineering/active/` in execution order.
6. Update `backlog/engineering/active/README.md` Active Sequence.
   - Verify queue/order updates remain consistent with `backlog/STATE_TRANSITION_CHECKLIST.md`.
7. Update `research/roadmap/PROGRAM_STATE_BOARD.md` queue counts and Now/Next items.
8. Ensure active stories include traceability metadata (`idea_id`, `phase`, `adr_refs`, metric field).
9. Update `prompts/active/next-agent-seed-prompt.md` only if special launch instructions are needed.
10. Commit refinement outputs and backlog state changes with an appropriate message.

## Constraints
- Preserve QA priority intent (`P0` highest urgency).
- Keep stories small, testable, and explicit.
- Do not implement fixes in PM mode.
- Do not skip commit after refinement decisions.
- Apply stage exit requirements in `docs/process/STAGE_EXIT_GATES.md`.
