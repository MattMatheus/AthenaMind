# Contributing

Thanks for helping test and improve Athena.

## Scope

This repository is a slim, user-facing distribution focused on:
- AthenaMind runtime behavior and docs
- AthenaWork user-facing guidance

Internal planning/research artifacts are intentionally archived outside this repo.

## Before You Open a PR

1. Read [README.md](README.md) and [knowledge-base/INDEX.md](knowledge-base/INDEX.md).
2. Reproduce the issue with explicit commands.
3. Run tests:
```bash
go test ./...
```
4. For docs-only changes, verify links and command accuracy.

## PR Requirements

- Clear problem statement and expected behavior.
- Reproduction steps and test evidence.
- CI must pass (`.github/workflows/ci.yml`).
- No secrets in commits or logs.
- Keep changes focused and reversible.

## Commit Guidance

Use clear, scoped messages, for example:
- `docs: clarify otlp auth examples`
- `retrieval: fix fallback determinism on latency degrade`

## High-Value Contributions

- Retrieval quality regressions and edge cases
- Governance and auditability correctness
- Telemetry/OTLP operational reliability
- User-facing docs accuracy for engineers and data scientists
