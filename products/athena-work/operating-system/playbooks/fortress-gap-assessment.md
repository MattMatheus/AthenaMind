# Fortress Gap Assessment Playbook

Run this review before promoting major workflow/security architecture changes.

## Control Checklist

### Data Lineage and Deletion
- [ ] End-to-end lineage exists for memory payload movement across retrieve, adjudicate, transfer, and receive.
- [ ] Deletion workflow can prove removal from primary stores, caches, artifacts, and log sinks.
- [ ] Retention rules are defined per data class and enforced automatically.

### Prompt and Policy Supply Chain
- [ ] Prompts, policies, and model config are versioned, reviewed, and linked to change records.
- [ ] Integrity checks (signing/verification or equivalent) are in place for promoted policy bundles.
- [ ] Drift detection alerts when runtime prompt/policy deviates from approved baseline.

### Egress-by-Intent
- [ ] Egress policy is tied to workflow purpose, not just destination allowlists.
- [ ] Secret-bearing workloads have stricter egress controls than non-secret workloads.
- [ ] Cross-domain transfer attempts produce deny events with actionable metadata.

### Model Risk Controls
- [ ] Model/dependency versions are pinned for governed lanes.
- [ ] Canary validation exists for policy-sensitive model or guard changes.
- [ ] Fallback behavior is defined for guard outage or elevated false positive rates.

### Adversarial and Insider Scenarios
- [ ] Abuse scenarios include malicious agent behavior, not only accidental leakage.
- [ ] Honey-token or deception controls exist for detection of suspicious exfiltration behavior.
- [ ] Incident response playbook includes containment and forensic collection steps.

### Segregation of Duties
- [ ] No single component can retrieve, adjudicate, transfer, and approve override alone.
- [ ] Break-glass approval requires independent authority and generates immutable audit evidence.
- [ ] High-risk policy changes require dual control before promotion.

### Evidence Automation
- [ ] Audit evidence bundle is generated continuously (not manually before reviews).
- [ ] Exception reports include owner, mitigation, and closure criteria.
- [ ] Compliance-ready control snapshots are reproducible from repo + telemetry sources.

### Rollback of Agent Changes (Required)
- [ ] Agent-authored code changes have a documented rollback path (`git revert` plan, release rollback plan, and affected scope).
- [ ] Prompt/policy changes have versioned rollback targets and a tested revert command path.
- [ ] Model/guard version changes have a pinned previous version and canary rollback switch.
- [ ] Secret broker, key, and credential rollout changes include emergency revoke + restore procedure.
- [ ] Runtime rollback drill has been executed at least once for:
  - code deployment rollback
  - policy bundle rollback
  - model/guard rollback
  - transfer channel deny-mode enforcement

## Exit Gate
- [ ] No open critical findings.
- [ ] Open high findings have approved owner, mitigation, and due trigger.
- [ ] Rollback controls are tested and evidenced.
