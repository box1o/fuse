# syntax=docker/dockerfile:1

FROM golang:1.24.3-alpine AS builder

WORKDIR /src

RUN apk add --no-cache ca-certificates git

COPY go.mod go.sum ./
RUN go mod download

COPY cmd ./cmd
COPY internal ./internal
COPY pkg ./pkg

RUN CGO_ENABLED=0 GOOS=linux go build \
    -buildvcs=false \
    -trimpath \
    -ldflags="-s -w" \
    -o /out/fuse ./cmd/api/main.go

FROM alpine:3.22 AS runtime

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app \
    && adduser -S -G app app

COPY --from=builder /out/fuse /app/fuse
COPY configs /app/configs

USER app

EXPOSE 3000

HEALTHCHECK --interval=10s --timeout=3s --start-period=15s --retries=5 \
    CMD wget --quiet --tries=1 --spider http://127.0.0.1:3000/health || exit 1

ENTRYPOINT ["/app/fuse"]
