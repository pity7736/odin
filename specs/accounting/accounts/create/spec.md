# Feature: Create Account

## Overview
Lets a user register one of their financial accounts so its balance and, later,
its transactions can be tracked. Registering an account is the starting point
for managing money in Odin: everything else builds on the accounts a user owns.

## User Story
As a user, I want to register a financial account with its name, type, currency
and current balance, so that I can start tracking my money in Odin.

## Acceptance Criteria
- A user must be logged in to create an account.
- An account always belongs to the user who creates it; the owner is never
  chosen by the user.
- The account name is required, cannot be blank, and can be at most 255
  characters long.
- Among that user's own accounts, the combination of name and currency must be
  unique: a user may reuse the same name across different currencies (e.g. a
  multi-currency wallet held as one account per currency), but not have two
  accounts with the same name in the same currency. Two different users may each
  have an account with the same name and currency.
- The account type is required and must be one of: savings, credit card, cash.
- The account currency is required and must be one of: COP, USD.
- The initial balance is required, cannot be negative, and may include decimals.
- At creation time the account's balance equals its initial balance.
- The system assigns the account a unique identifier.
- The system records the date the account was created.
- When creation succeeds, all of the account's data is returned to the user.

## Expected Behavior

### Successfully create an account
- Given a logged-in user with no account named "Bancolombia Ahorros"
- When the user creates an account with name "Bancolombia Ahorros", type
  savings, currency COP and initial balance 1500000.50
- Then the account is registered and belongs to that user
- And its balance equals its initial balance of 1500000.50
- And it is given a unique identifier and a creation date
- And all of the account's data is returned to the user

### Reject when the user is not logged in
- Given a user who is not logged in
- When the user tries to create an account
- Then the account is not created
- And the user is told they must be logged in

### Reject a blank name
- Given a logged-in user
- When the user tries to create an account whose name is empty or only spaces
- Then the account is not created
- And the user is told the name is required

### Reject a name that is too long
- Given a logged-in user
- When the user tries to create an account whose name is longer than 255
  characters
- Then the account is not created
- And the user is told the name is too long

### Reject a duplicate name in the same currency for the same user
- Given a logged-in user who already owns an account named "Global66" in USD
- When the user tries to create another account named "Global66" in USD
- Then the second account is not created
- And the user is told they already have an account with that name in that
  currency

### Allow the same name in a different currency for the same user
- Given a logged-in user who already owns an account named "Global66" in USD
- When the user creates an account named "Global66" in COP
- Then the account is created for that user

### Allow the same name and currency for different users
- Given two different users
- And the first user owns an account named "Global66" in USD
- When the second user creates an account named "Global66" in USD
- Then the account is created for the second user

### Reject a missing or invalid type
- Given a logged-in user
- When the user tries to create an account with no type, or with a type that is
  not one of savings, credit card or cash
- Then the account is not created
- And the user is told the type is invalid

### Reject a missing or invalid currency
- Given a logged-in user
- When the user tries to create an account with no currency, or with a currency
  that is not one of COP or USD
- Then the account is not created
- And the user is told the currency is invalid

### Reject a missing initial balance
- Given a logged-in user
- When the user tries to create an account without an initial balance
- Then the account is not created
- And the user is told the initial balance is required

### Reject a negative initial balance
- Given a logged-in user
- When the user tries to create an account with an initial balance below zero
- Then the account is not created
- And the user is told the initial balance cannot be negative

## Out of Scope
- Listing, viewing, editing or deleting accounts.
- Credit-card capacity/limit, and what a credit-card balance represents (amount
  owed vs. available credit). Deferred to the credit-card spending feature.
- Adding account types beyond savings, credit card and cash.
- A single account that holds balances in multiple currencies at once (a
  multi-currency wallet). For now such a wallet is represented as one account
  per currency. Deferred to a future feature.
- Adding currencies beyond COP and USD.
- Converting balances between currencies.
