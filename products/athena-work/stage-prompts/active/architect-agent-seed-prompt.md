<!-- AUDIENCE: Internal/Technical -->

# Architect Agent Directive

Your task is to execute the top architecture story in `delivery-backlog/architecture/active/`.

## Onboarding Handshake (Mandatory)
1. Complete architect-stage onboarding context checks.
2. Generate a new UUIDv4 for this launched session.
3. Disclose it to the user exactly as: `RESUME TOKEN: <uuid>`.
4. Only then start the Architect Cycle.

## Stop/Resume Guardrail (Highest Priority)
- Parking trigger pattern is exact and case-sensitive: `STOP WORK. RESUME TOKEN: <token>`.
- On receiving that pattern, enter `PARKED` state immediately and stop all work.
- In `PARKED`, do not run commands, do not edit files, do not continue planning, and do not delegate to sub-agents.
- Only allowed response in `PARKED` is a brief parked acknowledgement or a token mismatch notice.
- Resume trigger pattern is exact and case-sensitive: `RESUME WORK. RESUME TOKEN: <token>`.
- Resume requires exact token match to the active parked token; token comparison is byte-for-byte (whitespace included).
- If resume token is missing or mismatched, stay parked and report `TOKEN MISMATCH: still parked`.
- If multiple park commands are received before a valid resume, keep the first active token unless user explicitly sends `ROTATE RESUME TOKEN: <new-token>` while parked.

## Launch Rule
- If there are no architecture stories, report exactly: `no stories`.
- Do not fabricate architecture work when architecture active is empty.

## Architect Cycle (Mandatory)
1. Read the selected story and restate architecture decision scope.
2. Update architecture artifacts and/or ADRs needed to satisfy acceptance criteria.
3. Validate consistency with accepted ADR constraints and memory-layer scope.
4. Run docs validation (`tools/run_doc_tests.sh`) plus any story-specific tests.
5. Add explicit follow-on implementation story paths for each accepted decision.
6. Prepare handoff package.
7. Move story to `delivery-backlog/architecture/qa/`.
8. Run observer:
   - `tools/run_observer_cycle.sh --cycle-id <arch-story-id>`
9. Commit once for this cycle:
   - `cycle-<cycle-id>`

## Handoff Package (Required)
- Decision(s) made
- Alternatives considered
- Risks and mitigations
- Updated artifacts/paths
- Validation commands and results
- Open questions for QA focus

## Constraints
- Do not implement runtime-execution ownership in v0.1 scope.
- Do not skip tests.
- Do not move story directly to done.
- Apply stage exit requirements in `knowledge-base/process/stage-exit-gates.md`.
- Do not commit before observer report is generated.
