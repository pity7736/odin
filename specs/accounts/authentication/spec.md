# Feature: Authentication

## Overview

Users must identify themselves before accessing their financial data. Authentication controls who can enter the application and ensures that each user only sees their own information. The user's password never leaves their device — the application verifies identity without the server ever learning the password.

## User Stories

### Logging in
As a user, I want to log in with my email and password so that I can access my financial data.

### Receiving my encrypted data access
As a user, I want to receive everything I need to access my data after logging in so that my device can unlock my financial information in a single step.

### Staying logged in
As a user, I want to remain logged in while I actively use the application so that I don't have to re-enter my credentials constantly.

### Session expiration
As a user, I want my session to expire after a period of inactivity so that someone who finds my device can't access my finances.

### Logging out
As a user, I want to log out so that no one else can access my financial data from this device.

### Accessing protected features
As a user, I expect that all my financial data is only accessible after I log in.

## Acceptance Criteria

### Login
- I can log in with a valid email and password.
- If my email doesn't exist, I see an error message: "Correo o contraseña incorrectos".
- If my password is wrong, I see the same error message (no hint about which field is wrong).
- If I submit an empty email or password, I see an error message.
- My password never leaves my device.
- After logging in, I receive my session credentials, my encrypted data key, and the parameters my device needs to unlock it.

### Session
- After logging in, I stay logged in across interactions without re-entering my credentials.
- If I stop using the application for 30 days, I must log in again.
- Each time I use the application, the 30-day window resets from that moment.

### Logout
- I can log out from the application.
- Logging out ends the exact session I am currently using; my credentials stop working immediately.
- Logging out never reports success unless it actually ended my session.
- After logging out, I can no longer access any protected feature.

### Protected features
- If I try to access any financial feature without being logged in, I am rejected.
- The rejection for unauthorized access is consistent regardless of which feature I tried to reach.

## Expected Behavior

### Logging in successfully
- Given I have an account with email "julian@example.com"
- When I log in with the correct email and password
- Then I receive my session credentials, my encrypted data key, and the parameters to unlock it
- And I can access my financial data

### Logging in with wrong credentials
- Given I have an account with email "julian@example.com"
- When I log in with the wrong password
- Then I see "Correo o contraseña incorrectos"
- And I am not logged in

### Logging in with a non-existing email
- Given no account exists with email "unknown@example.com"
- When I try to log in with that email
- Then I see "Correo o contraseña incorrectos"
- And I am not logged in

### Logging in with missing or empty fields
- Given I am on the login screen
- When I submit without an email
- Then I see "El correo es obligatorio"
- When I submit without a password
- Then I see "La contraseña es obligatoria"

### Session stays active while I use the app
- Given I logged in 25 days ago
- When I use the application today
- Then I remain logged in
- And my session now expires 30 days from today, not from the original login

### Session expires after inactivity
- Given I logged in 31 days ago
- And I have not used the application since
- When I try to access my financial data
- Then I am rejected and must log in again

### Logging out
- Given I am logged in
- When I log out
- Then my session ends
- And I can no longer access my financial data

### Logging out actually ends the session I am using
- Given I am logged in
- When I log out
- Then the session I was using stops working immediately
- And trying to log out again does not report success, because there is no session left to end

### Accessing a protected feature without logging in
- Given I am not logged in
- When I try to access any financial feature
- Then I am rejected

## Out of Scope
- User registration (separate feature)
- Password recovery or reset
- Recovery keys
- Device sync and key sharing between devices
- Password change (involves re-encrypting the data key)
- Multiple simultaneous sessions per user
- Account lockout after failed attempts
