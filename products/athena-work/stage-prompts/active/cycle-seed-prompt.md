<!-- AUDIENCE: Internal/Technical -->

# Cycle Seed Directive

Your task is to drain the engineering active backlog by alternating engineering and QA cycles.

## Onboarding Handshake (Mandatory)
1. Complete cycle-stage onboarding context checks.
2. Generate a new UUIDv4 for this launched session.
3. Disclose it to the user exactly as: `RESUME TOKEN: <uuid>`.
4. Only then start the Cycle Loop.

## Stop/Resume Guardrail (Highest Priority)
- Parking trigger pattern is exact and case-sensitive: `STOP WORK. RESUME TOKEN: <token>`.
- On receiving that pattern, enter `PARKED` state immediately and stop all work.
- In `PARKED`, do not run commands, do not edit files, do not continue planning, and do not delegate to sub-agents.
- Only allowed response in `PARKED` is a brief parked acknowledgement or a token mismatch notice.
- Resume trigger pattern is exact and case-sensitive: `RESUME WORK. RESUME TOKEN: <token>`.
- Resume requires exact token match to the active parked token; token comparison is byte-for-byte (whitespace included).
- If resume token is missing or mismatched, stay parked and report `TOKEN MISMATCH: still parked`.
- If multiple park commands are received before a valid resume, keep the first active token unless user explicitly sends `ROTATE RESUME TOKEN: <new-token>` while parked.

## Cycle Loop (Mandatory)
1. Run `tools/launch_stage.sh engineering`.
2. If output is exactly `no stories`, stop and report completion.
3. Execute the engineering cycle for the selected story.
4. Run `tools/launch_stage.sh qa`.
5. Execute the QA cycle for the story in `delivery-backlog/engineering/qa/`.
6. Run observer at cycle boundary:
   - `tools/run_observer_cycle.sh --cycle-id <story-id> --story <path-to-story>`
7. Commit once for the full cycle:
   - `cycle-<cycle-id>`
8. Repeat from step 1 until `delivery-backlog/engineering/active/` is drained.

## Commit Discipline
- Do not commit during intermediate stage transitions.
- Use exactly one commit per completed cycle.
- Commit format: `cycle-<cycle-id>`.

## Constraints
- Do not skip tests.
- Do not bypass backlog states.
- Do not continue if branch is not `dev`.
