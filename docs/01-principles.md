# Odin: Development Principles

Odin is a personal finance management application that helps users track accounts, incomes, expenses, transfers, and categories. It provides a dual interface: HTMX-powered web pages and a REST API for mobile clients.

Our development process is founded on three core principles:

- **[Specification-Driven Development (SDD)](./03-sdd-workflow.md):** We define what we will build before we build it.
- **[Test-Driven Development (TDD)](./04-tdd-workflow.md):** We verify behavior with tests before implementing it.
- **[Modular Architecture](./02-architecture.md):** We create a decoupled, testable, and maintainable codebase with clear boundaries between components.

We use **Go 1.26** and **Fiber** for this project.

This documentation is organized into several parts:

- **[Architecture](./02-architecture.md)**
- **[Specification-Driven Development (SDD) Workflow](./03-sdd-workflow.md)**
- **[Test-Driven Development (TDD) Workflow](./04-tdd-workflow.md)**
- **[Code Standards](./05-code-standards.md)**
- **[Pillars of Quality](./06-quality-pillars.md)**
