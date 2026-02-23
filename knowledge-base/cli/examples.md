# CLI Examples

## Summary
Copy/paste examples for common CLI tasks.

## Write Entry
```bash
go run ./cmd/memory-cli write \
  --root memory \
  --id handoff-template \
  --title "Handoff Template" \
  --type instruction \
  --domain docs \
  --body "Always include risks, evidence, and next-state recommendation." \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "baseline docs quality" \
  --risk "low; reversible by git revert" \
  --notes "approved for docs baseline" \
  --embedding-endpoint http://localhost:11434
```

## Retrieve Entry
```bash
go run ./cmd/memory-cli retrieve \
  --root memory \
  --query "handoff instruction template" \
  --embedding-endpoint http://localhost:11434
```

## Run Evaluation
```bash
go run ./cmd/memory-cli evaluate \
  --root memory \
  --query-file cmd/memory-cli/testdata/eval-query-set-v1.json \
  --corpus-id memory-corpus-v1 \
  --query-set-id query-set-v1 \
  --embedding-endpoint http://localhost:11434
```

## Snapshot Create/List/Restore
```bash
go run ./cmd/memory-cli snapshot create \
  --root memory \
  --created-by clara \
  --reason "pre-release docs freeze"
```

```bash
go run ./cmd/memory-cli snapshot list --root memory
```

```bash
go run ./cmd/memory-cli snapshot restore \
  --root memory \
  --snapshot-id <snapshot-id> \
  --stage pm \
  --reviewer maya \
  --decision approved \
  --reason "rollback to known good set" \
  --risk "low; manifest/checksum verified" \
  --notes "approved restore"
```

## API Retrieve With Fallback
```bash
go run ./cmd/memory-cli api-retrieve \
  --root memory \
  --query "handoff instruction template" \
  --session-id docs-session-1 \
  --gateway-url http://127.0.0.1:8788
```
