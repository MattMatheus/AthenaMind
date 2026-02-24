# AthenaMind Quickstart

## Goal

Run a full governed write -> retrieve -> evaluate -> snapshot cycle.

If you installed a precompiled binary, replace `go run ./cmd/memory-cli` with `memory-cli` in each command.

## 1. Write One Entry

```bash
go run ./cmd/memory-cli write \
  --root memory \
  --id onboarding-guide \
  --title "Onboarding Guide" \
  --type prompt \
  --domain engineering \
  --body "Use deterministic fallbacks and always include source_path in retrieval outputs." \
  --stage planning \
  --reviewer maya \
  --decision approved \
  --reason "bootstrap baseline" \
  --risk "low" \
  --notes "approved"
```

## 2. Retrieve It

```bash
go run ./cmd/memory-cli retrieve \
  --root memory \
  --query "fallback and source path policy" \
  --domain engineering
```

Confirm response includes `selected_id`, `selection_mode`, and `source_path`.

## 3. Evaluate Quality Gate

```bash
go run ./cmd/memory-cli evaluate --root memory
```

Inspect `top1_useful_rate`, `fallback_determinism`, and latency metrics.

## 4. Create Snapshot

```bash
go run ./cmd/memory-cli snapshot create \
  --root memory \
  --created-by operator \
  --reason "post-quickstart checkpoint"
```

## 5. (Optional) Export OTel Traces

See [OTel/OTLP Setup](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/how-to/memory-cli-otel-setup.md).
