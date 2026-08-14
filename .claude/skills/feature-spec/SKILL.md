---
name: feature-spec
description: Spec-Driven Development workflow for Odin — covers starting a new feature, updating an existing one, AND fixing a bug in one. Use whenever creating or updating a feature's spec.md, design.md, or plan.md under specs/, or fixing a defect in a feature's behavior, before writing any feature code. Produces a business-focused spec (the WHAT), a per-change plan work order, and a durable technical design (the HOW/WHY) using the canonical templates. Triggers on requests to "new feature", "write a spec", "create a plan", "write a design", "update a feature", "fix a bug", "fix", or any work referencing specs/<module>/<feature>/.
---

# Feature Spec Workflow (SDD)

Every feature begins with **discovery** — questions to the user — because the
spec's content lives in the user's head, not in a template or the code. Only
after discovery come the **spec** (the WHAT) and the **plan** (the work order for
one change), each written and reviewed before any code. This skill operationalizes
that workflow. The human-facing narrative lives in `docs/03-sdd-workflow.md`; the
canonical file formats live in the three templates beside this file.

**Three documents per feature — each with ONE stable role:**
- **`spec.md`** — the WHAT, in business language. Living.
- **`design.md`** — the durable HOW and WHY. Living, and the **single source of
  truth** for anything durable. When design.md and a plan.md disagree, design.md
  wins.
- **`plan.md`** — the WORK ORDER for the *current* change (new feature, update,
  or bug fix). Disposable: **overwritten wholesale** by the next change; git
  history keeps every prior work order. It is never pruned or morphed — it is
  simply rewritten each time.

**Two rules keep this from degrading (do not skip either):**
- **Hydrate before merge.** Before a change merges, the durable decisions the
  work order introduced MUST be promoted into `design.md`. A decision that lives
  only in `plan.md` is invisible to anyone reading `design.md`.
- **Freeze after ship.** The moment a change merges, `plan.md` becomes a
  read-only historical record. You do not keep editing it — it sits untouched
  until the next change rewrites it from scratch. That is what stops it drifting
  back into a shape-shifting doc.

## Discovery (ALWAYS do this first)

Do not write, copy, or fill any file until discovery is complete. Ask **one
question at a time** and wait for the answer before the next. **Never assume** —
if something is unclear, ask.

1. **Q1 (branch):** Is this a **new feature**, an **update to an existing one**,
   or a **bug fix** in an existing one? The answer decides the path (see below).
2. **Q2:** What is the feature about, in your own words?
3. Then keep asking follow-ups until ALL of these are covered:
   - purpose / the benefit to the user
   - who uses it
   - the business rules
   - rejection & edge cases
   - what's out of scope
4. **For an update, also cover:** what exactly is changing, and why.
5. **For a bug fix, also cover:** what the user observes (the wrong behavior),
   what they expected instead, and how to reproduce it. Then decide whether the
   spec already describes the correct behavior (code diverges from spec — most
   bugs) or the spec itself is wrong/silent (a behavior problem). That decision
   picks the path below.

Only when the picture is complete do you proceed to write the spec.

## Path after discovery

- **New feature** → follow "Workflow" below.
- **Update to an existing feature** → follow "Updating an existing feature".
- **Bug fix** → follow "Fixing a bug in an existing feature".

## Files this skill owns

- `spec-template.md` — canonical structure for every `spec.md`
- `design-template.md` — canonical structure for every `design.md`
- `plan-template.md` — canonical structure for every `plan.md`

These templates are the single source of truth for spec/design/plan structure.
Do not reconstruct the format from memory — always start from the template.

## Location

All three files go in a feature-specific folder, organized by module:
`specs/<module>/<feature>/` — `spec.md`, `design.md`, `plan.md`. For atomic
per-operation features, add an operation subfolder:
`specs/accounting/accounts/create/`.

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

4. **Write the plan (work order)** from the agreed discussion. Copy
   `plan-template.md` to the feature folder as `plan.md`, **overwriting** any
   previous work order (git keeps it). Fill the Change section (what/why + the
   spec scenarios it satisfies), the delta architecture tree annotated
   CREATE/MODIFY/REGEN (only the files THIS change touches), Key Types &
   Signatures, Gaps / Bugs to Fix as checklist items, and the "Design decisions
   to hydrate into design.md" checklist (every durable decision this change
   introduces — this is the pre-merge hydrate list). The plan is the work order;
   the durable design lives in `design.md` and is written/updated at the hydrate
   gate (step 11), not now.

