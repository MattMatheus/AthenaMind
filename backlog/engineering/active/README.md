# Engineering Active Queue

Ordered execution queue for engineering implementation.

## Rules
- Top entry is the next engineering story to execute.
- If no engineering stories exist, launcher returns `no stories`.
- Queue ordering is product-first; process stories only outrank product work when process defects are blocking delivery (`docs/process/BACKLOG_WEIGHTING_POLICY.md`).

## Active Sequence
1. `STORY-20260223-launch-stage-memory-integration.md` - wires bootstrap/write-back into stage launch and observer workflows.
2. `STORY-20260223-add-ci-pipeline.md` - adds automated quality gate enforcement on `dev` and PR flows.
3. `STORY-20260223-readme-v01-alignment.md` - aligns top-level product messaging to implemented v0.1 behavior.
4. `STORY-20260223-sqlite-index-store.md` - starts v0.2 storage evolution once v0.1 queue-critical work is in motion.
5. `STORY-20260223-embedding-retrieval-via-ollama.md` - follows SQLite foundation for v0.2 retrieval-quality upgrade.
