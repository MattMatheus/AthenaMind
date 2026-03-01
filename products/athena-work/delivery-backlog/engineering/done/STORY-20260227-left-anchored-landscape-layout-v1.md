# Story: Shift workspace UI to left-anchored landscape layout with persistent nav

## Metadata
- `id`: STORY-20260227-left-anchored-landscape-layout-v1
- `owner_persona`: staff-personas/Software Engineer - Max.md
- `status`: done
- `idea_id`: direct
- `phase`: v0.2
- `adr_refs`: [ADR-0013, WSI-ADR-0002, WSI-ADR-0003]
- `success_metric`: On desktop landscape, first-screen orientation (where am I, what is next, where to navigate) is readable without centered-column scanning in <=10 seconds.
- `release_checkpoint`: required

## Problem Statement
- Current workspace composition feels centered like a social-feed layout; review target is landscape-first with clear left-anchored navigation and workspace content hierarchy.

## Scope
- In:
- Replace center-weighted shell with left-anchored layout for desktop/laptop breakpoints.
- Add persistent left navigation rail for primary workspace sections.
- Preserve mobile behavior with an accessible compact nav pattern.
- Out:
- New backend APIs or lane-state semantics.
- Major visual-theme redesign unrelated to layout/navigation structure.

## Acceptance Criteria
1. Desktop landscape view uses left-anchored information architecture (nav rail + workspace content), not a centered social-style column.
2. Primary navigation is persistently visible on desktop and clearly indicates active section/state.
3. Responsive behavior preserves usability on mobile/tablet without losing access to core sections.

## QA Checks
- Test coverage updated
- Tests pass
- No known regressions in touched scope

## Dependencies
- `products/athena-work/delivery-backlog/engineering/done/STORY-20260227-workspace-ui-read-only-board-v1.md`
- `products/athena-work/delivery-backlog/engineering/done/STORY-20260227-kanban-interaction-polish-and-accessibility-v1.md`

## Notes
- Reference direction from Matt: landscape workspace with left-anchored nav as primary orientation model.
