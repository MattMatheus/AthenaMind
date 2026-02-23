# Quickstart

## Summary
Create one memory entry, retrieve it, and confirm deterministic output fields.

## Preconditions
- Tooling validated via `installation.md`.
- Working directory at repository root.

## Main Flow
1. Write one approved memory entry:
```bash
go run ./cmd/memory-cli write \
  --root memory \
  --id getting-started \
  --title "Getting Started Prompt" \
  --type prompt \
  --domain onboarding \
  --body "Use this prompt to onboard a new engineer quickly." \
  --stage planning \
  --reviewer maya \
  --decision approved \
  --reason "onboarding baseline" \
  --risk "low; rollback by git revert" \
  --notes "approved in docs quickstart"
```
2. Retrieve using natural language:
```bash
go run ./cmd/memory-cli retrieve \
  --root memory \
  --query "onboard engineer prompt"
```
3. Confirm response contains:
   - `selected_id`
   - `selection_mode`
   - `source_path`

## Failure Modes
- `ERR_MUTATION_STAGE_INVALID`: use `--stage planning|architect|pm`.
- `ERR_MUTATION_EVIDENCE_REQUIRED`: provide `--reason --risk --notes`.
- `--query is required`: include `--query`.

## References
- `docs/cli/commands.md`
- `docs/workflows/write-retrieve-baseline.md`
