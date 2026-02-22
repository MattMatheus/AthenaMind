<!-- AUDIENCE: Internal/Technical -->

# Planning Session Directive

Your task is to run an interactive idea-generation session with the human operator before architecture or PM execution.

## Planning Cycle (Mandatory)
1. Start a structured conversation to capture goals, users, constraints, risks, and success metrics.
2. Record notes in a new file under `research/planning/sessions/` using `research/planning/PLANNING_SESSION_TEMPLATE.md`.
3. Classify each proposed idea:
   - Implementation work -> create story in `backlog/engineering/intake/` using `backlog/engineering/intake/STORY_TEMPLATE.md`.
   - Architecture/ADR decision work -> create story in `backlog/architecture/intake/` using `backlog/architecture/intake/ARCH_STORY_TEMPLATE.md`.
4. Provide a next-stage recommendation:
   - `architect` when architecture decisions are required first.
   - `pm` when intake is ready for refinement and ranking.
5. Commit planning notes and any created intake artifacts as: `plan-<plan-id>`.

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
