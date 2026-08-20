# Work Order: Create Chunk — new feature

**Feature design:** `specs/vault/chunks/create/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/vault/chunks/create/spec.md`

> Work order for: **building chunk creation**. Disposable — overwritten by the next
> change (git keeps the history). The living design is in design.md; hydrate it
> before this change merges, then freeze this file.

## Change

Introduce the `vault` module and the `EncryptedChunk` entity, and let a logged-in
user store one opaque encrypted chunk. The client sends a UUID `id` and an opaque
`content` blob (it packs `nonce‖ciphertext‖tag` itself); the server persists it
tied to the authenticated owner and never inspects it. Owner-scoped uniqueness:
if the owner already has that id, reject with 409. Minimal entity — `id`,
`ownerID`, `content`; `version`/`updated_at` arrive with the update/sync features.

Mirrors the auth feature's shape (handler → use case → repository port).
`NewFiberApplication` stays strict — all dependencies required (a `chunkRepository`
param is added) so production can never build with a nil repo. To stop constructor
changes from rippling across tests, the app-building test helpers route through one
shared **test app-builder** (`apptest`); future constructor changes touch only that
builder (and `main.go`), never the individual tests.

Satisfies spec scenarios: *Saving an item successfully*, *Saving without being
logged in*, *Saving an item whose identifier already exists*, *Saving without an
identifier*, *Saving without contents*.

## Architecture & Files (this change)
```
src/vault/domain/
├── chunkmodel/chunk.go                                   # CREATE  EncryptedChunk entity
└── repositories/chunk.go                                 # CREATE  ChunkRepository port

src/vault/application/use_cases/chunkcreator/
└── chunk_creator.go                                      # CREATE  use case

src/vault/infrastructure/
├── api/chunkhandler/
│   ├── chunk_handler.go                                  # CREATE  handler + response
│   └── body.go                                           # CREATE  CreateChunkBody + Validate
└── repositories/inmemory/chunk_repository.go             # CREATE  in-memory adapter

src/app/fiber_application.go                              # MODIFY  route + wiring + signature (+chunkRepository)
src/main.go                                               # MODIFY  wire in-memory chunk repo

.mockery.yaml                                             # MODIFY  add ChunkRepository
tests/unit/mocks/mock_ChunkRepository.go                  # REGEN

tests/unit/vault/domain/chunk_test.go                     # CREATE
tests/unit/vault/application/use_cases/create_test.go     # CREATE
tests/unit/vault/infrastructure/api/chunk_api_test/
├── create_test.go                                        # CREATE
└── body_test.go                                          # CREATE
tests/integration/vault/chunk_test.go                     # CREATE  (Phase 0 acceptance)

tests/builders/apptest/app.go                             # CREATE  shared test app-builder (absorbs constructor changes)

# MODIFY (one-time: route app construction through the apptest builder —
# future constructor changes won't touch these again):
tests/integration/accounts/auth_test.go                   # MODIFY
tests/unit/accounts/infrastructure/api/register_api_test/register_test.go  # MODIFY
tests/unit/app/error_handler_test.go                      # MODIFY
tests/unit/app/middleware_test.go                         # MODIFY
```

## Key Types & Signatures

```go
// src/vault/domain/chunkmodel — private fields; New validates: id is a valid UUID,
// ownerID non-empty, content non-empty.
type EncryptedChunk struct { /* id, ownerID, content string */ }
func New(id, ownerID, content string) (*EncryptedChunk, error)
func (self *EncryptedChunk) ID() string
func (self *EncryptedChunk) OwnerID() string
func (self *EncryptedChunk) Content() string

// src/vault/domain/repositories — owner-scoped
type ChunkRepository interface {
    Exists(ctx context.Context, ownerID, id string) (bool, error)
    Add(ctx context.Context, chunk *chunkmodel.EncryptedChunk) error
}

// src/vault/application/use_cases/chunkcreator
func New(id, ownerID, content string, chunkRepository repositories.ChunkRepository) ChunkCreator
func (self ChunkCreator) Create(ctx context.Context) (*chunkmodel.EncryptedChunk, error)

// src/vault/infrastructure/api/chunkhandler
type CreateChunkBody struct { ID string `json:"id"`; Content string `json:"content"` }
func (self CreateChunkBody) Validate() error   // presence of id + content
func New(chunkRepository repositories.ChunkRepository) chunkHandler
func (self chunkHandler) Create(ctx *fiber.Ctx) error   // owner from requestcontext.UserID()
type createChunkResponse struct { ID string `json:"id"` }

// src/app — stays strict/positional; gains a required chunkRepository so
// production can never build with a nil repo.
func NewFiberApplication(
    sessionRepository accountsrepos.SessionRepository,
    userRepository accountsrepos.UserRepository,
    authHasher authhasher.AuthHasher,
    chunkRepository vaultrepos.ChunkRepository,
) Application

// tests/builders/apptest — the ONE place that calls NewFiberApplication; defaults
// to in-memory repos + real bcrypt, with per-repo overrides. Absorbs future
// constructor changes so individual tests don't.
func New() *Builder
func (self *Builder) WithSessionRepository(r accountsrepos.SessionRepository) *Builder
func (self *Builder) WithUserRepository(r accountsrepos.UserRepository) *Builder
func (self *Builder) WithChunkRepository(r vaultrepos.ChunkRepository) *Builder
func (self *Builder) Build() app.Application
```

