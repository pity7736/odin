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

- [ ] Design the new server domain entities and sync API contract
- [ ] Auth redesign: key derivation support (Argon2id params), encrypted master key storage, session tokens
- [ ] Encrypted chunk storage: CRUD with version-based optimistic locking
- [ ] Sync endpoint: delta sync via `updated_at`, pagination for initial load
- [ ] Remove old domain code (Account, Category, Income entities, use cases, repositories, HTMX handlers/templates)

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
