# Feature: Saving an Item to the Vault

## Overview

Everything a user creates in the app — an account, an income, an expense — needs
to be saved beyond their device so it persists and can later reach their other
devices. This feature stores one such item on the service. The item is saved in a
form the service cannot read: the service keeps it safe and tied to its owner, but
never learns what it contains.

## User Stories

### Saving my data
As a user, I want the items I create in the app to be saved to the service, so
that they persist beyond the device I created them on.

### Keeping my data private from the service
As a user, I want the service to store my items without being able to read them,
so that my finances stay private even from the service itself.

### Owning my data
As a user, I want each item I save to be tied to me, so that no one else can reach
it.

## Acceptance Criteria

- I must be logged in to save an item; otherwise the request is rejected.
- Each item I save belongs to me.
- The item is stored exactly as my device provides it, and the service cannot read
  its contents.
- My device provides both the item's identifier and its contents.
- If the identifier is missing, the item is not saved.
- If the contents are missing, the item is not saved.
- If an item with the same identifier already exists, it is not saved again —
  saving never overwrites existing data.

## Expected Behavior

### Saving an item successfully
- Given I am logged in
- And no item exists with the identifier my device chose
- When I save an item with that identifier and its contents
- Then the item is stored and tied to me
- And I am told it was saved

### Saving without being logged in
- Given I am not logged in
- When I try to save an item
- Then the item is not saved
- And I am rejected

### Saving an item whose identifier already exists
- Given an item already exists with the identifier my device chose
- When I try to save another item with that same identifier
- Then the new item is not saved
- And the existing item is left unchanged
- And I am told an item with that identifier already exists

### Saving without an identifier
- Given I am logged in
- When I try to save an item without an identifier
- Then the item is not saved
- And I see an error

### Saving without contents
- Given I am logged in
- When I try to save an item with no contents
- Then the item is not saved
- And I see an error

## Out of Scope

- Reading a saved item back (separate feature).
- Updating or deleting a saved item (separate features).
- Syncing items across devices and paging the initial load (separate features).
- Encrypting the item — that happens on the user's device, before it is saved.
- Understanding or organizing the item's contents by type — the service never
  reads them; that happens on the device after retrieval.
