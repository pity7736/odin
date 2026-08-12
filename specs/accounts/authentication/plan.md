# Technical Plan: Authentication

**Corresponds to Spec:** `specs/accounts/authentication/spec.md`

## Overview

Expose login, logout, session management, and route protection for both HTMX (web) and REST (mobile) interfaces. Passwords are hashed with bcrypt. Session tokens are generated with `crypto/rand`. Sessions expire after 30 days of inactivity using a sliding window TTL.

Architectural decisions applied:
- Strategy pattern for handlers (inject rendering strategy + repositories)
- Individual repository injection (no factory interfaces)
- Test boundary at repository layer (handler tests go through real use cases with mock repos)
- Password hashing is an application concern via `PasswordHasher` interface, not a domain concern

## Architecture & Files Summary

```
src/shared/utils/
└── random_string.go                                         # CREATE

src/accounts/domain/
├── usermodel/
│   └── user.go                                              # CREATE
├── sessionmodel/
│   └── session.go                                           # CREATE
└── repositories/
    ├── user.go                                              # CREATE
    └── session.go                                           # CREATE

src/accounts/application/
├── passwordhasher/
│   └── password_hasher.go                                   # CREATE
└── use_cases/
    ├── sessionstarter/
    │   └── session_starter.go                               # CREATE
    └── sessionterminator/
        └── session_terminator.go                            # CREATE

src/accounts/infrastructure/
├── api/
│   ├── loginhandler/
│   │   ├── login_handler.go                                 # CREATE
│   │   └── body.go                                          # CREATE
│   ├── logouthandler/
│   │   └── logout_handler.go                                # CREATE
│   ├── htmx/
│   │   ├── htmxloginhandler/
│   │   │   └── handler.go                                   # CREATE
│   │   └── htmxlogouthandler/
│   │       └── handler.go                                   # CREATE
│   └── rest/
│       ├── restloginhandler/
│       │   └── handler.go                                   # CREATE
│       └── restlogouthandler/
│           └── handler.go                                   # CREATE
├── security/bcrypthasher/
│   └── bcrypt_hasher.go                                     # CREATE
└── repositories/pgrepositories/
    ├── user_repository.go                                   # CREATE
    └── session_repository.go                                # CREATE

src/shared/domain/
├── requestcontext/
│   └── context.go                                           # CREATE
└── odinerrors/                                               # (exists)

src/shared/infrastructure/api/
└── handler.go                                               # CREATE

src/app/
└── fiber_application.go                                     # CREATE

src/main.go                                                  # CREATE

tests/unit/accounts/domain/
├── user_test.go                                             # CREATE
└── session_test.go                                          # CREATE

tests/unit/accounts/application/use_cases/
├── login_test.go                                            # CREATE
└── logout_test.go                                           # CREATE

tests/unit/accounts/infrastructure/api/
├── login_api_test/
│   └── login_test.go                                        # CREATE
└── logout_api_test/
    └── logout_test.go                                       # CREATE

tests/unit/testrepositoryfactory/
└── factory.go                                               # CREATE

tests/builders/
├── request_builder.go                                       # CREATE
└── userbuilder/
    └── user.go                                              # CREATE

tests/unit/mocks/                                            # REGENERATED
```

## Domain Layer

### `User` entity (`src/accounts/domain/usermodel/user.go`)

```go
type User struct {
    id             string
    email          string
    hashedPassword string
}

func New(email, hashedPassword string) (*User, error)
func (self *User) ID() string
func (self *User) Email() string
func (self *User) HashedPassword() string
```

- `New` generates a UUIDv7 for the ID. Receives an already-hashed password — the entity has no knowledge of hashing.
- Password hashing and comparison are handled by `PasswordHasher` in the application layer.

### `Session` entity (`src/accounts/domain/sessionmodel/session.go`)

```go
const DefaultTTL = 30 * 24 * time.Hour

type Session struct {
    expiresAt time.Time
    createdAt time.Time
    token     string
    userID    string
}

func New(userID string, ttl time.Duration) (*Session, error)
func NewFromRepository(token, userID string, createdAt, expiresAt time.Time) *Session
func (self *Session) Token() string
func (self *Session) UserID() string
func (self *Session) CreatedAt() time.Time
func (self *Session) ExpiresAt() time.Time
func (self *Session) IsExpired() bool
func (self *Session) Extend(ttl time.Duration)
```

