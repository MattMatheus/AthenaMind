Agents - This is NOT part of your normal workflow.  You should ignore this file unless founder directs you to inspect or implement from this file.

# Session Briefing: AthenaMind Evaluation and Strategic Direction
**Date:** 2026-02-23
**Participants:** Matt Matheus (founder/operator), Claude (Copilot CLI)


---

## 1. Repository Evaluation

### Memory System (Product) — cmd/memory-cli

**What exists today:**
- Go CLI binary, ~1,800 lines in a single `cmd/memory-cli/main.go`
- `go.mod` targets Go 1.22 with zero external dependencies
- Subcommands: `write`, `retrieve`, `evaluate`, `snapshot` (create/list/restore), `serve-read-gateway`, `api-retrieve`
- File-based storage: markdown content files + JSON-serialized `.yaml` index/metadata
- Retrieval pipeline: token-overlap semantic scoring → confidence gate → deterministic fallback (exact-key → path-priority)
- Governance: review-first write policy, autonomous-run guards (`AUTONOMOUS_RUN` env var blocks writes), rejected-decision audit trails with rework requirements
- Telemetry: structured JSONL events with trace IDs, session IDs, scenario IDs, operator verdicts, latency tracking
- Constraint enforcement: cost budgets (`MEMORY_CONSTRAINT_COST_MAX_PER_RUN_USD`), traceability enforcement, reliability freeze gates, latency degradation fallbacks — all fail-closed
- Schema versioning: major/minor compatibility checks on index and metadata files
- Snapshot lifecycle: create with checksums/manifests, list, restore with compatibility blocking
- Evaluation harness: 50-query benchmark measuring top-1 useful rate, fallback determinism, selection mode reporting, source trace completeness
- API gateway: HTTP read endpoint with CLI-parity fallback when gateway is unavailable
- 871 lines of tests covering write/retrieve flows, governance enforcement, schema compat, telemetry emission, snapshot lifecycle, constraint enforcement, API parity

**Key strengths identified:**
- Governance-first design is genuinely differentiated — most memory systems ignore trust/auditability entirely
- Metric-driven development with ADR-0002 scorecard and KPI set is unusually disciplined
- "Precision over volume" retrieval philosophy directly addresses the most common real-world failure mode (context stuffing)
- Constraint enforcement (cost, traceability, reliability, latency) is production-grade thinking applied early

**Key concerns identified:**
- Single-file monolith: 1,822 lines with no packages despite ADR-0003 defining module boundaries
- No real semantic retrieval: token-overlap scoring, no embeddings despite README/ADRs describing FAISS
- JSON files named `.yaml` — confusing for operators inspecting the store
- Zero external dependencies means hand-rolling everything
- Tests are thorough but only exercise through CLI entry points — no unit tests for internal logic in isolation

### Agentic Engineering System (Process)

**What exists today:**
- Full SDLC encoded in markdown: Planning → Architect → Engineering → QA → PM Refinement → Observer → Commit
- 91 total commits, 33 engineering stories done, 11 architecture stories done, 25+ observer reports
- Stage launcher (`tools/launch_stage.sh`) enforces branch safety and emits scoped seed prompts
- Backlog state machine: intake → ready → active → qa → done (with blocked/archive lanes)
- Intake validation script with metadata enforcement
- 7 named agent personas with structured templates (Role/Mission/Scope/Inputs/Outputs/Workflow/Constraints)
- Observer pattern: deterministic diff/memory-delta capture at cycle boundary
- Atomic commit convention: one commit per cycle, all artifacts included
- Doc-test harness: shell-based tests validating process artifact structure
- 17 ADRs documenting architectural decisions
- 315 markdown files across the repo

**Critical reframe from founder:** Agents are the sole operators. The documentation density is intentional — it's the runtime environment for stateless workers. Every guardrail that a human team absorbs through osmosis must be explicit for agents.

**Key strengths:**
- O(1) orientation protocol: agent reads AGENTS.md → runs launcher → gets one story. Constant cost regardless of codebase size.
- Demonstrated 9-second agent self-orientation (branch check, queue query, launcher resolve, handoff confirm)
- Observer pattern creates accountability and traceability most agentic workflows lack
- "done ≠ shipped" distinction with release bundles prevents premature shipping
- The process investment is designed to compound: future agent sessions are more productive because the rails exist

**Key concerns:**
- 150:1 docs-to-code ratio (justified by agent-operator context, but needs to earn its keep in sprint velocity)
- Risk of legibility attractor: agents preferring process stories over harder product work
- Some artifact staleness (founder snapshot vs program board counts diverged)
- No CI/CD — all verification is manual shell scripts

---

## 2. Strategic Direction Delivered

### Sprint Sequence (8 stories created in delivery-backlog/engineering/intake/)

