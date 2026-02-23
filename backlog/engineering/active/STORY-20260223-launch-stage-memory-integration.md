# Story: Integrate memory bootstrap into launch_stage.sh

## Metadata
- `id`: STORY-20260223-launch-stage-memory-integration
- `owner_persona`: personas/Product Manager - Maya.md
- `status`: active
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0001, ADR-0002, ADR-0009]
- `success_metric`: agent sessions launched via launch_stage.sh include retrieved memory context in their prompt payload
- `release_checkpoint`: required

## Problem Statement
- The agent bootstrap protocol and episode write-back exist but are not wired into the agent launch flow. Without integration into `launch_stage.sh`, agents must manually invoke memory-cli, which they will not do without explicit instructions. This story closes the dogfooding loop: agents automatically consume memory at startup and produce episodes at cycle close.

## Scope
- In:
  - Modify `scripts/launch_stage.sh` to invoke `memory-cli bootstrap` before emitting the seed prompt
  - Append bootstrap payload to the seed prompt context (or output it alongside the prompt)
  - Modify `scripts/run_observer_cycle.sh` to invoke `memory-cli episode write` with cycle outcome data after observer report generation
  - Graceful degradation: if memory-cli binary is not available or bootstrap returns empty, proceed with existing behavior (no hard dependency)
- Out:
  - Changes to memory-cli behavior
  - Changes to seed prompt content (only appending retrieved context)
  - Requiring memory-cli for the system to function (soft integration)

## Acceptance Criteria
1. `launch_stage.sh` invokes `memory-cli bootstrap` and includes result in output when binary is available.
2. If `memory-cli` is unavailable, `launch_stage.sh` proceeds normally with a warning.
3. `run_observer_cycle.sh` invokes `memory-cli episode write` with cycle data when binary is available.
4. If `memory-cli` is unavailable, `run_observer_cycle.sh` proceeds normally with a warning.
5. Integration tested with at least one manual cycle run demonstrating bootstrap → work → episode write-back flow.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- STORY-20260223-agent-bootstrap-protocol
- STORY-20260223-episode-writeback

## Notes
- This is a soft integration. The memory system enhances the agent workflow but does not gate it. If the memory-cli binary doesn't exist yet (e.g., not compiled), the system continues to work as before.
- The integration point should be clearly documented so agents understand what memory context they are receiving and what episode data is being captured.
