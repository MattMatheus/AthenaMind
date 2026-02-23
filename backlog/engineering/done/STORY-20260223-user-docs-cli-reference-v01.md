# Story: Build Comprehensive User CLI Reference

## Metadata
- `id`: STORY-20260223-user-docs-cli-reference-v01
- `owner_persona`: Technical Writer - Clara.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.1
- `adr_refs`: [ADR-0009, ADR-0015]
- `success_metric`: CLI reference documents all supported v0.1 commands with args, examples, and expected outputs
- `release_checkpoint`: required

## Problem Statement
- There is no canonical user-facing CLI command reference, causing reliance on code reading and ad hoc tribal knowledge.

## Scope
- In:
  - Produce `docs/cli/` command reference pages.
  - Document command syntax, options, and examples.
  - Capture output/error contracts and telemetry-relevant notes where user-visible.
- Out:
  - New CLI command implementation.
  - Non-v0.1 command sets.

## Acceptance Criteria
1. Every user-supported v0.1 CLI command is documented.
2. Each command has example usage and expected output behavior.
3. Reference is traceable to implementation and tests.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `cmd/memory-cli/main.go`
- `cmd/memory-cli/main_test.go`
- `research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`

## Notes
- Favor concise reference tables and tested command examples.
