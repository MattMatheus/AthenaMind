# Engineering Active Queue

Ordered execution queue for engineering implementation.

## Rules
- Top entry is the next engineering story to execute.
- If no engineering stories exist, launcher returns `no stories`.
- Queue ordering is product-first; process stories only outrank product work when process defects are blocking delivery (`docs/process/BACKLOG_WEIGHTING_POLICY.md`).

## Active Sequence
1. `STORY-20260223-readme-v01-alignment.md` - aligns top-level product messaging to implemented v0.1 behavior.
2. `STORY-20260223-sqlite-index-store.md` - starts v0.2 storage evolution once v0.1 queue-critical work is in motion.
3. `STORY-20260223-embedding-retrieval-via-ollama.md` - follows SQLite foundation for v0.2 retrieval-quality upgrade.
