<!--
CANONICAL PLAN FORMAT. This template is the single source of truth for the
structure of a plan.md. Copy it into specs/<module>/<feature>/plan.md and fill
every section, then DELETE all HTML comments.

WHAT plan.md IS:
- The WORK ORDER for ONE change to the feature (new feature, update, or bug fix).
- Disposable and point-in-time. It guides the implementer for THIS change only.
- Overwritten wholesale by the next change to the feature. git history keeps
  every prior work order (`git log -- specs/<module>/<feature>/plan.md`), so
  nothing is lost — there is no need to keep old plans in the tree, and NO
  per-change filenames (plan-<change>.md). One plan.md, rewritten each time.

LIFECYCLE:
- WORK ORDER (while building): all sections present; it guides implementation.
- FROZEN (after the change ships): the moment the change merges, plan.md becomes
  a read-only historical record of that change. Do NOT keep editing it. It sits
  untouched until the NEXT change rewrites it from scratch.

THE HYDRATE GATE (do not skip):
- Before a change merges, the DURABLE decisions this work order introduced MUST
  be promoted ("hydrated") into design.md. design.md is the living source of
  truth; plan.md is a snapshot. If a decision lives only in plan.md, it is lost
  to future readers who read design.md. Hydrate first, then freeze.

RULES:
- design.md is the authority for anything durable. plan.md never restates the
  full design — it references design.md and records only what THIS change does.
- name the APPROACH, not the code: shapes/signatures to guide the implementer,
  not full bodies or SQL. These are transient and the source owns them.
-->

# Work Order: <feature name> — <this change in a few words>

**Feature design:** `specs/<module>/<feature>/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/<module>/<feature>/spec.md`

> Work order for: **<this change>**. Disposable — overwritten by the next change
> (git keeps the history). The living design is in design.md; hydrate it before
> this change merges, then freeze this file.

## Change
<!-- What this change does and WHY, in a few sentences. For a new feature: the
     scope being built. For an update: what behavior changes and why. For a bug
     fix: the observed wrong behavior, the expected behavior, and the root cause.
     Link the spec scenarios this change satisfies. -->

## Architecture & Files (this change)
<!-- ONLY the files this change touches, annotated CREATE / MODIFY / REGEN, laid
     out by layer. Not the whole feature tree (that lives in design.md) — just
     the delta. -->
```
src/<module>/...
└── ...                                    # CREATE | MODIFY | REGEN

tests/...
└── ...                                    # CREATE | MODIFY
```

## Key Types & Signatures
<!-- Interfaces/type shapes/signatures that guide the implementer for THIS
     change: ports, entity constructors, command shapes, repository methods.
     Shapes, not bodies. Transient — the source owns these once written. -->

## Gaps / Bugs to Fix
<!-- Checklist. Each item: what is wrong / to build, where (file:line), and the
     correct behavior per the spec. Every spec scenario and every gap needs a
     test. For a bug fix, list the FAILING reproduction test(s) first — one per
     rejection/edge path the defect touches. -->
- [ ] ...

## Design decisions to hydrate into design.md
<!-- The pre-merge checklist for the HYDRATE GATE. List every durable decision
     THIS change introduced or altered that must be promoted into design.md
     (Design Decisions & Rationale, Data Flow, Request & Response, Known
     Limitations, Quality Pillars). Tick each once it is in design.md. Empty
     only if the change genuinely altered nothing durable. -->
- [ ] ...
