---
name: feature-spec
description: Spec-Driven Development workflow for Odin — covers BOTH starting a new feature and updating an existing one. Use whenever creating or updating a feature's spec.md or plan.md under specs/, before writing any feature code. Produces a business-focused spec (the WHAT) and a technical plan (the HOW) using the canonical templates. Triggers on requests to "new feature", "write a spec", "create a plan", "update a feature", or any work referencing specs/<module>/<feature>/.
---

# Feature Spec Workflow (SDD)

Every feature begins with **discovery** — questions to the user — because the
spec's content lives in the user's head, not in a template or the code. Only
after discovery come the **spec** (the WHAT) and the **plan** (the HOW), each
written and reviewed before any code. This skill operationalizes that workflow.
The human-facing narrative lives in `docs/03-sdd-workflow.md`; the canonical file
formats live in the two templates beside this file.

## Discovery (ALWAYS do this first)

Do not write, copy, or fill any file until discovery is complete. Ask **one
question at a time** and wait for the answer before the next. **Never assume** —
if something is unclear, ask.

1. **Q1 (branch):** Is this a **new feature** or an **update to an existing
   one**? The answer decides the path (see below).
2. **Q2:** What is the feature about, in your own words?
3. Then keep asking follow-ups until ALL of these are covered:
   - purpose / the benefit to the user
   - who uses it
   - the business rules
   - rejection & edge cases
   - what's out of scope
4. **For an update, also cover:** what exactly is changing, and why.

Only when the picture is complete do you proceed to write the spec.

## Path after discovery

- **New feature** → follow "Workflow" below.
- **Update to an existing feature** → follow "Updating an existing feature".

## Files this skill owns

- `spec-template.md` — canonical structure for every `spec.md`
- `plan-template.md` — canonical structure for every `plan.md`

These templates are the single source of truth for spec/plan structure. Do not
reconstruct the format from memory — always start from the template.

## Location

Both files go in a feature-specific folder, organized by module:
`specs/<module>/<feature>/`. For atomic per-operation features, add an operation
subfolder: `specs/accounting/accounts/create/`.

## Workflow

Discovery (above) must be complete first.

1. **Write the spec** from the discovery answers. Copy `spec-template.md` to the feature folder as
   `spec.md`. Fill every section, delete the HTML-comment guidance. The spec is
   business language only — no technical terms (Fiber, HTMX, Postgres, HTTP,
   REST, API, endpoint, handler, repository, database). Describe the intended
   correct behavior, not what buggy code currently does.

2. **STOP for review.** The user reviews and approves the spec before any plan
   is written. Do not start the plan until they approve.

