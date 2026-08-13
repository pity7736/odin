# Specification & Plan-Driven Development Workflow

Every new feature begins with a **Specification** and a **Plan**. This ensures we build the right thing (the Spec) and build the thing right (the Plan).

- **Location:** Both `spec.md` and `plan.md` will be located in a feature-specific subfolder within `specs/`, organized by module. For example: `specs/accounting/accounts/`.
- **Process:**
    1. A `spec.md` and a `plan.md` are written and reviewed before development begins.
    2. Once approved, development can begin following TDD.

## Specification File (The WHAT)

- **Audience:** Product Managers, Stakeholders, Engineers.
- **Purpose:** To define the business requirements, user stories, and acceptance criteria for a feature in plain, non-technical language.
- **Specs describe intended behavior** — what the feature SHOULD do, not what the current code does. If current code has bugs, the spec is the source of truth for correct behavior.
- **Format:** The canonical `spec.md` structure lives in the feature-spec skill at `.claude/skills/feature-spec/spec-template.md`. Copy it and fill it in.

## Plan File (The HOW)

- **Audience:** Engineers.
- **Purpose:** To describe the high-level technical implementation plan for a corresponding specification.
- **Lifecycle — two phases:** A `plan.md` is a **work order** while the feature is being built (it carries the gaps to fix, the type/signature shapes, the TDD guidance). Once the feature ships it is **pruned into a living design doc** — the durable record of the *how* and, above all, the *why*. The transient work-order sections are deleted; the design decisions, architecture, data flow, contract, and Quality Pillars remain.
- **Source of truth:** After a feature ships, the **code** is the source of truth for *how it is built* and the **spec** for *what it does*. The pruned plan is the durable record of the design **rationale** — the reasoning the code cannot capture — not a prose mirror of the code.
- **For existing features:** The plan shows the touched files (real paths) in its architecture tree and flags bugs, missing tests, and gaps as items to fix.
- **What NOT to include:**
  - Full implementation of functions
  - Complete method bodies with business logic
  - Detailed error handling beyond error types
  - SQL / query text / any code-level implementation — name the approach, not the code
  - In the pruned living doc: prose copies of code the source owns (exact signatures, field lists), which only drift
- **Format:** The canonical `plan.md` structure lives in the feature-spec skill at `.claude/skills/feature-spec/plan-template.md`. Copy it and fill it in.

## Quality Pillars in Plans

Every plan **must** include a Quality Pillars section that addresses all four pillars from [docs/06-quality-pillars.md](./06-quality-pillars.md):

1. **Security**
2. **Reliability**
3. **Performance**
4. **Observability**

Each pillar must have at least one line stating what applies to this feature. **"Deferred"** is a valid answer when there is no production infrastructure to support it yet, but it must include a short justification. This ensures the decision to skip is conscious, not accidental.

## Updating a Feature

First confirm the change really is an update to *this* feature. A cross-cutting or infrastructural change (for example, swapping storage from in-memory to Postgres) is a **new feature** with its own `specs/` folder, not an update here. Because plans reference ports (interfaces), not concrete adapters, many infrastructure changes touch no feature plan at all.

When an existing feature's own behavior or implementation genuinely changes, there is still one `spec.md` and one `plan.md` per feature — edited in place; git history preserves prior versions. **Do not create per-change files** (`plan-<change-description>.md`).

1. **Edit the spec in place:** Modify `spec.md` to describe the feature as it should be after the change. The spec is always living.
2. **Re-open the plan into a work order:** The pruned `plan.md` (a living design doc) is re-opened for this change — add back the transient sections (types/signatures, gaps to fix) for the delta and update the design rationale.
3. **Prune it again after shipping:** Once the change ships, prune the plan back into a living design doc, exactly as for a new feature.
