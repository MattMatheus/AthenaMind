# Go Toolchain Setup for Memory CLI

## Purpose
Ensure engineering and QA environments can run `go version` and `go test ./...` for the memory CLI.

## Required Version
- Source of truth: `go.mod`
- Current minimum: `go 1.22`

## Install (macOS/Homebrew)
1. Install Go:
   - `brew install go`
2. Verify installation:
   - `go version`
3. Verify repo requirement and preflight:
   - `tools/check_go_toolchain.sh`
4. Run module tests:
   - `go test ./...`

## Install (Linux)
1. Install Go from your distro package manager or the official tarball.
2. Verify installation:
   - `go version`
3. Verify repo requirement and preflight:
   - `tools/check_go_toolchain.sh`
4. Run module tests:
   - `go test ./...`

## Troubleshooting
- If `go` is not found:
  - Ensure your shell PATH includes the Go binary location.
  - Restart shell and re-run `go version`.
- If preflight fails on version:
  - Upgrade Go to at least the version declared in `go.mod`.
- If tests fail after preflight passes:
  - This is a code/test issue; continue with normal engineering debug workflow.
