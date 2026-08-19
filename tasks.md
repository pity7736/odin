# Tasks

Odin is a zero-knowledge, E2E-encrypted personal finance app.
Two repos: `odin` (Go sync/auth server), `odin-android` (Kotlin finance app).

## Completed — Tooling

- [x] Upgrade Go to 1.26 and update all dependencies
- [x] Set up `.golangci.yml` with linters adapted from Gideon
- [x] Create `Makefile` with `lint`, `test`, `mocks`, `coverage`, `coverage-check`, `check` targets
- [x] Set up test coverage measurement and establish baseline threshold (91.9%)
- [x] Pre-commit hook (lint + test)
- [x] Fix `.mockery.yaml` to use mockery v3 config syntax

## Phase 1 — Server pivot (Go, `odin`)

Strip the current server down to a sync/auth API. New domain: User, EncryptedChunk, KeyParams, SyncToken.

- [x] Auth redesign: key derivation support (Argon2id params), encrypted master key storage, double-hash chain (client Argon2id → server bcrypt), app-scoped handlers, centralized error handling, self-validating request bodies
- [x] Remove old domain code (Account, Category, Income entities, use cases, repositories, HTMX handlers/templates)
- [ ] CLI + login verification: build a small CLI (`cmd/odin-cli`) that does client-side crypto (Argon2id key derivation), derives the auth hash for the seeded user, calls the login endpoint, and verifies it gets back the encrypted master key and key params. First real proof the crypto handshake works. The CLI grows with each subsequent task.
- [ ] Registration: accept email, auth_hash, encrypted_master_key, and key_params from the client. Server bcrypt-hashes the auth_hash and persists the user. Remove seed data — the CLI registers its own users from this point on. Full SDD cycle.
- [ ] Create encrypted chunk: the CLI encrypts plaintext with AES-256-GCM using the master key, sends the blob to a new server endpoint, reads it back, decrypts, and verifies the plaintext matches. Server-side: accept a client-encrypted blob (ciphertext + IV + auth tag) with a client-generated ID, persist it. This is where the EncryptedChunk entity gets designed — you can't create what you haven't defined. Full SDD cycle.
- [ ] Read encrypted chunk: retrieve a single encrypted chunk by ID for its owner. Server returns the opaque blob as-is — no decryption, no interpretation. CLI decrypts and verifies the round-trip.
- [ ] Update encrypted chunk: replace an existing chunk's ciphertext with a new version. Uses optimistic locking (version field) so concurrent edits from multiple devices fail cleanly instead of silently overwriting each other. CLI tests both paths: successful update with correct version, and rejection with stale version.
- [ ] Delete encrypted chunk: soft-delete a chunk by marking it as deleted (tombstone). Hard-delete would break sync — other devices need to learn the chunk is gone, not just stop seeing it. CLI verifies the chunk is gone on read after deletion.
- [ ] Sync endpoint: return all chunks changed since a given `updated_at` timestamp (delta sync). The client sends its last-known timestamp, the server returns everything newer. This is how devices stay in sync without downloading the full dataset every time. CLI creates/updates/deletes chunks, then syncs and verifies only the changes come back.
- [ ] Pagination for initial load: when a new device syncs for the first time (no `updated_at`), the full dataset could be large. Page the response so the client can load incrementally instead of waiting for one massive payload. CLI verifies paged responses assemble into the full dataset.

## Phase 2 — Crypto module (Kotlin, `odin-android`)

Isolated, well-defined interface. Replaceable with Rust + FFI later.

- [ ] Key derivation (Argon2id)
- [ ] AES-256-GCM encrypt/decrypt
- [ ] Master key generation and wrapping
- [ ] Interface design (`VaultCrypto`) for future Rust swap

## Phase 3 — Mobile app foundation (Kotlin, `odin-android`)

- [ ] Registration and login with E2E key setup
- [ ] Android Keystore integration for session persistence (biometric unlock)
- [ ] Local SQLite database
- [ ] Sync engine: upload/download encrypted chunks, version-based conflict resolution

## Phase 4 — Financial features (Kotlin, `odin-android`)

Core accounting, all client-side.

- [ ] Accounts (create, list, balances)
- [ ] Categories (create, list)
- [ ] Income
- [ ] Expenses

## Phase 5 — Advanced features (Kotlin, `odin-android`)

- [ ] Transfers
- [ ] Events (group transactions by trip, party, meeting, etc.)
- [ ] Multi-currency support
- [ ] Reports and dashboards

## Phase 6 — AI-powered expense entry (Kotlin + Gideon integration)

Voice/text expense creation to eliminate manual data entry friction — the retention differentiator.

- [ ] Sitia integration: mobile app authenticates with Sitia (all AI requests are routed through Sitia, Gideon doesn't handle public requests)
- [ ] Text-based expense parsing: natural language → structured expense (amount, currency, category, description)
- [ ] Gideon expense skill: mobile app sends plaintext to Sitia → Sitia routes to Gideon → Gideon parses → returns structured data → mobile encrypts locally (opt-in, zero-knowledge tradeoff with clear disclosure)
- [ ] Category matching: client sends user's category list alongside the text for accurate classification
- [ ] Audio clip support: speech-to-text → text parsing pipeline
- [ ] Multi-currency disambiguation (e.g. "20 dollars" when user has USD and EUR accounts)

## Phase 7 — Expand

- [ ] Web client
- [ ] iOS (extract crypto to Rust + UniFFI)
- [ ] Shared budgets / household access
