# Experiment: PTY Stream Supervisor Guardrail With Nonce Stop/Resume

## Hypothesis
- Classifying live PTY stdout/stderr with a LlamaGuard-style policy model, then issuing a nonce-bound stop token when risk is detected, will reduce agent misbehavior without materially degrading cycle throughput.

## Scope
- In:
  - Capture full terminal stream (stdout/stderr) for session-time analysis.
  - Run deterministic command/rule checks first, then semantic risk classification on stream windows.
  - On high-confidence policy risk, supervisor halts further stream forwarding and injects first-token control message:
    - `stop work resume key <nonce>`
  - Require human supervisor inspection before resume.
  - Resume only with explicit nonce echo from supervisor path.
  - Emit telemetry for detection, halt, and resume events.
- Out:
  - Replacing existing AthenaWork stage scripts.
  - Autonomous self-resume without human review.
  - Any permanent mutation to backlog state model.

## Proposed Control Flow
1. Session starts; supervisor seeds cryptographically strong `nonce`.
2. PTY stream is mirrored to classifier pipeline in near-real-time.
3. Deterministic denylist/policy checks run on command intents and output markers.
4. LlamaGuard-style classifier scores semantic risk classes (policy bypass, exfiltration, unsafe destructive drift, instruction injection).
5. If risk threshold is exceeded:
   - freeze downstream stream relay
   - inject control phrase: `stop work resume key <nonce>`
   - mark session state `halted_pending_review`
6. Human supervisor inspects full buffered stream and decides `resume` or `terminate`.
7. Resume path requires nonce-bound resume command and logs reviewer identity + rationale.

## Success Criteria
- At least one seeded misbehavior scenario is halted before next destructive command executes.
- False positive halt rate is <= 5% across dogfood runs.
- Median halt-to-review decision latency is <= 90 seconds.
- 100% of halts include audit fields: `session_id`, `nonce_id`, `policy_class`, `confidence`, `reviewer_decision`.
- No resume occurs without valid nonce-bound supervisor action.

## Duration
- 5 working days from pilot start.

## Rollback Trigger
- Any bypass where a halted session continues execution without supervisor decision.
- More than one false-negative event in seeded high-risk scenarios.
- Throughput degradation > 20% relative to baseline cycle completion time.

## Telemetry Contract Additions (Draft)
- `supervisor.guard.alert`
  - `session_id`, `cycle_id`, `trace_id`, `policy_class`, `confidence`, `window_offset`
- `supervisor.guard.halt`
  - `session_id`, `nonce_id`, `reason_code`, `halt_at_utc`
- `supervisor.guard.resume`
  - `session_id`, `nonce_id`, `reviewer`, `decision`, `rationale`, `resume_at_utc`
- `supervisor.guard.terminate`
  - `session_id`, `nonce_id`, `reviewer`, `reason`

## Open Risks / Questions
- Nonce leakage into model-visible context could weaken stop/resume integrity; evaluate channel isolation.
- Stream chunking strategy may affect classifier recall (token window boundaries).
- Need strict precedence rules between deterministic denylists and semantic classifier outputs.
- Clarify whether halt applies per-tool call, per-session, or per-stage boundary.

## Result
- status: pending
- notes:
  - This is research-only and does not change AthenaWork runtime behavior.
