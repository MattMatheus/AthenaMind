# Public Testing Program

## Objective

Collect high-quality real-world feedback from engineers and data scientists operating memory-assisted coding workflows.

## What We Need Tested

- Retrieval correctness under real repo contexts
- Deterministic fallback and explainability fields
- Governance gate behavior under failure modes
- Telemetry/OTLP export reliability
- Docs and operator UX clarity

## Program Workflow

1. Install via [knowledge-base/getting-started/binaries.md](knowledge-base/getting-started/binaries.md) or source setup docs
2. Run [knowledge-base/workflows/README.md](knowledge-base/workflows/README.md)
3. Configure traces via [knowledge-base/how-to/memory-cli-otel-setup.md](knowledge-base/how-to/memory-cli-otel-setup.md)
4. Submit findings via issue templates

## Data Handling Guidance

- Remove secrets from logs before posting.
- Use synthetic or scrubbed data when possible.
- For sensitive incidents, use private escalation in [SUPPORT.md](SUPPORT.md).

## Maintainer Cadence

- Triage new test reports daily.
- Reproduce and label severity within 48 hours.
- Publish weekly known issues and fixes in [CHANGELOG.md](CHANGELOG.md).
