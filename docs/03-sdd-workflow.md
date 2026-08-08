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
- **Example: `specs/accounting/accounts/spec.md`**

  ```markdown
  # Feature: Financial Accounts

  ## Overview
  Users need to manage their financial accounts to track where their money is.

  ## User Story
  As a user, I want to create and view my financial accounts so that I can
  organize my finances by institution or purpose (savings, checking, cash).

  ## Acceptance Criteria
  - I can create an account with a name and an initial balance.
  - I can see a list of all my accounts with their current balances.
  - I can view the details of a specific account.
  - Each account belongs to a single user and cannot be seen by others.

  ## Expected Behavior

  ### Creating an account
  - Given I am logged in
  - When I create an account with name "Savings" and initial balance of 100,000
  - Then the account is created with a balance equal to the initial balance

  ### Viewing my accounts
  - Given I have accounts "Savings" and "Cash"
  - When I view my accounts list
  - Then I see both accounts with their names, balances, and creation dates
  - And I do not see accounts belonging to other users

  ## Out of Scope
  - Deleting or archiving accounts
  - Editing account names after creation
  ```

## Plan File (The HOW)

- **Audience:** Engineers.
- **Purpose:** To describe the high-level technical implementation plan for a corresponding specification.
- **For existing features:** The plan references existing code that satisfies requirements and flags bugs, missing tests, and gaps as items to fix.
- **What to include:**
  - Component/package structure
  - Key interfaces and function signatures
  - Data flow between components
  - Implementation phases (as TDD Red-Green-Refactor steps)
  - Existing code that already satisfies requirements (with file paths)
  - Bugs or gaps to fix
  - Quality Pillars section (see below)
  - Files summary (CREATE / MODIFY)
- **What NOT to include:**
  - Full implementation of functions
  - Complete method bodies with business logic
  - Detailed error handling beyond error types
- **Example: `specs/accounting/accounts/plan.md`**

  ```markdown
  # Technical Plan: Financial Accounts

  **Corresponds to Spec:** `specs/accounting/accounts/spec.md`

  ## Overview
  Account creation and listing already exist. This plan documents the current
  implementation and identifies gaps.

  ## Existing Code
  - Domain: `src/accounting/domain/account/account.go`
  - Use Case: `src/accounting/application/use_cases/accountcreator/`
  - Handlers: REST and HTMX variants in `src/accounting/infrastructure/api/handlers/accounthandler/`
  - Repository: `src/accounting/infrastructure/repositories/pgrepositories/account_repository.go`

  ## Gaps
  - [ ] Missing REST handlers for GET /accounts and GET /accounts/:id
  - [ ] No input sanitization on account name

  ## Quality Pillars
  - **Security:** Ownership validation via requestContext — done
  - **Reliability:** Decimal arithmetic for balances, NOT_FOUND error tagging — done
  - **Performance:** Deferred — no bottleneck risk at current scale
  - **Observability:** Deferred — no production infrastructure yet

  ## Files Summary
  | Action | File |
  |--------|------|
  | MODIFY | `src/accounting/infrastructure/api/handlers/accounthandler/...` |
  | CREATE | `specs/accounting/accounts/spec.md` |
  | CREATE | `specs/accounting/accounts/plan.md` |
  ```

## Quality Pillars in Plans

Every plan **must** include a Quality Pillars section that addresses all four pillars from [docs/06-quality-pillars.md](./06-quality-pillars.md):

1. **Security**
2. **Reliability**
3. **Performance**
4. **Observability**

Each pillar must have at least one line stating what applies to this feature. **"Deferred"** is a valid answer when there is no production infrastructure to support it yet, but it must include a short justification. This ensures the decision to skip is conscious, not accidental.

## Updating a Feature

When an existing feature needs to change:

1. **Update the spec:** The spec always reflects the current intended behavior. Modify it to describe the feature as it should be *after* the change. Git history preserves previous versions.
2. **Write a new plan:** Create a new plan scoped to the change only, describing the delta from current implementation to the updated spec. Name it descriptively: `plan-<change-description>.md`.
3. **Keep the original plan:** It documents the initial implementation and remains useful for context.

For example:
- `specs/accounting/accounts/spec.md` — updated to include currency support
- `specs/accounting/accounts/plan.md` — original implementation plan (unchanged)
- `specs/accounting/accounts/plan-add-currency-support.md` — describes only what changes
