# Engineering Handoff: STORY-20260223-episode-writeback

## What Changed
- Added `memory-cli episode` command family with:
  - `episode write`
  - `episode list`
- Implemented `internal/episode` package for episode persistence and listing.
- Episode write now:
  - requires governance review evidence (same gate model as writes)
  - persists structured episode JSON under `<root>/episodes/<repo>/`
  - updates `<root>/episodes/<repo>/latest.json` for bootstrap consumption
  - indexes episodes into memory index as `type=episode` for retrieval round-trip
  - emits telemetry (`memory.episode.write`, `operation=episode_write`)
- Extended index validation/CRUD to support `type=episode` and `episodes/` path prefix.
- Added tests for episode write/list/retrieve and autonomous-run enforcement.
- Updated CLI docs for `episode` commands.

## Why It Changed
- Story requires a minimal Pillar-2 state-memory write-back path so completed cycles can be persisted and reused by bootstrap.

## Test Updates Made
- `internal/episode/episode_test.go`
- `cmd/memory-cli/main_test.go`:
  - `TestEpisodeWriteListAndRetrieveRoundTrip`
  - `TestEpisodeWriteBlockedDuringAutonomousRun`

## Test Run Results
- `go test ./...` passed.
- `scripts/run_doc_tests.sh` passed.

## Open Risks / Questions
- `episode write` currently stores `latest.json` under repo scope with scenario marker `episode`; later launch-stage integration may refine scenario mapping.

## Recommended QA Focus Areas
- Governance gating behavior for `episode write` (especially autonomous-run block).
- `episode list --repo` output correctness and ordering.
- Retrieval round-trip on indexed episode entries (`type=episode`).
- Telemetry event and operation fields for episode writes.

## New Gaps Discovered
- None.
