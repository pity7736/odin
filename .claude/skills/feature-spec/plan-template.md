<!--
CANONICAL PLAN FORMAT. This template is the single source of truth for the
structure of a plan.md. Copy it into specs/<module>/<feature>/plan.md and fill
every section, then DELETE all HTML comments.

LIFECYCLE — a plan.md has TWO phases:
- WORK ORDER (while building): all sections present; it guides implementation.
- LIVING DESIGN DOC (after shipping): PRUNE it into a durable design + rationale
  doc that matches the shipped feature. Pruning = delete the TRANSIENT sections
  (Key Types & Signatures, Gaps / Bugs to Fix) and strip the CREATE/MODIFY
  annotations from the tree. Keep the DURABLE sections (rationale, architecture
  shape, data flow, contract, Quality Pillars). Never keep a prose copy of code
  the source already owns (signatures, field lists) — that only drifts.
Each section below is marked DURABLE or TRANSIENT.

RULES (from docs/03-sdd-workflow.md):
- Audience: engineers. This is the HOW; the spec is the WHAT.
- For EXISTING features: show the touched files (with real paths, annotated
  MODIFY) in the Architecture & Files Summary tree, and flag bugs, missing
  tests, and gaps as checklist items in Gaps / Bugs to Fix.
- INCLUDE: Design Decisions & Rationale, an Architecture & Files Summary tree,
  data flow, Quality Pillars, and — when the feature has an external interface —
  a Request & Response contract. During the work-order phase also include Key
  Types & Signatures and Gaps / Bugs to Fix (both pruned after shipping).
- REFERENCE SWAPPABLE DEPENDENCIES BY THEIR PORT, NOT THEIR ADAPTER: a feature
  plan depends on the repository INTERFACE (domain port), never the concrete
  implementation (e.g. an in-memory or Postgres adapter). The concrete adapter
  is owned by that adapter's own feature/plan and wired at the composition root,
  so infra swaps do not touch feature plans.
- DO NOT INCLUDE: full function bodies, complete business logic, detailed error
  handling beyond naming the error types, OR any code-level implementation such
  as SQL / query text / concrete queries — name the APPROACH, not the code. The
  code is the single source of truth for the code; duplicating it here only
  guarantees drift.
- The Quality Pillars section is MANDATORY and must address all four pillars.
  "Deferred" is allowed ONLY with a short justification (see docs/06-quality-pillars.md).
-->

# Technical Plan: <feature name>

**Corresponds to Spec:** `specs/<module>/<feature>/spec.md`

## Overview
<!-- DURABLE. What this feature is, technically. For existing features, what
     already exists vs what this plan changes (the "changes" framing is trimmed
     to the durable description once the feature ships). -->

## Design Decisions & Rationale
<!-- DURABLE — the heart of the living doc. The non-obvious choices and WHY, the
     thing the code does NOT capture. One bullet per decision: the choice + the
     reason + the alternative rejected. E.g. "Currency is derived from Money, not
     stored on the account — avoids a second source of truth; rejected a separate
     currency field." -->
- ...

## Architecture & Files Summary
<!-- DURABLE structure (TRANSIENT annotations). An annotated file-tree of the
     packages/files this feature touches, laid out by layer (domain / application
     / infrastructure / tests), annotated CREATE / MODIFY / REGEN. Doubles as the
     architecture view AND the file manifest — no separate flat table. On pruning,
     keep the tree structure, strip the CREATE/MODIFY/REGEN annotations. -->
```
src/<module>/domain/
└── ...                                    # CREATE

src/<module>/application/
└── ...                                    # MODIFY

src/<module>/infrastructure/
└── ...                                    # MODIFY

tests/...
└── ...                                    # CREATE

specs/<module>/<feature>/
├── spec.md                                # CREATE
└── plan.md                                # CREATE
```

## Data Flow
<!-- DURABLE. How a request moves through the layers for this operation. Keep it
     to the sequence of components and what each one produces. -->

## Request & Response
<!-- DURABLE. ONLY when the feature has an external interface (a client sends data
     and/or receives a response). For a purely internal change (domain refactor,
     shared value object, background job) write "N/A — no external interface" and
     skip. Show the CONTRACT/shape, illustrative — not implementation. Cover BOTH
     interfaces where they differ: REST = JSON body/response; HTMX = form fields
     in, rendered fragment or redirect out. -->

**Request data** (fields the client provides):

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| ... | ... | ... | ... |

**REST** — `<METHOD> /api/v1/<path>`
```json
// request
{ ... }
// success <status>
{ ... }
// error <status>
{ ... }
```

**HTMX** — `<METHOD> /<path>`
- Form fields: `...`
- Success:
  - Rendered fragment: `<template>` (e.g. the new list row)
  - Target & swap: `<target element>`, `<strategy>` — for a list, choose
    append (`beforeend`) vs prepend (`afterbegin`); DECIDE per feature.
  - Out-of-band updates: `<e.g. reset the form, bump a counter>` (`hx-swap-oob`)
  - Or, instead of a swap: `HX-Redirect` / `HX-Trigger: <event>`
- Error: renders `<error template>` into `<target>`

## Key Types & Signatures
<!-- TRANSIENT — pruned after shipping (the code owns these). Interfaces/type
     shapes/signatures that guide the implementer. Ports, entity constructors,
     command shapes, repository methods. Shapes, not bodies. -->

## Gaps / Bugs to Fix
<!-- TRANSIENT — pruned after shipping (all closed). Checklist. Each item: what is
     wrong, where (file:line), and the correct behavior per the spec. -->
- [ ] ...

## Quality Pillars
<!-- DURABLE. MANDATORY — one line minimum per pillar. "Deferred" needs a
     justification. -->
- **Security:** ...
- **Reliability:** ...
- **Performance:** ...
- **Observability:** ...