- `New` generates token via `utils.RandomString` (`crypto/rand`), sets `createdAt = time.Now()`, `expiresAt = time.Now().Add(ttl)`.
- `NewFromRepository` reconstitutes from storage without generating a new token.
- `IsExpired` returns `time.Now().After(self.expiresAt)`.
- `Extend` sets `expiresAt = time.Now().Add(ttl)`.

### Repository interfaces (`src/accounts/domain/repositories/`)

**`user.go`:**

```go
type UserRepository interface {
    GetByEmail(ctx context.Context, email string) (*usermodel.User, error)
    Add(ctx context.Context, user *usermodel.User) error
}
```

**`session.go`:**

```go
type SessionRepository interface {
    Add(ctx context.Context, session *sessionmodel.Session) error
    Get(ctx context.Context, token string) (*sessionmodel.Session, error)
    Save(ctx context.Context, session *sessionmodel.Session) error
    Delete(ctx context.Context, token string) error
}
```

- `Get` returns `nil, nil` for not found or expired sessions.
- `Save` updates an existing session (sliding window extension).
- `Delete` removes a session (logout).

### Shared utils (`src/shared/utils/random_string.go`)

```go
func RandomString(length uint8) (string, error)
```

Uses `crypto/rand` to generate alphanumeric strings.

### Request context (`src/shared/domain/requestcontext/context.go`)

```go
type RequestContext struct {
    userID    string
    requestID string
}

func New(userID string) (*RequestContext, error)
func NewAnonymous() *RequestContext
func (self *RequestContext) UserID() string
func (self *RequestContext) RequestID() string
func (self *RequestContext) IsAuthenticated() bool
```

- `New` validates non-empty userID, generates UUIDv7 for requestID.
- `NewAnonymous` creates context with empty userID.
- `IsAuthenticated` returns `self.userID != ""`.

## Application Layer

### `PasswordHasher` interface (`src/accounts/application/passwordhasher/password_hasher.go`)

```go
type PasswordHasher interface {
    Compare(hashedPassword, password string) bool
}
```

- `Compare` checks a plain password against a hashed one.

### `SessionStarter` (`src/accounts/application/use_cases/sessionstarter/session_starter.go`)

```go
type SessionStarter struct {
    email             string
    password          string
    userRepository    repositories.UserRepository
    sessionRepository repositories.SessionRepository
    passwordHasher    passwordhasher.PasswordHasher
}

func New(email, password string, userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository, passwordHasher passwordhasher.PasswordHasher) *SessionStarter
func (self *SessionStarter) Start(ctx context.Context) (*sessionmodel.Session, error)
```

Steps:
1. `userRepository.GetByEmail(ctx, email)` — propagate error.
2. If user is nil or `passwordHasher.Compare(user.HashedPassword(), password)` fails → `odinerrors` with tag `DOMAIN`, message `"email or password are wrong"`, external `"Correo o contraseña incorrectos"`.
3. `sessionmodel.New(userID, DefaultTTL)` — propagate error.
4. `sessionRepository.Add(ctx, session)` — propagate error.
5. Return session.

### `SessionTerminator` (`src/accounts/application/use_cases/sessionterminator/session_terminator.go`)

```go
type SessionTerminator struct {
    sessionRepository repositories.SessionRepository
}

func New(sessionRepository repositories.SessionRepository) *SessionTerminator
func (self *SessionTerminator) Terminate(ctx context.Context, token string) error
```

Calls `sessionRepository.Delete(ctx, token)`.

## Infrastructure Layer

### Login handler orchestrator (`src/accounts/infrastructure/api/loginhandler/login_handler.go`)

```go
type LoginHandler interface {
    HandleResponse(session *sessionmodel.Session) error
    HandleBadRequest(err error) error
    ContentType() string
}

type loginHandler struct {
    userRepository    repositories.UserRepository
    sessionRepository repositories.SessionRepository
    passwordHasher    passwordhasher.PasswordHasher
    handler           LoginHandler
}

func New(userRepository repositories.UserRepository, sessionRepository repositories.SessionRepository, passwordHasher passwordhasher.PasswordHasher, handler LoginHandler) *loginHandler
func (self *loginHandler) Login(ctx *fiber.Ctx) error
```

