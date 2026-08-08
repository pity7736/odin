# Architecture

Odin is a personal finance application that serves two interfaces from a shared domain layer: HTMX-rendered web pages and a REST JSON API for mobile clients.

Odin **MUST** follow Clean Architecture. This separates concerns, enhances testability, and makes the system independent of external frameworks and services.

The architecture is divided into three main layers: **Domain**, **Application**, and **Infrastructure**.

## Directory Structure

```
.
├── go.mod
├── src/
│   ├── main.go
│   ├── app/
│   │   ├── application.go
│   │   └── fiber_application.go
│   ├── <module>/
│   │   ├── domain/
│   │   │   ├── <entity>/
│   │   │   │   └── <entity>.go
│   │   │   └── repositories/
│   │   │       └── repositories.go
│   │   ├── application/
│   │   │   └── use_cases/
│   │   │       └── <usecase>/
│   │   │           ├── command.go
│   │   │           └── <usecase>.go
│   │   └── infrastructure/
│   │       ├── api/
│   │       │   └── handlers/
│   │       └── repositories/
│   └── shared/
│       ├── domain/
│       │   ├── odinerrors/
│       │   └── requestcontext/
│       └── infrastructure/
│           ├── api/
│           └── templates/
├── tests/
│   ├── unit/
│   ├── integration/
│   └── builders/
├── docs/
├── specs/
├── CLAUDE.md
└── .golangci.yml
```

**Note:** Directories and packages are created only when needed — never pre-create empty packages.

## Modules

Odin has three modules:

- **`accounting`**: Core financial domain — accounts, categories, incomes, expenses, transfers, money.
- **`accounts`**: User identity — users, sessions, authentication.
- **`shared`**: Cross-cutting concerns — error handling (`odinerrors`), request context, templates, common interfaces.

## Layers

### a. Domain Layer (`src/<module>/domain`)

- **Purpose**: Contains the core business logic, entities, and rules. It is the heart of the application.
- **Contents**: Each entity is encapsulated in its own package with its struct, constructor, validation, and methods. Repository interfaces define contracts for data access.
- **Rules**:
  - Has **zero** dependencies on any other layer. It is pure business logic.
  - Does not know about Fiber, HTTP, databases, or any external framework.
  - Example: An `Account` entity defines rules about balances and ownership — without knowing about REST or HTMX.

### b. Application Layer (`src/<module>/application`)

- **Purpose**: Orchestrates the flow of data and commands between the domain and infrastructure. Contains use cases.
- **Contents**:
  - **Use Cases**: Encapsulate a specific operation (e.g., `AccountCreator`, `IncomeCreator`). Each use case has its own package containing the use case logic and its command.
  - **Commands** (`command.go`): Immutable DTOs used to carry data into the use cases. Co-located with their use case.
- **Rules**:
  - Depends only on the **Domain** layer.
  - Receives dependencies (like repository interfaces) via constructor injection.
  - Does not handle HTTP requests or framework-specific details.

### c. Infrastructure Layer (`src/<module>/infrastructure`)

- **Purpose**: Implements the details of how the application interacts with the outside world.
- **Contents**:
  - **API Handlers** (`infrastructure/api/handlers`): Expose use cases via HTTP. Handlers adapt incoming requests into application commands and format output. Each entity has handler variants for both HTMX and REST responses.
  - **Repositories** (`infrastructure/repositories`): Concrete implementations of domain repository interfaces (currently in-memory, planned migration to PostgreSQL).
- **Rules**:
  - Depends on the **Application** and **Domain** layers.
  - Handles all interactions with Fiber, databases, and external services.

## Dual Interface Pattern

Odin serves two types of clients from the same domain logic:

- **HTMX (Web)**: Handlers render `.gohtml` templates and return HTML fragments. Routes are at root level (e.g., `/accounts`, `/categories`).
- **REST (Mobile)**: Handlers return JSON responses. Routes are under `/api/v1/` (e.g., `/api/v1/accounts`).

Both handler types delegate to the same use cases and domain entities. The strategy pattern (e.g., `LoginHandler` interface) allows the same orchestration logic to produce different response formats.

## Data Flow

```
HTTP Request
    → Middleware [Infrastructure] (auth, session)
        → Handler [Infrastructure] (parse request, create command)
            → Use Case [Application] (orchestrate business logic)
                → Domain Entity (validate, enforce rules)
                → Repository Interface [Domain] → Repository Impl [Infrastructure]
            ← Result
        ← Response (HTML template or JSON)
```

## Error Handling

The `odinerrors` package provides tagged errors that map to HTTP status codes:

- `DOMAIN` → 400 Bad Request
- `NOT_FOUND` → 404 Not Found
- `RENDER` → 500 Internal Server Error (template rendering failures)
- `UNKNOWN` → 500 Internal Server Error

Errors have two message fields: `message` (internal, English, for logs) and `external` (user-facing, Spanish, for API responses).

## Design Principles

1. **Clean Architecture:** Dependencies point inward. Domain has zero external dependencies. Infrastructure implements domain interfaces.
2. **Dual Interface:** Same business logic serves both web (HTMX) and mobile (REST) clients.
3. **Stateless Handlers:** Handlers are thin adapters. Business logic lives in use cases and domain entities.
4. **Dependency Injection:** Components receive their dependencies through constructors, enabling testability.
5. **No Empty Packages:** Only create directories and packages when there is code to put in them.
