# Work Order: Authentication — E2E encryption pivot

**Feature design:** `specs/accounts/authentication/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/accounts/authentication/spec.md`

> Work order for: **adapting authentication for E2E encryption**. Disposable —
> overwritten by the next change (git keeps the history). The living design is
> in design.md; hydrate it before this change merges, then freeze this file.

## Change

Adapt the authentication flow for a zero-knowledge, E2E-encrypted architecture.
The server no longer receives the user's password — it receives a pre-derived
auth hash. The `User` entity gains two new fields: an encrypted master key (an
opaque blob the server stores and returns) and key derivation params (a domain
value object describing how the client should derive keys from the password).
The `hashedPassword` field is renamed to `authHashDigest` to reflect that the
server stores a digest of the auth hash, not a hashed password. The login
response expands to return the session token, encrypted master key, and key
params. The `PasswordHasher` port is renamed to `AuthHasher` to reflect what it
actually does. All HTMX references are removed from the design. The bcrypt
implementation stays — hashing a high-entropy auth hash server-side is
sufficient.

The LoginHandler/LogoutHandler strategy pattern is removed — with HTMX gone,
only REST remains. Handlers render JSON directly and are app-scoped values
(created once at startup, `*fiber.Ctx` passed as method parameter). Error
handling is centralized in Fiber's global `errorHandler`, which renders JSON
with external messages for all odin errors. `LoginBody` validates itself.
`loginhandler` builds the response directly from session + user (no intermediate
`LoginResult` struct). `keyParamsResponse` has a constructor that maps from the
domain `KeyParams` value object.

**Spec scenarios satisfied:** Logging in successfully (receives session
credentials + encrypted data key + key params), Logging in with wrong
credentials, Logging in with a non-existing email, Logging in with missing or
empty fields, Session stays active, Session expires, Logging out, Logging out
actually ends the session, Accessing a protected feature without logging in.

## Architecture & Files (this change)

```
src/accounts/domain/
├── usermodel/user.go                              # MODIFY — rename hashedPassword→authHashDigest, add encryptedMasterKey, keyParams fields
├── keyparams/key_params.go                        # CREATE — KeyParams value object (value type, NOT pointer)
└── repositories/user.go                           # no change

src/accounts/application/
├── authhasher/auth_hasher.go                      # CREATE (rename from passwordhasher/)
└── use_cases/
    └── sessionstarter/session_starter.go           # MODIFY — rename password→authHash, PasswordHasher→AuthHasher, HashedPassword()→AuthHashDigest()

src/accounts/infrastructure/
├── api/
│   ├── loginhandler/
│   │   ├── login_handler.go                       # MODIFY — app-scoped value type, no LoginHandler interface, no LoginResult, handlers return errors for global errorHandler, JSON response built directly from session+user
│   │   └── body.go                                # MODIFY — rename Password→AuthHash, add Validate() method
│   ├── logouthandler/
│   │   └── logout_handler.go                      # MODIFY — app-scoped value type, no LogoutHandler interface, inline JSON response
│   └── rest/
│       ├── restloginhandler/handler.go            # DELETE — inlined into loginhandler
│       └── restlogouthandler/handler.go           # DELETE — inlined into logouthandler
├── security/
│   └── bcrypthasher/bcrypt_hasher.go              # MODIFY — implement AuthHasher instead of PasswordHasher, rename params to storedHash/authHash
└── repositories/inmemory/
    ├── user_repository.go                         # RENAME from pgrepositories/, MODIFY — rename PGUserRepository→InMemoryUserRepository, seed users with encrypted master key and key params
    └── session_repository.go                      # RENAME from pgrepositories/ — rename PGSessionRepository→InMemorySessionRepository

src/shared/domain/requestcontext/context.go            # NO CHANGE — carries userID/requestID, unaware of auth hashes or key params

src/app/
└── fiber_application.go                           # MODIFY — AuthHasher instead of PasswordHasher, handlers created once at startup,
                                                   # global errorHandler renders JSON with external messages for odin errors,
                                                   # remove restloginhandler/restlogouthandler imports
                                                   # Bearer middleware, validateSession, loginRequired: NO CHANGE

src/main.go                                        # MODIFY — wire AuthHasher, use inmemory repositories

tests/unit/accounts/domain/
├── user_test.go                                   # MODIFY — test new User fields, renamed getter
└── key_params_test.go                             # CREATE — test KeyParams value object with Spanish external messages

tests/unit/accounts/application/use_cases/
└── login_test.go                                  # MODIFY — rename password→authHash, assert user returned

tests/unit/accounts/infrastructure/api/
├── login_api_test/
│   ├── login_test.go                              # MODIFY — rename password→authHash, assert expanded response, assert error responses exclude key data, odin errors render JSON body
│   └── body_test.go                               # CREATE — unit tests for LoginBody.Validate()
└── logout_api_test/logout_test.go                 # MODIFY — remove restlogouthandler wiring, wire logouthandler directly

tests/unit/app/
└── middleware_test.go                             # no change

tests/builders/
├── userbuilder/user.go                            # MODIFY — build users with encrypted master key and key params
└── request_builder.go                             # MODIFY — use AuthHasher, bearer-only auth (remove cookie fallback)

tests/integration/accounts/
└── auth_test.go                                   # MODIFY — add login integration tests (success + wrong credentials), wire AuthHasher, use inmemory repositories

tests/unit/mocks/
├── mock_SessionRepository.go                      # no change
└── mock_UserRepository.go                         # no change

src/accounts/application/passwordhasher/           # DELETE
```

