# Installation

## Summary
Install the minimum tooling and verify AthenaMind can run tests locally.

## Preconditions
- Repository cloned locally.
- macOS or Linux terminal access.

## Steps
1. Install Go 1.22+ (source of truth: `go.mod`).
2. Validate the toolchain:
   - `tools/check_go_toolchain.sh`
3. Run baseline tests:
   - `go test ./...`
4. Run documentation gates:
   - `tools/run_doc_tests.sh`

## Failure Modes
- `go` missing or old: follow `knowledge-base/how-to/GO_TOOLCHAIN_SETUP.md`.
- Failing tests: resolve code-level issues before using docs workflows.

## References
- `knowledge-base/how-to/GO_TOOLCHAIN_SETUP.md`
- `cmd/memory-cli/main_test.go`
