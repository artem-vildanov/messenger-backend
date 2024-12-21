FROM golang:1.23 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN go build -o ./bin/main ./cmd/service/main.go

FROM debian:bookworm

WORKDIR /app

COPY --from=builder /app/bin/main /app/bin/main
COPY --from=builder /app/config/.env /app/config/.env

EXPOSE 8080

CMD ["./bin/main"]