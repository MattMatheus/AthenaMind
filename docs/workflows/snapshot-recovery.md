# Snapshot Recovery Workflow

## Summary
Capture a full snapshot and restore it through policy-governed review.

## Preconditions
- Existing memory corpus.
- Reviewer available to approve restore evidence.

## Steps
1. Create snapshot:
```bash
go run ./cmd/memory-cli snapshot create \
  --root memory \
  --created-by clara \
  --reason "checkpoint before high-risk edits"
```
2. List snapshots:
```bash
go run ./cmd/memory-cli snapshot list --root memory
```
3. Restore snapshot:
```bash
go run ./cmd/memory-cli snapshot restore \
  --root memory \
  --snapshot-id <snapshot-id> \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "recover known-good corpus" \
  --risk "low; checksum validation enforced" \
  --notes "approved restore"
```

## Expected Behavior
- Restore blocks on schema incompatibility.
- Restore blocks on invalid checksums.
- Snapshot audit events are written to `memory/audits/`.

## References
- `backlog/engineering/done/STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
- `backlog/engineering/done/QA-RESULT-STORY-20260222-memory-snapshot-mvp-implementation-v02.md`
