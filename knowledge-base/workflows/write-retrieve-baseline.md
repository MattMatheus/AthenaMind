# Write/Retrieve Baseline Workflow

## Goal

Validate governed memory write and deterministic retrieval metadata.

## Steps

1. Write an approved entry.
2. Retrieve with natural language query.
3. Confirm output includes `selected_id`, `selection_mode`, `source_path`.
4. Repeat query to confirm deterministic fallback behavior on low-confidence paths.

## Success Criteria

- Entry files exist under memory root.
- Retrieval output includes required fields.
- No policy violations.

## Related Docs

- [Quickstart](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/getting-started/quickstart.md)
- [CLI Commands](https://github.com/MattMatheus/AthenaMind/blob/main/knowledge-base/cli/commands.md)