- Parse body, validate email/password presence (use `strings.Clone` on parsed values).
- Validation errors use `odinerrors` with tag `DOMAIN`:
  - Empty email → message `"email is required"`, external `"El correo es obligatorio"`.
  - Empty password → message `"password is required"`, external `"La contraseña es obligatoria"`.
  - Malformed body → message `"wrong body"`, external `"Datos de solicitud inválidos"`.
- Construct `SessionStarter` with repos, call `Start`.
- Error → 400 + `strategy.HandleBadRequest(err)`. Success → 201 + `strategy.HandleResponse(session)`.

### HTMX login strategy (`src/accounts/infrastructure/api/htmx/htmxloginhandler/handler.go`)

- `HandleResponse`: set `__Secure-odin-session` cookie (`Secure`, `HttpOnly`, `SameSite=Strict`), set `HX-Redirect` to `ctx.Query("next", "/")`.
- `HandleBadRequest`: render `login_error` template with the error's external message.

### REST login strategy (`src/accounts/infrastructure/api/rest/restloginhandler/handler.go`)

- `HandleResponse`: return JSON `{"token": "<token>", "error": ""}`.
- `HandleBadRequest`: return JSON `{"token": "", "error": "<external message>"}`.

### Logout handler orchestrator (`src/accounts/infrastructure/api/logouthandler/logout_handler.go`)

```go
type LogoutHandler interface {
    HandleResponse() error
    ContentType() string
}

type logoutHandler struct {
    sessionRepository repositories.SessionRepository
    handler           LogoutHandler
}

func New(sessionRepository repositories.SessionRepository, handler LogoutHandler) *logoutHandler
func (self *logoutHandler) Logout(ctx *fiber.Ctx) error
```

- Extract token from cookie (HTMX) or Bearer header (REST).
- Construct `SessionTerminator`, call `Terminate`.
- Delegate to `strategy.HandleResponse()`.

### HTMX logout strategy (`src/accounts/infrastructure/api/htmx/htmxlogouthandler/handler.go`)

- Clear `__Secure-odin-session` cookie (`MaxAge: -1`), set `HX-Redirect: /auth/login`.

### REST logout strategy (`src/accounts/infrastructure/api/rest/restlogouthandler/handler.go`)

- Return JSON `{"message": "session closed"}`.

### Bcrypt hasher (`src/accounts/infrastructure/security/bcrypthasher/bcrypt_hasher.go`)

```go
type BcryptHasher struct{}

func New() BcryptHasher
func (self BcryptHasher) Compare(hashedPassword, password string) bool
```

- `Compare` uses `bcrypt.CompareHashAndPassword`.
- Dependency: `golang.org/x/crypto/bcrypt`.

### In-memory repositories (`src/accounts/infrastructure/repositories/pgrepositories/`)

**`user_repository.go`:**
- Backed by `map[string]*usermodel.User` keyed by email.
- Constructor seeds initial users via `usermodel.New` with pre-hashed passwords (uses bcrypt directly — infrastructure layer).
- `GetByEmail` returns `nil, nil` for unknown emails.

**`session_repository.go`:**
- Backed by `map[string]*sessionmodel.Session` keyed by token.
- `Get` checks `IsExpired()` before returning. Returns `nil, nil` for expired or unknown tokens.
- `Save` updates the session in the map.
- `Delete` removes from the map.

## Middleware and Route Protection

### Cookie middleware (global)

1. Set `requestcontext.Key` to `NewAnonymous()`.
2. Read `__Secure-odin-session` cookie. If empty, continue.
3. `sessionRepository.Get(ctx, token)` — on error, return 500.
4. If session nil, continue as anonymous.
5. If valid, create `RequestContext`, set in `ctx.Locals`. Call `session.Extend(DefaultTTL)`, `sessionRepository.Save(ctx, session)`.

### Bearer token middleware (`/api/v1` group)

