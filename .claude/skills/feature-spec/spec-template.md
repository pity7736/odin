<!--
CANONICAL SPEC FORMAT. This template is the single source of truth for the
structure of a spec.md. Copy it into specs/<module>/<feature>/spec.md and fill
every section, then DELETE all HTML comments.

RULES (from docs/01-principles.md and CLAUDE.md):
- Business-focused, NOT technical. A Product Manager must understand it.
- FORBIDDEN words: Fiber, HTMX, Postgres, HTTP, REST, API, endpoint, handler,
  repository, database, JSON, form, request, response, status code, use case.
- Use business language: users, accounts, income, balance, log in, etc.
- Describe intended behavior — what the feature SHOULD do. If current code has
  bugs, the spec is the source of truth for correct behavior, not the code.
-->

# Feature: <feature name in plain language>

## Overview
<!-- 1-3 sentences: why this matters to a user, in business terms. -->

## User Stories
<!-- One or more. A feature almost always spans several stories (e.g. logging
     in, staying logged in, logging out); give each its own named block. -->

### <story name>
As a <user>, I want to <goal>, so that <benefit>.

## Acceptance Criteria
<!-- Bullet list of what must be true for the feature to be considered done.
     Each bullet is a business rule, not a technical step. -->
- ...

## Expected Behavior
<!-- One block per scenario, Given/When/Then. Cover the happy path AND every
     rejection/edge case the user can hit. -->

### <scenario name>
- Given <starting situation>
- When <the user does something>
- Then <the observable outcome>
- And <additional outcome, optional>

## Out of Scope
<!-- List what this feature deliberately does NOT cover, so reviewers know the
     boundaries. -->
- ...
