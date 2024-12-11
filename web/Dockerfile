FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git make gcc musl-dev

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o beacon-web ./cmd

# Use the official Alpine image for a small production image
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/beacon-web .
COPY --from=builder /app/config.toml .

EXPOSE 3000
ENTRYPOINT ["./beacon-web"]
