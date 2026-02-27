# AthenaMind Transfer Protocol v1

## Goal
Define the contract for a dedicated `athenamind-transfer` binary that handles inter-agent state transmission with policy, audit, and observability controls.

## Scope
- In:
  - Envelope schema for transmissions.
  - Binary command contract.
  - Adjudication and fail-closed behavior.
  - Protobuf + signature + replay-protection requirements.
- Out:
  - Retrieval/indexing semantics (owned by existing AthenaMind memory paths).

## Canonical Schema
- Protobuf file:
  - `products/athena-mind/proto/athenamind/transfer/v1/transfer.proto`
- Core entities:
  - `TransferEnvelope`
  - `TransferRequest`
  - `TransferResult`
  - `TransferService.Send`

## Adjudication Model
- Guard model is advisory only.
- Deterministic policy engine is final authority (`allow|deny|redact`).
- If adjudication, audit write, signature, or verification fails: deny transfer.

## Binary Contract (`athenamind-transfer`)

### `athenamind-transfer send`
- Purpose: validate/adjudicate/sign/transmit an envelope.
- Required inputs:
  - `--request-id`
  - `--session-id`
  - `--cycle-id`
  - `--source-agent-id`
  - `--target-agent-id`
  - `--domain`
  - `--classification`
  - `--policy-stage`
  - `--memory-id` (repeatable)
  - one of:
    - `--payload-file <path>`
    - `--payload-ref <opaque-reference>`
- Optional:
  - `--story-id`
  - `--ttl-seconds` (default 300)
  - `--route`
  - `--audit-attr key=value` (repeatable)
  - `--emit-envelope <path>` (write signed protobuf envelope)

### `athenamind-transfer receive`
- Purpose: verify/decrypt/re-evaluate policy for destination and accept or reject.
- Required inputs:
  - `--envelope-file <path>` or `--envelope-b64 <value>`
  - `--target-agent-id`
- Optional:
  - `--emit-result <path>`

### `athenamind-transfer verify`
- Purpose: cryptographic and policy preflight checks without delivery.
- Inputs:
  - `--envelope-file <path>`
- Outputs:
  - validation result + failure reason.

### `athenamind-transfer inspect`
- Purpose: metadata-only inspection for operator debugging.
- Inputs:
  - `--envelope-file <path>`
- Rule:
  - never emit raw payload unless explicit break-glass mode is enabled.

## Security Requirements
- Signature required for every envelope.
- Replay protection required (`nonce`, `issued_at_unix`, `expires_at_unix`, dedupe cache).
- Envelope classification required and policy-evaluated at send and receive.
- Payload encryption required in transit and at rest.
- No raw payload in default logs/traces.

## OpenTelemetry Requirements
- Root span: `memory.transfer`
- Child spans:
  - `memory.policy_eval`
  - `memory.sign`
  - `memory.transmit`
  - `memory.verify`
  - `memory.deliver`
- Required attrs:
  - `request_id`, `session_id`, `cycle_id`, `source_agent_id`, `target_agent_id`
  - `classification`, `decision`, `risk_level`, `payload_bytes`

## Failure Semantics
- Deny by default.
- Hard fail conditions:
  - missing required metadata
  - invalid signature
  - expired TTL
  - nonce replay
  - unresolved destination allowlist
  - audit sink unavailable

## Interop Rule
- Governed lanes may only exchange state through this protocol.
- Non-canonical channels are blocked unless explicitly enabled in experiment namespace.

## Related
- `knowledge-base/architecture/AGENT_STATE_TRANSFER_AND_OBSERVABILITY.md`
- `knowledge-base/architecture/ARGO_MULTI_AGENT_CONTROL_PLANE_RESEARCH.md`
- `products/athena-work/operating-system/decisions/WSI-ADR-0002-argo-multi-agent-delegation-model.md`
