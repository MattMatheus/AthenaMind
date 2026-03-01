<!-- AUDIENCE: Internal/Technical -->

# Planning Session Directive

Your task is to run an interactive idea-generation session with the human operator before architecture or PM execution.

## Onboarding Handshake (Mandatory)
1. Complete planning-stage onboarding context checks.
2. Generate a new UUIDv4 for this launched session.
3. Disclose it to the user exactly as: `RESUME TOKEN: <uuid>`.
4. Only then start the Planning Cycle.

## Stop/Resume Guardrail (Highest Priority)
- Parking trigger pattern is exact and case-sensitive: `STOP WORK. RESUME TOKEN: <token>`.
- On receiving that pattern, enter `PARKED` state immediately and stop all work.
- In `PARKED`, do not run commands, do not edit files, do not continue planning, and do not delegate to sub-agents.
- Only allowed response in `PARKED` is a brief parked acknowledgement or a token mismatch notice.
- Resume trigger pattern is exact and case-sensitive: `RESUME WORK. RESUME TOKEN: <token>`.
- Resume requires exact token match to the active parked token; token comparison is byte-for-byte (whitespace included).
- If resume token is missing or mismatched, stay parked and report `TOKEN MISMATCH: still parked`.
- If multiple park commands are received before a valid resume, keep the first active token unless user explicitly sends `ROTATE RESUME TOKEN: <new-token>` while parked.

## Planning Cycle (Mandatory)
1. Start a structured conversation to capture goals, users, constraints, risks, and success metrics.
2. Record notes in a new file under `product-research/planning/sessions/` using `product-research/planning/PLANNING_SESSION_TEMPLATE.md`.
3. Classify each proposed idea:
   - Implementation work -> create story in `delivery-backlog/engineering/intake/` using `delivery-backlog/engineering/intake/STORY_TEMPLATE.md`.
   - Architecture/ADR decision work -> create story in `delivery-backlog/architecture/intake/` using `delivery-backlog/architecture/intake/ARCH_STORY_TEMPLATE.md`.
4. Ensure all created intake items include traceability metadata (`idea_id`, `phase`, `adr_refs`, metric fields).
5. Provide a next-stage recommendation:
   - `architect` when architecture decisions are required first.
   - `pm` when intake is ready for refinement and ranking.
6. Set planning session status to `finalized` once intake artifacts are created and linked.
7. Run observer and capture cycle delta:
   - `tools/run_observer_cycle.sh --cycle-id <plan-id>`
8. Commit once for this cycle:
   - `cycle-<cycle-id>`

## Session Output Requirements
- Problem framing and target outcomes
- Assumptions and constraints
- Candidate ideas/options considered
- Decision/gap list by owner lane (engineering vs architecture)
- Concrete intake items created (paths + ids)
- Recommended next stage and rationale

## Constraints
- Do not implement production changes in planning mode.
- Do not skip writing session notes.
- Do not place architecture decision work in engineering intake.
- Do not commit before observer report is generated.
