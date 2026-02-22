# External Sources (2026-02-22)

## Purpose
Primary-source references used to inform module-boundary decisions.

## Sources
- Podman documentation (daemonless model, rootless mode, pod primitives):
  - https://docs.podman.io/en/v4.8.0/
  - https://docs.podman.io/en/stable/markdown/podman-pod-create.1.html
  - https://docs.podman.io/en/v4.6.1/markdown/podman-run.1.html
- SQLite WAL documentation (concurrency/durability behavior for local-first state):
  - https://www.sqlite.org/wal.html
- Language Server Protocol spec (semantic navigation contract):
  - https://microsoft.github.io/language-server-protocol/
- OpenTelemetry traces and semantic conventions (audit/trace design baseline):
  - https://opentelemetry.io/docs/concepts/signals/traces/
  - https://opentelemetry.io/docs/specs/semconv/general/trace/
- FAISS wiki (index options and selection tradeoffs for local semantic retrieval):
  - https://github.com/facebookresearch/faiss/wiki
  - https://github.com/facebookresearch/faiss/wiki/The-index-factory

## Design-Relevant Takeaways
- Podman: rootless and pod-level isolation support the working-memory boundary as a first-class runtime module.
- SQLite WAL: viable default for local-first concurrent writes with predictable operational model.
- LSP: standardized protocol justifies a separate semantic-navigation module rather than language-specific tight coupling.
- OpenTelemetry: trace/span + semantic attributes provide a concrete schema foundation for recall/write auditability.
- FAISS: index strategy should be modular and swappable as data scale changes.
