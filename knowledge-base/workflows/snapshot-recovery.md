# Snapshot Recovery Workflow

## Goal

Create a safe rollback point and restore if memory corpus quality regresses.

## Flow

1. Create snapshot.
2. List snapshots.
3. Restore with governance evidence.

## Commands

```bash
go run ./cmd/memory-cli snapshot create --root memory --created-by operator --reason "checkpoint"
go run ./cmd/memory-cli snapshot list --root memory
go run ./cmd/memory-cli snapshot restore --root memory --snapshot-id <id> --stage pm --reviewer maya --decision approved --reason "recover" --risk "low" --notes "approved"
```

## Validation

- Snapshot manifest exists.
- Restore succeeds only with valid policy evidence and compatibility.
