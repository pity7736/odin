# Technical Plan: Authentication

**Corresponds to Spec:** `specs/accounts/authentication/spec.md`

## Overview

Login, logout, session management, and route protection across both the web
(cookie-based) and mobile (bearer-token) interfaces. Passwords are verified with
bcrypt; session tokens come from `crypto/rand`; sessions use a 30-day
sliding-window TTL that is extended on every authenticated request.

This update changes how the REST login **response** is shaped and how a failed
login is **classified**: the response now carries only the field relevant to the
outcome, and wrong credentials are reported as an authentication failure (401)
distinct from malformed input (400), consistently on both interfaces.

## Design Decisions & Rationale

- **REST login returns only the outcome-relevant field.** Success →
  `{"token": ...}`, failure → `{"error": ...}`, via `omitempty` on a single
  `response` struct. Rejected the prior `{token, error}` shape that always
  carried an empty companion field (dead data), and rejected two separate
  structs — the fields are mutually exclusive and each is never empty on its own
  path, so `omitempty` expresses the contract without extra types.
- **Wrong credentials are 401, malformed/missing input is 400.** These are
  different failures — "your credentials are wrong" vs "you sent a bad request" —
  and conflating both under `Domain`/400 hid that. A new `Unauthorized`
  odinerror tag lets the layers distinguish them.
- **Status is chosen once, in the shared orchestrator.** `login_handler` maps
  the error tag → status, so REST and HTMX stay consistent from a single source
  of truth. Rejected per-strategy status ownership (two places to keep in sync,
  more surgery).
- **The web htmx config gains a `401 → swap:true` rule.** `base.gohtml`
  deliberately makes `400` swap so the login-error fragment renders on the web;
  moving wrong-credentials to 401 would otherwise fall into the `[45].. → no
  swap` bucket and silently stop rendering the error. The new rule mirrors the
  existing 400 rule.
- **Only 4xx produces a client error body; everything else is a 500.**
  `login_handler` renders an error body for `Unauthorized`/`Domain` only; any
  other error is returned up to the central `errorHandler`. This removes the
  prior latent nil-pointer panic (a propagated non-odin error was fed to
  `HandleBadRequest`, which dereferences a nil `*odinerrors.Error`).
- **Strategy pattern for handlers.** A shared orchestrator (`login_handler`,
  `logout_handler`) runs the use case; a per-interface strategy renders the
  result (cookie + redirect for web, JSON for mobile).
- **Password hashing is an application concern, not a domain one.** The `User`
  entity stores an already-hashed password and knows nothing about hashing;
  verification goes through the `PasswordHasher` port.
- **Sliding-window sessions.** Each validated request calls `Extend`, resetting
  expiry to now + TTL, so activity keeps a session alive and inactivity lets it
  lapse. Expired sessions are deleted on read.
- **Two auth carriers, one enforcement point.** Web sends the
  `__Secure-odin-session` cookie (`Secure`, `HttpOnly`, `SameSite=Strict`);
  `/api/v1` sends a bearer token. `loginRequired` centralizes enforcement —
  401 for `/api/`, redirect to `/auth/login?next=<path>` for the web — with no
  inline auth checks in feature handlers.
- **Field-agnostic credential error.** Both a wrong email and a wrong password
  return the identical message and status, leaking nothing about which field
  was wrong.
- **Logout deletes the exact session the middleware validated.** The auth
  middleware already resolves and validates the request's session, so it stashes
  that session's token in the request (`handler.SessionTokenKey`, alongside the
  `RequestContext` it already stores). Logout reads that token and terminates it,
  doing no cookie/bearer extraction of its own. This replaces a shared
  `extractToken` that preferred the cookie over the bearer token — a bug where a
  REST logout deleted a stray cookie's session (when present) while leaving the
  caller's bearer session alive and still returning success, so the caller was
  never actually logged out. Rejected both the shared cookie-first extractor and
  a per-interface `Token()` method on `LogoutHandler`: the latter still assumes
  the interface matches the credential (a cookie-authenticated REST call would
  delete nothing and false-succeed), whereas deleting the validated token is
  correct regardless of which credential authenticated the request.

## Architecture & Files Summary

