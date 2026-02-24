# AthenaWork Quickstart

## Summary

AthenaWork is the staged operating workflow (planning/engineering/qa/pm). In this slim repo, the full operator pack is archived; this page explains how to run it and where to find full docs.

## Where AthenaWork Lives

Archived operator pack:
- `/Users/foundry/Experiments/Archived/AthenaMind-internal-2026-02-24/products/athena-work`

## Core Stages

- `planning`: clarify/shape work.
- `engineering`: implement top-ranked story.
- `qa`: validate gates and regressions.
- `pm`: refine queues and update priorities.
- `cycle`: continuous engineering+qa loop.

## Typical Run (Operator Pack)

From the archived AthenaWork root:

```bash
./tools/launch_stage.sh engineering
./tools/launch_stage.sh qa
./tools/run_observer_cycle.sh --cycle-id <id>
```

## Why It Matters To AthenaMind Users

AthenaWork provides the disciplined workflow around AthenaMind changes: queueing, gate checks, QA evidence, and cycle closure.

## Full AthenaWork Guide

See [AthenaWork Product Guide](/Users/foundry/Experiments/Current/AthenaMind/knowledge-base/product/athenawork.md).
