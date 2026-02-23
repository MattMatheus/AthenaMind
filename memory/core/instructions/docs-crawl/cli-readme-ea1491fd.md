# Cli Readme Ea1491fd

# CLI Reference

## Summary
Reference entry point for user-facing memory CLI commands, options, and expected outputs.

## Intended Audience
- Users running AthenaMind memory CLI flows directly.
- QA/operators validating command behavior and evidence capture.

## Preconditions
- Go toolchain and CLI build/test path are working.
- User understands baseline concepts and lifecycle model.

## Main Flow
1. Use `commands.md` for full command and flag reference.
2. Use `examples.md` for copy/paste workflows.
3. Validate behavior against `cmd/memory-cli/main_test.go`.
4. Cross-check with `knowledge-base/concepts/README.md` contracts.

## Failure Modes
- Command examples drift from implemented behavior.
- Missing argument documentation causes invalid requests.
- Output contracts change without docs updates.

## References
- `cmd/memory-cli/main.go`
- `cmd/memory-cli/main_test.go`
- `product-research/architecture/MEMORY_CLI_V1_ARCHITECTURE.md`
- `delivery-backlog/engineering/done/STORY-20260222-memory-cli-v1-write-retrieve-baseline.md`

## Pages
- `commands.md`
- `examples.md`
