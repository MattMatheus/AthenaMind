# AthenaMind Skill

Use this skill when the task is about memory indexing, retrieval quality, embeddings, snapshots, or memory CLI behavior.

## Workspace Focus
- `products/athena-mind`
- Compatibility links at repo root currently point here.

## Typical Commands
```bash
ATHENA_MEMORY_ROOT="${ATHENA_MEMORY_ROOT:-./memory/core}"
go run ./products/athena-mind/cmd/memory-cli retrieve --root "$ATHENA_MEMORY_ROOT" --query "memory lifecycle"
go run ./products/athena-mind/cmd/memory-cli verify embeddings --root "$ATHENA_MEMORY_ROOT"
go run ./products/athena-mind/cmd/memory-cli verify health --root "$ATHENA_MEMORY_ROOT" --query "memory lifecycle"
```
