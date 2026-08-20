# Technical Design: Create Chunk

**Corresponds to Spec:** `specs/vault/chunks/create/spec.md`

## Overview

The `vault` module stores a user's financial data as opaque, client-encrypted
records — `EncryptedChunk`s. This feature is the first operation: a logged-in user
stores one chunk. The server persists it tied to its owner and never inspects the
content. Each domain record (an account, an expense…) becomes one chunk; the whole
relational model lives inside the encrypted payload and is reconstructed only on
the client after decryption.

## Design Decisions & Rationale

- **`EncryptedChunk` is minimal: `id`, `ownerID`, `content`.** `version`
  (optimistic locking) and `updated_at` (delta sync) are deliberately deferred to
  the update/sync features that need them — the entity grows with its features
  rather than speculating now. Fields are private; UUID validation lives in the
  constructor.
- **Client-generated, UUID-validated `id`.** Offline-first clients need a stable id
  at creation time (to store locally and to reference between records) before any
  server round-trip, and it makes retries idempotent. The id is validated as a
  well-formed UUID at the boundary so a future Postgres native `uuid` column stays
  valid; UUIDv7 is the client convention (time-ordered → good index locality). The
  server can't cleanly enforce the version, so it only parse-checks.
- **`content` is one opaque blob.** The client packs `nonce‖ciphertext‖tag`
  (AES-256-GCM) into a single base64 string. The server stores and returns it
  verbatim and never separates or reads it — splitting IV/tag server-side would
  buy nothing.
- **Existence check via `Exists(ownerID, id) bool`, not a fetch.** A duplicate
  check needs a boolean, not the (potentially large) content blob — cheaper and
  explicit intent. A fetch method is added by the read feature when it actually
  returns data. Chosen over `GetByID` returning `(nil, nil)` because a bool has no
  nil/NotFound ambiguity.
- **Ownership is owner-scoped everywhere.** `Exists` takes the owner; uniqueness is
  per-user, never global — full data isolation and no cross-user id oracle.
  Duplicate id → `AlreadyExists` (→ 409); a create never overwrites (that's
  update's job).
- **Auth-required; owner from the request context, never the body.** The route sits
  behind `loginRequired`; the handler reads the owner from
  `requestcontext.UserID()`. A client cannot spoof ownership by sending an owner in
  the payload. Response is `201 {id}` — a pure acknowledgement, since the client
  already knows the id.
- **`NewFiberApplication` stays strict/positional** (all deps required — production
  never builds with a nil repo) but gained a required `chunkRepository`. To stop
  constructor changes from rippling across tests, a shared `apptest` builder is now
  the single caller of the constructor; future modules touch only it.

## Architecture & Files Summary

```
src/vault/domain/
├── chunkmodel/chunk.go            # EncryptedChunk entity (id/ownerID/content, UUID-validated)
└── repositories/chunk.go          # ChunkRepository port: Exists, Add (owner-scoped)

src/vault/application/use_cases/chunkcreator/
└── chunk_creator.go               # duplicate check (Exists) → create → persist

src/vault/infrastructure/
├── api/chunkhandler/
│   ├── chunk_handler.go           # handler + createChunkResponse; owner from request context
│   └── body.go                    # CreateChunkBody + Validate (presence)
└── repositories/inmemory/chunk_repository.go   # owner→id map

src/app/fiber_application.go       # POST /api/v1/chunks (loginRequired) + chunkRepository dep
src/main.go                        # wires the in-memory chunk repo

tests/builders/apptest/app.go      # single caller of NewFiberApplication (in-memory defaults + overrides)
tests/unit/vault/...               # domain, application, body, handler tests
tests/integration/vault/chunk_test.go

specs/vault/chunks/create/
├── spec.md
├── design.md
└── plan.md
```

## Data Flow

**Create (`POST /api/v1/chunks`, auth-required):** `loginRequired` admits the
request; the handler reads the owner from the request context, parses the JSON
body, and delegates presence validation to `CreateChunkBody.Validate()`
(`strings.Clone` on parsed strings). It constructs `ChunkCreator` and calls
`Create`, which: `Exists(ownerID, id)` → true ⇒ `AlreadyExists` (409); else
`chunkmodel.New` (UUID validation → Domain/400 on a bad id) → `Add`. On success the
handler responds `201 {id}`. On failure the error returns to Fiber's global
`errorHandler`, which maps the odin tag to a status and renders `{"error": …}`.

## Request & Response

**Create** — `POST /api/v1/chunks` (Bearer token required)
```json
// request
{ "id": "<client-uuid-v7>", "content": "<opaque base64 nonce‖ciphertext‖tag>" }
// success 201
{ "id": "<client-uuid-v7>" }
// duplicate id 409
{ "error": "El elemento ya existe" }
// missing/invalid id or content, or malformed body 400
{ "error": "Datos de solicitud inválidos" }
// no / invalid bearer token 401  (empty body, from loginRequired)
```

## Known Limitations

- **In-memory repository** — data lost on restart; storage adapter TBD. Concurrent
  registrations/creates on the shared map have no mutex (same as the auth repos);
  duplicate safety under concurrency will rely on the future store's unique
  constraint.
- **No content size limit** — the server accepts an arbitrarily large opaque blob;
  a cap belongs with the persistent store / transport limits later.
- **id trust** — the server parse-checks the id is a UUID but does not enforce v7;
  ordering benefits depend on the client honoring the convention.

## Quality Pillars

- **Security:** auth-required (`loginRequired`); owner taken from the validated
  request context, never the body, so ownership can't be spoofed; content stored
  opaquely and never inspected; per-user isolation (owner-scoped `Exists`, no
  cross-user oracle); `strings.Clone` on parsed body values.
- **Reliability:** duplicate id rejected cleanly (409); repository errors propagate
  to the global handler; invalid UUID rejected before persistence.
- **Performance:** existence check returns a bool, not the blob. Otherwise deferred
  — in-memory store, no hot path in this change.
- **Observability:** Deferred — no production telemetry yet; `odinerrors` location
  tracking aids debugging.
