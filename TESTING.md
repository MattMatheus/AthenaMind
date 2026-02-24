# Testing Guide

## Baseline

Run full suite:
```bash
go test ./...
```

## Public Test Loop

1. Install and verify dependencies.
2. Execute quickstart workflows.
3. Run evaluation gate with your query set.
4. Capture telemetry and trace evidence.
5. Report findings using templates.

## Recommended Telemetry Stack (Local)

Use Maple OSS locally via Podman for telemetry-first testing.

1. Start Maple:
```bash
git clone https://github.com/makisuo/maple.git /tmp/maple
cp /tmp/maple/.env.example /tmp/maple/.env
podman compose -f /tmp/maple/docker-compose.yml up --build -d
```
2. Configure AthenaMind OTLP export:
```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:3474"
export OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"
export OTEL_EXPORTER_OTLP_HEADERS="x-maple-ingest-key=<your-ingest-key>"
export ATHENA_OTEL_SAMPLE_RATE="1.0"
```
3. Run AthenaMind commands under test, then verify telemetry in Maple UI (`http://localhost:3471`).

See `knowledge-base/how-to/memory-cli-otel-setup.md` for full details and alternatives.

## Suggested Matrix

- OS: macOS, Linux
- Embedding endpoint: local Ollama, remote endpoint
- Retrieval mode: `classic`, `hybrid`
- Retrieval backend: `sqlite`, `qdrant`, `neo4j` (when available)
- Constraint conditions: normal, forced latency degrade, reliability freeze

## Evidence Expectations

For every bug report include:
- exact command(s)
- full stderr/stdout snippet
- expected vs actual result
- root path/config flags/env used
- whether issue is deterministic