```
src/shared/domain/odinerrors/
└── tags.go                                                 # MODIFY (add Unauthorized)

src/shared/utils/
└── random_string.go

src/accounts/domain/
├── usermodel/user.go
├── sessionmodel/session.go
└── repositories/
    ├── user.go
    └── session.go

src/accounts/application/
├── passwordhasher/password_hasher.go
└── use_cases/
    ├── sessionstarter/session_starter.go                   # MODIFY (Unauthorized tag)
    ├── sessionterminator/session_terminator.go
    └── sessionvalidator/session_validator.go

src/accounts/infrastructure/
├── api/
│   ├── loginhandler/
│   │   ├── login_handler.go                                # MODIFY (tag→status mapping)
│   │   └── body.go
│   ├── logouthandler/logout_handler.go                     # MODIFY (read validated token)
│   ├── htmx/
│   │   ├── htmxloginhandler/handler.go
│   │   └── htmxlogouthandler/handler.go
│   └── rest/
│       ├── restloginhandler/handler.go                     # MODIFY (single struct + omitempty)
│       └── restlogouthandler/handler.go
├── security/bcrypthasher/bcrypt_hasher.go
└── repositories/pgrepositories/
    ├── user_repository.go
    └── session_repository.go

src/shared/domain/requestcontext/context.go
src/shared/infrastructure/api/handler.go                     # MODIFY (SessionTokenKey const)

src/app/
└── fiber_application.go                                     # MODIFY (errorHandler: Unauthorized→401; validateSession stashes token)

src/shared/infrastructure/templates/
└── base.gohtml                                             # MODIFY (htmx 401 → swap:true)

tests/unit/accounts/application/use_cases/
└── login_test.go                                           # MODIFY (Unauthorized tag)

tests/unit/accounts/infrastructure/api/login_api_test/
└── login_test.go                                           # MODIFY (401 + field presence)

specs/accounts/authentication/
├── spec.md
└── plan.md                                                 # MODIFY
```

Repositories are referenced by their domain ports (`UserRepository`,
`SessionRepository`); the `pgrepositories` adapters are wired at the composition
root and owned by their own concern — swapping them does not touch this plan.

## Data Flow

**Login (REST — `POST /api/v1/auth/login`):** orchestrator parses the body and
validates presence (`strings.Clone` on parsed values) → `SessionStarter.Start`
looks up the user by email and compares the password via `PasswordHasher` →
on success it creates a session (`crypto/rand` token, TTL) and persists it →
strategy returns `{"token": ...}` at 201. On failure the orchestrator maps the
error tag to a status and the strategy returns `{"error": ...}`.

**Login (HTMX — `POST /auth/login`):** same orchestrator and use case; the web
strategy sets the session cookie and an `HX-Redirect` to `next` on success, or
renders the `login_error` fragment on failure.

**Request authentication (middleware):** the cookie middleware (global) and the
bearer middleware (`/api/v1`) both read the token, call `SessionValidator.Validate`
(unknown/expired → anonymous; other errors → 500), extend the session on success,
and place a `RequestContext` in `ctx.Locals`. `loginRequired` then admits
authenticated requests and rejects the rest (401 for `/api/`, redirect for web).

**Logout:** the auth middleware has already validated the request's session and
stashed its token in the request; logout reads that token and
`SessionTerminator.Terminate` deletes that session. The web strategy then clears
the cookie and redirects; the mobile strategy returns a confirmation message.
Because the token deleted is the one that authenticated the request, a second
logout with the same credential is rejected by the middleware (401).

## Request & Response

**Request data** (login):

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| email | string | yes | empty → 400 "El correo es obligatorio" |
| password | string | yes | empty → 400 "La contraseña es obligatoria" |

**REST** — `POST /api/v1/auth/login`
```json
// request
{ "email": "some@email.com", "password": "secret" }
// success 201
{ "token": "<token>" }
// wrong credentials 401
{ "error": "Correo o contraseña incorrectos" }
// missing/empty field or malformed body 400
{ "error": "El correo es obligatorio" }
```

**HTMX** — `POST /auth/login`
- Form fields: `email`, `password` (query `next` carries the post-login target)
- Success: set `__Secure-odin-session` cookie (`Secure`, `HttpOnly`,
  `SameSite=Strict`) + `HX-Redirect: <next>`
- Error: render `login_error` fragment into `#login_error`. Wrong credentials
  now return 401; the htmx `responseHandling` config swaps on `400` and `401`
  so the fragment renders for both validation and credential failures.

