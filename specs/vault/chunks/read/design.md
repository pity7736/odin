# Technical Design: Read Chunk

**Corresponds to Spec:** `specs/vault/chunks/read/spec.md`

## Overview

The `vault` module stores a user's financial data as opaque, client-encrypted
`EncryptedChunk`s. This feature is the read counterpart to create: a logged-in
user retrieves one chunk by its id and gets it back exactly as stored (`id` +
opaque `content`) so the client can decrypt it. The server returns the content
verbatim and never inspects it. Retrieval is owner-scoped, so a chunk is
reachable only by the user that created it.

## Design Decisions & Rationale

- **Owner-scoped fetch that signals absence explicitly:
  `ChunkRepository.Get(ownerID, id) → (*EncryptedChunk, error)`.** On a miss the
  repository returns a `NotFound` odin error, not `(nil, nil)` — a nil/NotFound
  bool has no ambiguity and mirrors the user repository precedent (the auth
  feature's `GetByEmail`). `Exists` (a bool, used by create's duplicate check)
  stays; `Get` is added because read actually needs the content blob. The two
  coexist by intent: `Exists` answers "is it there?", `Get` returns the data.
- **Not-found is uniform across nonexistent / malformed / other-user ids — all
  404.** The lookup is keyed by owner, so another user's chunk is simply absent
  from this owner's space and yields the same `NotFound` as a nonexistent id.
  This gives full per-user isolation and, deliberately, no cross-user id oracle:
  the response never distinguishes "exists but not yours" from "does not exist".
- **The read path does NOT validate the id as a UUID.** Create parse-checks the
  id and rejects a malformed one as `Domain`/400; read does the opposite — a
  malformed id is passed straight to the owner-scoped lookup, misses, and returns
  404. Validating on read would leak a "well-formed but absent" vs "malformed"
  distinction and contradict the spec, which treats an invalid id as not found.
- **`ChunkGetter` propagates the repository's `NotFound` verbatim.** The use case
  is a thin passthrough to `Get`; it does not translate the error. This contrasts
  with login's `SessionStarter`, which turns a `NotFound` from the user lookup
  into an `Unauthorized` (to avoid a user-existence oracle). For read, not-found
  IS the desired outcome, so the `NotFound` tag flows to the global error handler
  and maps to 404. The use case is kept (rather than calling the repo from the
  handler) for symmetry with `ChunkCreator` and as a home for future read logic.
- **The `NotFound` external message is minted in the infrastructure repository.**
  The Spanish user-facing string (`"El elemento no existe"`) lives in the
  in-memory adapter, consistent with the user repository. Message-ownership
  location is an existing codebase convention; this feature follows it rather
  than relocating it.
- **Auth-required; owner from the request context, never the input.** The route
  sits behind `loginRequired`; the handler reads the owner from
  `requestcontext.UserID()` and the id from the path param (`strings.Clone`d
  against fasthttp buffer reuse). A client cannot reach another user's data by
  supplying an owner. `NewFiberApplication`'s signature is unchanged — the read
  route reuses the `chunkRepository` dependency the create feature already added.
- **Success returns the full stored chunk, `200 {id, content}`.** Read returns
  the content (unlike create's `201 {id}` acknowledgement) because the client's
  whole purpose is to decrypt it.

## Architecture & Files Summary

```
src/vault/domain/
└── repositories/chunk.go            # ChunkRepository port: Get added (owner-scoped), alongside Exists/Add

src/vault/application/use_cases/chunkgetter/
└── chunk_getter.go                  # ChunkGetter: passthrough to Get, propagates NotFound verbatim

src/vault/infrastructure/
├── api/chunkhandler/chunk_handler.go        # Get method + getChunkResponse; owner from context, id from path
└── repositories/inmemory/chunk_repository.go # Get: owner-scoped lookup, NotFound on miss

src/app/fiber_application.go         # GET /api/v1/chunks/:id (loginRequired); constructor unchanged

tests/unit/vault/application/use_cases/read_test.go
tests/unit/vault/infrastructure/api/chunk_api_test/read_test.go
tests/integration/vault/chunk_test.go        # read round-trip + cross-owner isolation
tests/unit/mocks/mock_ChunkRepository.go      # includes Get

.bruno/chunks/                       # get chunk, not found, invalid id, unauthorized
cmd/odin-cli/main.go                 # get-chunk command: fetch + decrypt (read round-trip)

specs/vault/chunks/read/
├── spec.md
├── design.md
└── plan.md          # current work order (see plan-template.md)
```

## Data Flow

**Read (`GET /api/v1/chunks/:id`, auth-required):** `loginRequired` admits the
request. The handler reads the owner from the request context and the id from the
path param (`strings.Clone`), builds `ChunkGetter`, and calls `Get`, which
delegates to `chunkRepository.Get(ownerID, id)`. A hit returns the chunk → the
handler responds `200 {id, content}`. A miss returns a `NotFound` odin error,
propagated verbatim to Fiber's global `errorHandler`, which maps the tag to 404
and renders `{"error": …}`. No UUID validation happens anywhere on this path.

## Request & Response

**Read** — `GET /api/v1/chunks/:id` (Bearer token required)
```json
// success 200
{ "id": "<uuid>", "content": "<opaque base64 nonce‖ciphertext‖tag>" }
// not found / another user's chunk / malformed id — all 404
{ "error": "El elemento no existe" }
// no / invalid bearer token 401  (empty body, from loginRequired)
```

Web/HTMX: N/A — this operation is REST-only (client retrieves and decrypts).

## Known Limitations

- **In-memory repository** — data lost on restart; the persistent storage adapter
  is TBD and owned by its own feature. Read is a single map lookup with no mutex
  (same as the other in-memory repos).
- **No cross-user oracle by design** — the trade-off is that a legitimate owner
  cannot distinguish "I mistyped the id" from "that chunk was never created"; both
  are 404. This is intentional and preferred over leaking existence.

## Quality Pillars

- **Security:** auth-required (`loginRequired`); owner taken from the validated
  request context, never client input; per-user isolation via owner-scoped lookup
  with no cross-user id oracle (nonexistent / malformed / other-user all 404);
  content returned opaque and never inspected; `strings.Clone` on the path param.
- **Reliability:** absence mapped cleanly to 404 via the `NotFound` tag; non-
  NotFound repository errors propagate to the global handler; no UUID validation
  means no false 400s on the read path.
- **Performance:** a single owner-scoped map lookup; no fetch-then-check. Otherwise
  deferred — in-memory store, no hot path in this change.
- **Observability:** Deferred — no production telemetry yet; `odinerrors` location
  tracking aids debugging, consistent with the rest of the module.
