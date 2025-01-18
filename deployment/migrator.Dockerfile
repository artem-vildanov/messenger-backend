FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git build-base

WORKDIR /app
ARG ENV=dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .
COPY ./config/.env.${ENV} ./config/.env

RUN go build -o ./bin/migrator ./cmd/migrator/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/bin/migrator ./bin/migrator
COPY --from=builder /app/config/.env ./.env
COPY --from=builder /app/migrations ./migrations

CMD ["./bin/migrator"]