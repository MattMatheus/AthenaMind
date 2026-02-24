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
