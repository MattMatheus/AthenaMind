# External Sources: MVP Constraints (2026-02-22)

## Sources
- Google SRE Workbook (SLO + error budget policy)
  - https://sre.google/workbook/implementing-slos/
  - https://sre.google/workbook/error-budget-policy/
- Podman rootless/runtime limits
  - https://docs.podman.io/en/v4.0.0/markdown/podman.1.html
  - https://docs.podman.io/en/stable/markdown/podman-pod-create.1.html
- SQLite WAL behavior
  - https://sqlite.org/wal.html
- OpenTelemetry HTTP spans/metrics semantic guidance
  - https://opentelemetry.io/docs/specs/semconv/http/http-spans/
  - https://opentelemetry.io/docs/specs/semconv/http/http-metrics/

## Design-Relevant Takeaways
- Reliability targets should be expressed as explicit SLOs with a written error budget policy.
- Runtime safety/cost bounds must be enforceable at container runtime (CPU/memory/time) not just advisory.
- Local-first state persistence can rely on SQLite WAL with known single-host constraints.
- Traceability constraints should align span/metric semantics to avoid ambiguous latency or incomplete audits.