**REST** — `DELETE /api/v1/auth/logout` → `{"message": "session closed"}`
(bearer token identifies the session).

**HTMX** — `POST /auth/logout` → clears the cookie + `HX-Redirect: /auth/login`.

## Key Types & Signatures

- `odinerrors.Unauthorized` — new `Tag` constant.
- `sessionstarter.start` — the wrong-credentials error is built with
  `WithTag(odinerrors.Unauthorized)` (message English, external Spanish
  unchanged).
- `loginHandler.login` — on use-case error, select status by tag:
  `Unauthorized`→401 and `Domain`→400 both delegate to `HandleBadRequest`; any
  other error is returned up (no `HandleBadRequest` call) so `errorHandler`
  yields 500. `validateRequestBody` failures remain 400.
- `restloginhandler.response` — one struct, `omitempty` on both fields:
  `Token string json:"token,omitempty"`, `Error string json:"error,omitempty"`.
- `errorHandler` — add `case odinerrors.Unauthorized: code = 401`.
- `handler.SessionTokenKey` — new fiber `Locals` key. `validateSession` sets
  `c.Locals(handler.SessionTokenKey, session.Token())` after validating the
  session. `logoutHandler.Logout` reads it (`token, _ := ctx.Locals(...).(string)`)
  and terminates that session; the cookie-first `extractToken` and its `strings`
  import are removed. The `LogoutHandler` interface is unchanged.

## Gaps / Bugs to Fix

- [ ] Add `Unauthorized` to `odinerrors/tags.go`.
- [ ] Tag the wrong-credentials error in `sessionstarter/session_starter.go:49`
      as `Unauthorized`.
- [ ] `login_handler.login`: map error tag → status; return non-4xx errors up
      the stack instead of forcing 400 through `HandleBadRequest` (fixes the
      nil-pointer panic path).
- [ ] `restloginhandler/handler.go`: collapse to a single `response` struct with
      `omitempty` on `token` and `error`.
- [ ] `base.gohtml`: add `{"code":"401","swap":true}` before the `[45]..` rule.
- [ ] `fiber_application.go` `errorHandler`: map `Unauthorized` → 401.
- [ ] Tests: REST and HTMX wrong-credentials assert 401; REST success asserts no
      `error` key; REST failure asserts no `token` key; use-case test asserts the
      `Unauthorized` tag.
- [ ] Tests: login where the use case returns a non-4xx error — (a) a raw
      non-odin error and (b) an odinerror with a non-4xx tag — assert the
      response is 500 and no panic. Covers `handleError`'s two error-forwarding
      branches (the nil-panic fix), required by the 100% business-logic rule.
- [ ] Mocks: repository interfaces are unchanged — no regeneration.
- [ ] Add `SessionTokenKey` const to `shared/infrastructure/api/handler.go`;
      `validateSession` stashes `session.Token()` under it after validating.
- [ ] `logout_handler.go`: read the token from `ctx.Locals(SessionTokenKey)` and
      terminate it; delete the cookie-first `extractToken` and its `strings`
      import. `LogoutHandler` interface unchanged (no mock regeneration).
- [ ] Tests: REST logout terminates the session that authenticated the request
      (the bearer session) even when a `__Secure-odin-session` cookie is also
      present (regression for the cookie-first bug); HTMX logout terminates the
      cookie session; the middleware stashes the validated token.

## Quality Pillars

- **Security:** bcrypt verification, `crypto/rand` tokens, 30-day sliding TTL,
  cookie hardening (`Secure`/`HttpOnly`/`SameSite=Strict`), `strings.Clone` on
  parsed body values, centralized `loginRequired`. Wrong email and wrong
  password remain indistinguishable (identical 401 + message); the 400/401 split
  only distinguishes malformed input from an authentication failure, not which
  credential field was wrong.
- **Reliability:** nil-safe bearer/cookie middleware; propagated non-4xx errors
  now surface as 500 instead of nil-panicking on a 400 path; expired sessions
  are rejected and deleted; no panics in the auth flow.
- **Performance:** Deferred — in-memory repositories, minimal user set; no
  hot-path concern in this change.
- **Observability:** Deferred — no production telemetry yet; `odinerrors`
  location tracking provides debugging context.
