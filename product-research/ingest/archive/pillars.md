This three-pillar approach—combining procedural Markdown storage, Document DB memory, and isolated workspace pods for user-side code—creates a highly resilient and scalable "brain" for Project Athena.
Your idea to clone repositories directly into a pod aligns perfectly with Athena's architecture, specifically its built-in sandboxing capabilities and semantic tooling. Here is how we can seamlessly integrate this third pillar into Athena:
1. Cloned Repositories via Sandbox Routing
Athena already features a PolicyAwareExecutionBackend that supports routing agent runs into isolated environments. Instead of the agent generating code blindly, you can configure Athena to spin up a temporary Kubernetes pod (or Docker container) specifically for that run.
• The Execution Flow: When a task involves a user's codebase, Athena can clone the repository into this secure sandbox. The agent can then safely read, modify, and actually execute tests or scripts without risking the host system.
• Configuration: You can enforce this isolation using environment variables like ATHENA_RUNTIME_ISOLATION_DEFAULT_PROFILE and ATHENA_SANDBOX_ENABLED=true.
2. Codebase-Specific Memory & Preferences
To handle user-preference data specific to a given repository (like branch naming conventions, commit message formats, or architectural gotchas), you have two distinct locations to store this knowledge based on its scope:
• In-Repo Markdown (Local Scope): The most effective way to store repository-specific directives is to treat this context as part of the codebase itself. Just as tools like Claude Code use a CLAUDE.md or .claude/MEMORY.md file, you can instruct Athena to always look for a local .athena/preferences.md file immediately after cloning a repo. This ensures that any agent touching the code instantly adopts the project's specific commit and branching rules.
• Document DB (Global Scope): If the agent learns a developer's preference that applies across multiple repositories (e.g., "This developer always prefers rebase over merge"), the agent can write that to the Document DB under a USER.md or user profile schema.
3. Leveraging LSP for Codebase "Understanding"
Once the repository is cloned into the pod, Athena doesn't have to rely on token-heavy file dumps or basic grep searches to understand it. Athena includes built-in Language Server Protocol (LSP) and symbolic navigation tools.
• Because the codebase is local to the pod, the agent can use tools like athena_lsp_definition, athena_lsp_references, and athena_lsp_symbols to semantically navigate the code.
• This drastically reduces the number of turns the agent needs to take and limits context-window bloat, as the tools return only relevant code snippets and signatures rather than entire files.
The Complete 3-Pillar Strategy
By finalizing this architecture, Athena's memory and execution model becomes incredibly robust:
1. Procedural Memory (Markdown Blob Share): The agent's core instructions, standard operating procedures, and "Skills" are mounted to the orchestrator.
2. State & Semantic Memory (Document DB): The database tracks the orchestrator's DAG workflow states, episodic logs, and global user preferences to prevent context drift across sessions.
3. Working Environment (Isolated Pods): Target repositories are cloned into isolated Kubernetes pods where the agent can securely execute code, read local repository preferences, and perform deep semantic analysis via LSP.
