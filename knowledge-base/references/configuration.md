# Configuration Reference

## Memory Roots

- `--root`: explicit memory root for CLI commands.
- Default examples use `memory` in repo root.

## Retrieval

- `--mode`: `classic|hybrid`
- `--top-k`: candidate trace count
- `--retrieval-backend`: `sqlite|qdrant|neo4j`
- `--embedding-endpoint`: explicit embedding endpoint used by retrieval/indexing

## Constraint Env Vars

- `MEMORY_CONSTRAINT_COST_MAX_PER_RUN_USD`
- `MEMORY_CONSTRAINT_FORCE_TRACE_MISSING`
- `MEMORY_CONSTRAINT_RELIABILITY_FREEZE`
- `MEMORY_CONSTRAINT_FORCE_LATENCY_DEGRADED`
- `MEMORY_CONSTRAINT_LATENCY_P95_RETRIEVAL_MS`

## Evaluation Env Vars

- `MEMORY_EVAL_MAX_P95_MS` (explicit eval override)

## Backend Env Vars

### Qdrant

- `ATHENA_QDRANT_URL`
- `ATHENA_QDRANT_COLLECTION`
- `ATHENA_QDRANT_API_KEY`

### Neo4j

- `ATHENA_NEO4J_HTTP_URL`
- `ATHENA_NEO4J_USER`
- `ATHENA_NEO4J_PASSWORD`
- `ATHENA_NEO4J_DATABASE`

## OTel/OTLP Env Vars

- `ATHENA_OTEL_SAMPLE_RATE`
- `ATHENA_OTEL_STDOUT`
- `ATHENA_ENV`
- `OTEL_EXPORTER_OTLP_ENDPOINT`
- `OTEL_EXPORTER_OTLP_PROTOCOL`
- `OTEL_EXPORTER_OTLP_HEADERS`

## Telemetry Files

- events: `<root>/telemetry/events.jsonl`
- retrieval metrics: `<root>/telemetry/retrieval-metrics.jsonl`
