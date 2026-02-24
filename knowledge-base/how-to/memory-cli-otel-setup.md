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

## Local Maple (Podman)

Maple open source provides a local OTLP ingest gateway and UI.

```bash
git clone https://github.com/makisuo/maple.git /tmp/maple
cp /tmp/maple/.env.example /tmp/maple/.env
podman compose -f /tmp/maple/docker-compose.yml up --build -d
```

Maple local endpoints:
- UI: `http://localhost:3471`
- API: `http://localhost:3472`
- Ingest gateway: `http://localhost:3474`

Send AthenaMind traces to Maple ingest:

```bash
export OTEL_EXPORTER_OTLP_ENDPOINT="http://localhost:3474"
export OTEL_EXPORTER_OTLP_PROTOCOL="http/protobuf"
export OTEL_EXPORTER_OTLP_HEADERS="x-maple-ingest-key=<your-ingest-key>"
export ATHENA_OTEL_SAMPLE_RATE="1.0"
go run ./cmd/memory-cli retrieve --root memory --query "maple local test"
```

Notes:
- Use the ingestion key from Maple Settings.
- `Authorization=Bearer <key>` is also accepted by Maple ingest.
- For direct collector testing (bypassing Maple ingest), point to `http://localhost:4318`.

## Security Guidance

- Never commit collector tokens.
- Prefer environment injection at runtime.
- Rotate credentials regularly.
