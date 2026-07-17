# syntax=docker/dockerfile:1

FROM golang:1.24.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -buildvcs=false \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/fuse \
    ./cmd/api/main.go

FROM alpine:3.22

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/fuse /app/fuse
COPY --from=builder /src/configs /app/configs

USER app

EXPOSE 3000

ENTRYPOINT ["/app/fuse"]