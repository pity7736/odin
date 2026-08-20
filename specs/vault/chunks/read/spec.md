# Feature: Reading an Item from the Vault

## Overview

Once a user has saved an item to the vault, they need to get it back so their
device can decrypt it and show them their data. This feature returns one saved
item, exactly as it was stored. The service hands back the item without ever
reading it: the contents stay unreadable to the service, and only the owner can
retrieve them.

## User Stories

### Getting my data back
As a user, I want to retrieve an item I previously saved, so that my device can
decrypt it and show me my data.

### Reaching only my own data
As a user, I want to retrieve only the items that belong to me, so that no one
else's data is ever exposed to me and mine is never exposed to anyone else.

## Acceptance Criteria

- I must be logged in to retrieve an item; otherwise the request is rejected.
- I identify the item I want by its identifier, one item at a time.
- When I retrieve an item that belongs to me, I get it back exactly as it was
  stored, including everything my device needs to decrypt it.
- If no item exists with that identifier, I am told it was not found.
- If the item exists but belongs to another user, I am told it was not found —
  I am never told it exists.
- If the identifier is not a valid identifier at all, I am told it was not found.

## Expected Behavior

### Retrieving my item successfully
- Given I am logged in
- And an item I saved exists with the identifier I provide
- When I retrieve the item by that identifier
- Then I get the item back exactly as it was stored

### Retrieving without being logged in
- Given I am not logged in
- When I try to retrieve an item
- Then the item is not returned
- And I am rejected

### Retrieving an item that does not exist
- Given I am logged in
- And no item exists with the identifier I provide
- When I try to retrieve the item by that identifier
- Then I am told the item was not found

### Retrieving an item that belongs to another user
- Given I am logged in
- And an item with the identifier I provide exists but belongs to another user
- When I try to retrieve the item by that identifier
- Then I am told the item was not found
- And I am never told the item exists

### Retrieving with an invalid identifier
- Given I am logged in
- When I try to retrieve an item using an identifier that is not valid
- Then I am told the item was not found

## Out of Scope

- Decrypting the item — that happens on the user's device, after it is returned.
- Retrieving more than one item at once, listing, or searching items (separate
  features).
- Updating or deleting a saved item (separate features).
- Syncing items across devices and paging the initial load (separate features).