3. **Plan investigation & discussion** — do this BEFORE writing any plan file.
   The plan's input comes from the code and architecture, not the user's head,
   so investigate first, then align:
   a. **Read the relevant code.** For an existing feature, the feature's code.
      For a new one, `docs/02-architecture.md`, `docs/05-code-standards.md`, and
      the closest existing feature to mirror its patterns.
   b. **Surface findings.** What already exists, bugs, gaps, and architectural
      concerns. Raise concerns and challenge decisions — do not stay quiet.
   c. **Discuss** with the user until aligned. Be STRICT, not agreeable. Do not
      default to agreement — if the user is wrong, say so plainly and explain
      why; if the user is right, say why with real technical arguments, not
      praise. Every position (yours or the user's) must be backed by an
      argument. Sycophancy here produces bad plans.
   d. **GATE:** explicitly ask "are you good with the discussion?" Do NOT write
      the plan until the user says yes. This is a separate approval from the
      plan review in step 5.

4. **Write the plan** from the agreed discussion — at this stage it is a WORK
   ORDER (it gets pruned to a living design doc after shipping, step 8). Copy
   `plan-template.md` to the feature folder as `plan.md`. Capture the Design
   Decisions & Rationale (the WHY), show touched files by real path in the
   architecture tree, list bugs/gaps as checklist items, and complete the
   mandatory Quality Pillars section (all four; "Deferred" only with
   justification).

5. **STOP for review.** The user approves the written plan before implementation.

6. **Implement** in a FRESH session or subagent that works only from `spec.md`
   and `plan.md` — not from the design conversation. If it cannot build from the
   spec and plan alone, the plan was incomplete; that is useful signal, so stop
   and fix the plan rather than filling gaps from memory. Follow TDD
   (Red-Green-Refactor): implement test-first in dependency order
   (domain → application → infrastructure); every spec scenario and every gap in
   the plan MUST have a test; regenerate mocks after the repository interfaces
   are final. End with `make check` (lint + test + coverage) GREEN — never hand
   red code to a reviewer.
   **STOP-ON-DEVIATION:** the plan was agreed together, so any departure from it
   is decided together. The moment reality diverges from the plan — a test
   fails, the code is not shaped as the plan assumed, an approach will not work,
   anything unexpected — STOP. Do not improvise, and above all do not resolve it
   by deleting tests, dropping functionality, or changing an agreed design on
   your own. Explain what you found and why it does not fit the plan, then
   discuss the next move and get agreement before acting. This applies to ANY
   deviation, not only destructive ones.

7. **Review** in a SEPARATE fresh reviewer session/subagent — not the one that
   implemented. The reviewer reads the real git diff and files (not a summary)
   and checks:
   - **Correctness/quality:** run `/code-review`.
   - **Odin standards:** conformance to `docs/05-code-standards.md` and
     `CLAUDE.md` (self receiver, no source comments, 100% coverage on business
     logic, descriptive names, internal errors English / external Spanish,
     `strings.Clone` on Fiber body data, constructor-after-struct ordering).
   - **Plan conformance:** the implementation built what the plan describes, and
     NO tests or functionality were removed unless the plan called for it.
   Report findings, discuss, and fix. After fixes, re-run `make check` GREEN.
   (Ad-hoc subagents for now; revisit dedicated implementer/reviewer agents with
   fixed model+effort later if feedback warrants.)

   NOTE: `make check` is a GATE, not a one-time step — it must be GREEN every
   time code changes (end of implementation, after review fixes, after manual-
   review fixes). Never proceed past a red check.

8. **Manual code review** by the user, back in the main session (NOT the reviewer
   subagent — do not simulate this). This is a DISCUSSION, exactly like Plan
   investigation & discussion: the agent WAITS for the user's feedback and
   questions and ANSWERS them — every question gets a real, direct answer, no
   deflecting, no exceptions. Do NOT change any code during the discussion;
   making a change requires the user's explicit permission — we are discussing
   now, changes come later. Once the discussion settles and the user approves the
   changes, apply them, then re-run `make check` GREEN. When all doubts are
   resolved and agreed changes applied, the agent asks "is everything OK?" and
   waits for the user's approval. GATE.

9. **Manual test** by the user, who runs the application themselves. The agent
   waits. What the user finds routes by KIND:
   - **A bug** (the code does not match the spec) → fix in the code, re-run
     `make check` GREEN.
   - **A behavior problem** (the spec itself specifies the wrong thing) → this is
     NOT a code patch. Loop back: update `spec.md` → re-review the spec → adjust
     the plan → implement. Never silently patch code to a behavior the spec does
     not describe; that makes the spec lie.
   GATE on the user's approval.

10. **Prune the plan into a living design doc** once everything above passes and
   the feature ships. `plan.md` stops being a work order and becomes the durable record of
   HOW and WHY. Delete the TRANSIENT sections (Key Types & Signatures, Gaps /
   Bugs to Fix) and strip the CREATE/MODIFY/REGEN annotations from the
   architecture tree. Keep the DURABLE sections (Overview, Design Decisions &
   Rationale, architecture shape, Data Flow, Request & Response, Quality
   Pillars). Do NOT keep a prose copy of code the source owns (signatures, field
   lists) — that only drifts. The rationale is the point: it is what the code
   cannot tell a future reader.

## Updating an existing feature

First confirm it really IS an update to THIS feature — a change to this
feature's own behavior or implementation. If the change is cross-cutting or
infrastructural (e.g. swapping storage from in-memory to Postgres), it is a NEW
feature with its own `specs/` folder, not an update here. And because plans
reference ports, not adapters, many infrastructure changes touch no feature
plan at all — check before assuming this path applies.

Discovery (above) must be complete first — including what is changing and why.

There is one `spec.md` and one `plan.md` per feature; git preserves prior
versions. Do NOT create per-change files (`plan-<change>.md`). The `spec.md` is
always living. The `plan.md` at rest is a living design doc; an update re-opens
it into a work order, then it is pruned back.

- Edit `spec.md` in place to reflect the new intended behavior. STOP for review
  before the plan.
- Run the same **Plan investigation & discussion** step as the Workflow (read
  the existing code, surface findings, discuss, and pass the "are you good with
  the discussion?" gate) before writing anything.
- Re-open `plan.md` into a WORK ORDER for this change: add back the transient
  Key Types & Signatures and Gaps / Bugs to Fix for the delta, and update the
  Design Decisions & Rationale.
- STOP for review before implementation.
- Then implement, verify, review, manually review/test, and PRUNE exactly as
  Workflow steps 6–10 (fresh implementer; `make check` gate; separate reviewer;
  your manual code review and manual test; prune back to a living design doc).

## Checklist before handing a spec for review

- No technical terms anywhere in the spec.
- User Story present (As a… I want… so that…).
- Expected Behavior covers the happy path AND every rejection/edge case.
- Out of Scope section present.

## Checklist before handing a plan for review

- Corresponds-to-Spec link present.
- Design Decisions & Rationale present (the WHY behind non-obvious choices, not a
  copy of the code).
- Gaps/bugs listed as actionable checklist items.
- All four Quality Pillars addressed; any "Deferred" is justified.
- Architecture & Files Summary tree present, laid out by layer, annotated
  CREATE / MODIFY / REGEN.
- Request & Response contract present (REST + HTMX), or explicitly "N/A — no
  external interface". For HTMX, the swap is specified: target, strategy
  (append/prepend decided), any out-of-band updates, or redirect/trigger.
