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

## Priority 1 — Core Accounting

Existing features that need specs to document intended behavior and fix known bugs.
Each plan should address relevant structural issues (error handling, auth, routing) as they arise.

- [ ] `specs/accounts/authentication/` — Login, session management, loginRequired consolidation, fix REST API nil dereference on invalid Bearer token
- [ ] `specs/accounting/accounts/` — Create, list, get account
- [ ] `specs/accounting/categories/` — Create, list categories. Plan must address: category handler Create error propagation (application errors vs domain errors). BUG: validation errors (empty name, empty type, invalid type) return 500 instead of 400 — commented-out tests in `tests/unit/accounting/infrastructure/api/category_api_test/category_test.go`
- [ ] `specs/accounting/income/` — Create income (known bug: balance subtracts instead of adds). Plan must address: ignored Render errors in htmxcreateincomehandler, ignored Add/Save errors in incomecreator (needs transaction design)

## Tech Debt

- [ ] Extend `RequestBuilder` to support partial/broken requests (raw cookies, raw headers) and remove duplicate utility functions in `tests/testutils/requests.go`

## Priority 2 — New Features

Features from the original roadmap that need specs before implementation.

- [ ] `specs/accounting/expenses/` — Create, list expenses
- [ ] `specs/accounting/transfers/` — Create, list transfers
- [ ] `specs/shared/postgres-repositories/` — Real DB implementation for all entities
