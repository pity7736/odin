# Feature: Listing All Items in the Vault

## Overview

When a user opens the app on a device, they need all of their stored data at
once so the device can rebuild their accounts, income, and expenses. This feature
returns every item the user has saved, each exactly as it was stored. The service
hands the items back without ever being able to read them, and only ever returns
the requesting user's own items.

## User Stories

### Loading all my data
As a user, I want to retrieve everything I have saved in one request, so that my
device can rebuild my full picture of my finances.

### Reaching only my own data
As a user, I want to receive only the items that belong to me, so that no one
else's data is ever exposed to me and mine is never exposed to anyone else.

## Acceptance Criteria

- I must be logged in to list my items; otherwise the request is rejected.
- I receive every item that belongs to me, and only items that belong to me.
- Each item comes back exactly as it was stored, including everything my device
  needs to work with it.
- The items are ordered with the most recently created first.
- If I have not saved any items, I receive an empty collection — this is a
  success, not an error.

## Expected Behavior

### Listing my items successfully
- Given I am logged in
- And I have saved several items
- When I list all my items
- Then I receive all of them, each exactly as it was stored
- And they are ordered with the most recently created first

### Listing when I have no items
- Given I am logged in
- And I have not saved any items
- When I list all my items
- Then I receive an empty collection
- And the request succeeds

### Listing returns only my own items
- Given I am logged in
- And other users have saved items of their own
- When I list all my items
- Then I receive only my items
- And no item belonging to another user is ever included

### Listing without being logged in
- Given I am not logged in
- When I try to list items
- Then no items are returned
- And I am rejected

## Out of Scope

- Loading items in pages for very large collections (separate feature).
- Returning only the items that changed since a given moment (separate feature).
- Searching or filtering items by their contents.
