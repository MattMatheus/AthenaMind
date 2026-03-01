<!-- AUDIENCE: Internal/Technical -->

# PM Refinement Directive

Run this cycle when QA has created intake bugs or when `delivery-backlog/engineering/active/` is empty.

## Onboarding Handshake (Mandatory)
1. Complete PM-stage onboarding context checks.
2. Generate a new UUIDv4 for this launched session.
3. Disclose it to the user exactly as: `RESUME TOKEN: <uuid>`.
4. Only then start the PM Refinement Cycle.

## Stop/Resume Guardrail (Highest Priority)
- Parking trigger pattern is exact and case-sensitive: `STOP WORK. RESUME TOKEN: <token>`.
- On receiving that pattern, enter `PARKED` state immediately and stop all work.
- In `PARKED`, do not run commands, do not edit files, do not continue planning, and do not delegate to sub-agents.
- Only allowed response in `PARKED` is a brief parked acknowledgement or a token mismatch notice.
- Resume trigger pattern is exact and case-sensitive: `RESUME WORK. RESUME TOKEN: <token>`.
- Resume requires exact token match to the active parked token; token comparison is byte-for-byte (whitespace included).
- If resume token is missing or mismatched, stay parked and report `TOKEN MISMATCH: still parked`.
- If multiple park commands are received before a valid resume, keep the first active token unless user explicitly sends `ROTATE RESUME TOKEN: <new-token>` while parked.

## PM Refinement Cycle (Mandatory)
1. Review new items in `delivery-backlog/engineering/intake/`.
2. Run intake validation before moving any items:
   - `tools/validate_intake_items.sh`
3. If validation fails:
   - fix missing metadata fields or invalid status values in-place;
   - move misfiled `ARCH-*` items to `delivery-backlog/architecture/intake/`;
   - move misfiled `STORY-*`/`BUG-*` items to `delivery-backlog/engineering/intake/`;
   - rerun `tools/validate_intake_items.sh` until it passes.
4. Refine each item into clear, actionable stories/bug cards.
5. Rank and move refined items to `delivery-backlog/engineering/active/` in execution order.
   - Apply product-first weighting from `knowledge-base/process/backlog-weighting-policy.md`.
   - Process stories may outrank product work only when a broken process is blocking delivery or gate enforcement.
6. Update `delivery-backlog/engineering/active/README.md` Active Sequence.
   - Verify queue/order updates remain consistent with `delivery-backlog/STATE_TRANSITION_CHECKLIST.md`.
7. Update `product-research/roadmap/PROGRAM_STATE_BOARD.md` queue counts and Now/Next items.
8. Ensure active stories include traceability metadata (`idea_id`, `phase`, `adr_refs`, metric field).
9. Update `stage-prompts/active/next-agent-seed-prompt.md` only if special launch instructions are needed.
10. Run observer:
   - `tools/run_observer_cycle.sh --cycle-id PM-<date>-<slug>`
11. Commit once for this cycle:
   - `cycle-<cycle-id>`

## Constraints
- Preserve QA priority intent (`P0` highest urgency).
- Keep stories small, testable, and explicit.
- Do not implement fixes in PM mode.
- Apply stage exit requirements in `knowledge-base/process/stage-exit-gates.md`.
- Do not commit before observer report is generated.
