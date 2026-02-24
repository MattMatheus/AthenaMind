# Security Policy

## Supported Surface

This repository includes user-facing runtime and docs for AthenaMind and AthenaWork guidance.

## Reporting a Vulnerability

Do not open public issues for suspected vulnerabilities.

Please report privately with:
- impact summary
- affected file/path and version/commit
- reproducible proof-of-concept
- suggested remediation (if available)

Use the private escalation path in [SUPPORT.md](SUPPORT.md).

## Sensitive Data Rules

- Never commit API keys, tokens, or credentials.
- Treat telemetry payloads as potentially sensitive.
- Redact identifiers in public bug reports unless needed for reproduction.

## OTLP/Telemetry Notes

Remote collectors may require auth headers or mTLS.
Store credentials in environment variables only, never in repo files.
