# Story: Implement episode write-back at cycle boundary

## Metadata
- `id`: STORY-20260223-episode-writeback
- `owner_persona`: staff-personas/Product Manager - Maya.md
- `status`: qa
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0001, ADR-0002, ADR-0006, ADR-0009]
- `success_metric`: each completed cycle persists a structured episode record retrievable by future bootstrap calls
- `release_checkpoint`: required

## Problem Statement
- The memory system has write and retrieve but no mechanism for agents to record what happened during a work session. Without episode records, the bootstrap protocol has no cross-session context to retrieve. Episode write-back is the minimal Pillar 2 (State Memory) implementation needed to close the dogfooding loop.

## Scope
- In:
  - Add `memory-cli episode write` subcommand that accepts structured episode data:
    - `--repo`, `--session-id`, `--cycle-id`, `--story-id`
    - `--outcome` (success | partial | blocked)
    - `--summary` (brief text or file)
    - `--files-changed` (comma-separated list)
    - `--decisions` (brief text or file)
  - Store episodes as structured entries in the memory store (type: `episode`, domain: repo-derived)
  - Episode writes go through existing governance gates (review required per ADR-0006)
  - Add `memory-cli episode list` to list episodes for a repo/scenario
  - Add tests for episode write and retrieval round-trip
- Out:
  - Automatic/autonomous episode writes (must be reviewed per governance policy)
  - Episode analytics or aggregation
  - Changes to existing memory entry types

## Acceptance Criteria
1. `memory-cli episode write` persists a structured episode record in the memory store.
2. Episode records are retrievable via `memory-cli retrieve` with episode-relevant queries.
3. Episode records are retrievable via `memory-cli episode list` filtered by repo.
4. Episode writes respect existing governance gates (blocked during autonomous runs).
5. Telemetry event emitted for episode write operations.
6. Episode schema documented in metadata.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260223-agent-bootstrap-protocol (bootstrap consumes episodes)

## Notes
- Episode records are the simplest form of Pillar 2 (State Memory). Keep the schema minimal: who, what, when, outcome, and a brief summary.
- The observer report already captures similar data at cycle boundary — consider whether the observer script should invoke episode write-back as part of its flow.
