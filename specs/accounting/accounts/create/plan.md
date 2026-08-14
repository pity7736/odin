# Technical Plan: Create Account

**Corresponds to Spec:** `specs/accounting/accounts/create/spec.md`

> **WORK ORDER (re-opened):** fix — REST rejections return the correct status but
> an empty response body, so the user is never told why creation failed
> (spec.md: every rejection scenario ends "the user is told …", and criterion
> "the user is told the reason"). See **Gaps / Bugs to Fix**. Prune back to a
> living design doc once shipped.

## Overview

Registering a financial account: the user provides a name, an account type
(savings / credit card / cash), a currency (COP / USD), and an initial balance;
the system stores the account under the authenticated owner with a generated id,
a creation date, and a starting balance equal to the initial balance. Served over
both HTMX (web) and REST (mobile).

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
  `AccountType`. Rejected: a `*Currency`, where `COP() == COP()` is false
  (identity, not code) and `Money{}` could hold a nil currency.
- **Uniqueness is `(owner, name, currency)`, case-insensitive on name** — choice:
  same name allowed across currencies, blocked within one. Reason: models a
  multi-currency wallet (e.g. Global66) as one account per currency without
  stuffing the currency into the name string. Rejected: `(owner, name)`
  uniqueness (blocks the wallet case) and case-sensitive compare (allows
  near-duplicate look-alikes).
- **Owner comes from the authenticated session, never the request** — choice:
  read the owner from `RequestContext`. Reason: a client-supplied owner lets one
  user create accounts under another's id.
- **Errors are created once at their origin and propagated unchanged** — choice:
  a layer returns the error it received (`return nil, err`); it wraps only to
  introduce a *genuinely new* error (e.g. a render failure → `Render` tag), and
  such a wrap sets an internal English `message` only, never a second `external`.
  Reason: the domain error already carries its final tag and Spanish external;
  re-wrapping to add English external prefixes builds a mixed-language chain
  (`"error creating account: validation error: El nombre es obligatorio"`) via
  `ExternalError()` and tells the user nothing new. Rejected: per-layer wrapping,
  and translating the prefixes (which only turns clean redundancy into Spanish
  redundancy). A single generic Spanish fallback at the presentation boundary
  covers errors that carry no external at all.
