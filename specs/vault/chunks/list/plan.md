# Work Order: List Chunks — new feature

**Feature design:** `specs/vault/chunks/list/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/vault/chunks/list/spec.md`

> Work order for: **building bulk retrieval of all a user's chunks**. Disposable —
> overwritten by the next change (git keeps the history). The living design is in
> design.md; hydrate it before this change merges, then freeze this file.

## Change

Let a logged-in user retrieve every `EncryptedChunk` they own in one request,
each returned exactly as stored (`id` + opaque `content`), ordered newest-first.
This is the base full-fetch primitive that the future delta-sync and pagination
features refine; it is the one-time bootstrap a fresh device uses to rebuild its
local copy.

Retrieval is owner-scoped: only the requester's own chunks come back, never
another user's. An owner with no chunks receives an empty collection with a
success status — empty is not an error. Ordering is newest-first, achieved by
sorting on the id descending: the entity carries no timestamp (deferred to the
sync feature), and the UUIDv7 client convention makes descending id ≈ most
recently created first. When sync adds a real created/updated timestamp, ordering
should move to it.

This change is additive to the read feature and needs no new plumbing. It adds a
`GetAll` method to the port + adapter, a thin `ChunkLister` use case, a `List`
method on the existing `chunkHandler`, and a route. `NewFiberApplication`'s
signature does NOT change. Response is a wrapper object `{"chunks": [...]}` (not a
bare array) so the coming pagination/delta features can add top-level fields
(cursor, sync token) without a breaking change.

Satisfies spec scenarios: *Listing my items successfully*, *Listing when I have
no items*, *Listing returns only my own items*, *Listing without being logged in*.

## Architecture & Files (this change)
```
src/vault/domain/
└── repositories/chunk.go                                 # MODIFY  add GetAll to ChunkRepository port

src/vault/application/use_cases/chunklister/
└── chunk_lister.go                                       # CREATE  list use case

src/vault/infrastructure/
├── api/chunkhandler/chunk_handler.go                     # MODIFY  add List method + listChunksResponse (+ reuse chunk item shape)
└── repositories/inmemory/chunk_repository.go             # MODIFY  implement GetAll (owner-scoped, sorted id desc, empty slice on none)

src/app/fiber_application.go                              # MODIFY  route GET /api/v1/chunks (loginRequired)

tests/unit/mocks/mock_ChunkRepository.go                  # REGEN   (GetAll added to port)

tests/unit/vault/application/use_cases/list_test.go       # CREATE
tests/unit/vault/infrastructure/api/chunk_api_test/list_test.go   # CREATE
tests/integration/vault/chunk_test.go                     # MODIFY  add list acceptance scenario(s)
```

## Key Types & Signatures

```go
// src/vault/domain/repositories — owner-scoped bulk fetch added; Exists/Add/Get
// unchanged. GetAll returns the owner's chunks newest-first (id descending); an
// owner with no chunks yields an empty slice and a nil error (NOT NotFound).
type ChunkRepository interface {
    Exists(ctx context.Context, ownerID, id string) (bool, error)
    Add(ctx context.Context, chunk *chunkmodel.EncryptedChunk) error
    Get(ctx context.Context, ownerID, id string) (*chunkmodel.EncryptedChunk, error)
    GetAll(ctx context.Context, ownerID string) ([]*chunkmodel.EncryptedChunk, error)
}

// src/vault/application/use_cases/chunklister — thin, mirrors chunkgetter.
func New(ownerID string, chunkRepository repositories.ChunkRepository) ChunkLister
func (self ChunkLister) List(ctx context.Context) ([]*chunkmodel.EncryptedChunk, error)

// src/vault/infrastructure/api/chunkhandler — method on the EXISTING chunkHandler.
func (self chunkHandler) List(ctx *fiber.Ctx) error   // owner from requestcontext; 200 wrapper
type listChunksResponse struct { Chunks []getChunkResponse `json:"chunks"` }
```

Ordering is the repository's responsibility: the in-memory adapter sorts the
owner's chunks by id descending before returning; a future Postgres adapter does
`ORDER BY id DESC`. The use case returns what the repository gives. Success →
`200 {"chunks": [...]}`; empty owner → `200 {"chunks": []}`. Unauthenticated →
401 via `loginRequired` (no body).

## Implementation Phases (TDD)

Double-loop: Phase 0 adds the failing happy-path list acceptance test first; it
stays RED until Phase 3 wires the route. Inner phases go green in order.
`make check` green is the exit criterion.

