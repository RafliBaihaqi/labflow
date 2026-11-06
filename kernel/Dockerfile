# syntax=docker/dockerfile:1.6
ARG GO_VERSION=1.26.0
ARG ALPINE_VERSION=3.21
ARG MIGRATE_VERSION=v4.18.1

# Build stage
FROM golang:${GO_VERSION}-alpine AS builder
WORKDIR /app

# Cache go mod downloads
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -trimpath -o main .

# Fetch the golang-migrate CLI binary
FROM alpine:${ALPINE_VERSION} AS migrate
ARG MIGRATE_VERSION
RUN apk --no-cache add curl tar && \
    curl -fsSL "https://github.com/golang-migrate/migrate/releases/download/${MIGRATE_VERSION}/migrate.linux-amd64.tar.gz" \
        | tar xz -C /usr/local/bin migrate && \
    chmod +x /usr/local/bin/migrate

# Runtime stage
FROM alpine:${ALPINE_VERSION}

LABEL org.opencontainers.image.title="labflow-backend"
LABEL org.opencontainers.image.description="Laboratory Information Management System - Go backend"
LABEL org.opencontainers.image.source="https://github.com/RafliBaihaqi/labflow"

WORKDIR /app

RUN apk --no-cache add ca-certificates wget curl && \
    addgroup -g 1000 -S appgroup && \
    adduser -u 1000 -S appuser -G appgroup

COPY --from=builder --chown=appuser:appgroup /app/main .
COPY --from=migrate /usr/local/bin/migrate /usr/local/bin/migrate
COPY --chown=appuser:appgroup sql/migrations /app/sql/migrations
COPY --chown=appuser:appgroup entrypoint.sh /app/entrypoint.sh
RUN chmod +x /app/entrypoint.sh

USER appuser

EXPOSE 8080

# Use GET not HEAD -- Gin GET route does not handle HEAD via wget --spider (which sends HEAD)
HEALTHCHECK --interval=10s --timeout=10s --start-period=60s --retries=8 \
    CMD curl -f http://localhost:8080/health || wget -qO- http://localhost:8080/health | grep -q ok || exit 1

ENTRYPOINT ["/app/entrypoint.sh"]
