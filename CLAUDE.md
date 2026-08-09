# Odin - Personal Finance

Personal finance management application with dual interface: HTMX for web, REST API for mobile. Built with **Go 1.26** and **Fiber**.

## Quick Reference

| Topic | Documentation |
|-------|---------------|
| Core Principles | [docs/01-principles.md](docs/01-principles.md) |
| Architecture | [docs/02-architecture.md](docs/02-architecture.md) |
| SDD Workflow | [docs/03-sdd-workflow.md](docs/03-sdd-workflow.md) |
| TDD Workflow | [docs/04-tdd-workflow.md](docs/04-tdd-workflow.md) |
| Code Standards | [docs/05-code-standards.md](docs/05-code-standards.md) |
| Quality Pillars | [docs/06-quality-pillars.md](docs/06-quality-pillars.md) |

## Critical Rules

### Communication

- **When the user asks a question, ONLY answer the question.** Do not take action, do not make changes, do not run fixes. Just answer. This is a strict rule with no exceptions.
- **NEVER make assumptions.** If something is unclear or ambiguous, ASK the user instead of assuming. This is a strict rule with no exceptions.
- **Raise concerns and challenge decisions.** When you identify potential issues or have concerns about a technical decision, speak up.
- **Correct the user immediately when they are wrong.** Do not let them stay wrong to avoid conflict. State the correction first, no softening.

### Token Efficiency

- **Prefer Explore agent for codebase searches.** Use the Agent tool with `subagent_type=Explore` for open-ended searches, understanding patterns, or exploring multiple files. This saves context tokens.
- **Use direct tools only when necessary.** Grep and Read should be used for specific, targeted queries where you know exactly what you're looking for.

### Development Flow

1. **READ DOCS FIRST:** Before implementing ANY feature, READ these files - This is MANDATORY:
   - `docs/02-architecture.md` - Understand the component structure
   - `docs/05-code-standards.md` - Follow coding conventions
2. **Spec First:** Create `specs/<feature>/spec.md` and `plan.md` before coding
   - **Specs must be business-focused, NOT technical** - Product managers must understand them
   - **NO technical terms in specs**: No Fiber, HTMX, Postgres, HTTP, API endpoints, etc.
   - Use business language: "users", "accounts", "income", "balance", etc.
   - Specs must include "Expected Behavior" section with Given/When/Then scenarios
   - **ALL technical details go in plan.md ONLY**
3. **TDD Always:** Red-Green-Refactor for all business logic
4. **Clean Architecture:** Domain -> Application -> Infrastructure (dependencies point inward)

### Code Standards

- Private fields by default, public only for DTOs
- Struct names MUST be descriptive
- **Variable names must describe role, never type and never abbreviated.** When a single instance exists, the type name may double as the role (e.g. `account`). When multiple instances of the same type exist in the same scope, each must be named by what it represents, not what it is (e.g. `savingsAccount` / `checkingAccount`, never `account1` / `account2`).
- Receiver name: `self`
- 100% test coverage for business logic
- **No comments in source code** - code must be self-documenting
- **Code Organization:** Constructors MUST be defined immediately after the struct definition. Methods MUST be defined immediately after the constructor(s).
- All files end with a trailing newline
- **Error messages:** Internal (`message`) in English, external/user-facing (`external`) in Spanish

### Tooling

- **Always use Makefile targets.** Do not run manual commands when a Makefile target exists.
- **Tests:** `make test` to run all tests
- **Coverage:** `make coverage` to see coverage, `make coverage-check` to verify threshold
- **Lint:** `make lint` to run linters
- **Full check:** `make check` runs lint + test + coverage-check
- **Mocks:** mockery — configured in `.mockery.yaml`, run with `go run github.com/vektra/mockery/v3`

### Testing

See `docs/05-code-standards.md` section 3 for all testing rules.

### Security

- Validate input at handlers
- Never hardcode secrets
- Use `strings.Clone()` when storing data parsed from Fiber request bodies to prevent fasthttp buffer reuse corruption
