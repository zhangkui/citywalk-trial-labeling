# CityWalk BUG-003 Delivery

## Project

CityWalk is a city-walk and cultural guide platform built with Go, MySQL, Redis, Vue 3, and Docker Compose.

## Standard Commands

```bash
docker compose up -d
go test ./scripts/verify -count=1 -run '^TestBug003_BusinessRegression$'
go build ./...
go vet ./...
bash build_benzhi_docker.sh
```

## Expected Behavior

- Route list and total use identical city, theme, and minimum-rating filters.
- An invalid minRate query returns HTTP 400 instead of silently disabling filtering.
- Valid combined filters and pagination remain available.
