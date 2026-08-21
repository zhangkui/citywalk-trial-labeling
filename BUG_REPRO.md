# BUG-002 Reproduction

## Issue

The route favorite endpoint can report success while the persisted favorite relation and the route display aggregate remain inconsistent.

## Steps

1. Create a route with `favorite_count` set to zero.
2. Insert one route favorite relation for a user, leaving the aggregate stale.
3. Call the route favorite endpoint again for the same user and route.
4. Read the route aggregate and favorite relation count.
5. Repeat the request when both values are already correct.

## Actual Result

- The stale aggregate is not reconciled when the relation already exists.
- Repeated repository mutations can change the aggregate independently of whether the relation changed.
- The route detail page can retain an old displayed count after a successful request.

## Expected Result

- Every successful favorite operation must reconcile `favorite_count` from the actual relation rows.
- Duplicate requests must leave exactly one relation and the correct aggregate.
- The route detail display must refresh from the authoritative backend state.
