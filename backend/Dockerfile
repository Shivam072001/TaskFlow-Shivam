FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.mod ./
COPY . .
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /server ./cmd/server
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /seeder ./scripts/seed

FROM alpine:3.19

RUN apk add --no-cache ca-certificates curl

ARG MIGRATE_VERSION=4.17.0
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz \
    | tar xz -C /usr/local/bin

WORKDIR /app

COPY --from=builder /server /app/server
COPY --from=builder /seeder /app/seeder
COPY migrations/ /app/migrations/
COPY entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

EXPOSE 8080

ENTRYPOINT ["/app/entrypoint.sh"]