5. **STOP for review.** The user approves the written plan before implementation.

6. **Implement** in a FRESH session or subagent that works only from `spec.md`,
   the current `plan.md` work order, and (for an update or bug fix) the feature's
   `design.md` — not from the design conversation. If it cannot build from those
   files alone, the plan was incomplete; that is useful signal, so stop
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
   - **Plan conformance:** the implementation built what the `plan.md` work order
     describes, and NO tests or functionality were removed unless the plan called
     for it.
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

9. **Update the Bruno collection** so the user can exercise the API without
   hand-building requests. REST only: web/HTMX is tested in the browser, not
   Bruno. **One request per endpoint is NOT enough** — a lone happy-path request
   only proves the happy path and is close to useless for finding defects. Add a
   SEPARATE, named request for every meaningful scenario the endpoint has: the
   success case AND each rejection and edge case (wrong credentials, missing /
   empty / malformed fields, unauthorized access, not-found, conflicts, etc.).
   Mirror the spec's Expected Behavior scenarios one-to-one — if a scenario has
   no request, it will not get tested. Create/update the
   `.bruno/<area>/<request>.yml` files and keep shared values (`host`, auth
   token, secrets) in the environment file, not inline in requests.

10. **Manual test** by the user, who runs the application themselves — the API
   with the Bruno requests from step 9, the web in the browser. The agent
   waits. What the user finds routes by KIND:
   - **A bug** (the code does not match the spec) → fix in the code, re-run
     `make check` GREEN.
   - **A behavior problem** (the spec itself specifies the wrong thing) → this is
     NOT a code patch. Loop back: update `spec.md` → re-review the spec → adjust
     the plan → implement. Never silently patch code to a behavior the spec does
     not describe; that makes the spec lie.
   GATE on the user's approval.

11. **Hydrate `design.md`, then freeze `plan.md`** once everything above passes,
   BEFORE the change merges. This is the hydrate gate — the one bit of discipline
   the disposable-plan model depends on.
   - **Hydrate:** work through the plan's "Design decisions to hydrate into
     design.md" checklist. For a NEW feature, create `design.md` from
     `design-template.md`, filled from the work order and the shipped code. For an
     update or bug fix, edit the existing `design.md` in place so every durable
     decision, data-flow change, contract change, and new limitation this change
     introduced is reflected there. `design.md` must match the shipped code when
     you are done. Do NOT copy code the source owns (signatures, field lists) —
     capture the rationale, the thing the code cannot tell a future reader.
   - **Freeze:** leave `plan.md` exactly as it is — it is now the historical work
     order for this change. Do NOT prune it, do NOT edit it further. The next
     change to this feature overwrites it from scratch; git keeps this one.
   If a decision lives only in `plan.md` after this step, it is effectively lost —
   that is the failure this gate exists to prevent.

## Updating an existing feature

First confirm it really IS an update to THIS feature — a change to this
feature's own behavior or implementation. If the change is cross-cutting or
infrastructural (e.g. swapping storage from in-memory to Postgres), it is a NEW
feature with its own `specs/` folder, not an update here. And because plans
reference ports, not adapters, many infrastructure changes touch no feature
plan at all — check before assuming this path applies.

Discovery (above) must be complete first — including what is changing and why.

There is one `spec.md`, one `design.md`, and one `plan.md` per feature. `spec.md`
and `design.md` are living; `plan.md` is overwritten per change. Do NOT create
per-change files (`plan-<change>.md`) — git history holds prior work orders.

- Edit `spec.md` in place to reflect the new intended behavior. STOP for review
  before the plan.
- Run the same **Plan investigation & discussion** step as the Workflow — and
  read the existing `design.md` first, since it is the living record of how the
  feature works today. Surface findings, discuss, and pass the "are you good with
  the discussion?" gate before writing anything.
- Write a fresh `plan.md` work order for this change from `plan-template.md`,
  **overwriting** the previous work order (git keeps it). Fill the Change,
  delta architecture tree, Key Types & Signatures, Gaps / Bugs to Fix, and the
  "Design decisions to hydrate into design.md" checklist for the delta.
- STOP for review before implementation.
- Then implement, verify, review, manually review, update the Bruno collection,
  manually test, and — at the hydrate gate — update `design.md` in place and
  freeze `plan.md`, exactly as Workflow steps 6–11 (fresh implementer; `make
  check` gate; separate reviewer; your manual code review; Bruno REST requests;
  your manual test; hydrate `design.md`, freeze `plan.md`).

