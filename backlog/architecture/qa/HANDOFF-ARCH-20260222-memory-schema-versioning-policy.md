# Architecture Handoff: ARCH-20260222-memory-schema-versioning-policy

## Decision(s) Made
- Accepted a canonical schema versioning policy for memory artifacts using `schema_version` with `MAJOR.MINOR` semantics.
- Defined compatibility behavior for newer minor versions and fail-fast behavior for unsupported major versions.
- Specified required fields and validation rules for `index.yaml` and metadata files, plus migration expectations for incompatible changes.

## Alternatives Considered
- No explicit versioning contract.
  - Rejected due to compatibility drift risk across CLI upgrades.
- Lockstep strict schema with no additive extensions.
  - Rejected because it slows low-risk iterative improvements.

## Risks and Mitigations
- Risk: migration overhead increases implementation complexity.
  - Mitigation: additive minor policy and deterministic major migration rules reduce disruption.
- Risk: invalid legacy artifacts block retrieval/write paths.
  - Mitigation: explicit validation errors and migration-needed signals.

## Updated Artifacts
- `research/decisions/ADR-0010-memory-schema-versioning-policy.md`
- `research/architecture/MEMORY_SCHEMA_AND_VERSIONING_POLICY.md`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `backlog/architecture/qa/ARCH-20260222-memory-schema-versioning-policy.md`

## Validation Commands and Results
- `scripts/run_doc_tests.sh` -> PASS

## Open Questions for QA Focus
- Are schema validation/error contracts specific enough for implementation tests?
- Is the major/minor compatibility policy enforceable without ambiguity in CLI behavior?
