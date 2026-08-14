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
  - [x] `specs/accounting/accounts/create/` — Create account
- [ ] `specs/accounting/categories/` — Create, list categories. Plan must address: category handler Create error propagation (application errors vs domain errors). BUG: validation errors (empty name, empty type, invalid type) return 500 instead of 400 — commented-out tests in `tests/unit/accounting/infrastructure/api/category_api_test/category_test.go`
- [ ] `specs/accounting/income/` — Create income (known bug: balance subtracts instead of adds). Plan must address: ignored Render errors in htmxcreateincomehandler, ignored Add/Save errors in incomecreator (needs transaction design)

## Priority 2 — New Features

Features from the original roadmap that need specs before implementation.

- [ ] `specs/accounting/expenses/` — Create, list expenses
- [ ] `specs/accounting/transfers/` — Create, list transfers
- [ ] `specs/shared/postgres-repositories/` — Real DB implementation for all entities
- [ ] `specs/shared/i18n/` — Translation/localization system for all user-facing strings (labels, account type names, error messages) and localized long-form dates. Deferred until a second locale is needed; today the UI hardcodes Spanish. Would let clients (and the accounts-list view) share one source of localized labels instead of duplicating per view.

## Tech Debt

- [x] Fix `.mockery.yaml` to use mockery v3 config syntax (`with-expecter` and `outpkg` are v2 keys rejected by v3) so `make mocks` works correctly and mocks are never hand-edited again
- [ ] Extend `RequestBuilder` to support partial/broken requests (raw cookies, raw headers) and remove duplicate utility functions in `tests/testutils/requests.go`
- [ ] Review HTMX implementation: verify correct use of hx-* attributes, swap strategies, response handling config, and overall patterns across templates and handlers
- [ ] Robust error-handling design across the API. Current problems: (1) the tag→HTTP-status mapping is duplicated in `loginhandler.handleError` and `app.errorHandler`, so a tag change must be kept in sync in two places; (2) the central `errorHandler` returns only a status code with an **empty body**, dropping the external (Spanish) message — so accounts/categories give clients no error text, and each handler that wants to surface a message (login) must format its own body, which is why the mapping is duplicated in the first place; (3) `loginhandler.handleError`'s `default` branch drops the external message for any odin tag other than `Domain`/`Unauthorized` (latent — currently unreachable from the login path); (4) `restloginhandler`/`htmxloginhandler` `HandleBadRequest` discard the `errors.As` result and would nil-deref on a non-odin error (latent — the only caller pre-filters). Design a single source of truth for tag→status AND a consistent error-body contract (REST JSON `{"error": ...}`, HTMX fragment) so handlers can return the error up without losing the user-facing message.

## Dev Process Improvements

Durable improvements to how we build and verify the app. Unlike the Priority
sections, this stays after every feature has its spec and plan.

- [ ] Automate the Bruno collection in CI so API requests self-verify and run headless. Today the `.bruno` requests only carry `docs:` — running them proves transport, not correctness. Add an `assert` block (or Tests script) to every request asserting the HTTP status AND the response body (REST: `res.body.error` for rejections, the account fields for success), mirroring each spec's Expected Behavior. Add a login-first step (or `--env-var`) that captures the session token into `{{token}}` so a full run authenticates itself. Then wire `bru run --env local` into a make target + CI with a JSON/JUnit reporter. Caveat: this collection uses the newer opencollection `.yml` format (`opencollection: 1.0.0`) — confirm the installed `@usebruno/cli` version runs it before relying on it in CI. Note: Go unit tests remain the primary regression guard; Bruno automation is end-to-end smoke testing of the live app.

