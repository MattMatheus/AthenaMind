# Installation

## Summary
Install the minimum tooling and verify AthenaMind can run tests locally.

## Preconditions
- Repository cloned locally.
- macOS or Linux terminal access.

## Steps
1. Install Go 1.22+ (source of truth: `go.mod`).
2. Configure writable memory storage:
   - Default runtime path is `~/.athena/memory/<repo-id>`.
   - In Codex sandbox, add `~/.athena` to writable roots in Codex config.
   - Fallback option: set `ATHENA_MEMORY_ROOT` to a writable workspace path.
3. Validate the toolchain:
   - `tools/check_go_toolchain.sh`
4. Run baseline tests:
   - `tools/run_stage_tests.sh --mode push`
5. Run documentation gates:
   - `tools/run_doc_tests.sh`

## Failure Modes
- `go` missing or old: follow `knowledge-base/how-to/GO_TOOLCHAIN_SETUP.md`.
- Memory path permission denied in Codex: add `~/.athena` to writable roots or export `ATHENA_MEMORY_ROOT` to a writable workspace path.
- Failing tests: resolve code-level issues before using docs workflows.

## References
- `knowledge-base/how-to/GO_TOOLCHAIN_SETUP.md`
- `cmd/memory-cli/main_test.go`
