# Write/Retrieve Baseline Workflow

## Summary
Create a memory artifact and retrieve it with predictable selection metadata.

## Preconditions
- `docs/getting-started/installation.md` completed.

## Steps
1. Write entry with approved review evidence.
2. Retrieve by natural language query.
3. Confirm `selection_mode`, `selected_id`, and `source_path` are present.
4. Re-run query to confirm deterministic behavior when fallback is used.

## Expected Evidence
- Created files under `memory/prompts|instructions` and `memory/metadata`.
- Updated `memory/index.yaml`.
- JSON retrieval output with required fields.

## References
- `backlog/engineering/done/STORY-20260222-memory-cli-v1-write-retrieve-baseline.md`
- `backlog/engineering/done/QA-RESULT-STORY-20260222-memory-cli-v1-write-retrieve-baseline.md`
