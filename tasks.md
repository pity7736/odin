# Tasks

Each feature needs a `specs/<module>/<feature>/spec.md` and `plan.md` before implementation.
Specs describe intended behavior. Plans reference existing code and flag gaps/bugs to fix.

## Priority 0 — Tooling

Safety net that protects everything else. Must be in place before feature work.

- [x] Upgrade Go to 1.26 and update all dependencies
- [x] Set up `.golangci.yml` with linters adapted from Gideon
- [x] Create `Makefile` with `lint`, `test`, `mocks`, `coverage`, `coverage-check`, `check` targets
- [x] Set up test coverage measurement and establish baseline threshold (91.9%)
- [x] Pre-commit hook (lint + test)

## Priority 1 — Foundation

Fix structural issues before specifying features.

- [ ] `specs/shared/error-handling/` — Centralize error-to-HTTP mapping, fix empty response bodies, complete NOT_FOUND handling, error response format for REST and HTMX. Includes: add test cases for ignored Render errors in `htmxcreateincomehandler`, add test cases for ignored Create errors in category handler calls (`fiber_application.go`), and define how application-level validation errors (not domain errors) propagate to HTTP status codes
- [ ] `specs/accounts/authentication/` — Consolidate auth (loginRequired everywhere, fix REST API nil dereference on invalid Bearer token, cookie + Bearer token handling)
- [ ] `specs/shared/routing/` — Extract route registration from the god function in fiber_application.go, consistent auth middleware for all routes

## Priority 2 — Core Accounting

Existing features that need specs to document intended behavior and fix known bugs.

- [ ] `specs/accounting/accounts/` — Create, list, get account (known bugs: none after buffer fix)
- [ ] `specs/accounting/categories/` — Create, list categories
- [ ] `specs/accounting/income/` — Create income (known bug: balance subtracts instead of adds)

## Priority 3 — New Features

Features from the original roadmap that need specs before implementation.

- [ ] `specs/accounting/expenses/` — Create, list expenses
- [ ] `specs/accounting/transfers/` — Create, list transfers
- [ ] `specs/shared/postgres-repositories/` — Real DB implementation for all entities

## Priority 4 — Hardening

- [ ] `specs/shared/input-validation/` — Unique names for accounts and categories, payload validation
- [ ] `specs/accounts/user-management/` — User invitations
- [ ] `specs/accounting/financial-institutions/` — Create, get, link to accounts
