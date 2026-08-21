# CityWalk BUG-001 Delivery

## Project

CityWalk is a city-walk and cultural guide platform built with Go, MySQL, Redis, Vue 3, and Docker Compose.

## Standard Commands

Start dependencies and services:

```bash
docker compose up -d
```

Run the BUG-001 verification command:

```bash
go test ./scripts/verify -count=1 -run '^TestBug001_BusinessRegression$'
```

Build the project locally:

```bash
go build ./...
go vet ./...
```

Build the delivery Docker image:

```bash
bash build_benzhi_docker.sh
```

## BUG-001 Reproduction

1. Edit an existing route that already has waypoints.
2. Submit a new title and waypoint list where one waypoint is invalid.
3. Observe that the request fails, but the route title changes and the old waypoints disappear.
4. Submit a valid update with a custom waypoint order.
5. Observe that the saved order does not match the order sent by the client.

## Expected Behavior

- A failed route edit must leave the route and its waypoints unchanged.
- A successful route edit must preserve the client-provided waypoint order.