## Fixing a bug in an existing feature

A bug is code that does not do what the spec says. It is a close cousin of an
update, with one decisive difference: the intended behavior usually already
lives in `spec.md`, so there is nothing to change there. Do NOT invent a new
spec or rewrite the old one to match the buggy code — that makes the spec lie.

First confirm it really is a bug in THIS feature (the code diverges from this
feature's own spec), not a missing capability. A missing capability is an
update or a new feature, not a bug fix.

**MANDATORY: failing tests that reproduce the bug come first.** Before any fix
is written, there MUST be at least one test that exercises the broken behavior
and FAILS for the reason the user reported (Red). This is non-negotiable — it
proves the bug is real, pins down the exact defect, and guards against
regression. The fix is complete only when those tests pass (Green) with no other
test broken. A bug fix without a test that failed before it is not done.

**One reproduction test is rarely enough — cover EVERY scenario the defect
touches.** A single root cause usually manifests across many cases; each is a
distinct test. Enumerate them before writing the fix and do not stop at the
first. For example, "rejections return an empty response body" is not one case —
it is every rejection path: empty name, duplicate name+currency, invalid type,
invalid currency, not-found, unauthorized, and so on. Missing one scenario
leaves that path unguarded and free to regress. Add each reproduction test to
the `plan.md` work order's Gaps / Bugs to Fix as its own checklist item.

Put these tests at the RIGHT LEVEL — the layer where the defect actually lives
(see docs/04-tdd-workflow.md). A bug in domain or application logic is a unit
test against that unit; a bug in a handler or in shared HTTP wiring (error
mapping, response serialization) is a unit test of that component. Reach for an
integration test only when the defect is genuinely in the wiring BETWEEN
components and cannot be shown at the unit level. Do not default to an
integration test because it is the easiest place to hit the endpoint.

Discovery (above) must be complete first — including the observed wrong
behavior, the expected behavior, and how to reproduce it.

Then split by KIND (the same split as Workflow step 10):

- **Code diverges from the spec** (the spec is already correct — most bugs) →
  the spec does NOT change. Confirm `spec.md` already covers the correct
  behavior and its Expected Behavior has a scenario for the case that broke; if
  that scenario is missing, add it (this is a spec gap, review it as in Workflow
  step 2). Otherwise go straight to **Plan investigation & discussion**: read the
  feature's `design.md` and code, find the root cause, and discuss it. Write a
  fresh `plan.md` work order (overwriting the previous one) whose Gaps / Bugs to
  Fix names the actual defect and root cause, and whose "Design decisions to
  hydrate into design.md" checklist captures any decision the fix changes. Then
  implement, verify, review, manually review, update Bruno, manually test, and —
  at the hydrate gate — update `design.md` and freeze `plan.md`, exactly as
  Workflow steps 6–11. The implementer works test-first: a FAILING test that
  reproduces the bug comes before the fix (Red), then the fix makes it green.

- **The spec itself is wrong or silent** (a behavior problem — the code does
  what the spec says, but the spec specifies the wrong thing) → this is not a
  pure bug fix. Follow "Updating an existing feature": edit `spec.md` to the
  correct intended behavior, STOP for review, then the rest of the flow.

## Checklist before handing a spec for review

- No technical terms anywhere in the spec.
- User Stories present — one or more, each As a… I want… so that…
- Expected Behavior covers the happy path AND every rejection/edge case.
- Out of Scope section present.

## Checklist before handing a plan (work order) for review

- Links to `design.md` and `spec.md` present.
- Change section describes what this change does and why, and the spec scenarios
  it satisfies.
- Delta architecture tree present — ONLY the touched files, laid out by layer,
  annotated CREATE / MODIFY / REGEN.
- Gaps / Bugs to Fix listed as actionable checklist items; every spec scenario
  and every gap has a test (for a bug, a FAILING reproduction test per path).
- "Design decisions to hydrate into design.md" checklist present — every durable
  decision this change introduces (empty only if nothing durable changed).

## Checklist before merge (the hydrate gate)

- `design.md` exists (created for a new feature) and matches the shipped code.
- Every item on the plan's "Design decisions to hydrate into design.md" checklist
  is now reflected in `design.md`: Design Decisions & Rationale, Data Flow,
  Request & Response, Known Limitations, Quality Pillars (all four).
- `design.md` contains no code the source owns (signatures, field lists) — only
  the durable shape and the rationale.
- `plan.md` is left frozen as the historical work order — not pruned, not edited
  further.
