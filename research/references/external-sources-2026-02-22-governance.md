# External Sources: Governance and Human-in-the-Loop (2026-02-22)

## Sources
- NIST AI Risk Management Framework 1.0
  - https://www.nist.gov/itl/ai-risk-management-framework
- OWASP Top 10 for LLM Applications
  - https://owasp.org/www-project-top-10-for-large-language-model-applications/
- OpenTelemetry semantic conventions (for auditable policy events)
  - https://opentelemetry.io/docs/specs/semconv/general/trace/

## Design-Relevant Takeaways
- Governance controls should be explicit, measurable, and continuously monitored (not one-time policies).
- Human oversight is mandatory for higher-risk operations and low-confidence outputs.
- Prompt/data handling risks require policy checks before memory promotion or autonomous execution.
- Audit trails must capture decision rationale and policy path for post-incident analysis.
