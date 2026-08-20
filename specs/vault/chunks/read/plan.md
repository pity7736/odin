# Work Order: Read Chunk — new feature

**Feature design:** `specs/vault/chunks/read/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/vault/chunks/read/spec.md`

> Work order for: **building chunk retrieval by id**. Disposable — overwritten by
> the next change (git keeps the history). The living design is in design.md;
> hydrate it before this change merges, then freeze this file.

## Change

Let a logged-in user retrieve one stored `EncryptedChunk` by its id. The client
sends the id in the path; the server returns the chunk exactly as stored (`id` +
opaque `content`) so the client can decrypt it. Retrieval is **owner-scoped**: the
lookup is keyed by owner, so a chunk that does not exist, is malformed, or belongs
to another user is all indistinguishable — each yields *not found* (404), giving
full isolation and no cross-user id oracle.

This change is additive to the create feature. The plumbing already exists: the
`NotFound`→404 mapping, `loginRequired`, the request-context owner, the `apptest`
builder, and the `chunkRepository` constructor param. This adds a fetch method to
the port + adapter, a thin `chunkgetter` use case, a `Get` method on the existing
`chunkHandler`, and a route. `NewFiberApplication`'s signature does NOT change.

The repository signals absence with an explicit `NotFound` odin error, mirroring
the user repository (commit 67727d7) — not `(nil, nil)`, and not the `Exists`
bool. The `chunkgetter` use case propagates that `NotFound` verbatim (unlike
login, which translates `NotFound`→`Unauthorized`), since here not-found is the
desired outcome. The read path does NOT UUID-parse-check the id: a malformed id
simply misses the owner-scoped lookup → `NotFound`, matching the spec.

Satisfies spec scenarios: *Retrieving my item successfully*, *Retrieving without
being logged in*, *Retrieving an item that does not exist*, *Retrieving an item
that belongs to another user*, *Retrieving with an invalid identifier*.

## Architecture & Files (this change)
```
src/vault/domain/
└── repositories/chunk.go                                 # MODIFY  add Get to ChunkRepository port

src/vault/application/use_cases/chunkgetter/
└── chunk_getter.go                                       # CREATE  read use case

src/vault/infrastructure/
├── api/chunkhandler/chunk_handler.go                     # MODIFY  add Get method + getChunkResponse
└── repositories/inmemory/chunk_repository.go             # MODIFY  implement Get (NotFound on miss)

src/app/fiber_application.go                              # MODIFY  route GET /api/v1/chunks/:id (loginRequired)

tests/unit/mocks/mock_ChunkRepository.go                  # REGEN   (Get added to port)

tests/unit/vault/application/use_cases/read_test.go       # CREATE
tests/unit/vault/infrastructure/api/chunk_api_test/read_test.go   # CREATE
tests/integration/vault/chunk_test.go                     # MODIFY  add read acceptance scenario(s)
```

## Key Types & Signatures

```go
// src/vault/domain/repositories — owner-scoped fetch added; Exists/Add unchanged.
// Get returns a NotFound odin error when the owner has no chunk with that id.
type ChunkRepository interface {
    Exists(ctx context.Context, ownerID, id string) (bool, error)
    Add(ctx context.Context, chunk *chunkmodel.EncryptedChunk) error
    Get(ctx context.Context, ownerID, id string) (*chunkmodel.EncryptedChunk, error)
}

// src/vault/application/use_cases/chunkgetter — thin, mirrors chunkcreator.
func New(id, ownerID string, chunkRepository repositories.ChunkRepository) ChunkGetter
func (self ChunkGetter) Get(ctx context.Context) (*chunkmodel.EncryptedChunk, error)

// src/vault/infrastructure/api/chunkhandler — method on the EXISTING chunkHandler.
func (self chunkHandler) Get(ctx *fiber.Ctx) error   // owner from requestcontext, id from ctx.Params("id"), strings.Clone
type getChunkResponse struct { ID string `json:"id"`; Content string `json:"content"` }
```

The read path performs NO UUID validation (invalid id → NotFound, not Domain/400).
Absence external message (Spanish): `"El elemento no existe"`, internal English
`"chunk not found"`, tag `NotFound` (→ 404). Unauthenticated → 401 via
`loginRequired` (no body). Success → `200 {id, content}`.

## Implementation Phases (TDD)

