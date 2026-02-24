# Memory CLI OTel/OTLP Setup

## Summary

AthenaMind CLI emits OpenTelemetry spans and can export them to local or remote OTLP collectors.

## Basic Controls

- `ATHENA_OTEL_SAMPLE_RATE` (default `1.0`): trace sampling ratio `0.0-1.0`.
- `ATHENA_OTEL_STDOUT` (`1|true|yes|on`): enable stdout span exporter for local debugging.
- `ATHENA_ENV`: optional environment tag.

## OTLP Export (Collector-Managed)

Use standard OTel env vars supported by your collector setup:

- `OTEL_EXPORTER_OTLP_ENDPOINT`
- `OTEL_EXPORTER_OTLP_PROTOCOL` (`grpc` or `http/protobuf`)
- `OTEL_EXPORTER_OTLP_HEADERS`

Example:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="https://collector.example.com:4317"
export OTEL_EXPORTER_OTLP_PROTOCOL="grpc"
export OTEL_EXPORTER_OTLP_HEADERS="authorization=Bearer <token>"
export ATHENA_OTEL_SAMPLE_RATE="1.0"
```

## Does OTLP Require A Secret?

- Local/open collectors: often no secret.
- Managed remote collectors: usually yes (token/API key header) or mTLS certs.

## Local Debug Mode

```bash
export ATHENA_OTEL_STDOUT=1
go run ./cmd/memory-cli retrieve --root memory --query "trace check"
```

## Security Guidance

- Never commit collector tokens.
- Prefer environment injection at runtime.
- Rotate credentials regularly.
