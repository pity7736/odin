# Technical Design: Authentication

**Corresponds to Spec:** `specs/accounts/authentication/spec.md`

## Overview

Login, logout, session management, and route protection for a zero-knowledge,
E2E-encrypted architecture. The server never sees the user's password — it
receives a pre-derived auth hash from the client. The double-hashing chain
(client Argon2id → server bcrypt) ensures neither a database leak nor a transit
interception compromises the vault. Sessions use a 30-day sliding-window TTL
extended on every authenticated request. The login response returns the session
token, the user's encrypted master key, and key derivation params so the client
can unlock the vault.

## Design Decisions & Rationale

- **AuthHasher replaces PasswordHasher.** The server hashes auth hashes (already
  high-entropy output of client-side Argon2id), not passwords. The port is named
  to reflect what it actually does. Bcrypt is sufficient server-side because
  the input is already high-entropy.
- **`hashedPassword` renamed to `authHashDigest`.** The server never sees a
  password; it stores a digest of the client-derived auth hash. The field name
  reflects the actual data.
- **KeyParams is a domain value object and a VALUE TYPE (not pointer).** It
  holds algorithm, iterations, memory, parallelism, and salt. It validates
  structure (non-empty, positive values) with Spanish external messages, but not
  specific algorithms — the domain doesn't know which algorithms are valid;
  that's an infrastructure/crypto concern. Every user must have key params; the
  zero value is the error case, not nil. Struct field ordering follows alignment
  rule 2.4: strings (16 bytes) before ints (8 bytes). Getters use value
  receivers — it's a value object, not an entity.
- **User entity gains encryptedMasterKey and keyParams.** Both stored at
  registration, returned at login. `encryptedMasterKey` is an opaque string
  (the server stores it but cannot decrypt it). `keyParams` is a KeyParams
  value type.
- **Login response contract.** Success: `token` + `encrypted_master_key` +
  `key_params` (algorithm, iterations, memory, parallelism, salt). Error
  responses exclude key data via `omitempty`.
- **SessionStarter.Start returns (Session, User, error).** The handler needs
  both session (token) and user (encrypted master key, key params) to build the
  login response. Three return values is a known smell — consider refactoring to
  an Authenticator use case with `AuthResult{Session, User}` when the next auth
  change lands.
- **No handler strategy pattern.** With HTMX removed, only REST remains.
  Handlers are app-scoped value types created once at startup; `*fiber.Ctx` is
  passed as a method parameter. No intermediate strategy interfaces, no
  per-request allocation.
- **LoginBody validates itself.** Field validation (`Validate()`) lives on the
  body struct, not in the handler — the body knows its own rules.
- **Centralized error handling.** Fiber's global `errorHandler` renders JSON
  with external messages for all odin errors. Handlers just return errors —
  no local error rendering or status-code mapping. Uses `defer` for
  `ctx.Status(code)` to avoid duplicated status-setting. Non-odin errors get
  500 with no body.
- **No LoginResult struct.** The handler builds `loginResponse` directly from
  session + user. An intermediate struct added no value with only one consumer.
- **`newKeyParamsResponse` constructor.** Encapsulates the domain-to-response
  mapping for KeyParams, keeping the transformation inside the response type.
- **REST login returns only the outcome-relevant fields.** Success →
  `{token, encrypted_master_key, key_params}`, failure → `{"error": ...}`,
  via `omitempty` on a single `loginResponse` struct.
- **Wrong credentials are 401, malformed/missing input is 400.** Different
  failures — "your credentials are wrong" vs "you sent a bad request" — with
  distinct status codes via the `Unauthorized` odinerror tag.
- **Sliding-window sessions.** Each validated request calls `Extend`, resetting
  expiry to now + TTL. Expired sessions are deleted on read.
- **Bearer-only auth.** REST API uses bearer tokens. Cookie-based auth removed
  with the HTMX/web interface.
- **Field-agnostic credential error.** Both a wrong email and a wrong password
  return the identical message and status, leaking nothing about which field
  was wrong.
- **Logout deletes the exact session the middleware validated.** The auth
  middleware resolves and validates the request's session, stashing that
  session's token in the request (`handler.SessionTokenKey`). Logout reads that
  token and terminates it, doing no token extraction of its own.
- **Repositories named for what they are.** `inmemory/InMemoryUserRepository`
  and `InMemorySessionRepository` — honest naming for in-memory
  implementations. Referenced by domain ports; swapping adapters doesn't touch
  this design.
- **User repository fetch contract: absence is a `NotFound` error, not `(nil,
  nil)`.** `UserRepository.GetByEmail` is a pure fetch that returns an
  `odinerrors` `NotFound` when the user is absent — consistent with
  `SessionRepository.GetByID`, and removing the nil-sentinel footgun (callers no
  longer nil-check; a real adapter signals absence the same way). Existence
  questions use a dedicated `UserRepository.Exists(email) (bool, error)` — no
  whole-user fetch, explicit intent — mirroring the chunk repo. `SessionStarter`
  translates a `NotFound` from `GetByEmail` into the field-agnostic
  wrong-credentials 401 (detected via `errors.As` + `Tag() == NotFound`, the same
  pattern `SessionValidator` uses); other errors propagate; a found user is
  compared. `UserRegistrar`'s duplicate check uses `Exists`. Behavior is
  unchanged — a wrong email still returns the identical 401.

