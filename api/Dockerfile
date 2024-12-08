FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o ./beacon-web -v ./cmd/...

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/beacon-web .
COPY --from=builder /app/config.toml .

ENTRYPOINT ["./beacon-web"]
