# Technical Plan: Create Account

**Corresponds to Spec:** `specs/accounting/accounts/create/spec.md`

## Overview

Registering a financial account: the user provides a name, an account type, a
currency, and an initial balance; the system stores the account under the
authenticated owner with a generated id, a creation date, and a starting balance
equal to the initial balance. Served over both HTMX (web) and REST (mobile).

Account creation already exists end-to-end (domain entity, use case, repository
port, in-memory adapter, HTMX + REST handlers, routing); this plan closes the
gaps between that code and the spec: account **type**, selectable **currency**
(COP/USD), stronger **name** validation (trim, non-empty, ≤255), **uniqueness**
on `(owner, name, currency)`, and user-facing error messages in Spanish.

## Design Decisions & Rationale

- **Type is a domain value object in its own package** — choice: a struct value
  object parallel to `Currency` with `Savings/CreditCard/Cash` + `NewFromString`.
  Reason: type and currency are siblings (constrained descriptors that arrive
  together on the account), so they should share a shape. Rejected: a bare int8
  enum copied from `constants.CategoryType`, which uses plain `fmt.Errorf` (no
  Spanish external, wrong error tag) and a non-exhaustive `String()`.
- **Currency is derived from Money, not stored on the account** — choice: no
  `currency` field; the account's currency is its balance's currency via a
  `Currency()` getter. Reason: `initialBalance`/`balance` already carry currency,
  so a separate field is a second source of truth that can disagree. Rejected: an
  explicit `currency` field that must be kept in sync with the money.
- **Currency is a value type, not a pointer** — choice: `Currency` passed/returned
  by value. Reason: it's a small immutable descriptor; value semantics make `==`
  compare by code and eliminate the nil-currency hazard; consistent with
  `AccountType`. Rejected: the current `*Currency`, where `COP() == COP()` is
  false (identity, not code) and `Money{}` can hold a nil currency.
- **Uniqueness is `(owner, name, currency)`, case-insensitive on name** — choice:
  same name allowed across currencies, blocked within one. Reason: models a
  multi-currency wallet (e.g. Global66) as one account per currency without
  stuffing the currency into the name string. Rejected: `(owner, name)`
  uniqueness (blocks the wallet case) and case-sensitive compare (allows
  near-duplicate look-alikes).
- **Owner comes from the authenticated session, never the request** — choice:
  read the owner from `RequestContext` (already the case). Reason: a
  client-supplied owner lets one user create accounts under another's id.
- **Errors are created once at their origin and propagated unchanged** — choice:
  a layer returns the error it received (`return nil, err`); it wraps only to
  introduce a *genuinely new* error (e.g. a render failure → `Render` tag), and
  such a wrap sets an internal English `message` only, never a second `external`.
  Reason: the domain error already carries its final tag and Spanish external;
  re-wrapping to add English external prefixes builds a mixed-language chain
  (`"error creating account: validation error: El nombre es obligatorio"`) via
  `ExternalError()` and tells the user nothing new. Rejected: the current
  per-layer wrapping, and translating the prefixes (which only turns clean
  redundancy into Spanish redundancy). A single generic Spanish fallback at the
  presentation boundary covers errors that carry no external at all.
- **Credit-card capacity/semantics deferred** — choice: initial balance ≥ 0 for
  all types; no credit-limit field now. Reason: a limit only does work at
  spending time, which doesn't exist yet; adding a field with no consumer is
  speculative. Rejected: modeling credit-card debt/available-credit during
  creation.
- **Display is product-facing, mapped in a reusable handler view model** —
  choice: the Spanish type label (Ahorros / Tarjeta de crédito / Efectivo) and
  the ISO `yyyy-mm-dd` creation date are built in a small reusable view model in
  the handler layer; the created row renders from it, and the template stays
  dumb. Reason: these are product decisions (in the spec), and mapping them in
  Go keeps the domain free of localization and the mapping unit-testable —
  unlike branching in an untestable template. `AccountType.String()` and the REST
  response keep the English *code* AND the machine-readable RFC3339 `created_at`
  (no mobile app yet; a future client localizes) — only the web view model
  produces the Spanish label + ISO date. The view model is written
  to be reused, but the **accounts-list rendering is NOT modified here** (listing
  is a separate feature) — so a newly-created row shows Spanish/ISO while existing
  list rows stay as-is until the list feature adopts the view model. Rejected:
  branching in the template, and a Spanish label on the domain value object.
  Full i18n (a translation system, localized long-form dates) remains a deferred
  separate feature.

## Architecture & Files Summary