Double-loop: Phase 0 adds the failing happy-path read acceptance test first; it
stays RED until Phase 3 wires the route. Inner phases go green in order.
`make check` green is the exit criterion.

### Phase 0: Acceptance (outer loop) — failing happy-path integration test
**Red:** in `tests/integration/vault/chunk_test.go`, add a scenario: an
authenticated `POST /api/v1/chunks` then `GET /api/v1/chunks/:id` for the same
owner → `200 {id, content}` matching what was stored. Fails now (no route). Also
add the cross-owner scenario: owner A stores a chunk, owner B `GET`s that id →
`404`. Do not touch again until Phase 3.
**Green:** achieved by Phases 1–3.

### Phase 1: Domain — port + mock
**Red:** none new at this layer beyond the mock the app/handler tests need.
**Green:** add `Get(ctx, ownerID, id) (*chunkmodel.EncryptedChunk, error)` to the
`ChunkRepository` port; regenerate the mock (`go run github.com/vektra/mockery/v3`)
so `MockChunkRepository` has `Get` for the use-case and handler tests.

### Phase 2: Application — chunkgetter use case
**Red:** `read_test.go` (MockChunkRepository): existing chunk (`Get(owner,id)` →
chunk) → returns that chunk; not found (`Get` → `NotFound` error) → propagates the
`NotFound` error verbatim (tag preserved); a non-NotFound repo error also
propagates.
**Green:** implement `ChunkGetter`: `New(id, ownerID, chunkRepository)` and
`Get(ctx)` delegating to `chunkRepository.Get(ctx, ownerID, id)`, returning its
`(chunk, err)` unchanged.

### Phase 3: Infrastructure — inmemory Get, handler Get, route
**Red:** handler `read_test.go` (authenticated via RequestBuilder;
MockChunkRepository): success → `200` `{id, content}` matching stored; not found
(`Get` → NotFound) → `404` "El elemento no existe"; a chunk owned by another user
(repo returns NotFound for this owner) → `404`; a non-UUID / malformed id → `404`
(no 400 — assert the handler passes the raw id straight to `Get` without
validating); unauthenticated (anonymous session) → `401`.
**Green:**
- Implement `InMemoryChunkRepository.Get`: look up `chunks[ownerID][id]`; on miss
  return a `NotFound` odin error (external `"El elemento no existe"`, internal
  `"chunk not found"`), mirroring the user repo.
- Add `chunkHandler.Get`: read owner via
  `ctx.Locals(requestcontext.Key).(*requestcontext.RequestContext).UserID()`,
  read `strings.Clone(ctx.Params("id"))`, build `chunkgetter.New`, call `Get`, on
  success respond `200 getChunkResponse{id, content}`; errors returned to the
  global `errorHandler`.
- Register `apiV1.Get("/chunks/:id", func(ctx){ loginRequired(ctx, chunk.Get) })`.
  The Phase 0 acceptance tests now pass. Finish: `make check` GREEN.

## Design decisions to hydrate into design.md
- [ ] Read is auth-required (`loginRequired`); owner from the request context;
  route `GET /api/v1/chunks/:id`; success `200 {id, content}` (chunk returned
  exactly as stored).
- [ ] `ChunkRepository.Get(ownerID, id)` added — **owner-scoped fetch** that
  signals absence with an explicit `NotFound` odin error (mirrors the user
  repository, 67727d7), chosen over `(nil, nil)` for no nil/NotFound ambiguity.
  `Exists` stays for create's duplicate check.
- [ ] Not-found is **indistinguishable across nonexistent / malformed / other-user
  ids** — all 404, no cross-user id oracle, full isolation (owner-scoped lookup).
- [ ] Read does **NOT** UUID-validate the id (unlike create's 400 on a bad id): a
  malformed id misses the lookup → 404, per spec.
- [ ] `chunkgetter` use case is a thin passthrough that **propagates** the repo's
  `NotFound` verbatim (contrast: login translates `NotFound`→`Unauthorized`).
- [ ] The `NotFound` external Spanish message (`"El elemento no existe"`) is minted
  in the infrastructure repo, consistent with the user repository (message
  ownership location is a known codebase convention, not relitigated here).
- [ ] Quality Pillars — Security (auth-required, owner from context, per-user
  isolation / no oracle, `strings.Clone` on the path param, content returned
  opaque), Reliability (NotFound mapped cleanly to 404, repo errors propagate),
  Performance/Observability (deferred — in-memory, no telemetry).
