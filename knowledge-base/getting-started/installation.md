# Installation

## Summary

Install runtime prerequisites and verify the slim AthenaMind distribution is healthy.

## Prerequisites

- macOS or Linux shell.
- Go 1.22+.
- Network access if pulling new Go modules.

## Install Steps

1. Choose install method:
   - Precompiled binary: [Precompiled Binaries](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/getting-started/binaries.md)
   - Build from source: continue below
2. Clone and enter repo.
3. Verify Go:
```bash
go version
```
4. Download modules:
```bash
go mod download
```
5. Run full tests:
```bash
go test ./...
```

## Optional Services

### Ollama embeddings
```bash
podman compose -f docker-compose.ollama.yml up -d
```

### Retrieval backend experiments
```bash
podman compose -f docker-compose.experiment.yml up -d
```

## Memory Storage Defaults

- Default root for examples: `memory`
- Recommended for real use: set explicit root per environment and persist on durable storage.

## Observability Defaults

- Local telemetry events: `<root>/telemetry/events.jsonl`
- Retrieval metrics: `<root>/telemetry/retrieval-metrics.jsonl`
- OTel tracing enabled in CLI runtime, with OTLP/collector config via env vars.

## Verify Runtime

```bash
go run ./cmd/memory-cli retrieve --root memory --query "health check"
```

If index is empty, you should see a clear error indicating no entries exist yet.
