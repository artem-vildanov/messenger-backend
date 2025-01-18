FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .
COPY ./config/.env.test ./config/.env

RUN go build -o ./bin/main ./cmd/service/main.go
RUN go test -c -o ./bin/tests ./tests

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/main ./bin/main
COPY --from=builder /app/bin/tests ./bin/tests
COPY --from=builder /app/config/.env ./tests/.env
COPY --from=builder /app/config/.env ./.env

CMD ["./bin/main"]
