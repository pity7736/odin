# Work Order: Authentication — repository fetch contract (NotFound + Exists)

**Feature design:** `specs/accounts/authentication/design.md` (the living source of truth)
**Corresponds to Spec:** `specs/accounts/authentication/spec.md`

> Work order for: **making the user repository signal absence explicitly**.
> Disposable — overwritten by the next change (git keeps history). Hydrate
> design.md before merge, then freeze this file.

## Change

Refactor the `UserRepository` fetch contract so absence is explicit, matching what
`SessionRepository` already does. **Behavior is unchanged** — `spec.md` does not
change (login with a wrong email still returns the same field-agnostic 401).

- `GetByEmail` becomes a pure **fetch**: returns the user, or an `odinerrors`
  `NotFound` when absent (was `(nil, nil)`).
- New `Exists(email) (bool, error)` for existence checks — no fetch of the whole
  user, explicit intent (mirrors the chunk repo's `Exists`).
- `sessionstarter` (login) treats a `NotFound` from `GetByEmail` as
  wrong-credentials; any other error propagates; a found user is compared (the old
  `user != nil` nil-check disappears).
- `userregistrar` (registration) switches its duplicate check from `GetByEmail`
  to `Exists` — an available email no longer fetches a user, and it no longer
  depends on a nil return.

Satisfies the same auth spec scenarios (login success / wrong password / non-
existing email / registration duplicate) — with no behavior change.

## Architecture & Files (this change)
```
src/accounts/domain/repositories/
└── user.go                                               # MODIFY  add Exists; GetByEmail doc-contract (NotFound on absent)

src/accounts/application/use_cases/
├── sessionstarter/session_starter.go                     # MODIFY  branch on NotFound
└── userregistrar/user_registrar.go                       # MODIFY  duplicate check via Exists

src/accounts/infrastructure/repositories/inmemory/
└── user_repository.go                                    # MODIFY  GetByEmail → NotFound on absent; add Exists

.mockery.yaml                                             # (unchanged — UserRepository already listed)
tests/unit/mocks/mock_UserRepository.go                   # REGEN  (interface gained Exists)

tests/unit/accounts/application/use_cases/
├── login_test.go                                         # MODIFY  not-found case → GetByEmail returns NotFound
└── register_test.go                                      # MODIFY  GetByEmail expectations → Exists
tests/unit/accounts/infrastructure/api/
├── login_api_test/login_test.go                          # MODIFY  not-found case → NotFound
└── register_api_test/register_test.go                    # MODIFY  GetByEmail expectations → Exists
tests/integration/accounts/auth_test.go                   # VERIFY  unknown-email login still 401, duplicate register still 409 (likely no change)
```

## Key Types & Signatures

```go
// src/accounts/domain/repositories — GetByEmail is now a pure fetch (NotFound on
// absent); Exists added for existence checks.
type UserRepository interface {
    GetByEmail(ctx context.Context, email string) (*usermodel.User, error) // NotFound when absent
    Exists(ctx context.Context, email string) (bool, error)
    Add(ctx context.Context, user *usermodel.User) error
}

// sessionstarter.Start: GetByEmail → NotFound ⇒ wrong-credentials (401);
// other error ⇒ propagate; found ⇒ compare auth hash.
// userregistrar.Register: Exists → true ⇒ AlreadyExists (409); false ⇒ create.
// NotFound detection: errors.As(&odinErr) && odinErr.Tag() == odinerrors.NotFound
// (mirrors sessionvalidator).
```

## Implementation Phases (TDD)

### Phase 1: Domain port + mock
**Red:** n/a (interface).
**Green:** add `Exists` to the `UserRepository` interface; regenerate the mock
(`go run github.com/vektra/mockery/v3@latest`) so downstream tests have
`Exists` on `MockUserRepository`.

### Phase 2: Application — sessionstarter + userregistrar
**Red:**
- `login_test.go`: when `GetByEmail` returns an `odinerrors` `NotFound`, `Start`
  returns the wrong-credentials 401 (Unauthorized tag), NOT the NotFound; a
  non-NotFound error still propagates; the happy path (user found) is unchanged.
- `register_test.go`: `Exists` → true ⇒ `AlreadyExists`, `Add` NOT called;
  `Exists` → false ⇒ `Add` called; `Exists` error ⇒ propagate.
**Green:** `sessionstarter.Start` branches on the `NotFound` tag (mirror
`sessionvalidator`); drop the `user != nil` check. `userregistrar` calls `Exists`
instead of `GetByEmail`.

### Phase 3: Infrastructure — in-memory repo + handler tests
**Red:**
- in-memory user repo test: `GetByEmail` on a missing email returns a
  `NotFound`-tagged error; `Exists` returns false for absent, true after `Add`.
- `login_api_test`: the not-found login case mocks `GetByEmail` → `NotFound` (was
  `nil, nil`) → still 401. Non-NotFound repo errors (existing cases) still propagate.
- `register_api_test`: duplicate/available cases mock `Exists` (was `GetByEmail`).
**Green:** in-memory `GetByEmail` returns `NotFound` on absent; implement `Exists`.
Confirm `tests/integration/accounts/auth_test.go` still passes unchanged
(unknown-email login → 401, duplicate register → 409). Finish: `make check` GREEN.

## Design decisions to hydrate into design.md
- [ ] `UserRepository` fetch contract: `GetByEmail` returns `NotFound` on absent
  (was `(nil, nil)`), consistent with `SessionRepository`; existence checks use a
  dedicated `Exists` (no whole-user fetch), consistent with the chunk repo.
- [ ] `sessionstarter` translates `NotFound` → wrong-credentials (behavior
  unchanged, field-agnostic 401); `userregistrar` duplicate check uses `Exists`.
- [ ] Reliability pillar: absence is an explicit error, not a nil sentinel — removes
  the nil-check footgun; a future adapter can signal not-found the same way.
