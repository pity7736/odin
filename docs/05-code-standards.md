# Code Standards

## 1. Tooling

- **Language:** Go 1.26
- **Web Framework:** Fiber v2
- **Template Engine:** gohtml templates with HTMX
- **UUID Generation:** `github.com/google/uuid` (UUIDv7). UUIDs stored and passed as `string`.
- **Decimal Arithmetic:** `github.com/govalues/decimal` for money values — never use floats.
- **Dependency Mocking:** `github.com/vektra/mockery` — always use `EXPECT()` to set up mock calls, never `.On()` directly.
- **Formatting:** All code must be formatted with `gofmt`.
- **Dependency Management:**
  - When importing a new package, add it with `go get <package>`.
  - Always run `go mod tidy` after adding dependencies.

## 2. Code Style

### 2.1. Encapsulation

- **Private Fields by Default:** Structs must use private fields to enforce encapsulation.
- **Judicious Use of Getters/Setters:** Only expose when necessary for the public API.
- **Exception for DTOs:** Structs for data transfer and serialization (e.g., API request/response bodies) may have public fields.

### 2.2. Package Structure

- All application packages live under `src/` following Clean Architecture layers (domain, application, infrastructure).
- Each package has a single, clear responsibility.
- Constructor functions should be named `New()` since the package name already provides context (e.g., `accountcreator.New(command, repository)`).
- Struct names MUST be descriptive.
- **Code Organization:**
    - Constructors MUST be defined immediately after the struct definition.
    - Methods MUST be defined immediately after the constructor(s).

### 2.3. Instantiation and Memory

- **Avoid Primitive Pointers:** Do not use `*string`, `*int` unless clearly justified for optionality.
- **Context as First Argument:** All functions that do I/O must accept `context.Context` as their first argument.

### 2.4. Struct Field Alignment

Order struct fields from largest to smallest to minimize padding and reduce memory usage.

| Size (Bytes) | Go Types |
| :--- | :--- |
| 24 | `slice` |
| 16 | `string`, `interface`, `any` |
| 8 | `int`, `uint`, `int64`, `uint64`, `float64`, `pointer`, `map`, `chan`, `func` |
| 4 | `int32`, `uint32`, `float32` |
| 2 | `int16`, `uint16` |
| 1 | `int8`, `uint8`, `byte`, `bool` |

### 2.5. General Idioms

- **Self-Documenting Code:** Code must be self-documenting through clear naming and proper structure. **Comments are strictly prohibited in source code.** If you feel a comment is necessary, the code is either too complex or poorly named and MUST be refactored until it is self-explanatory.
- **Blank Lines:**
  - A single blank line MUST separate function and method declarations.
  - Inside functions, blank lines are strictly forbidden. Never add blank lines after variable declarations, error checks, or any closing bracket.
  - If you feel the need for logical separation within a function, you MUST extract that logic into a new function instead of adding a blank line.
  - A blank line after a close bracket (`}`) of a control structure within a function is NOT allowed.
  - **Exception:** A blank line is allowed between the closing `)` of a multi-line function signature and the first statement in the body.
- **Trailing Newline:** All files MUST end with a single newline character.
- **Function Ordering:** Called functions and methods MUST be defined below **all** their callers. This ensures a top-down reading flow.
- **Use `any`:** Always use `any` instead of `interface{}`.
- **Naming Conventions:**
  - Follow Go's `PascalCase` and `camelCase` conventions.
  - Receiver name: `self`.
  - **Initialisms (ID, HTTP, API, URL, etc.):**
    - Must follow visibility rules. If they start a private field or local variable, use lowercase (e.g., `id`, `url`, `apiClient`).
    - If they are exported or in the middle/end of a name, use uppercase (e.g., `ID`, `CompanyID`, `companyID`, `getURL`).
  - **No Abbreviations:** Use complete words (e.g., `account` not `acct`, `category` not `cat`, `repository` not `repo`). Exceptions: `err`, and loop indices `i`, `j`, `k`.
  - **Conciseness:** Prefer the simplest full word that is clear in context. Do NOT add redundant domain prefixes if the scope is already clear.
  - Do NOT use suffixes like `Entity` (e.g., use `account`, not `accountEntity`).
  - `ctx` is reserved for `context.Context`.

### 2.6. Error Handling

- Return errors explicitly; do not panic at runtime.
- Use the `odinerrors` builder pattern for domain and application errors:
  - `message` field: English (for developers and logs).
  - `external` field: Spanish (user-facing messages in the app).
  - Tag with the appropriate error type (`DOMAIN`, `NOT_FOUND`, `RENDER`).
- Wrap errors with context using `odinerrors.NewErrorBuilder("context").WithWrapped(err).Build()`.

### 2.7. Fiber-Specific

- **Buffer Reuse:** Always use `strings.Clone()` when storing string values parsed from Fiber request bodies. Fiber's fasthttp reuses request buffers, which can silently corrupt stored strings across sequential requests.
- **Value Objects:** Use `MustNew()` constructors only in test code. In production code, always handle the error from `New()`.

## 3. Testing Guidelines

### 3.1. Test Naming and Structure

Tests must use the `Should` suffix in the function name and describe the behavior in `t.Run` without the word "should". No underscores in test names.

```go
func TestCreateAccountHandlerShould(t *testing.T) {
    t.Run("be able to create an account", func(t *testing.T) {
        // ...
    })
    t.Run("return error when initial balance is not valid", func(t *testing.T) {
        // ...
    })
}
```

### 3.2. Test Simplicity

Tests must be simple and declarative. One behavior per test. No conditionals or loops in tests.

### 3.3. Test Coverage

- 100% test coverage for all logic-bearing code (Domain and Application layers).
- A lower literal coverage (90-95%) is acceptable for the entire project to account for trivial boilerplate like getters and setters.
- Behavioral correctness is the priority; trivial functions are implicitly verified through higher-level tests.

### 3.4. Production Code vs Test Code

Production code must never be added solely to satisfy test requirements. Tests adapt to production code, not the other way around.

### 3.5. Test Builders

Use the builder pattern in `tests/builders/` to construct test data. Builders provide sensible defaults and fluent configuration:

- `.Build()` creates a domain object directly (for unit tests).
- `.Create(repository)` creates and persists via the use case layer (for integration tests).