- **Each interface handler formats its own error body; the global `errorHandler`
  owns only the `tag → status` mapping** — choice: on a failed create, the HTMX
  handler renders the Spanish message as an HTML fragment and the REST handler
  serializes it as `{"error": <external>}`; both then `return err`, and the
  shared `errorHandler` (`src/app/fiber_application.go`) sets the HTTP status
  from the error tag without writing a body. Reason: the *error* must be the
  same across interfaces (same status, same Spanish external message), only the
  *format* differs — HTML for the browser, JSON for the client — so formatting
  belongs in each interface's own handler while the status mapping stays shared
  and identical. This mirrors the pattern the HTMX handler already used.
  Rejected: content-negotiation inside the global `errorHandler` (couples the
  shared handler to a feature's template and OOB target) and letting the global
  handler write a JSON body for everyone (would put JSON on HTMX error
  responses). The Spanish fallback for a message-less error
  (`"No se pudo crear la cuenta"`) is shared by both handlers via a single
  `externalOrFallback` helper.
- **Credit-card capacity/semantics deferred** — choice: initial balance ≥ 0 for
  all types; no credit-limit field. Reason: a limit only does work at spending
  time, which doesn't exist yet; adding a field with no consumer is speculative.
  Rejected: modeling credit-card debt/available-credit during creation.
- **Display is product-facing, mapped in a reusable handler view model** —
  choice: the Spanish type label (Ahorros / Tarjeta de crédito / Efectivo) and
  the ISO `yyyy-mm-dd` creation date are built in a small reusable view model in
  the handler layer; the created row renders from it and the template stays dumb.
  Reason: these are product decisions (in the spec), and mapping them in Go keeps
  the domain free of localization and the mapping unit-testable — unlike branching
  in an untestable template. `AccountType.String()` and the REST response keep the
  English *code* and the machine-readable RFC3339 `created_at` (no mobile app yet;
  a future client localizes) — only the web view model produces the Spanish label
  + ISO date. Rejected: branching in the template, and a Spanish label on the
  domain value object. Full i18n (a translation system, localized long-form dates)
  remains a deferred separate feature.

## Architecture & Files Summary

```
src/accounting/domain/
├── accounttypemodel/
│   └── account_type.go
├── money/
│   ├── currency.go
│   └── money.go
├── account/
│   └── account.go
└── repositories/
    └── repositories.go

src/accounting/application/use_cases/accountcreator/
├── command.go
└── account_creator.go

src/accounting/infrastructure/
├── repositories/pgrepositories/
│   └── account_repository.go
└── api/handlers/accounthandler/
    ├── accountviewmodel/account_view_model.go
    ├── createaccounthandler/handler.go
    ├── restcreateaccounthandler/handler.go          # MODIFY (write JSON error body)
    └── htmxcreateaccounthandler/handler.go          # MODIFY (use shared externalOrFallback)

src/shared/infrastructure/api/
└── external_error.go                                # CREATE (package handler: ExternalOrFallback)

src/shared/infrastructure/templates/pages/
└── accounts.gohtml

tests/unit/accounting/infrastructure/handlers/accounthandlers/
└── rest_handler_test.go                             # MODIFY (assert error body on rejections)
```

## Data Flow

```
POST /accounts  |  POST /api/v1/accounts   (loginRequired → RequestContext in Locals)
  → createaccounthandler.Handle
      parse body → name, raw initial balance, type, currency (strings.Clone each)
      if body is not valid JSON → Domain error "Datos de solicitud inválidos" (400)
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
  → HTMX: build an account view model (Spanish type label, ISO date); render the
    created row prepended at the top of the table via an out-of-band swap; on
    error render the Spanish message into the #account-error container (also OOB),
    then return the error so errorHandler sets the status.
  → REST: on success, JSON with the type code, currency code, and RFC3339
    created_at; on error, JSON { "error": <Spanish external | fallback> }, then
    return the error so errorHandler sets the status. The error body is written
    by the REST handler itself — errorHandler writes no body.
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
// success 201
{ "id": "...", "name": "Global66", "initial_balance": "1500000.50",
  "balance": "1500000.50", "type": "savings", "currency": "COP",
  "user_id": "...", "created_at": "<RFC3339>" }
// error 400 — Spanish external message, written by the REST handler;
//            status set from the error tag by the global errorHandler
{ "error": "El saldo inicial es obligatorio" }
```
The type and currency stay as codes and `created_at` as RFC3339; a client
localizes for display. Every rejection (blank/too-long name, duplicate
name+currency, missing/invalid type, missing/invalid currency, missing/negative
initial balance) returns the same `{ "error": <message> }` shape; a
message-less error falls back to `"No se pudo crear la cuenta"`.

**HTMX** — `POST /accounts` (form in `accounts.gohtml`)
- Form fields: `name`, `initial_balance`, `type` (select), `currency` (select).
- Success: renders the created row from the account view model (Spanish type
  label, ISO date), inserted at the top of the accounts table via an out-of-band
  swap.
- Error: renders the Spanish external message into the `#account-error` container
  (out-of-band), so the reason is shown to the user.
- Existing list rows are rendered by the separate list feature and are not
  affected here (see Known Limitations).

## Known Limitations

- Existing accounts-list rows (owned by the separate list feature) still show the
  English type code and long-form date; only the newly-created row uses the
  Spanish-label / ISO-date view model. The list feature will adopt the shared
  view model when it is next worked on.
- Full i18n (a translation system and localized long-form dates) is deferred to
  its own feature; the current code hardcodes the single Spanish locale in the
  view model, consistent with the rest of the app.
- The uniqueness check is not atomic in the in-memory adapter (check-then-add);
  acceptable for the single-user in-memory setup, to revisit with the Postgres
  adapter.
- `PGAccountRepository` is in-memory and its `Save` panics; account creation uses
  `Add`, so this does not affect the feature.
- `odinerrors.ExternalError()` emits a leading `": "` when a wrapping node has an
  empty external; not triggered here because errors are created once at origin and
  not chained, so left as-is.

## Key Types & Signatures

Shared fallback helper, moved out of the HTMX handler package so both handlers
use one implementation:

```
package handler  // src/shared/infrastructure/api (dir already hosts package handler)
func ExternalOrFallback(err error, fallback string) string
```
- Returns the error's Spanish `ExternalError()` when present, else `fallback`.
- HTMX handler passes fallback `"No se pudo crear la cuenta"`; REST handler
  passes the same fallback.

REST handler error path (`restcreateaccounthandler/handler.go`):
```
account, err := self.handler.Handle(ctx)
if err != nil {
    _ = ctx.JSON(fiber.Map{"error": api.ExternalOrFallback(err, "No se pudo crear la cuenta")})
    return err                 // errorHandler sets the status from the tag
}
```
Content-type is already `application/json` (set before delegating). `errorHandler`
is unchanged — it maps `tag → status` and writes no body.

## Gaps / Bugs to Fix

- [ ] **REST rejections return an empty body.** `restcreateaccounthandler/handler.go`
      returns the error without writing a response body; the global `errorHandler`
      (`src/app/fiber_application.go:204-205`) sets the status and writes nothing,
      so every rejected create yields the right status with no message. Fix: the
      REST handler writes `{ "error": ExternalOrFallback(err, …) }` before
      returning the error.
- [ ] **Extract `externalOrFallback` to `src/shared/infrastructure/api`** as
      `ExternalOrFallback(err, fallback)`; update the HTMX handler to use it
      (interim — ahead of the broader error-handling task).
- [ ] **Reproduction tests (Red first) in `rest_handler_test.go`** — one per
      rejection, asserting the response body carries the Spanish external
      message: blank name → `El nombre es obligatorio`; duplicate name+currency →
      `Ya tienes una cuenta con ese nombre en esa moneda`; missing initial
      balance → `El saldo inicial es obligatorio`; invalid type → `Tipo de cuenta
      inválido`; invalid currency → `Moneda inválida`; negative initial balance →
      `El saldo inicial no puede ser negativo`; and the fallback branch
      (message-less error, e.g. malformed body) → `No se pudo crear la cuenta`.
      Extend the existing "return error when name is empty" and "initial balance
      is not valid" tests with a body assertion rather than duplicating them.
- [ ] **Delete `tests/unit/app/account_error_response_test.go`** — a premature
      app-level file from the wrong first diagnosis; the defect is observable at
      the REST handler's own unit level, so the tests belong in `rest_handler_test.go`.
- [ ] **Malformed body returned 500, not 400** (found in manual test).
      `createaccounthandler.createCommand` returned the raw `ctx.BodyParser` error
      untagged, so the global `errorHandler` fell through to its 500 default. Fix:
      wrap it as a `Domain`-tagged `odinError` with external "Datos de solicitud
      inválidos" (mirrors `loginhandler.validateRequestBody`), preserving the parse
      cause via `WithWrapped`. Guarded by the "return error when body is not valid"
      test asserting `Tag() == Domain` and the message in the body.

## Quality Pillars

- **Security:** All inputs validated at the domain/handler boundary (name length
  and content, type and currency against closed sets); owner taken from the
  authenticated `RequestContext`, never the body; `strings.Clone` on every parsed
  body string (name, initial balance, type, currency) to avoid fasthttp buffer
  reuse; 255-char cap bounds name size.
- **Reliability:** Uniqueness check prevents duplicate `(user, name, currency)`
  accounts; value objects make invalid type/currency unrepresentable; no panics
  in the create path; no ignored errors (render errors handled and tested).
  Every rejection returns its Spanish reason in the response body (HTML for HTMX,
  JSON for REST), guarded by a reproduction test per rejection path so no path
  can silently regress to an empty body. Known non-atomicity of the in-memory
  uniqueness check is documented above.
- **Performance:** Deferred — single user, in-memory repository. The uniqueness
  check is an O(n) scan over the user's accounts, negligible at this scale;
  revisit with an indexed lookup in the Postgres feature.
- **Observability:** Deferred — no production telemetry yet. `odinerrors`
  location tracking on domain errors provides debugging context in logs.
