# syntax=docker/dockerfile:1

# Builder
FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -o /app

# Runner
FROM gcr.io/distroless/static-debian11 AS runner

WORKDIR /app

COPY --from=builder /app /

CMD ["/meilisearch-syncer"]
