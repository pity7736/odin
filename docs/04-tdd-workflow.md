# Test-Driven Development (TDD) Workflow

We follow the "Red-Green-Refactor" cycle.

1. **Red:** Write a failing test in `tests/unit/` or `tests/integration/` that expresses a single piece of required functionality. This test should not compile or should fail.

2. **Green:** Write the simplest possible production code in `src/` to make the test pass. Do not add extra functionality.

3. **Refactor:** Clean up the code (both test and production) while keeping the tests green. Improve naming, remove duplication, and enhance clarity.

All business logic **must** be developed with TDD, starting with unit tests.

## What to Test

- **Domain entities:** Validation rules, business logic, state transitions (e.g., account balance after income).
- **Use cases:** Orchestration flow, error propagation, repository interactions.
- **Handlers:** Request parsing, response formatting, error cases, content type negotiation.
- **HTMX handlers:** Verify rendered HTML contains expected elements (using `strings.Contains` on response body).

## What NOT to Unit Test

- Database connectivity (integration test when Postgres is implemented).
- Configuration loading from environment (trivial).
- Template file existence (verified by integration tests).

## Test Infrastructure

Odin uses several testing tools and patterns:

- **mockery:** Auto-generated mocks with `EXPECT()` style. Configured in `.mockery.yaml`.
- **gomonkey:** Patching `uuid.NewV7()` and `time.Now()` for deterministic tests.
- **testify:** `assert` for assertions, `mock` for mock verification.
- **Builders:** Test data builders in `tests/builders/` (`AccountBuilder`, `FiberContextBuilder`, `CategoryBuilder`, `UserBuilder`) with fluent API and `.Build()` / `.Create(repository)` methods.
- **Test utilities:** `tests/testutils/` provides helpers like `IsUUIDv7()`, `IsTimeClose()`.

## Test File Location

- Unit tests: `tests/unit/<module>/<layer>/<feature>/`
- Integration tests: `tests/integration/<module>/`
- Mocks: `tests/unit/mocks/`
- Builders: `tests/builders/`
