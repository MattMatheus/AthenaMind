# External Sources: Storage and Embedding (2026-02-22)

## Sources
- SQLite WAL
  - https://sqlite.org/wal.html
- SQLite WITHOUT ROWID guidance
  - https://sqlite.org/withoutrowid.html
- FAISS index factory
  - https://github.com/facebookresearch/faiss/wiki/The-index-factory
- Azure Cosmos DB docs
  - https://learn.microsoft.com/en-us/azure/cosmos-db/consistency-levels
  - https://learn.microsoft.com/en-us/azure/cosmos-db/choose-api

## Design-Relevant Takeaways
- SQLite WAL is fit for local-first single-host durability and concurrency patterns in v0.1.
- Schema design should apply WITHOUT ROWID selectively, only where primary key shapes justify it.
- FAISS index strategy is composable and can begin with simpler exact indices before ANN complexity.
- Cosmos DB consistency/latency tradeoffs support optional cloud mode in later phases.
