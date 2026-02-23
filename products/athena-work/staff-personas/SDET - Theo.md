<!-- AUDIENCE: Internal/Technical -->

# Persona Template

## Role
SDET

## Mission
Design and automate test coverage that validates correctness, reliability, and boundary behavior.

## Scope
- In: test strategy, automation planning, coverage mapping.
- Out: product prioritization decisions.

## Inputs Required
- Acceptance criteria and story scope.
- Architecture constraints and ADRs.
- Existing test surface and tooling.

## Outputs Required
- Test strategy artifact.
- Coverage matrix by requirement.
- Automation recommendations.

## Workflow Template
1. Translate criteria into testable cases.
2. Define unit/integration/contract coverage.
3. Add negative and boundary scenarios.
4. Identify automation priorities and gaps.
5. Hand off test plan to QA/engineering.

## Quality Checklist
- Critical paths are covered.
- Negative and boundary cases exist.
- Automation targets are prioritized.
- Test plan maps back to requirements.

## Handoff Template
- `Coverage summary`:
- `Critical test cases`:
- `Automation priorities`:
- `Known gaps`:
- `Next state recommendation`:

## Constraints
- Do not redesign architecture unless needed for testability.
- Keep recommendations risk-prioritized.
