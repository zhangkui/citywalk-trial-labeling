# CityWalk BUG-002 Delivery

## Project

CityWalk is a city-walk and cultural guide platform built with Go, MySQL, Redis, Vue 3, and Docker Compose.

## Standard Commands

Start dependencies and services:

```bash
docker compose up -d
```

Run the BUG-002 verification command:

```bash
go test -race ./scripts/verify -count=10 -run '^TestBug002_BusinessRegression$'
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

## BUG-002 Reproduction

1. Create a route whose favorite aggregate is zero.
2. Insert an existing favorite relation for the route without changing the aggregate.
3. Repeat the favorite request for the same user and route.
4. Observe that the API succeeds while the displayed aggregate remains stale.
5. Repeat an already-correct favorite request and verify that neither the relation nor aggregate is duplicated.

## Expected Behavior

- A repeated favorite request must reconcile the aggregate from actual favorite relations.
- An already-correct repeated request must remain idempotent.
- Concurrent or repeated execution must not create duplicate relations or aggregate drift.