```
src/accounting/domain/
├── accounttypemodel/
│   └── account_type.go                                     # CREATE
├── money/
│   ├── currency.go                                         # MODIFY (value type: COP/USD/CurrencyFromString/Equals)
│   └── money.go                                            # MODIFY (currency value field, Currency getter, ...Currency, Spanish parse error)
├── account/
│   └── account.go                                          # MODIFY (type field, name trim/≤255, getters, Spanish)
└── repositories/
    └── repositories.go                                     # MODIFY (ExistsByNameAndCurrency)

src/accounting/application/use_cases/accountcreator/
├── command.go                                              # MODIFY (accountType)
└── account_creator.go                                      # MODIFY (uniqueness check, pass type, propagate errors)

src/accounting/infrastructure/
├── repositories/pgrepositories/
│   └── account_repository.go                               # MODIFY (implement ExistsByNameAndCurrency)
└── api/handlers/accounthandler/
    ├── accountviewmodel/account_view_model.go              # CREATE (reusable: Spanish type label + ISO date)
    ├── createaccounthandler/handler.go                     # MODIFY (body: type/currency, missing-balance check, Clone)
    ├── restcreateaccounthandler/handler.go                 # MODIFY (response: type/currency)
    └── htmxcreateaccounthandler/handler.go                 # MODIFY (drop external wrap, generic fallback, render created row from view model)

src/shared/infrastructure/templates/pages/
└── accounts.gohtml                                         # MODIFY (form type/currency inputs, row+header cols, OOB error target, created row from view model — list render untouched)

tests/
├── builders/account_builder.go                             # MODIFY (WithType/WithCurrency)
├── unit/accounting/domain/
│   ├── accounttype/account_type_test.go                    # CREATE
│   ├── account/account_test.go                             # MODIFY
│   └── money/*                                             # MODIFY (currency/money tests)
├── unit/accounting/application/use_cases/account/
│   └── create_account_test.go                              # MODIFY
├── unit/accounting/infrastructure/handlers/accounthandlers/
│   ├── account_view_model_test.go                          # CREATE (label + ISO date mapping)
│   ├── rest_handler_test.go                                # MODIFY
│   └── htmx_handler_test.go                                # MODIFY
├── integration/accounting/accounts_test.go                 # MODIFY
└── unit/mocks/mock_AccountRepository.go                    # REGEN

specs/accounting/accounts/create/
├── spec.md                                                 # CREATE
└── plan.md                                                 # CREATE
```

## Data Flow

```
POST /accounts  |  POST /api/v1/accounts   (loginRequired → RequestContext in Locals)
  → createaccounthandler.Handle
      parse body → name, raw initial balance, type, currency (strings.Clone each)
      if rawInitialBalance == "" → Domain error "El saldo inicial es obligatorio"  (missing)
      currency := moneymodel.CurrencyFromString(currencyCode)  → Domain error if not {COP,USD}
      accountType := accounttypemodel.NewFromString(typeCode)  → Domain error if not {savings,credit_card,cash}
      initialBalance := moneymodel.New(rawInitialBalance, currency)  → Domain error (Spanish) if non-numeric (invalid)
      command := accountcreator.NewCreateAccountCommand(name, initialBalance, accountType)
    → AccountCreator.Create(ctx)                               (owner = RequestContext.UserID())
        repository.ExistsByNameAndCurrency(ctx, name, currency) → Domain error if duplicate
        accountmodel.New(name, userID, initialBalance, accountType) → validate (trim, ≤255, non-empty)
        repository.Add(ctx, account)
    ← account (id, name, initialBalance, balance, type, currency, userID, createdAt)
  → HTMX strategy: render account_created / create_account_error
  → REST strategy: JSON response including type and currency
```

## Request & Response

**Request data** (fields the client provides):

| Field | Type | Required | Notes |
|-------|------|----------|-------|
| name | string | yes | non-empty after trim, ≤ 255 chars |
| initial_balance | string (decimal) | yes | ≥ 0, decimals allowed; empty → "required" |
| type | string | yes | one of `savings`, `credit_card`, `cash` (case-insensitive) |
| currency | string | yes | one of `COP`, `USD` (case-insensitive) |

Owner is NOT a request field — it comes from the authenticated session.

