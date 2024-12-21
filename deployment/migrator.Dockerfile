FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o ./bin/migrator ./cmd/migrator/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/bin/migrator /app/bin/migrator
COPY --from=builder /app/config/.env /app/config/.env
COPY --from=builder /app/migrations /app/migrations

EXPOSE 8080

CMD ["./bin/migrator"]