1. Read `Authorization` header. If not `Bearer <token>` format, continue.
2. `sessionRepository.Get(ctx, token)` — on error, return 500.
3. If session nil, continue as anonymous.
4. If valid, create `RequestContext`, set in `ctx.Locals`. Call `session.Extend(DefaultTTL)`, `sessionRepository.Save(ctx, session)`.

### `loginRequired` function

```go
func loginRequired(ctx *fiber.Ctx, handler handler.Handler) error
```

- Authenticated → delegate to handler.
- Unauthenticated + JSON → 401.
- Unauthenticated + HTML → redirect to `/auth/login?next=<path>`.
- All protected routes use `loginRequired`. No inline auth checks.

### Route registration

**`NewFiberApplication` signature** receives repositories and the password hasher:

```go
func NewFiberApplication(
    accountingRepositoryFactory accountingrepositoryfactory.RepositoryFactory,
    sessionRepository accountsrepos.SessionRepository,
    userRepository accountsrepos.UserRepository,
    passwordHasher passwordhasher.PasswordHasher,
) Application
```

**Auth routes (no login required):**
- `GET /auth/login` — render login form.
- `POST /auth/login` — HTMX login.
- `POST /api/v1/auth/login` — REST login.

**Logout routes (login required):**
- `POST /auth/logout` — HTMX logout.
- `DELETE /api/v1/auth/logout` — REST logout.

**Protected routes (all use `loginRequired`):**
- `POST /categories`, `GET /categories`
- `POST /api/v1/categories`, `GET /api/v1/categories`
- `POST /accounts`, `GET /accounts`, `GET /accounts/:accountID`
- `POST /api/v1/accounts`
- `POST /accounts/:accountID/incomes`

### Entry point (`src/main.go`)

Create each repository and the bcrypt hasher, pass to `NewFiberApplication`.

## Implementation Phases (TDD)

### Phase 1: Domain — User

**Red:** `New` generates a UUIDv7 ID. Entity stores email and pre-hashed password.

**Green:** Implement `User` with UUIDv7. No hashing dependency.

### Phase 2: Domain — Session

**Red:** `New` returns token of expected length with `expiresAt ≈ now + ttl`. `IsExpired` returns false for new session, true when past. `Extend` resets `expiresAt`.

**Green:** Implement `Session` with `crypto/rand`.

### Phase 3: Use Cases

**Red:** Successful login returns session. Wrong password/email returns `"Correo o contraseña incorrectos"`. Repo errors propagate. Successful logout deletes session. Logout repo error propagates.

**Green:** Implement `SessionStarter` and `SessionTerminator`.

### Phase 4: Handlers

**Red:** REST/HTMX login with valid credentials. REST/HTMX login with wrong credentials returns Spanish errors. Validation errors return Spanish messages. REST/HTMX logout clears session. Accessing protected route after logout fails.

**Green:** Implement login/logout handlers and strategies.

### Phase 5: Middleware and Route Protection

**Red:** Cookie middleware sets `RequestContext` for valid session. Cookie middleware returns anonymous for expired session. Bearer middleware sets `RequestContext` for valid token. Bearer middleware returns anonymous for invalid token (no panic). Both middlewares return 500 on session lookup error. Both middlewares extend session on valid request. `loginRequired` delegates for authenticated users. `loginRequired` returns 401 for JSON, redirect for HTML. All category routes enforce auth via `loginRequired`.

**Green:** Implement middlewares, `loginRequired`, wire all routes.

### Phase 6: Mock Regeneration

Run `go run github.com/vektra/mockery/v3` after all repository interfaces are final.

## Quality Pillars

- **Security:** bcrypt password hashing, `crypto/rand` for session tokens, 30-day sliding window TTL, cookie attributes (`Secure`, `HttpOnly`, `SameSite=Strict`), consistent auth enforcement via `loginRequired`, `strings.Clone` on parsed request body values.
- **Reliability:** nil-safe Bearer token middleware, session lookup errors handled (not silenced), expired sessions rejected in both middlewares, no panics in auth flow.
- **Performance:** Deferred — single user, in-memory repos.
- **Observability:** Deferred — no production infrastructure. `odinerrors` location tracking provides debugging context.