**REST** — `POST /api/v1/accounts`
```json
// request
{ "name": "Global66", "initial_balance": "1500000.50", "type": "savings", "currency": "COP" }
// success 201  (createAccountResponse)
{ "id": "...", "name": "Global66", "initial_balance": "1500000.50",
  "balance": "1500000.50", "type": "savings", "currency": "COP",
  "user_id": "...", "created_at": "<RFC3339>" }
// error 400  (Spanish external via the global errorHandler; match its existing shape)
{ "error": "El saldo inicial es obligatorio" }
```
The error body shape is whatever the existing global `errorHandler` in
`fiber_application.go` already emits — confirm and match it; do not introduce a
new error shape.

**HTMX** — `POST /accounts` (form in `accounts.gohtml`, `hx-swap="none"`)
- Form fields: `name`, `initial_balance`, `type`, `currency` (the last two must be
  ADDED to `create_account_form` — inputs/selects for type and currency).
- Success: renders `account_created`, a `<tr hx-swap-oob="afterbegin:#accounts-view">`
  that **prepends** the new row to the top of the accounts table (`#accounts-view`).
  This is a live-updating list via out-of-band swap; the form's `hx-swap="none"`
  is intentional because insertion happens out-of-band.
  - The row currently shows Name / InitialBalance / Balance / CreatedAt. Add
    **type and currency** to the row AND matching header columns in the table.
- Error: renders `create_account_error` with the error's Spanish external message.
  See the error-display gap below — as wired, this fragment is not shown.

## Key Types & Signatures

### `AccountType` value object (`src/accounting/domain/accounttypemodel/account_type.go`)

```go
type AccountType struct {
    code string
}

func Savings() AccountType
func CreditCard() AccountType
func Cash() AccountType
func NewFromString(value string) (AccountType, error)
func MustNewFromString(value string) AccountType
func (self AccountType) String() string
func (self AccountType) Equals(other AccountType) bool
```

- Canonical wire codes: `savings`, `credit_card`, `cash`. `NewFromString`
  lowercases input and rejects anything else with `odinerrors.Domain`
  (external Spanish, e.g. `"Tipo de cuenta inválido"`).
- `MustNewFromString` is test-only (panics), matching the `MustNew` convention.
- `String()` maps all three explicitly — no silent fallback.

### Currency additions (`src/accounting/domain/money/currency.go`)

`Currency` becomes a **value type**, not a pointer (see the conversion gap
below). `COP()` also changes to return a value.

```go
func COP() Currency
func USD() Currency
func CurrencyFromString(code string) (Currency, error)
func (self Currency) Equals(other Currency) bool
```

- `CurrencyFromString` accepts only `COP`/`USD` (case-insensitive), rejects
  others with `odinerrors.Domain` (external Spanish, e.g. `"Moneda inválida"`).
  Named `CurrencyFromString` (not `NewFromString`) to avoid confusion with
  `moneymodel.New`, which builds `Money` in the same package.
- `Equals` compares by `code`, mirroring `AccountType.Equals`, so the adapter's
  uniqueness check has a comparison mechanism.

### Money getter (`src/accounting/domain/money/money.go`)

```go
func (self Money) Currency() Currency
```

- Also Spanish-ize the parse error `moneymodel.New` returns for a non-numeric
  value (`money.go:30`), keeping the `Domain` tag.

### `Account` changes (`src/accounting/domain/account/account.go`)

```go
func New(name, userID string, initialBalance moneymodel.Money, accountType accounttypemodel.AccountType) (*Account, error)
func NewFromRepository(id, name, userID string, initialBalance, balance moneymodel.Money, accountType accounttypemodel.AccountType, createdAt time.Time) (*Account, error)
func (self *Account) Type() accounttypemodel.AccountType
func (self *Account) Currency() moneymodel.Currency   // returns self.balance.Currency()
```

- Add the `accountType` field (respect struct field alignment ordering).
- `validateData`: trim name; reject blank/whitespace-only
  (external `"El nombre es obligatorio"`); reject length > 255
  (external `"El nombre es demasiado largo"`). Existing negative-balance checks
  keep behavior (≥ 0 allowed) but get Spanish externals.

### Command changes (`.../accountcreator/command.go`)

```go
func NewCreateAccountCommand(name string, initialBalance moneymodel.Money, accountType accounttypemodel.AccountType) CreateAccountCommand
func (self CreateAccountCommand) AccountType() accounttypemodel.AccountType
```

- Currency needs no command field — it rides inside `initialBalance` (Money).

### Repository port change (`src/accounting/domain/repositories/repositories.go`)

```go
type AccountRepository interface {
    Add(ctx context.Context, account *accountmodel.Account) error
    ExistsByNameAndCurrency(ctx context.Context, name string, currency moneymodel.Currency) (bool, error)
    GetAll(ctx context.Context) ([]*accountmodel.Account, error)
    GetByID(ctx context.Context, id string) (*accountmodel.Account, error)
    Save(ctx context.Context, account *accountmodel.Account) error
}
```

