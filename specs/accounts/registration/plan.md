# Work Order: Registration — new feature

**Feature design:** `specs/accounts/registration/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/accounts/registration/spec.md`

> Work order for: **building registration**. Disposable — overwritten by the next
> change (git keeps the history). The living design is in design.md; hydrate it
> before this change merges, then freeze this file.

## Change

Build account registration. A new user is created from `email`, `auth_hash`,
`encrypted_master_key`, and `key_params` supplied by the client. The server
bcrypt-hashes the `auth_hash`, stores the user (with the opaque
`encrypted_master_key` and the `key_params` value object), and returns a minimal
`201 {id, email}`. Registration does **not** start a session — login is separate.
Seed data is removed; from now on users are created through registration.

Mirrors the existing login feature (`loginhandler` / `sessionstarter`).

Satisfies spec scenarios: *Creating an account successfully*, *Registering with
an email that already exists*, *Registering with a missing or empty email*,
*Registering with a missing or empty password*, *Registering with incomplete
unlock information*, *Registering does not log me in*.

## Architecture & Files (this change)
```
src/shared/domain/odinerrors/
└── tags.go                                                   # MODIFY  add AlreadyExists tag

src/accounts/application/
├── authhasher/auth_hasher.go                                 # MODIFY  add Hash to port
└── use_cases/userregistrar/user_registrar.go                 # CREATE  registration use case

src/accounts/infrastructure/
├── security/bcrypthasher/bcrypt_hasher.go                    # MODIFY  implement Hash
├── api/registerhandler/register_handler.go                   # CREATE  handler + response
├── api/registerhandler/body.go                               # CREATE  RegisterBody + Validate
└── repositories/inmemory/user_repository.go                  # MODIFY  remove seed (start empty)

src/app/
└── fiber_application.go                                      # MODIFY  409 mapping + register route/wiring

tests/unit/app/
└── error_handler_test.go                                     # CREATE  AlreadyExists -> 409 mapping
tests/unit/accounts/infrastructure/security/
└── bcrypt_hasher_test.go                                     # CREATE  Hash/Compare round-trip
tests/unit/accounts/application/use_cases/
└── register_test.go                                          # CREATE  use case unit tests
tests/unit/accounts/infrastructure/api/register_api_test/
├── register_test.go                                          # CREATE  handler tests
└── body_test.go                                              # CREATE  body validation tests
tests/unit/accounts/infrastructure/api/login_api_test/
└── login_test.go                                             # MODIFY  arrange user via builder (seed gone)
tests/integration/accounts/
└── auth_test.go                                              # MODIFY  register->login round-trip, duplicate, builder-arranged login
```

No mock regeneration: no repository interface changed (`Add` already exists);
`AuthHasher` is not mocked (mockery config covers only the two repositories).

## Key Types & Signatures

```go
// src/shared/domain/odinerrors/tags.go  — new tag, mapped to 409 in errorHandler
AlreadyExists Tag

// src/accounts/application/authhasher/auth_hasher.go  — port grows one method
type AuthHasher interface {
    Compare(storedHash, authHash string) bool
    Hash(authHash string) (string, error)
}

// src/accounts/application/use_cases/userregistrar — mirrors sessionstarter
func New(
    email, authHash, encryptedMasterKey string,
    keyParams keyparams.KeyParams,
    userRepository repositories.UserRepository,
    authHasher authhasher.AuthHasher,
) UserRegistrar
func (self UserRegistrar) Register(ctx context.Context) (*usermodel.User, error)

// src/accounts/infrastructure/api/registerhandler/body.go
type RegisterBody struct {
    Email              string          `json:"email"`
    AuthHash           string          `json:"auth_hash"`
    EncryptedMasterKey string          `json:"encrypted_master_key"`
    KeyParams          keyParamsBody   `json:"key_params"`
}
type keyParamsBody struct {
    Algorithm   string `json:"algorithm"`
    Salt        string `json:"salt"`
    Iterations  int    `json:"iterations"`
    Memory      int    `json:"memory"`
    Parallelism int    `json:"parallelism"`
}
func (self RegisterBody) Validate() error   // presence of email/auth_hash/encrypted_master_key

// src/accounts/infrastructure/api/registerhandler/register_handler.go
func New(userRepository repositories.UserRepository, authHasher authhasher.AuthHasher) registerHandler
func (self registerHandler) Register(ctx *fiber.Ctx) error
type registerResponse struct {
    ID    string `json:"id"`
    Email string `json:"email"`
}
```

External (Spanish) messages: duplicate → `"El correo ya está registrado"`
(AlreadyExists); empty email → `"El correo es obligatorio"`; empty auth_hash →
`"La contraseña es obligatoria"`; empty encrypted master key → internal English
`"encrypted master key is required"`, external `"Datos de solicitud inválidos"`
(it is a client-generated, user-invisible field — no bespoke user-facing text);
malformed body → `"Datos de solicitud inválidos"`; key-param errors come from
`keyparams.New` (already Spanish; shared with login, unchanged here).

## Implementation Phases (TDD)

Double-loop TDD: Phase 0 writes the failing happy-path **acceptance test** (the
outer loop / north star). It stays RED through Phases 1–5 and goes GREEN only when
the wiring lands in Phase 5. The inner unit loops (Phases 1–4) go green phase by
phase. Full-green `make check` is the exit criterion for the session — the outer
test is the one knowingly-red test until the end.