## Key Types & Signatures

```go
// Domain value object — src/accounts/domain/keyparams/key_params.go
// VALUE TYPE, not pointer. Every user MUST have key params.
// Struct field order: strings (16 bytes) first, then ints (8 bytes) — per 2.4 alignment rule.
type KeyParams struct {
    algorithm   string
    salt        string
    iterations  int
    memory      int
    parallelism int
}
func New(algorithm string, iterations, memory, parallelism int, salt string) (KeyParams, error)
// Returns value, not pointer. Validation errors include Spanish external messages.
// Getters: Algorithm(), Iterations(), Memory(), Parallelism(), Salt()
// Getter receivers: value receiver (KeyParams, not *KeyParams) — it's a value object.

// User entity changes — src/accounts/domain/usermodel/user.go
// Rename: hashedPassword → authHashDigest
// New fields: encryptedMasterKey string, keyParams keyparams.KeyParams (VALUE, not pointer)
// Constructor: New(email, authHashDigest, encryptedMasterKey string, keyParams keyparams.KeyParams) (*User, error)
// Renamed getter: HashedPassword() → AuthHashDigest()
// New getters: EncryptedMasterKey(), KeyParams()

// Renamed port — src/accounts/application/authhasher/auth_hasher.go
type AuthHasher interface {
    Compare(storedHash, authHash string) bool
}

// Login body — src/accounts/infrastructure/api/loginhandler/body.go
// Body validates itself via Validate() — field validation is not a handler concern.
type LoginBody struct {
    Email    string `json:"email"`
    AuthHash string `json:"auth_hash"`
}
func (self LoginBody) Validate() error

// Login handler — src/accounts/infrastructure/api/loginhandler/login_handler.go
// App-scoped value type (created once at startup, *fiber.Ctx passed as method param).
// No LoginHandler interface, no LoginResult struct.
// Handler returns errors — global errorHandler in fiber_application.go renders JSON.
// Success path builds loginResponse directly from session + user.
type loginHandler struct { ... }
func New(...) loginHandler
func (self loginHandler) Login(ctx *fiber.Ctx) error

// Login response structs (private, in loginhandler)
type loginResponse struct {
    Token              string             `json:"token,omitempty"`
    EncryptedMasterKey string             `json:"encrypted_master_key,omitempty"`
    KeyParams          *keyParamsResponse `json:"key_params,omitempty"`
    Error              string             `json:"error,omitempty"`
}
type keyParamsResponse struct { ... }
func newKeyParamsResponse(params keyparams.KeyParams) *keyParamsResponse

// Logout handler — src/accounts/infrastructure/api/logouthandler/logout_handler.go
// App-scoped value type. No LogoutHandler interface.
type logoutHandler struct { ... }
func New(sessionRepository repositories.SessionRepository) logoutHandler
func (self logoutHandler) Logout(ctx *fiber.Ctx) error

// Global error handler — src/app/fiber_application.go
// Centralized: maps odin error tags to HTTP status codes and renders JSON
// with external message. Non-odin errors get 500 with no body.
// Uses defer for ctx.Status(code) to avoid duplicated status-setting.
```