- `ExistsByNameAndCurrency` scopes to the context user (like `GetAll`/`GetByID`),
  compares name case-insensitively and currency via `Currency.Equals`. Adapter implements it in
  `.../pgrepositories/account_repository.go`. Mocks regenerated.

### Use case change (`.../accountcreator/account_creator.go`)

- Before constructing the account, call
  `repository.ExistsByNameAndCurrency(ctx, command.Name(), command.InitialBalance().Currency())`.
  If it returns true → create (at origin) an `odinerrors.Domain` error, external
  `"Ya tienes una cuenta con ese nombre en esa moneda"`.
- Pass `command.AccountType()` into `accountmodel.New`.
- On the `accountmodel.New` / repository errors, **propagate** (`return nil, err`)
  — drop the current `WithExternalMessage("validation error")` wrap.

### Handler DTO changes

- `createAccountBody` (inner handler): add `Type string` and `Currency string`
  (`json`/`form` tags); `strings.Clone` name, initial balance, type, currency.
- Presence check for a **missing initial balance** (empty raw string) →
  `odinerrors.Domain`, external `"El saldo inicial es obligatorio"`, raised in the
  handler before `moneymodel.New` (mirrors the auth login handler's email/password
  presence checks).
- `createAccountResponse` (REST handler): add `type` and `currency` fields.

## Gaps / Bugs to Fix

- [ ] **No account type.** `Account` (`account.go`) has no type field and no type
      value object exists. Create `accounttypemodel` and add the field.
- [ ] **Currency not selectable; USD missing.** `moneymodel.New` defaults to
      `COP()` and only `COP()` exists (`currency.go`); `newCurrency` is
      unexported. Add `USD()`, a public `CurrencyFromString` restricted to
      {COP, USD}, a `Currency.Equals`, and a `Money.Currency()` getter so the
      account can expose and compare its currency.
- [ ] **`Currency` is a pointer; make it a value.** `COP()` returns `*Currency`
      and `Money.currency` is `*Currency` (`currency.go`, `money.go:12`). Pointer
      identity makes `COP() == COP()` false and allows a nil currency. Convert
      `Currency` to a value type: `COP`/`USD`/`CurrencyFromString` return
      `Currency`; `Money.currency` becomes `Currency`; `New`/`MustNew` take
      `...Currency`; update `Subtract` and every other `*Currency` usage. `Equals`
      becomes `self == other`. Mirrors `AccountType`'s value semantics.
- [ ] **Weak name validation.** `validateData` (`account.go`) only checks
      `name == ""`. Add: trim surrounding whitespace, reject blank/whitespace-only,
      reject length > 255. Store the trimmed name.
- [ ] **No uniqueness.** `AccountCreator.Create` performs no duplicate check and
      `AccountRepository` has no lookup method. Add
      `ExistsByNameAndCurrency` to the port, implement it in the in-memory
      adapter (case-insensitive name, scoped to the context user), and enforce it
      in the use case.
- [ ] **English user-facing errors at the origin.** `validateData` (`account.go`)
      sets `external` equal to the English `message`, AND `moneymodel.New`
      (`money.go:30`) builds an English external (`"<x> is not valid money
      value"`) on the create path for an unparseable balance. Rewrite these
      origin `external` messages in Spanish (keep internal `message` English).
- [ ] **Stop re-wrapping errors in the layers (mixed-language chain).** The use
      case (`account_creator.go:24-29`) wraps the domain error with
      `WithExternalMessage("validation error")`, and the HTMX handler
      (`htmxcreateaccounthandler/handler.go`) wraps with
      `WithExternalMessage("error creating account")`; `ExternalError()` chains
      these into `"error creating account: validation error: El nombre es
      obligatorio"`. Propagate instead: use case returns `return nil, err`; HTMX
      handler drops its `"error creating account"` external. KEEP the HTMX
      handler's render-failure wrap (`Render` tag) — that is a genuinely new
      error. Do NOT touch `odinerrors` (see known limitation).
- [ ] **Generic Spanish fallback for external-less errors.** A technical failure
      (repo/render error) can reach the user with an empty `external` and show a
      blank message. Provide one default at the presentation boundary (e.g.
      `"No se pudo crear la cuenta"`) when the error has no external.