### Phase 0: Acceptance (outer loop) — failing happy-path integration test
**Red:** `auth_test.go` (integration): register a brand-new user via
`POST /api/v1/auth/register` → `201` with `{id, email}`; then log in via
`POST /api/v1/auth/login` with the same email + auth_hash → `201` returning
`token`, `encrypted_master_key`, and `key_params`. Fails now (no register route);
must pass once Phase 5 completes. Do NOT touch it again until then.
**Green:** achieved by Phases 1–5 — nothing implemented in this phase.

### Phase 1: Shared — AlreadyExists tag → 409
**Red:** `error_handler_test.go` asserts an odinerror tagged `AlreadyExists`
renders `409` with `{"error": <external>}`; a `Domain` error still renders `400`
(guard against regression).
**Green:** add `AlreadyExists` to `tags.go`; add its `case` in `errorHandler`
mapping to `http.StatusConflict`.

### Phase 2: Application/Infra — AuthHasher.Hash
**Red:** `bcrypt_hasher_test.go`: `Hash(x)` returns `(h, nil)` where
`Compare(h, x)` is true and `Compare(h, "other")` is false.
**Green:** add `Hash(authHash string) (string, error)` to the `AuthHasher`
interface; implement in `BcryptHasher` via `bcrypt.GenerateFromPassword` at
`bcrypt.DefaultCost`.

### Phase 3: Application — UserRegistrar use case
**Red:** `register_test.go` (MockUserRepository, real bcrypthasher):
- new email (`GetByEmail` → nil): `Add` is called once with a user whose email
  matches, whose `AuthHashDigest` satisfies `Compare(digest, authHash)`, and
  whose encrypted master key + key params are set; returns `(user, nil)`.
- duplicate (`GetByEmail` → existing user): returns an `AlreadyExists` error,
  `Add` is NOT called.
- `GetByEmail` error and `Add` error each propagate.
**Green:** implement `UserRegistrar`: `GetByEmail` → non-nil ⇒ AlreadyExists;
else `authHasher.Hash(authHash)` → `usermodel.New(email, digest,
encryptedMasterKey, keyParams)` → `userRepository.Add`; return the user.

### Phase 4: Infrastructure — body + handler
**Red:**
- `body_test.go`: empty email → "El correo es obligatorio"; empty auth_hash →
  "La contraseña es obligatoria"; empty encrypted_master_key → external "Datos de
  solicitud inválidos" (internal "encrypted master key is required"); all present
  → nil.
- `register_test.go` (handler, builder-arranged repo): success → `201`,
  `{id, email}` present, user persisted; duplicate email → `409` "El correo ya
  está registrado"; empty email → `400`; empty auth_hash → `400`; empty
  encrypted_master_key → `400`; invalid key params (e.g. empty algorithm / zero
  iterations) → `400` with the keyparams Spanish message; malformed body → `400`
  "Datos de solicitud inválidos".
**Green:** implement `RegisterBody.Validate()` (presence only); implement
`registerHandler.Register`: `BodyParser` (→ malformed 400), `Validate()`, build
`keyparams.New(...)` (its error returns as 400), `strings.Clone` every parsed
string (email, auth_hash, encrypted_master_key, algorithm, salt), construct
`userregistrar`, `Register`, respond `201 {id, email}`. Errors returned to the
global `errorHandler`.

### Phase 5: Infrastructure — wiring + seed removal
**Red:** `auth_test.go` (integration): register duplicate email → `409`. Rearrange
the existing login integration test(s) — and the login **unit** test in
`login_test.go` — to arrange their user via `userbuilder.New().WithEmail(...).
Create(userRepository)` and post `userbuilder.DefaultPassword` (seed is gone).
When the wiring below lands, the Phase 0 acceptance test goes GREEN.
**Green:** make `NewInMemoryUserRepository()` start with an empty map (remove the
two seeded users and the `hashPassword` helper); in `NewFiberApplication` build
`registerhandler.New(userRepository, authHasher)` and register
`apiV1.Post("/auth/register", register.Register)` (unauthenticated — no
`loginRequired`).

Finish: `make check` GREEN.

## Design decisions to hydrate into design.md

- [ ] Registration is its own feature/use case (`userregistrar`), mirroring
  `sessionstarter`; it does not start a session — login stays separate.
- [ ] `AlreadyExists` odinerror tag added (domain-named, not "Conflict") mapped
  to HTTP 409; duplicate email is the first user.
- [ ] `AuthHasher` port gains `Hash`; the server bcrypt-hashes the client-derived
  `auth_hash` at registration (double-hash chain: client Argon2id → server
  bcrypt). Bcrypt is sufficient because the input is already high-entropy.
- [ ] Handler builds the `KeyParams` value object from the nested body so the use
  case receives a valid value object; presence validation lives on `RegisterBody`.
- [ ] `encrypted_master_key` is stored opaquely and never echoed back; success
  response is minimal `{id, email}`.
- [ ] Seed data removed — users now exist only via registration; tests arrange
  users via `userbuilder`. In-memory repositories still lose data on restart.
- [ ] Registration is reachable unauthenticated (bearer middleware sets anonymous
  context; no `loginRequired` wrapper).
- [ ] Quality Pillars: Security (password never sent; server bcrypt over the
  auth_hash; `strings.Clone` on parsed body; opaque master key; field-agnostic
  presence errors), Reliability (duplicate rejected cleanly, repo errors
  propagate to global handler), Performance/Observability (deferred — in-memory,
  no telemetry yet).
