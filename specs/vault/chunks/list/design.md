# Technical Design: List Chunks

**Corresponds to Spec:** `specs/vault/chunks/list/spec.md`

## Overview

The `vault` module stores a user's financial data as opaque, client-encrypted
`EncryptedChunk`s. This feature returns every chunk a logged-in user owns in one
request, each exactly as stored (`id` + opaque `content`), ordered newest-first.
It is the base full-fetch primitive that the future delta-sync and pagination
features refine — the one-time bootstrap a fresh device uses to rebuild its local
copy of the vault. The server returns the blobs verbatim and never reads them.

## Design Decisions & Rationale

- **Owner-scoped bulk fetch: `ChunkRepository.GetAll(ownerID) → ([]EncryptedChunk,
  error)`.** Returns only the requester's chunks; another user's data is never
  reachable. Added alongside `Exists`/`Get`; the collection read needs its own
  method rather than overloading the single-id `Get`.
- **Empty is a success, not a `NotFound`.** An owner with no chunks yields an
  empty slice and a nil error — distinct from single-read, where a miss is a
  `NotFound` (404). Listing an empty vault is a valid state (a brand-new user), so
  it returns `200 {"chunks": []}`, never an error.
- **Newest-first via id descending, and ordering is the repository's job.** The
  entity carries no timestamp — `updated_at` was deliberately deferred to the sync
  feature — so the only field to order by today is the id. The UUIDv7 client
  convention makes descending id ≈ most-recently-created-first. The in-memory
  adapter sorts the slice by id descending; a future Postgres adapter does `ORDER
  BY id DESC`. The server cannot enforce v7, so this ordering is correct only as
  far as clients honor the convention; when sync introduces a real created/updated
  timestamp, ordering should move to it. Sorting lives in the adapter so the use
  case stays a pure passthrough.
- **Response is a wrapper object `{"chunks": [...]}`, not a bare array.** This is
  the one shape deviation from single-read's bare `{id, content}`. Pagination and
  delta sync are the immediate next features and both need to add top-level fields
  (a page cursor, a sync token); a bare top-level array cannot grow without a
  breaking change, a wrapper can. Each item reuses the single-read item shape
  (`id` + `content`).
- **`ChunkLister` is a thin passthrough to `GetAll`.** It mirrors `ChunkGetter`
  and is kept for layer symmetry and as a future home for list-level logic, even
  though it currently only forwards the call.
- **Auth-required; owner from the request context, never input.** The route sits
  behind `loginRequired`; the handler reads the owner from `requestcontext`. There
  is no body and no path param, so nothing is parsed from the Fiber request and no
  `strings.Clone` is needed. `NewFiberApplication`'s signature is unchanged — the
  route reuses the existing `chunkRepository` dependency.

## Architecture & Files Summary

```
src/vault/domain/
└── repositories/chunk.go            # ChunkRepository port: GetAll added (owner-scoped), alongside Exists/Add/Get

src/vault/application/use_cases/chunklister/
└── chunk_lister.go                  # ChunkLister: passthrough to GetAll

src/vault/infrastructure/
├── api/chunkhandler/chunk_handler.go        # List method + listChunksResponse wrapper (reuses getChunkResponse item)
└── repositories/inmemory/chunk_repository.go # GetAll: owner-scoped gather, empty slice on none, sorted id desc

src/app/fiber_application.go         # GET /api/v1/chunks (loginRequired); constructor unchanged

tests/unit/vault/application/use_cases/list_test.go
tests/unit/vault/infrastructure/api/chunk_api_test/list_test.go
tests/integration/vault/chunk_test.go        # list: newest-first ordering, owner isolation, empty collection
tests/unit/mocks/mock_ChunkRepository.go      # includes GetAll

.bruno/chunks/                       # list chunks, list chunks empty, list chunks unauthorized
cmd/odin-cli/main.go                 # list-chunks command: fetch all + decrypt each (list round-trip)

specs/vault/chunks/list/
├── spec.md
├── design.md
└── plan.md          # current work order (see plan-template.md)
```

## Data Flow

**List (`GET /api/v1/chunks`, auth-required):** `loginRequired` admits the
request. The handler reads the owner from the request context, builds
`ChunkLister`, and calls `List`, which delegates to
`chunkRepository.GetAll(ownerID)`. The repository gathers the owner's chunks
(empty slice if the owner has none), sorts them by id descending, and returns
them. The handler maps each chunk to the item shape and responds
`200 {"chunks": [...]}`. Repository errors propagate to Fiber's global
`errorHandler`.

## Request & Response

**List** — `GET /api/v1/chunks` (Bearer token required)
```json
// success 200
{ "chunks": [ { "id": "<uuid>", "content": "<opaque base64 nonce‖ciphertext‖tag>" }, ... ] }
// owner with no chunks 200
{ "chunks": [] }
// no / invalid bearer token 401  (empty body, from loginRequired)
```

Order: newest-first (id descending). Web/HTMX: N/A — REST-only (client retrieves
and reconstructs locally).

## Known Limitations

- **Unbounded response** — this returns every chunk the owner has, in one payload.
  Safe now (in-memory store, small datasets, CLI testing), but for a large vault
  it is heavy; **pagination is the next feature** and exists precisely to bound
  this. Steady-state clients avoid re-fetching everything by persisting locally and
  using delta sync; this full list is the one-time cold bootstrap.
- **Ordering depends on the UUIDv7 client convention** — the server cannot enforce
  v7, so newest-first is only as correct as the client's id generation. A real
  created/updated timestamp (arriving with sync) should replace id-ordering.
- **In-memory repository** — data lost on restart; the persistent adapter is TBD.
  `GetAll` gathers and sorts in memory with no mutex, like the other in-memory
  repos.

## Quality Pillars

- **Security:** auth-required (`loginRequired`); owner taken from the validated
  request context, never input; strict per-user isolation — only the owner's
  chunks are gathered, no cross-user data ever returned; content returned opaque
  and never inspected.
- **Reliability:** an empty vault is handled as a success (`{"chunks": []}`), not
  an error; repository errors propagate to the global handler; ordering is
  deterministic (id descending).
- **Performance:** a single owner-scoped gather plus an in-memory sort. Unbounded
  for now — bounding it is the explicit job of the upcoming pagination feature, so
  performance work is deferred there rather than duplicated here.
- **Observability:** Deferred — no production telemetry yet; `odinerrors` location
  tracking aids debugging, consistent with the rest of the module.