### Phase 0: Acceptance (outer loop) — failing happy-path integration test
**Red:** in `tests/integration/vault/chunk_test.go`, add a scenario: an
authenticated owner stores several chunks via `POST /api/v1/chunks`, then
`GET /api/v1/chunks` → `200 {"chunks": [...]}` containing all of them, ordered
newest-first (id descending), each with its stored id+content. Add the isolation
scenario: owner A stores chunks, owner B `GET /api/v1/chunks` → only B's chunks
(and B with none → `{"chunks": []}`). Fails now (no route). Do not touch again
until Phase 3.
**Green:** achieved by Phases 1–3.

### Phase 1: Domain — port + mock
**Red:** none new at this layer beyond the mock the app/handler tests need.
**Green:** add `GetAll(ctx, ownerID) ([]*chunkmodel.EncryptedChunk, error)` to the
`ChunkRepository` port; regenerate the mock (`make mocks`) so `MockChunkRepository`
has `GetAll` for the use-case and handler tests.

### Phase 2: Application — chunklister use case
**Red:** `list_test.go` (MockChunkRepository): owner with chunks (`GetAll(owner)`
→ slice) → returns that slice unchanged; owner with none (`GetAll` → empty slice)
→ returns empty slice, no error; a repo error propagates.
**Green:** implement `ChunkLister`: `New(ownerID, chunkRepository)` and `List(ctx)`
delegating to `chunkRepository.GetAll(ctx, ownerID)`, returning its `(slice, err)`
unchanged.

### Phase 3: Infrastructure — inmemory GetAll, handler List, route
**Red:** handler `list_test.go` (authenticated via RequestBuilder;
MockChunkRepository): several chunks → `200` `{"chunks":[...]}` in newest-first
order, each {id, content} matching stored; empty → `200` `{"chunks":[]}`;
unauthenticated (anonymous session) → `401`. Repo-level assertion (in the
inmemory test or the integration test): `GetAll` returns only the owner's chunks,
sorted id descending, and an empty slice for an unknown owner.
**Green:**
- Implement `InMemoryChunkRepository.GetAll`: gather `chunks[ownerID]` values into
  a slice (empty slice if the owner has no map entry — NOT NotFound), sort by
  `ID()` descending, return it.
- Add `chunkHandler.List`: read owner via
  `ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext).UserID()`,
  build `chunklister.New`, call `List`, map each chunk to `getChunkResponse`, and
  respond `200 listChunksResponse{Chunks: ...}`; errors returned to the global
  `errorHandler`.
- Register `apiV1.Get("/chunks", func(ctx){ loginRequired(ctx, chunk.List) })`.
  Fiber distinguishes it from `GET /chunks/:id`. The Phase 0 acceptance tests now
  pass. Finish: `make check` GREEN.

## Design decisions to hydrate into design.md
- [ ] List is auth-required (`loginRequired`); owner from the request context;
  route `GET /api/v1/chunks`; success `200 {"chunks":[{id, content}, ...]}` (each
  chunk returned exactly as stored).
- [ ] `ChunkRepository.GetAll(ownerID)` added — **owner-scoped bulk fetch** that
  returns only the owner's chunks; an owner with none yields an **empty slice +
  nil error** (empty is success, NOT NotFound). `Exists`/`Get` unchanged.
- [ ] **Ordering is newest-first via id descending**, and ordering is the
  repository's responsibility (in-memory sort; future Postgres `ORDER BY id
  DESC`). Relies on the UUIDv7 client convention because the entity has no
  timestamp yet; when sync adds a created/updated timestamp, ordering moves to it.
- [ ] Response is a **wrapper object `{"chunks": [...]}`**, not a bare array —
  chosen so the coming pagination/delta features can add top-level fields (cursor,
  sync token) without breaking the shape. This is the base full-fetch primitive
  those features refine; it is the one-time device bootstrap.
- [ ] `ChunkLister` use case is a thin passthrough to `GetAll` (mirrors
  `ChunkGetter`), kept for layer symmetry and a future home for list logic.
- [ ] Quality Pillars — Security (auth-required, owner from context, per-user
  isolation — only the owner's chunks, content returned opaque), Reliability
  (empty handled as success, repo errors propagate), Performance (single
  owner-scoped gather + sort; unbounded for now — pagination is the next feature),
  Observability (deferred — in-memory, no telemetry).
