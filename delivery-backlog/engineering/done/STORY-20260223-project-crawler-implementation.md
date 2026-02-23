# Story: Project Crawler Implementation

## Metadata
- `id`: STORY-20260223-project-crawler-implementation
- `owner_persona`: Software Architect - Ada.md
- `status`: active
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0005, ADR-0007]
- `success_metric`: Single command to index all markdown files in a codebase
- `release_checkpoint`: required

## Problem Statement
Onboarding new agents to a codebase is manual. We need a way to automatically discover and index core knowledge files.

## Scope
- In: `crawl` command in `memory-cli`, recursive file discovery, batch embedding integration.
- Out: Non-knowledge files (binaries, large datasets), external URL crawling.

## Acceptance Criteria
1. `memory-cli crawl --dir <path>` recursively finds `.md` files.
2. Discovered files are batched for embedding using the new batch logic.
3. Successfully indexes a target directory into the memory root.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- Azure OpenAI Batch Embedding Support (Completed)

## Notes
- Integrated with the newly added Azure batching logic for efficiency.
