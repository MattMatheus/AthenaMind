# Story: Define and implement agent bootstrap protocol

## Metadata
- `id`: STORY-20260223-agent-bootstrap-protocol
- `owner_persona`: personas/Product Manager - Maya.md
- `status`: qa
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0001, ADR-0002, ADR-0007, ADR-0009]
- `success_metric`: agent sessions start with relevant memory context retrieved via memory-cli, measurable as non-empty bootstrap payload in at least one dogfooding scenario
- `release_checkpoint`: required

## Problem Statement
- The memory CLI exists but nothing consumes it. The cross-pillar bootstrap order (load procedural → retrieve semantic → execute → persist) is documented in `research/architecture/FOUNDATION_PILLARS.md` but not implemented. Without a concrete bootstrap protocol, the memory system cannot prove its core value proposition: instant resume and cross-session continuity.

## Scope
- In:
  - Define a `memory-cli bootstrap` subcommand that:
    - Accepts `--repo`, `--session-id`, `--scenario` flags
    - Retrieves relevant procedural memories for the current context
    - Retrieves last episode summary for the repo/scenario (if episode store exists)
    - Outputs a structured context payload (JSON) suitable for injection into agent prompts
  - Define the bootstrap payload schema (memory entries + episode context + metadata)
  - Add tests for bootstrap with empty memory store (graceful zero-result case)
  - Add tests for bootstrap with seeded memory entries
- Out:
  - Episode write-back (separate story)
  - Integration with launch_stage.sh (separate story)
  - Changes to existing write/retrieve behavior

## Acceptance Criteria
1. `memory-cli bootstrap` subcommand exists and returns structured context payload.
2. Payload includes retrieved memory entries with selection metadata (id, mode, source_path).
3. Bootstrap with empty memory store returns valid empty payload without error.
4. Bootstrap with seeded entries returns relevant matches.
5. Telemetry event emitted for bootstrap operation.
6. Bootstrap payload schema is documented.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260223-split-main-go-into-packages (recommended but not blocking)

## Notes
- This is the concrete implementation of "instant resume." Keep the payload minimal and useful rather than comprehensive.
- The bootstrap command is the primary integration surface for agent consumers.
