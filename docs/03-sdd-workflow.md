# Specification & Design-Driven Development Workflow

Every feature is described by three documents, each with **one stable role**: a
**Spec** (the WHAT), a **Design** (the durable HOW and WHY), and a **Plan** (the
work order for the current change). This ensures we build the right thing (the
Spec), record why we built it the way we did (the Design), and drive each change
deliberately (the Plan).

- **Location:** `spec.md`, `design.md`, and `plan.md` live together in a
  feature-specific subfolder within `specs/`, organized by module. For example:
  `specs/accounting/accounts/create/`.
- **Process:**
    1. A `spec.md` is written and reviewed, then a `plan.md` work order is
       written and reviewed, before development begins.
    2. Once approved, development follows TDD.
    3. Before the change merges, `design.md` is hydrated from the plan and the
       plan is frozen.

The operational, step-by-step version of this workflow lives in the feature-spec
skill (`.claude/skills/feature-spec/SKILL.md`). This document is the narrative.

## Specification File (The WHAT)

- **Audience:** Product Managers, Stakeholders, Engineers.
- **Purpose:** To define the business requirements, user stories, and acceptance criteria for a feature in plain, non-technical language.
- **Specs describe intended behavior** — what the feature SHOULD do, not what the current code does. If current code has bugs, the spec is the source of truth for correct behavior.
- **Living:** `spec.md` is edited in place whenever the intended behavior changes.
- **Format:** The canonical `spec.md` structure lives in the feature-spec skill at `.claude/skills/feature-spec/spec-template.md`. Copy it and fill it in.

## Design File (The durable HOW and WHY)

- **Audience:** Engineers.
- **Purpose:** The durable, living record of how the feature works and — above all — *why* it is built that way: the reasoning the code cannot capture.
- **Living, and the single source of truth for anything durable.** When `design.md` and a `plan.md` disagree, `design.md` wins. It is kept accurate to the shipped code for the life of the feature: every change **hydrates** it (see below) before merging.
- **Source of truth boundary:** the **code** is the source of truth for *how it is built*, the **spec** for *what it does*; `design.md` holds the design **rationale** and shape — not a prose mirror of the code.
- **What NOT to include:** full function bodies, complete business logic, detailed error handling beyond naming error types, SQL / query text, or prose copies of code the source owns (exact signatures, field lists) — all of which only drift. Name the approach and the shape, not the code.
- **Format:** The canonical `design.md` structure lives in the feature-spec skill at `.claude/skills/feature-spec/design-template.md`. Copy it and fill it in.

## Plan File (The work order for one change)

- **Audience:** Engineers (and the implementer subagent).
- **Purpose:** To guide the implementer through **one change** — a new feature, an update, or a bug fix. It carries the change description, the delta architecture (touched files, annotated CREATE/MODIFY/REGEN), the type/signature shapes, the **Implementation Phases (TDD)** — the ordered work in dependency order, each phase with its Red (tests first) and Green (implementation) — and the list of durable decisions to hydrate into `design.md`.
- **Disposable and point-in-time.** There is one `plan.md` per feature, **overwritten wholesale** by the next change. git history keeps every prior work order (`git log -- specs/<module>/<feature>/plan.md`), so nothing is lost. **Do not create per-change files** (`plan-<change>.md`).
- **Frozen after ship.** The moment a change merges, `plan.md` becomes a read-only historical record of that change. It is not pruned and not edited further — it sits untouched until the next change rewrites it from scratch.
- **For existing features:** the plan shows only the touched files (real paths) in its delta tree and flags bugs, missing tests, and gaps as items to fix.
- **Format:** The canonical `plan.md` structure lives in the feature-spec skill at `.claude/skills/feature-spec/plan-template.md`. Copy it and fill it in.

## The Hydrate Gate

Because `plan.md` is disposable and `design.md` is the durable record, the one
piece of discipline the model depends on is **hydration**: before a change
merges, every durable decision the work order introduced MUST be promoted into
`design.md`. A decision that lives only in `plan.md` is invisible to anyone
reading `design.md` and is effectively lost once the next change overwrites the
plan.

At the gate, before merge:

1. **Hydrate `design.md`.** For a new feature, create it from the template,
   filled from the work order and the shipped code. For an update or bug fix,
   edit it in place so every durable decision, data-flow change, contract change,
   and new limitation is reflected. `design.md` must match the shipped code.
2. **Freeze `plan.md`.** Leave it exactly as it is — the historical work order
   for this change. Do not prune it, do not edit it further.

## Quality Pillars in Designs

Every `design.md` **must** include a Quality Pillars section that addresses all four pillars from [docs/06-quality-pillars.md](./06-quality-pillars.md):

1. **Security**
2. **Reliability**
3. **Performance**
4. **Observability**

Each pillar must have at least one line stating what applies to this feature. **"Deferred"** is a valid answer when there is no production infrastructure to support it yet, but it must include a short justification. This ensures the decision to skip is conscious, not accidental.

## Updating a Feature (and Fixing a Bug)

First confirm the change really is an update to *this* feature. A cross-cutting or infrastructural change (for example, swapping storage from in-memory to Postgres) is a **new feature** with its own `specs/` folder, not an update here. Because designs reference ports (interfaces), not concrete adapters, many infrastructure changes touch no feature design at all.

When an existing feature's own behavior or implementation genuinely changes:

1. **Edit the spec in place** (only if intended behavior changes). For a bug where the spec is already correct, the spec does not change — the code diverges from it.
2. **Read the existing `design.md`** — it is the living record of how the feature works today — then investigate the code and discuss the change.
3. **Write a fresh `plan.md` work order** for this change, overwriting the previous one (git keeps it). For a bug fix, its Implementation Phases open with a FAILING reproduction test (Phase 1) per path the defect touches.
4. **Implement, review, and manually test** following TDD.
5. **At the hydrate gate:** update `design.md` in place to reflect the change, then freeze `plan.md`.