## Implementation Phases (TDD)

### Phase 1: Domain — KeyParams value object

**Red:** Create `tests/unit/accounts/domain/key_params_test.go`:
- Test that `New` creates a valid KeyParams with all fields accessible via getters.
- Test that `New` rejects empty algorithm → returns error with Spanish external message.
- Test that `New` rejects non-positive iterations → returns error with Spanish external message.
- Test that `New` rejects non-positive memory → returns error with Spanish external message.
- Test that `New` rejects non-positive parallelism → returns error with Spanish external message.
- Test that `New` rejects empty salt → returns error with Spanish external message.
- Assert the returned value is the zero value of `KeyParams` (not nil — it's a value type).

**Green:** Create `src/accounts/domain/keyparams/key_params.go`:
- `KeyParams` as a value type (not pointer).
- Struct field order: `algorithm string`, `salt string`, `iterations int`, `memory int`, `parallelism int` (strings first per alignment rule 2.4).
- Constructor returns `(KeyParams, error)`, not `(*KeyParams, error)`.
- All validation errors use `odinerrors` with both internal (English) and external (Spanish) messages.
- Getters use value receiver `(self KeyParams)`, not pointer receiver.

### Phase 2: Domain — User entity update

**Red:** Update `tests/unit/accounts/domain/user_test.go`:
- Test that `New` accepts `authHashDigest`, `encryptedMasterKey`, and `keyParams` (value type), accessible via getters.
- Test `AuthHashDigest()` getter (renamed from `HashedPassword()`).
- Test existing behavior still works (id, email).

**Green:** Modify `src/accounts/domain/usermodel/user.go`:
- Rename field `hashedPassword` → `authHashDigest`.
- Rename getter `HashedPassword()` → `AuthHashDigest()`.
- Add `encryptedMasterKey string` and `keyParams keyparams.KeyParams` fields (value type, not pointer).
- Update `New` constructor signature: `New(email, authHashDigest, encryptedMasterKey string, keyParams keyparams.KeyParams)`.
- Add `EncryptedMasterKey()` and `KeyParams()` getters.

### Phase 3: Application — Rename PasswordHasher to AuthHasher

**Red:** Update `tests/unit/accounts/application/use_cases/login_test.go`:
- Rename all references from password/PasswordHasher to authHash/AuthHasher.
- Update `HashedPassword()` → `AuthHashDigest()` references.
- Assert that `Start` returns `(*sessionmodel.Session, *usermodel.User, error)`.
- On success, assert the returned user matches.
- On failure, assert both session and user are nil.

**Green:**
- Create `src/accounts/application/authhasher/auth_hasher.go` with `AuthHasher` interface.
- Modify `src/accounts/application/use_cases/sessionstarter/session_starter.go`:
  - Rename `password` → `authHash`, `passwordHasher` → `authHasher`.
  - Change `Start` return type to `(*sessionmodel.Session, *usermodel.User, error)`.
  - Use `user.AuthHashDigest()` instead of `user.HashedPassword()`.
  - Return `(session, user, nil)` on success so the handler can build the response.
- Delete `src/accounts/application/passwordhasher/`.

### Phase 4: Infrastructure — bcrypt adapter

**Red:** No new tests — the bcrypt implementation is exercised through integration tests and the existing login tests.

**Green:** Modify `src/accounts/infrastructure/security/bcrypthasher/bcrypt_hasher.go` to implement `AuthHasher` instead of `PasswordHasher`. Rename params to `storedHash`/`authHash`. No structural change needed — the `Compare` method signature is identical.

### Phase 5: Infrastructure — Login handler and body validation

**Red:**
- Create `tests/unit/accounts/infrastructure/api/login_api_test/body_test.go`:
  - Test `LoginBody.Validate()` returns nil for valid body.
  - Test empty email returns odin error with Spanish external message.
  - Test empty auth hash returns odin error with Spanish external message.
- Update `tests/unit/accounts/infrastructure/api/login_api_test/login_test.go`:
  - Rename `password` → `auth_hash` in all request payloads.
  - Validation error for empty auth hash stays as `"La contraseña es obligatoria"` — the message is user-facing and users enter a password; the client derives the auth hash.
  - For the success case, assert the response contains `token`, `encrypted_master_key`, and `key_params` (with algorithm, iterations, memory, parallelism, salt).
  - For ALL error cases (wrong credentials, missing fields, malformed body, odin errors with any tag), assert `encrypted_master_key` and `key_params` are absent from the response. Odin errors render JSON body with external message regardless of tag.

**Green:**
- Modify `src/accounts/infrastructure/api/loginhandler/body.go`: rename `Password` → `AuthHash`, json tag `"password"` → `"auth_hash"`, add `Validate()` method with field validation and Spanish external messages.
- Modify `login_handler.go`:
  - App-scoped value type: `New` returns `loginHandler` (value, not pointer), no `*fiber.Ctx` in struct.
  - `Login(ctx *fiber.Ctx)` receives context as parameter.
  - Parsing error handled inline, field validation delegated to `body.Validate()`.
  - Handler returns errors to Fiber's global `errorHandler` — no local `handleError`/`renderError`.
  - Success path builds `loginResponse` directly from session + user, using `newKeyParamsResponse()`.
- Add `newKeyParamsResponse(params keyparams.KeyParams) *keyParamsResponse` constructor to encapsulate domain-to-response mapping.
- Delete `restloginhandler/` package entirely.

### Phase 6: Infrastructure — Logout handler and strategy removal

**Red:** Update `tests/unit/accounts/infrastructure/api/logout_api_test/logout_test.go`:
- Remove `restlogouthandler` import and wiring — wire `logouthandler` directly.

**Green:**
- Modify `src/accounts/infrastructure/api/logouthandler/logout_handler.go`:
  - App-scoped value type: `New` returns `logoutHandler` (value, not pointer), no `*fiber.Ctx` in struct.
  - `Logout(ctx *fiber.Ctx)` receives context as parameter.
  - Remove `LogoutHandler` interface, inline the JSON response (`{"message": "session closed"}`).
- Delete `src/accounts/infrastructure/api/rest/restlogouthandler/` package entirely.

### Phase 7: Centralized error handling

**Red:** Update `tests/unit/accounts/infrastructure/api/login_api_test/login_test.go`:
- Odin errors with non-4xx tags (e.g. Render) now render JSON body with external message (previously returned empty body). Update assertion from `assert.Zero(t, response.ContentLength)` to assert JSON error body.

**Green:**
- Modify `src/app/fiber_application.go` global `errorHandler`:
  - Render JSON `{"error": externalMessage}` for all odin errors, not just Domain/Unauthorized.
  - Use `defer func() { ctx.Status(code) }()` to avoid duplicated status-setting.
  - Non-odin errors: 500 with no body (unchanged).

### Phase 8: Wiring and test infrastructure

**Red:** Update `tests/integration/accounts/auth_test.go`:
- Add login integration test: successful login asserts token, encrypted master key, and key params.
- Add login integration test: wrong credentials asserts error message, no key data leaked.
- Update `tests/unit/app/middleware_test.go` to wire `AuthHasher` and `inmemory` repositories.

**Green:**
- Modify `src/app/fiber_application.go`: handlers created once at startup (`login := loginhandler.New(...)`, `logout := logouthandler.New(...)`), route handlers call `login.Login` and `logout.Logout` directly.
- Modify `src/main.go`: wire `bcrypthasher.New()` as `AuthHasher`, use `inmemory` repositories.
- Rename `pgrepositories/` → `inmemory/`, `PGUserRepository` → `InMemoryUserRepository`, `PGSessionRepository` → `InMemorySessionRepository`. Update all references.
- Modify `tests/builders/userbuilder/user.go`: build users with `authHashDigest`, encrypted master key, and key params (value type).
- Modify `tests/builders/request_builder.go`: use `AuthHasher`, bearer-only auth (remove cookie fallback), remove unused `handler` import.
- Modify user repository seed data with encrypted master key and key params.

### Phase 9: Mock regeneration

Run `go run github.com/vektra/mockery/v3` to regenerate mocks if any repository interface changed. Then `make check` to verify everything is green.

## Design decisions to hydrate into design.md

- [ ] AuthHasher replaces PasswordHasher — the server hashes auth hashes (already high-entropy), not passwords. Bcrypt is sufficient server-side.
- [ ] `hashedPassword` renamed to `authHashDigest` — the server never sees a password, it stores a digest of the client-derived auth hash.
- [ ] KeyParams is a domain value object and a VALUE TYPE (not pointer) — algorithm, iterations, memory, parallelism, salt. Validates structure with Spanish external messages, not specific algorithms. Every user must have key params; nil is not valid.
- [ ] KeyParams struct field ordering follows alignment rule 2.4: strings (16 bytes) before ints (8 bytes).
- [ ] KeyParams getters use value receivers — it's a value object, not an entity.
- [ ] User entity gains encryptedMasterKey (string) and keyParams (KeyParams value type) — stored at registration, returned at login.
- [ ] Login response contract: token + encrypted_master_key + key_params on success. Error responses exclude key data.
- [ ] HTMX / web interface removed entirely — REST-only, client-agnostic.
- [ ] SessionStarter.Start returns (Session, User, error) so the handler can build the response with the user's key data. Three return values is a smell — consider refactoring to an Authenticator use case with AuthResult{Session, User} when the next auth change lands.
- [ ] The server never sees the user's password — it receives and hashes a pre-derived auth hash. The double-hashing chain (client Argon2id → server bcrypt) ensures neither a database leak nor a transit interception compromises the vault.
- [ ] Repositories renamed from pgrepositories/PG* to inmemory/InMemory* — honest naming for in-memory implementations.
- [ ] LoginHandler/LogoutHandler strategy pattern removed — with HTMX gone, only REST remains. Handlers are app-scoped value types created once at startup, `*fiber.Ctx` passed as method parameter. No intermediate strategy interfaces or per-request allocation.
- [ ] LoginBody validates itself via Validate() — field validation is not a handler concern.
- [ ] Centralized error handling — Fiber's global errorHandler renders JSON with external messages for all odin errors. Handlers just return errors. Uses defer for ctx.Status(code) to avoid duplication.
- [ ] No LoginResult struct — handler builds loginResponse directly from session + user. The intermediate struct added no value with only one consumer.
- [ ] newKeyParamsResponse constructor encapsulates domain-to-response mapping for KeyParams.
- [ ] Request builder uses bearer-only auth (cookie fallback removed with HTMX).
- [ ] Remove Data Flow and Request & Response sections for HTMX.
- [ ] Update Architecture & Files Summary to reflect deleted HTMX handlers, new files, and renamed repositories.
- [ ] Update Quality Pillars: password never leaves device, auth hash double-hashing.
