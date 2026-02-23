Engineering Exploration: Adaptive Agent Memory & Cognitive Architectures

Subject: Transitioning from Context-Stuffed LLMs to Structured External Knowledge Graphs

Framework: Synthesis of Structured Knowledge (Ars Contexta) & Meta-Context Engineering (Koylan)

Date: February 2026
1. Executive Summary

Current LLM implementations suffer from Context Rot—the degradation of reasoning quality as the context window fills with non-essential tokens. This document explores a unified strategy that replaces linear chat history with a Skill Graph: a persistent, traversable, and self-healing markdown-based memory system optimized through Meta-Context Engineering (MCE).
2. Core Architectural Primitives

The system relies on 15 "Kernel Primitives" categorized into three operational phases:
A. The Ingestion Phase (Memory Write)

    Atomic Extraction: Identifying discrete, reusable facts from raw dialogue rather than saving transcripts.

    Salience Filtering: A heuristic-driven process to discard "conversational fluff" and retain only high-utility data.

    Consolidation: Comparing new information against existing notes to merge, update, or resolve contradictions.

B. The Retrieval Phase (Memory Read)

    Topology over Vector Similarity: Moving beyond simple RAG. Instead of using cosine similarity, the agent uses Wiki-links to traverse related concepts.

    Maps of Content (MOCs): Indexing files that act as "hubs," allowing the agent to enter a broad domain and narrow down to specific technical notes through explicit links.

C. The Maintenance Phase (Self-Healing)

    Link Integrity: Automated checks for broken references within the markdown vault.

    Pruning: Periodically removing outdated instructions or deprecated code snippets to maintain a high signal-to-noise ratio.

3. Integrated Strategy: The Self-Organizing Memory OS

The following table outlines how these divergent research paths merge into a single, unified system architecture.
Layer	Component	Integrated Function
Storage (L1)	Temporal Skill Graph	A markdown vault where every node (fact/skill) has a timestamp and a "Salience Score" that decays over time.
Logic (L2)	Kernel Hooks	Automated scripts that run whenever the agent writes to memory, performing Salience Filtering and Link Validation.
Meta (L3)	Skill Evolution Loop	A Meta-Agent that reviews performance logs to perform Agentic Crossover, rewriting the "Kernel Hooks" for better efficiency.
Hardware (L4)	KV-Cache Stabilization	Strict prompt formatting that keeps stable system instructions at the top, ensuring fast token processing and lower latency.
4. Key Technical Techniques
Meta Context Engineering (MCE)

MCE uses two distinct layers to evolve memory rather than just managing it:

    The Meta-Agent (Skill Evolution): Analyzes the history of past tasks, successes, and failures to synthesize new instructions and tools (skills).

    The Base-Agent (Context Optimization): Executes the evolved skills. It treats context as "flexible code" and adjusts its own file structures based on what the Meta-Agent has learned.

Progressive Disclosure (L1 → L2 → L3)

To maximize the LLM's attention budget, context is loaded in layers:

    L1: The agent only loads "Skill Names" and descriptions at startup.

    L2: It loads the full "Skill Logic" (Markdown) only when triggered.

    L3: It loads specific "Data Files" only when the logic requires them.
    Result: Reduces token usage by up to 87% compared to traditional "stuff-it-all-in" methods.

The "Cognitive Eraser" Pattern

When an agent reads a long file, it is instructed to:

    Extract core conclusions.

    Erase the original text from its context.

    Replace it with a Reference ID. This keeps the context window lean while maintaining "traceability."

5. Implementation Directives

    Develop the Bipartite Graph: Structure memory where nodes are either Entities (Facts) or Procedures (Skills). Use Wiki-links to bridge the two.

    Execute "Agentic Crossover" Sprints: Every 100 tasks, trigger a Meta-Agent routine to audit the "Skill" directory. If two skills overlap, consolidate them.

    KV-Cache Optimization: Ensure prefix stability by placing stable data (System Prompts/Tool Definitions) at the beginning of the prompt to ensure 10x speed improvement in high-frequency loops.

6. Sources & Research References

    Ars Contexta: Agent Memory Article Series (X)

    Muratcan Koylan: Meta-Context Engineering & Agentic Skill Evolution (X)

    Video: Deep Dive into Meta-Context Engineering