# BUG-001 Reproduction

## Issue

Editing a route can leave partial data behind when waypoint replacement fails, and successful updates can lose the waypoint order provided by the client.

## Steps

1. Create or use a route that already has existing waypoints.
2. Submit an edit request with a new title and a waypoint list containing one invalid waypoint.
3. Query the route details after the request fails.
4. Submit another valid edit request whose waypoint `order` values are intentionally not the same as the input array index.
5. Query the route details again.

## Actual Result

- After the failed edit, the route title is already updated and the previous waypoints are gone.
- After the successful edit, the persisted waypoint order does not match the client request.

## Expected Result

- If any waypoint insert fails, the full edit must roll back without changing the route or its existing waypoints.
- If the request succeeds, the stored waypoint order must match the client-provided `order`.
