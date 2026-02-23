# Troubleshooting Readme 999aea86

# Troubleshooting

## Summary
Diagnosis and recovery paths for common user-facing setup and workflow failures.

## Intended Audience
- Users blocked while running AthenaMind CLI workflows.
- QA/operators triaging reproducible documentation or behavior issues.

## Preconditions
- User can provide command output and environment context.
- Baseline setup and workflow docs have been followed.

## Main Flow
1. Find your error in `common-errors.md`.
2. Apply the fix and rerun the command.
3. Validate with tests and doc gates.
4. Open intake bug if unresolved.

## Failure Modes
- Missing logs/evidence prevents root-cause isolation.
- Recovery steps are applied out of order.
- Known issue is rediscovered because docs are stale.

## References
- `knowledge-base/how-to/GO_TOOLCHAIN_SETUP.md`
- `knowledge-base/cli/README.md`
- `delivery-backlog/engineering/intake/BUG_TEMPLATE.md`
- `knowledge-base/process/STAGE_EXIT_GATES.md`

## Pages
- `common-errors.md`