## Architecture & Files Summary

```
src/accounts/domain/
├── usermodel/user.go
├── keyparams/key_params.go
├── sessionmodel/session.go
└── repositories/
    ├── user.go
    └── session.go

src/accounts/application/
├── authhasher/auth_hasher.go
└── use_cases/
    ├── sessionstarter/session_starter.go
    ├── sessionterminator/session_terminator.go
    └── sessionvalidator/session_validator.go

src/accounts/infrastructure/
├── api/
│   ├── loginhandler/
│   │   ├── login_handler.go
│   │   └── body.go
│   └── logouthandler/logout_handler.go
├── security/bcrypthasher/bcrypt_hasher.go
└── repositories/inmemory/
    ├── user_repository.go
    └── session_repository.go

src/shared/domain/
├── odinerrors/
├── requestcontext/context.go
└── ...

src/shared/infrastructure/api/handler.go

src/app/
└── fiber_application.go

tests/unit/accounts/domain/
├── user_test.go
└── key_params_test.go

tests/unit/accounts/application/use_cases/
└── login_test.go

tests/unit/accounts/infrastructure/api/
├── login_api_test/
│   ├── login_test.go
│   └── body_test.go
└── logout_api_test/logout_test.go

tests/unit/app/
└── middleware_test.go

tests/integration/accounts/
└── auth_test.go

tests/builders/
├── userbuilder/user.go
└── request_builder.go
```

## Data Flow

**Login (`POST /api/v1/auth/login`):** handler parses the JSON body and
delegates field validation to `LoginBody.Validate()` (`strings.Clone` on parsed
values) → `SessionStarter.Start` looks up the user by email and compares the
auth hash via `AuthHasher` → on success it creates a session (`crypto/rand`
token, TTL), persists it, and returns both session and user → handler builds
`loginResponse` with token, encrypted master key, and key params (via
`newKeyParamsResponse`) at 201. On failure the handler returns the error to
Fiber's global `errorHandler`, which maps the odin error tag to a status code
and renders `{"error": externalMessage}`.

**Request authentication (middleware):** the bearer middleware (`/api/v1`) reads
the token from `Authorization: Bearer <token>`, calls
`SessionValidator.Validate` (unknown/expired → anonymous; other errors → 500),
extends the session on success, and places a `RequestContext` in `ctx.Locals`.
`loginRequired` then admits authenticated requests and rejects the rest (401).

**Logout (`DELETE /api/v1/auth/logout`):** the auth middleware has already
validated the request's session and stashed its token in the request; logout
reads that token and `SessionTerminator.Terminate` deletes that session,
returning `{"message": "session closed"}`. Because the token deleted is the one
that authenticated the request, a second logout with the same credential is
rejected by the middleware (401).

## Request & Response

**Request data** (login):

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| email | string | yes | empty → 400 "El correo es obligatorio" |
| auth_hash | string | yes | empty → 400 "La contraseña es obligatoria" |

**Login** — `POST /api/v1/auth/login`
```json
// request
{ "email": "some@email.com", "auth_hash": "<client-derived-auth-hash>" }
// success 201
{
  "token": "<session-token>",
  "encrypted_master_key": "<opaque-blob>",
  "key_params": {
    "algorithm": "argon2id",
    "iterations": 3,
    "memory": 65536,
    "parallelism": 4,
    "salt": "<base64-salt>"
  }
}
// wrong credentials 401
{ "error": "Correo o contraseña incorrectos" }
// missing/empty field or malformed body 400
{ "error": "El correo es obligatorio" }
```

**Logout** — `DELETE /api/v1/auth/logout` (bearer token identifies the session)
```json
// success 200
{ "message": "session closed" }
```

## Known Limitations

- In-memory repositories — data is lost on restart. Storage adapter TBD.
- Session token is the only auth mechanism — no refresh tokens, no device
  management.
- No rate limiting on login attempts.

## Quality Pillars

- **Security:** The password never leaves the device. The client derives an auth
  hash (Argon2id) and sends that to the server; the server bcrypt-hashes the
  auth hash for storage. The double-hashing chain ensures neither a database
  leak nor a transit interception reveals the password. `crypto/rand` tokens,
  30-day sliding TTL, `strings.Clone` on parsed body values, centralized
  `loginRequired`. Wrong email and wrong password remain indistinguishable
  (identical 401 + message). Error responses never include key data.
- **Reliability:** absence is an explicit `NotFound` error, not a nil sentinel —
  callers can't forget to nil-check; nil-safe bearer middleware; propagated
  non-4xx errors surface as 500 via the global error handler (no nil-pointer
  panics); expired sessions are rejected and deleted; no panics in the auth flow.
- **Performance:** Deferred — in-memory repositories, minimal user set; no
  hot-path concern in this change.
- **Observability:** Deferred — no production telemetry yet; `odinerrors`
  location tracking provides debugging context.