**Sprint 3 — Foundation Hardening:**
1. `STORY-20260223-split-main-go-into-packages` — extract internal/index, governance, retrieval, telemetry, snapshot, gateway packages; reduce main.go to thin CLI wiring
2. `STORY-20260223-add-ci-pipeline` — GitHub Actions running `go test ./...` + `tools/run_doc_tests.sh` on push to dev
3. `STORY-20260223-readme-v01-alignment` — rewrite README to reflect v0.1 reality, move vision to knowledge-base/product/VISION.md

**Sprint 4 — Dogfooding Loop (highest leverage):**
4. `STORY-20260223-agent-bootstrap-protocol` — `memory-cli bootstrap` subcommand returning structured context payload for agent sessions (depends on #1)
5. `STORY-20260223-episode-writeback` — `memory-cli episode write/list` for structured cycle records (depends on #4)
6. `STORY-20260223-launch-stage-memory-integration` — wire bootstrap into launch_stage.sh, episode write into observer script; soft integration with graceful degradation (depends on #4, #5)

**Sprint 5 — Real Retrieval:**
7. `STORY-20260223-sqlite-index-store` — v0.2, SQLite WAL-mode adapter replacing JSON file index, auto-migration, pure Go driver (depends on #1)
8. `STORY-20260223-embedding-retrieval-via-ollama` — v0.2, local embeddings via Docker-hosted Ollama, quality gate against 50-query benchmark (depends on #7, #1)

**Done signals per phase:**
- Phase 1: `go test ./...` passes in CI with split packages
- Phase 2: agent sessions start with retrieved memory context and write episodes back
- Phase 3: embedding retrieval beats token-overlap on the 50-query benchmark

**What NOT to do next:**
- No more process/meta-engineering stories until Phase 2 is complete
- No Pillar 3 (workspace-runtime/pods) — stay memory-layer-only per ADR-0007
- No cloud infra spend — Docker/K8s locally simulates everything needed
- No multi-operator governance yet — single-operator is fine through v0.2
- No UI — CLI and API are sufficient for agent consumers

### Docker/K8s Local Stack (planned for Sprint 5+)
- `memory-api`: Go read gateway service
- `ollama`: local embedding model (e.g., nomic-embed-text)
- SQLite embedded in memory-api container
- Optional: prometheus + grafana for KPI dashboards
- Runs on laptop, costs nothing, mirrors eventual cloud topology

---

## 3. Product Vision (captured in product-research/product/PRODUCT_VISION_V2.md)

### Core Insight
Agents waste 25%+ of context window on orientation. AthenaMind reduces this to near-zero through externalized, governed memory.

### Memory Hierarchy / Skill Pack Taxonomy
- **Core skills** — always loaded at bootstrap (cycle protocol, handoff format). 500-2000 tokens, stable.
- **Domain skills** — loaded by context match (Go patterns, SRE runbooks). Retrieved semantically.
- **Repo skills** — loaded per codebase (build commands, architecture). Generated by first-agent exploration.
- **Episode memory** — loaded on demand (what happened last session, what failed, what was decided).

### Key Insight: Two Consumers, One Memory
Same memory store serves both machines and humans:
- Agent asks "How does this org do DI in C#?" → gets executable skill pack
- Human asks same question → gets explanation with examples from actual codebase
- Both grounded in how the org actually does it, not Stack Overflow

### Organizational Knowledge as Side Effect
Every agent session that writes an episode, every decision record, every codebase bootstrap adds to the knowledge base. Agents document the organization as a side effect of working.

### Enterprise Engagement: The Memory Engineer
Services role that configures ingest pipeline for an org:
1. Points AthenaMind at knowledge sources (repos, wikis, runbooks, Confluence, Slack, ADRs, postmortems)
2. Specialized process crawls, sorts, indexes, embeds, vectorizes
3. Human review gates flag low-confidence entries
4. Org gets governed, auditable, queryable knowledge base
5. Audit trail proves exactly what was ingested — "the receipts"

### Revenue Model (if ever pursued)
- Open core: local single-operator CLI
- Team/org tier: centralized memory with shared retrieval and role-based governance
- Enterprise services: memory engineer engagement, custom ingest, compliance reporting

### v0.1 Primitives → Long-term Vision Mapping
| v0.1 Primitive | Long-term Role |
| --- | --- |
| File-based memory store | Adapter interface for any backend |
| Governance gates | Org-level write policy |
| Audit/telemetry | Compliance reporting |
| Retrieval pipeline | Skill pack retrieval |
| Bootstrap protocol | Agent onboarding at any scale |
| Episode write-back | Compounding knowledge loop |
| Evaluation harness | Retrieval quality SLA for enterprise |

---

## 4. Steelman Objections (strongest counterarguments)

1. **Context windows are growing faster than you can ship.** 200K+ tokens today. The 25% orientation tax may shrink to 2% through brute force before the product ships.

2. **Governance moat may be a governance wall.** Enterprise wants governance but hates friction. A competitor with "good enough" governance and zero-friction writes could win adoption first.

3. **One person building enterprise infrastructure.** The distance from working CLI to org-scale knowledge platform is enormous. Why won't Microsoft/Google/Anthropic add memory to their agent platforms natively?

4. **Dual-consumer model is two products pretending to be one.** Agents want structured machine-parseable packs. Humans want natural language explanations. Serving both well from one store means awkward format for one consumer or double the rendering surface.

5. **Ingest is where this dies.** Enterprise knowledge is a nightmare: outdated Confluence, Slack decisions never documented, three wikis claiming source of truth. Crawling reliably is harder than the memory system itself.

6. **Dogfooding in a vacuum.** Current testing is on a codebase designed for agents. Real test is a messy 500-service enterprise monorepo. (Founder countered: has access to org of thousands of engineers + dev friends for beta testing.)

7. **Skill pack taxonomy assumes stable knowledge.** Orgs change. Stale memory served confidently is worse than no memory. Staleness detection is in "trial" bucket, should probably be baseline.

8. **STRONGEST: Platform players will add memory natively.** GitHub, Google, Anthropic are all building agent ecosystems. Governance depth is a feature, not a platform. Features get copied. What prevents Copilot from adding audit trails to its own memory and making AthenaMind redundant?

### Founder Response
Not trying to build a company or outrun platform players. Goal is open source contribution, shared learning, and practical application at current employer. The "product" framing keeps engineering disciplined. If major players solve the problem natively, that's a win — the problem gets solved.

---

## 5. Engineering Insights for Blog Content

### Three strong post topics identified:

**Post 1: "I built a process system where AI agents run the entire SDLC"**
- The cycle protocol: query state → launch stage → execute → handoff → observe → commit
- Observer pattern for agent accountability
- 9-second self-orientation demonstration
- 91 commits, 33 done stories, all agent-operated
- Key insight: stateless agents need explicit process artifacts the way humans need culture and memory

**Post 2: "My agents waste 25% of context on orientation — here's how I got it to near-zero"**
- O(1) orientation protocol vs O(n) with codebase size
- AGENTS.md → launcher → one story (constant cost regardless of repo size)
- Skill pack hierarchy: core / domain / repo / episode
- First agent pays exploration cost once, every agent after pays near-zero
- Measured pain point, not theory

**Post 3: "Good engineering practice works for agents too"**
- Tight story scoping = context window management
- Structured handoffs = agent coordination protocol
- Interface-driven design = reduced orientation cost for agents
- Governance = trust in shared memory
- The 150 markdown files aren't documentation — they're a queryable database for stateless workers
- The files are Pillar 2 (state); the skill file is Pillar 1 (procedural); the current session is Pillar 3 (working) — the architecture was built intuitively before being named

### Writing strategy:
- Show work, don't theorize. Link to real repo with real code and artifacts.
- Twitter/X for discovery (short threads with one concrete insight + link to longer writeup)
- Personal site for depth (full writeups)
- Org mentorship role as organic amplification
- No hype, no flashy Tailwind site, no VC stalking

---

## 6. Artifacts Created This Session

1. `product-research/product/PRODUCT_VISION_V2.md` — long-term product vision, explicitly marked out-of-scope for v0.1
2. `delivery-backlog/engineering/intake/STORY-20260223-split-main-go-into-packages.md`
3. `delivery-backlog/engineering/intake/STORY-20260223-add-ci-pipeline.md`
4. `delivery-backlog/engineering/intake/STORY-20260223-readme-v01-alignment.md`
5. `delivery-backlog/engineering/intake/STORY-20260223-agent-bootstrap-protocol.md`
6. `delivery-backlog/engineering/intake/STORY-20260223-episode-writeback.md`
7. `delivery-backlog/engineering/intake/STORY-20260223-launch-stage-memory-integration.md`
8. `delivery-backlog/engineering/intake/STORY-20260223-sqlite-index-store.md`
9. `delivery-backlog/engineering/intake/STORY-20260223-embedding-retrieval-via-ollama.md`
10. This briefing document.

All intake stories validated against metadata schema (9/9 checks passing per story).

---

## 7. Next Actions

1. Run `./tools/launch_stage.sh pm` to have PM refinement rank intake stories into active queue.
2. Sprint 3 stories (split, CI, README) have no dependencies and can be ranked immediately.
3. Sprint 4 stories depend on split completion; Sprint 5 stories depend on Sprint 4.
4. Begin writing blog post #1 while agents execute Sprint 3.
