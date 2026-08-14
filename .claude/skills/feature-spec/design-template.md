<!--
CANONICAL DESIGN FORMAT. This template is the single source of truth for the
structure of a design.md. Copy it into specs/<module>/<feature>/design.md and
fill every section, then DELETE all HTML comments.

WHAT design.md IS:
- The DURABLE, living technical record of a feature: the HOW and the WHY.
- The single source of truth for anything durable about the feature. When a
  plan.md work order and design.md disagree, design.md wins.
- Living: it is kept accurate to the shipped code for the life of the feature.
  Every change to the feature HYDRATES design.md — the durable decisions the
  change introduced are promoted here before the change merges.

WHAT design.md IS NOT:
- Not a work order. The per-change, disposable "what I'm about to do" lives in
  plan.md (see plan-template.md). design.md never holds Implementation Phases,
  CREATE/MODIFY annotations, or literal code signatures the source owns.

RULES (from docs/03-sdd-workflow.md):
- Audience: engineers. This is the HOW; the spec is the WHAT.
- REFERENCE SWAPPABLE DEPENDENCIES BY THEIR PORT, NOT THEIR ADAPTER: depend on
  the repository INTERFACE (domain port), never a concrete adapter (in-memory /
  Postgres). Adapters are owned by their own feature and wired at the
  composition root, so infra swaps do not touch a feature's design.
- DO NOT duplicate code the source owns: no full function bodies, no field-by-
  field type dumps, no SQL/query text. Name the APPROACH and the shape, not the
  code — duplicated code only drifts. The rationale is the point: capture what
  the code CANNOT tell a future reader.
- The Quality Pillars section is MANDATORY and must address all four pillars.
  "Deferred" is allowed ONLY with a short justification (see docs/06-quality-pillars.md).
-->

# Technical Design: <feature name>

**Corresponds to Spec:** `specs/<module>/<feature>/spec.md`

## Overview
<!-- What this feature is, technically, in a few sentences. The durable
     description of the shipped design — not a changelog. -->

## Design Decisions & Rationale
<!-- The heart of the doc. The non-obvious choices and WHY — the thing the code
     does NOT capture. One bullet per decision: the choice + the reason + the
     alternative rejected. E.g. "Currency is derived from Money, not stored on
     the account — avoids a second source of truth; rejected a separate currency
     field." Every change that alters a decision updates the relevant bullet
     here. -->
- ...

## Architecture & Files Summary
<!-- The packages/files this feature owns, laid out by layer (domain /
     application / infrastructure / tests). Structure only — NO CREATE/MODIFY
     annotations (those are per-change and live in plan.md). -->
```
src/<module>/domain/
└── ...

src/<module>/application/
└── ...

src/<module>/infrastructure/
└── ...

tests/...
└── ...

specs/<module>/<feature>/
├── spec.md
├── design.md
└── plan.md          # current work order (see plan-template.md)
```

## Data Flow
<!-- How a request moves through the layers for this operation. The sequence of
     components and what each one produces. -->

## Request & Response
<!-- ONLY when the feature has an external interface. For a purely internal
     change write "N/A — no external interface" and skip. Show the CONTRACT/shape,
     illustrative — not implementation. Cover BOTH interfaces where they differ:
     REST = JSON body/response; HTMX = form fields in, rendered fragment or
     redirect out (name the target, swap strategy, any out-of-band updates). -->

## Known Limitations
<!-- Durable caveats a future engineer must know: non-atomic checks, deferred
     concerns, sharp edges left in place on purpose. -->

## Quality Pillars
<!-- MANDATORY — one line minimum per pillar. "Deferred" needs a justification. -->
- **Security:** ...
- **Reliability:** ...
- **Performance:** ...
- **Observability:** ...
