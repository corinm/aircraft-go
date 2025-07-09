FROM golang:1.23.10-alpine AS builder

RUN apk add --no-cache git
WORKDIR /app

RUN go install github.com/pressly/goose/v3/cmd/goose@latest

FROM alpine:3.20

RUN apk add --no-cache ca-certificates postgresql-client

COPY --from=builder /go/bin/goose /usr/local/bin/goose

WORKDIR /migrations

COPY ./data/sql/migrations ./
COPY .env ./

ENTRYPOINT ["goose", "up"]
