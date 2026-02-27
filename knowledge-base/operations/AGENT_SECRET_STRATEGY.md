# Agent Secret Strategy

## Goal
Define secret handling controls for agent execution environments beyond basic post-use rotation.

## Baseline
- Rotation after use remains mandatory.
- Secret values must not be stored in repo state, workflow metadata, or default logs.

## Required Controls Beyond Rotation

## 1. Secret Scope Minimization
- Issue per-agent, per-workflow, and per-purpose credentials.
- Use shortest feasible TTL.
- Bind credentials to destination scope (service, namespace, environment).

## 2. Just-In-Time Minting
- Mint secrets at step start, not at workflow start.
- Revoke immediately on step completion or failure.
- Avoid long-lived static credentials in pods.

## 3. Identity First
- Prefer workload identity (OIDC/IRSA/GKE WI/etc.) over injected static secrets.
- Use federated identity to obtain ephemeral service tokens.

## 4. Non-Bypassable Delivery Path
- Inject secrets only through a controlled secret broker path.
- Block ad hoc secret delivery via env files, annotations, or side channels.

## 5. Usage-Bound Policy
- Secret policy should encode:
  - who can use it
  - where it can be used
  - what actions are allowed
  - when it expires
- Deny if usage context mismatches policy claims.

## 6. Egress and Destination Controls
- Enforce network allowlists per workflow lane.
- Limit secret-bearing workloads to approved destination endpoints.

## 7. Runtime Memory Hygiene
- Avoid writing secrets to disk.
- Zero memory buffers where feasible for sensitive material.
- Disable core dumps in secret-bearing containers.

## 8. Audit + Trace Correlation
- Audit every secret issue, read, and revoke with `request_id`/`cycle_id`.
- Track secret usage events as OTel spans with metadata only (never values).

## 9. Detection and Guardrails
- Add policy checks for suspicious secret access frequency and destination drift.
- Alert on reuse outside TTL/window or by unexpected identity.

## 10. Break-Glass Protocol
- Require explicit operator approval for emergency secret override.
- Time-box, log, and auto-revoke break-glass credentials.

## 11. Supply Chain Boundaries
- Sign and verify agent images.
- Restrict secret access to trusted image attestations and approved runtime classes.

## 12. Testing and Chaos Validation
- Run periodic secret leakage tests (logs, traces, artifacts, crash paths).
- Run revocation chaos tests to ensure workflows fail closed correctly.

## Implementation Pattern for AthenaWork
1. Identity token from workload identity.
2. Broker mints short-lived secret for exact step scope.
3. Agent uses secret via in-memory handle.
4. `athenamind-transfer` emits usage audit metadata.
5. Secret revoked on completion/failure/timeout.

## Minimal Acceptance Checks
- No secret value appears in logs, traces, artifacts, or backlog docs.
- Expired/revoked secret usage is denied in runtime.
- Secret issuance and transfer events are correlated by trace id.
- Break-glass actions are visible and auto-expire.
