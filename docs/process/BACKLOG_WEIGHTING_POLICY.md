# Backlog Weighting Policy

Product-first ranking policy for PM refinement.

## Policy
1. Product and engineering stories are ranked above process-improvement stories by default.
2. Process-improvement stories are promoted ahead of product work only when one of the following is true:
   - a stage gate is failing or unenforceable,
   - a safety/quality contract is broken (for example CI, validation, or state-transition controls),
   - delivery is blocked without the process fix.
3. PM must record the reason whenever a process story is ranked above product work.

## Decision Rule
Use this order during PM queue ranking:
1. Blocked-by-broken-process product work
2. Product value/risk/dependency sequencing
3. Non-blocking process improvements

## Required Trace
When a process story is prioritized over product work, include in the story or queue notes:
- broken gate/contract reference,
- unblock condition,
- expected return to product-first order.
