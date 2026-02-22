<!-- AUDIENCE: Internal/Technical -->

# PM Refinement Directive

Run this cycle when QA has created intake bugs or when `backlog/engineering/active/` is empty.

## PM Refinement Cycle (Mandatory)
1. Review new items in `backlog/engineering/intake/`.
2. Refine each item into clear, actionable stories/bug cards.
3. Rank and move refined items to `backlog/engineering/active/` in execution order.
4. Update `backlog/engineering/active/README.md` Active Sequence.
   - Verify queue/order updates remain consistent with `backlog/STATE_TRANSITION_CHECKLIST.md`.
5. Update `prompts/active/next-agent-seed-prompt.md` only if special launch instructions are needed.
6. Commit refinement outputs and backlog state changes with an appropriate message.

## Constraints
- Preserve QA priority intent (`P0` highest urgency).
- Keep stories small, testable, and explicit.
- Do not implement fixes in PM mode.
- Do not skip commit after refinement decisions.
