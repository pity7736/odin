# Technical Design: Registration

**Corresponds to Spec:** `specs/accounts/registration/spec.md`

## Overview

Account creation for a zero-knowledge, E2E-encrypted architecture. A client
supplies an email, a client-derived `auth_hash`, an opaque `encrypted_master_key`,
and the `key_params` used to derive the auth hash. The server bcrypt-hashes the
`auth_hash`, persists the user with the opaque master key and key params, and
returns a minimal `201 {id, email}`. Registration never sees the password and
never starts a session — logging in is a separate step. It mirrors the login
feature (`registerhandler`/`userregistrar` parallel to `loginhandler`/
`sessionstarter`).

## Design Decisions & Rationale

- **Registration is its own use case (`userregistrar`) and does not start a
  session.** The client already holds all key material at registration time (it
  generated the master key and derived the auth hash), so returning session/key
  data buys nothing; keeping login separate also preserves login as an
  independent, testable handshake. Rejected auto-login (couples two concerns for
  a once-ever UX saving).
- **The endpoint is a resource: `POST /api/v1/users`.** Registration is creating
  a user resource, so it is modeled as a resource, not an RPC verb. Rejected
  `/auth/register` (verb path) and `/accounts/users` (leaks the internal Go
  module name into the URL, and "accounts" is ambiguous with financial accounts).
  Note the Go packages keep the business vocabulary (`registerhandler`,
  `userregistrar`) while the REST surface uses the resource noun — the feature is
  "registration"; the resource is "users".
- **`AlreadyExists` odinerror tag → HTTP 409.** A duplicate email is rejected
  with a domain-named tag (not "Conflict") that the global error handler maps to
  409. Domain vocabulary over transport vocabulary; the clear "already registered"
  message is intentional (users need to know to log in), enumeration hardening
  deferred.
- **`AuthHasher` port gains `Hash`; the server bcrypt-hashes the client-derived
  `auth_hash`.** Completes the double-hash chain (client Argon2id → server
  bcrypt). Bcrypt is sufficient server-side because the input is already the
  high-entropy output of client Argon2id. Previously only `Compare` existed
  (hashing lived only in the removed seed helper).
- **The handler builds the `KeyParams` value object; the body validates only
  presence.** `RegisterBody.Validate()` checks that email / auth_hash /
  encrypted_master_key are present; the handler constructs `keyparams.New(...)`
  so the use case receives a valid value object rather than loose fields, and the
  domain constructor (shared with login) enforces key-param validity with its own
  Spanish messages.
- **`encrypted_master_key` is opaque and never echoed back; success is minimal
  `{id, email}`.** The server stores but cannot read the master key. A missing
  `encrypted_master_key` returns the generic `"Datos de solicitud inválidos"`
  (internal `"encrypted master key is required"`) rather than a bespoke message —
  it is a client-generated, user-invisible field, unlike email/password which get
  specific user-facing messages.
- **Seed data removed — users exist only through registration.** The in-memory
  user repository now starts empty; tests arrange users via the `userbuilder`
  test builder instead of relying on seeded accounts.
- **Registration is reachable unauthenticated.** It sits on the `/api/v1` group
  whose bearer middleware sets an anonymous context when no token is present; no
  `loginRequired` wrapper guards it.

## Architecture & Files Summary

```
src/accounts/domain/
├── usermodel/user.go            # reused: New(email, authHashDigest, encryptedMasterKey, keyParams)
├── keyparams/key_params.go      # reused: validates key-param values
└── repositories/user.go         # port: UserRepository (GetByEmail, Add)

src/accounts/application/
├── authhasher/auth_hasher.go              # port: Compare + Hash
└── use_cases/userregistrar/user_registrar.go

src/accounts/infrastructure/
├── api/registerhandler/
│   ├── register_handler.go      # handler + registerResponse, wired at POST /api/v1/users
│   └── body.go                  # RegisterBody + keyParamsBody + Validate
├── security/bcrypthasher/bcrypt_hasher.go # implements Hash
└── repositories/inmemory/user_repository.go # seed removed (starts empty)

src/shared/domain/odinerrors/tags.go        # AlreadyExists tag
src/app/fiber_application.go                 # 409 mapping + POST /users wiring

tests/unit/accounts/application/use_cases/register_test.go
tests/unit/accounts/infrastructure/api/register_api_test/{register_test.go,body_test.go}
tests/unit/accounts/infrastructure/security/bcrypt_hasher_test.go
tests/unit/app/error_handler_test.go
tests/integration/accounts/auth_test.go      # register->login round-trip

specs/accounts/registration/
├── spec.md
├── design.md
└── plan.md
```

