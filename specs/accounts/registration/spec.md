# Feature: Registration

## Overview

A new person needs an account before they can use Odin. Registration creates that
account from an email and a password, so the person can later log in and access
their financial data. As with logging in, the password never leaves the person's
device — the app creates the account without the server ever learning the
password.

## User Stories

### Creating an account
As a new user, I want to create an account with my email and password, so that I
can start using the application and later log in to reach my financial data.

### Keeping my password private
As a new user, I want my password to never leave my device when I register, so
that no one — not even the service — can learn it.

### Being told my email is already taken
As a new user, I want to be clearly told when my email already has an account, so
that I know to log in instead of registering again.

## Acceptance Criteria

- I can create an account with a valid email and password.
- My password never leaves my device during registration.
- When I create my account, the app also stores everything it needs to unlock my
  financial data on future logins, so that my next step is simply to log in.
- If I register with an email that already has an account, I am clearly told the
  email is already registered.
- If I submit an empty email, I see an error message.
- If I submit an empty password, I see an error message.
- If the information my device provides to unlock my data is incomplete or
  invalid, my account is not created.
- Creating an account does not log me in — after registering, I log in as a
  separate step.

## Expected Behavior

### Creating an account successfully
- Given no account exists with email "julian@example.com"
- When I register with that email and a password
- Then my account is created
- And I am told the account was created successfully
- And I can now log in with that email and password

### Registering with an email that already exists
- Given an account already exists with email "julian@example.com"
- When I try to register again with that email
- Then my account is not created
- And I am told "El correo ya está registrado"

### Registering with a missing or empty email
- Given I am creating an account
- When I submit without an email
- Then my account is not created
- And I see "El correo es obligatorio"

### Registering with a missing or empty password
- Given I am creating an account
- When I submit without a password
- Then my account is not created
- And I see "La contraseña es obligatoria"

### Registering with incomplete unlock information
- Given I am creating an account
- When the information my device provides to unlock my data is missing or invalid
- Then my account is not created
- And I see an error message

### Registering does not log me in
- Given I successfully register a new account
- When registration completes
- Then I am not logged in
- And I must log in as a separate step to access my financial data

## Out of Scope

- Email verification (confirming ownership of the email by sending a message to
  it) — a separate, later feature.
- Logging in after registering — handled by the Authentication feature.
- Password recovery or reset.
- Recovery keys.
- Device sync and key sharing between devices.
- Protection against email enumeration (e.g. rate limiting) — a later hardening
  concern.
- The creation of the unlock information on the device (key derivation, data-key
  generation) — performed on the user's device, outside the service.
