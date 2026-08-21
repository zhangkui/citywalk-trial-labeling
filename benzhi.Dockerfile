FROM golang:1.22-bookworm

WORKDIR /workspace

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 8080

CMD go run ./cmd/api
