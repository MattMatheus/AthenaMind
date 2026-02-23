# Product Vision F02a617c

# AthenaMind Product Vision (Long-Term)

## Status
Preserved long-term vision artifact. This document is not implementation scope for v0.1 engineering cycles.

## Scope Boundary
This vision defines future direction beyond the current memory-layer delivery. Current execution scope remains v0.1 memory-layer only (ADR-0007).

## Core Insight
Agents spend significant context capacity on orientation, and that cost repeats every session. AthenaMind externalizes operating knowledge into a governed memory layer so agents and humans can access the right context quickly and consistently.

## Memory Hierarchy
- Core skills: always-loaded operating basics
- Domain skills: context-matched expertise
- Repo skills: codebase-specific operating knowledge
- Episode memory: session/cycle outcomes and decisions

## O(1) Orientation Protocol
1. Load core operating skill
2. Query launcher for current work
3. Retrieve repo/domain skills
4. Retrieve episode context from prior cycles

## Two Consumers, One Memory
The same governed memory store supports:
- Agent execution guidance
- Human explanation and onboarding guidance

## Onboarding Thesis
Organizational knowledge is fragmented across repos and operational systems. A governed memory layer compresses onboarding time by surfacing real patterns and decisions from actual project artifacts.

## Layered Policy Hierarchy (Future Team/Org Scale)
- Org level policies
- Team level conventions
- Repo level standards
- Session level context

## Governance as Product Moat
Auditability, policy-gated writes, review evidence, and traceability are required for shared organizational memory at scale.

## Enterprise Service Direction
Longer-term model includes memory-engineering services to configure ingestion, review gates, indexing, and compliance evidence for organizational knowledge sources.

## Revenue Direction
- Open core: local single-operator memory CLI
- Team/org tier: centralized governed memory
- Enterprise services: ingestion setup, skill-pack curation, compliance reporting

## v0.1 Building Blocks Toward Vision
| v0.1 Primitive | Long-term Role |
| --- | --- |
| File-based memory store | Adapter baseline for future backends |
| Governance gates | Foundation for org-level write policy |
| Audit and telemetry | Compliance and trust evidence |
| Retrieval pipeline | Skill-pack retrieval path |
| Bootstrap protocol | Agent onboarding flow |
| Episode write-back | Compounding memory loop |
| Evaluation harness | Retrieval quality SLA baseline |

## References
- `product-research/product/PRODUCT_VISION_V2.md`
- `product-research/product/VISION_WORKSHOP_2026-02-22.md`
- `product-research/roadmap/PHASED_IMPLEMENTATION_PLAN_V01_V03.md`
- `product-research/decisions/ADR-0007-memory-layer-scope-refinement.md`
