# Accessibility Profiles

## Purpose
Provide reusable workflow defaults that improve operating conditions for humans with different accessibility needs while preserving deterministic agent execution.

## Principles
- Make the next action obvious.
- Keep one active focus by default.
- Reduce avoidable context switching.
- Use neutral, respectful language in artifacts and UI.
- Preserve auditability and policy gates.

## Profile: Low-Vision Default
### Planning/UI defaults
- High contrast color pairs only.
- Minimum body text equivalent of 16px.
- Minimum line-height of 1.5.
- Persistent headings and section anchors.
- Avoid color-only status signaling; include text labels.

### Workflow defaults
- Show one stage and one top queue item first.
- Keep command outputs concise with explicit pass/fail markers.
- Require plain-language error summaries before raw logs.

## Profile: High-Variability Attention
### Planning/UI defaults
- Single-task board mode (one primary task, collapsed secondary tasks).
- Short checklist slices (max 5 items per active checklist).
- Stable layout with minimal motion and no auto-refresh jumps.
- Prominent completion signals (`done`, `blocked`, `needs human decision`).

### Workflow defaults
- Enforce explicit `next_stage` and `next_actions` in handoffs.
- Prefer deterministic, machine-readable errors with recovery steps.
- Promote one item at a time unless human explicitly requests batching.

## Inclusive Communication Standard
- Preferred terminology: `high-variability attention`, `attention-support mode`.
- Avoid stigmatizing labels in docs, prompts, and UI text.
- Record accommodation choices as workflow settings, not personal judgments.

## Operationalization Checklist
1. Select one or more accessibility profiles at cycle start.
2. Record selected profile(s) in cycle metadata and handoff notes.
3. Validate stage prompts and UI artifacts against selected profile rules.
4. Confirm low-vision defaults are present for all human-facing views.
5. Include QA evidence for accessibility conformance before cycle closure.

## Evidence Hooks
- Decision artifacts: note selected profile and deviations.
- Handoff artifacts: include pass/fail accessibility checklist.
- Observer reports: include accessibility drift findings if present.
