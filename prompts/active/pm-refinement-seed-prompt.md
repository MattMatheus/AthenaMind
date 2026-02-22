<!-- AUDIENCE: Internal/Technical -->

# PM Refinement Directive

Run this cycle when QA has created intake bugs or when `backlog/active/` is empty.

## PM Refinement Cycle (Mandatory)
1. Review new items in `backlog/intake/`.
2. Refine each item into clear, actionable stories/bug cards.
3. Rank and move refined items to `backlog/active/` in execution order.
4. Update `backlog/active/README.md` Active Sequence.
5. Update `prompts/active/next-agent-seed-prompt.md` only if special launch instructions are needed.

## Constraints
- Preserve QA priority intent (`P0` highest urgency).
- Keep stories small, testable, and explicit.
- Do not implement fixes in PM mode.