## Data Flow

**Register (`POST /api/v1/users`):** the handler parses the JSON body and delegates
presence validation to `RegisterBody.Validate()` (`strings.Clone` on every parsed
string). It builds a `KeyParams` value object via `keyparams.New` (its validation
errors surface as 400). It constructs `UserRegistrar` and calls `Register`, which:
looks up the email via the `UserRepository` port — a non-nil result yields an
`AlreadyExists` error (409); otherwise it `Hash`es the auth hash via the
`AuthHasher` port, builds the user (`usermodel.New`, UUIDv7 id), and persists it
via `Add`. On success the handler responds `201 {id, email}`. On failure it
returns the error to Fiber's global `errorHandler`, which maps the odin tag to a
status code and renders `{"error": externalMessage}`.

## Request & Response

**Register** — `POST /api/v1/users` (unauthenticated)
```json
// request
{
  "email": "julian@example.com",
  "auth_hash": "<client-derived-auth-hash>",
  "encrypted_master_key": "<opaque-blob>",
  "key_params": {
    "algorithm": "argon2id",
    "iterations": 3,
    "memory": 65536,
    "parallelism": 4,
    "salt": "<base64-salt>"
  }
}
// success 201
{ "id": "<uuid>", "email": "julian@example.com" }
// duplicate email 409
{ "error": "El correo ya está registrado" }
// missing/invalid field or malformed body 400
{ "error": "El correo es obligatorio" }
```

## Known Limitations

- **bcrypt 72-byte limit is handled asymmetrically.** Register's `Hash`
  (`bcrypt.GenerateFromPassword`) errors on an `auth_hash` longer than 72 bytes,
  while login's `Compare` silently truncates to 72. The current convention
  (Argon2id 32-byte output, base64 ≈ 44 bytes) stays well under 72, so the
  contract is: `auth_hash` must be ≤ 72 bytes. A SHA-256 pre-hash on both sides
  would remove the limit and the asymmetry if longer hashes are ever needed.
- **No atomic uniqueness under concurrency (in-memory adapter).** The duplicate
  check (`GetByEmail`) and the insert (`Add`) are separate calls, so concurrent
  registrations of the same email can both pass the check (TOCTOU race), and the
  in-memory map has no mutex (concurrent writes panic). Duplicate safety relies on
  the future persistent store's unique constraint; the in-memory repository is a
  throwaway adapter.
- **Non-odin errors return an empty-body 500.** Infrastructure failures in the
  register path (uuid generation, hash, or repository errors) are not
  `*odinerrors.Error`, so the global `errorHandler` falls through to a 500 with no
  body. This is the same fall-through that makes unmatched routes return 500
  instead of 404 (tracked as a bug in `tasks.md`).
- **In-memory repositories lose all data on restart.** Storage adapter TBD.

## Quality Pillars

- **Security:** The password never leaves the device; the server bcrypt-hashes the
  client-derived `auth_hash` (double-hash chain) and stores the `encrypted_master_key`
  opaquely, never echoing it back. `strings.Clone` guards against fasthttp buffer
  reuse on every parsed body value. Presence errors do not leak which stored data
  exists. Duplicate email is knowingly distinguishable (409) — an accepted
  enumeration tradeoff; rate limiting deferred.
- **Reliability:** Duplicate emails are rejected cleanly (`AlreadyExists` → 409)
  and repository errors propagate to the global handler. Known gaps: the
  non-atomic duplicate check and missing mutex on the in-memory adapter, and the
  empty-body 500 for non-odin errors (see Known Limitations).
- **Performance:** Deferred — in-memory store, one bcrypt hash per registration (a
  rare operation); no hot path in this change.
- **Observability:** Deferred — no production telemetry yet; `odinerrors` location
  tracking provides debugging context.
