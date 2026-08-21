# CityWalk BUG-004 Delivery

## Standard Commands

```bash
docker compose up -d
go test ./scripts/verify -count=1 -run '^TestBug004_BusinessRegression$'
go build ./...
go vet ./...
bash build_benzhi_docker.sh
```

## Expected Behavior

- Story content and all replacement media update atomically.
- A failed media insert leaves the original story and media unchanged.
- Client-provided media order is persisted and returned unchanged.