- [ ] **Balance path must produce three distinct spec messages.** The spec
      defines *missing* → `"El saldo inicial es obligatorio"` and *negative* →
      `"El saldo inicial no puede ser negativo"`, plus the implicit *invalid*
      (non-numeric) case. Today an empty balance flows into `moneymodel.New("")`
      and yields the generic English "not valid money value" — so the specific
      "required" wording is never produced. Split the outcomes by layer:
      - **Missing** (empty raw string): presence check at the create handler →
        Domain error `"El saldo inicial es obligatorio"`, before `moneymodel.New`.
      - **Invalid** (non-empty, non-numeric): `moneymodel.New` → Domain error
        with a Spanish external.
      - **Negative** (`< 0`): `validateData` in `account.go` → Domain error
        `"El saldo inicial no puede ser negativo"`.
- [ ] **Body/response omit type and currency.** `createAccountBody` and
      `createAccountResponse` (inner/REST handlers) must carry both; apply
      `strings.Clone` to the new parsed strings.
- [ ] **Templates (`accounts.gohtml`) must surface type and currency.** The
      `account_created` row shows only Name/InitialBalance/Balance/CreatedAt and
      the table has four headers — add **type and currency** as new row cells AND
      header columns. Also add **type and currency inputs** to
      `create_account_form` (currently only name + initial_balance).
- [ ] **HTMX error message is never shown (bug).** `create_account_form` uses
      `hx-swap="none"` and `create_account_error` renders a bare `<p>` with no
      `hx-swap-oob` target, so HTMX discards it — the user is never told why
      creation failed, violating every rejection scenario in the spec. Give the
      error fragment an out-of-band target (a dedicated error container in
      `accounts.gohtml`) so rejections are displayed.
- [ ] **Test builder gap.** `tests/builders/account_builder.go` has no
      `WithType`/`WithCurrency` — add them.
- [ ] **Audit ignored errors** in the create-path HTMX handler render calls; any
      `_ =` suppression must be handled and covered by a test (see project
      convention on ignored errors).
- [ ] **Type shows the English code; date shows English long-form.** The created
      row renders `{{ .Type }}` (the code `savings`) and
      `CreatedAt.Format("Monday, _2 January 2006")` (English month/weekday). The
      spec requires the type in Spanish and the date as ISO `yyyy-mm-dd`.
      Introduce a small **reusable account view model** in the accounting handler
      layer that maps the type code → Spanish label (Ahorros / Tarjeta de crédito
      / Efectivo) and formats `CreatedAt` as `2006-01-02`; render the created row
      from it and keep the template dumb. Keep `AccountType.String()` and the REST
      response as the English code. Cover the mapping with a unit test. Do NOT
      modify the accounts-list handler or its template (see known limitation).

### Known limitations (flagged, NOT fixed here)
- The accounts-**list** rows (owned by the separate list feature) are not updated
  to the Spanish-label / ISO-date view model, so existing rows keep the English
  code/long-form date until the list feature adopts the reusable view model. Only
  the newly-created row reflects the new display.
- Full i18n (a translation system and localized long-form dates) is deferred to
  its own future feature; this change hardcodes the single Spanish locale in the
  view model, consistent with the rest of the app.
- The uniqueness check is **not atomic** in the in-memory adapter (check-then-add).
  Acceptable for single-user/in-memory; revisit with the Postgres feature.
- `PGAccountRepository.Save` panics and the repo is in-memory only. Create uses
  `Add`, so this does not block the feature; out-of-scope debt.
- `odinerrors.ExternalError()` (`error.go:40-52`) emits a leading `": "` when a
  wrapping node has an empty `external` (the separator keys off the next node).
  A latent bug in shared code; NOT triggered once we stop chaining externals
  (errors created once at origin), so left out of scope for this feature.

## Quality Pillars

- **Security:** All inputs validated at the domain/handler boundary (name length
  and content, type and currency against closed sets); owner taken from the
  authenticated `RequestContext`, never the body; `strings.Clone` on every parsed
  body string (name, initial balance, type, currency) to avoid fasthttp buffer
  reuse; 255-char cap bounds name size.
- **Reliability:** Uniqueness check prevents duplicate `(user, name, currency)`
  accounts; value objects make invalid type/currency unrepresentable; no panics
  in the create path; no ignored errors (render errors handled and tested).
  Known non-atomicity of the in-memory uniqueness check is documented above.
- **Performance:** Deferred — single user, in-memory repository. The uniqueness
  check is an O(n) scan over the user's accounts, negligible at this scale;
  revisit with an indexed lookup in the Postgres feature.
- **Observability:** Deferred — no production telemetry yet. `odinerrors`
  location tracking on domain errors provides debugging context in logs.