Error messages: id/content are client-generated, user-invisible fields, so their
external messages are the generic `"Datos de solicitud inválidos"` (internal
English: `"id is required"`, `"id must be a valid uuid"`, `"content is required"`,
`"owner id is required"`). Duplicate → `AlreadyExists` tag (→ 409), external
`"El elemento ya existe"`. Unauthenticated → 401 via `loginRequired` (no body).

## Implementation Phases (TDD)

Double-loop: Phase 0 writes the failing happy-path acceptance test first; it stays
RED until Phase 4 wires the route. Inner phases go green in order. `make check`
green is the exit criterion.

### Phase 0: Acceptance (outer loop) — failing happy-path integration test
**Red:** `tests/integration/vault/chunk_test.go`: an authenticated
`POST /api/v1/chunks` with `{id: <uuid>, content: <blob>}` → `201 {id}`, and the
chunk is persisted for that owner (assert via the in-memory repo instance). Fails
now (no route). Do not touch again until Phase 4.
**Green:** achieved by Phases 1–4.

### Phase 1: Domain — entity + port + mock
**Red:** `chunk_test.go`: `New` builds a chunk and exposes id/ownerID/content;
rejects a non-UUID id, empty ownerID, empty content (each a Domain error).
**Green:** implement `chunkmodel.EncryptedChunk` + `New` (UUID parse-check via
`github.com/google/uuid`); define the `ChunkRepository` port; add `ChunkRepository`
to `.mockery.yaml` and regenerate (`go run github.com/vektra/mockery/v3`) so
downstream tests have `MockChunkRepository`.

### Phase 2: Application — chunkcreator use case
**Red:** `create_test.go` (MockChunkRepository): new id (`Exists(owner,id)` → false)
→ `Add` called once with a chunk matching id/owner/content, returns the chunk;
duplicate (`Exists` → true) → `AlreadyExists` error, `Add` NOT called;
`Exists` error and `Add` error each propagate.
**Green:** implement `ChunkCreator`: owner-scoped `Exists` → true ⇒
AlreadyExists; else `chunkmodel.New` → `Add`.

### Phase 3: Infrastructure — body + handler
**Red:** `body_test.go`: empty id / empty content → Domain error. `create_test.go`
(handler, authenticated via RequestBuilder; MockChunkRepository): success → `201`
`{id}`, chunk persisted; duplicate → `409` "El elemento ya existe"; empty id →
`400`; empty content → `400`; non-UUID id → `400`; malformed body → `400`;
unauthenticated (anonymous session) → `401`.
**Green:** `CreateChunkBody.Validate()` (presence); `chunkHandler.Create` reads the
owner via `ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext).UserID()`,
`strings.Clone`s parsed body fields, builds `chunkcreator`, responds `201 {id}`;
errors returned to the global `errorHandler`.

### Phase 4: Infrastructure — test app-builder + wiring
**Red:** create the shared `apptest` builder and route the existing app-building
test helpers (`newIntegrationApp`, `newApplication`, middleware/error-handler
builders) through it; the Phase 0 acceptance test now passes.
**Green:** add `chunkRepository` as a required param to `NewFiberApplication`
(production stays strict — no optional/nil deps) and update `main.go`. Add the
in-memory `ChunkRepository` (owner-scoped map, keyed by owner then id) and register
`apiV1.Post("/chunks", func(ctx){ loginRequired(ctx, chunk.Create) })`. The
`apptest` builder defaults the chunk repo to in-memory. Finish: `make check` GREEN.

## Design decisions to hydrate into design.md
- [ ] New `vault` module; `EncryptedChunk` entity is minimal (`id`, `ownerID`,
  `content`) — `version`/`updated_at` deferred to the update/sync features that need them.
- [ ] Client-generated **UUID** id (v7 by client convention) — validated at the
  boundary so a future Postgres native `uuid` column stays valid and gets v7 insert
  locality. Server treats content as one opaque blob (`nonce‖ciphertext‖tag`), never inspected.
- [ ] Ownership is **owner-scoped** everywhere: `Exists(ownerID, id)` (a bool — no
  fetch of the content blob, explicit intent; a fetch method is added by the read
  feature when it needs one); duplicate
  check is per-user (not global) — isolation + no cross-user id oracle. Duplicate → 409 `AlreadyExists`.
- [ ] Create is auth-required (`loginRequired`); owner comes from the request
  context, never the body. Route `POST /api/v1/chunks`. Response `201 {id}` (id is client-known; pure ack).
- [ ] `NewFiberApplication` stays strict (all deps required; `chunkRepository`
  added) — production never builds with a nil repo. Test churn is solved in the
  tests: a shared `apptest` builder is the single caller, so future constructor
  changes touch only it. Chunk uses an in-memory adapter; data lost on restart.
- [ ] Quality Pillars — Security (auth-required, owner from context not body, opaque
  content, `strings.Clone`, per-user isolation), Reliability (duplicate rejected
  cleanly, repo errors propagate), Performance/Observability (deferred — in-memory, no telemetry).
