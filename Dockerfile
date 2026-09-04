# syntax=docker/dockerfile:1
ARG GO_VERSION=1.23
# Note: User had 1.25.5-alpine, but 1.26.1 matches local go.mod.
ARG IMAGE_TAG=1.24-alpine

################################################################################
# Base stage for shared dependencies
FROM golang:${IMAGE_TAG} AS base
ENV GOTOOLCHAIN=auto
RUN apk add --no-cache git make
WORKDIR /app
COPY go.mod go.sum ./
# Implement build cache for go modules
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

################################################################################
# Tools stage for installing dependencies
FROM base AS tools
# Implement build cache for go installs
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go install github.com/air-verse/air@latest && \
    go install github.com/swaggo/swag/cmd/swag@latest

################################################################################
# Development stage
FROM tools AS dev
COPY . .
# Docs generation for development
# RUN swag init -g main.go --parseDependency --parseInternal \
#     --dir ./cmd/api,./internal/features/auth,./internal/features/users,./internal/features/appointments,./internal/features/slips,./internal/features/students \
#     --output ./docs/internal_docs --instanceName internal
# RUN swag init -g main.go --parseDependency --parseInternal \
#     --dir ./cmd/api,./internal/features/students/integrations,./internal/features/m2mclients \
#     --output ./docs/integrations --instanceName integrations
CMD ["air"]

################################################################################
# Builder stage for production
FROM tools AS builder
COPY . .
RUN swag init -g main.go --parseDependency --parseInternal \
    --dir ./cmd/api,./internal/features/auth,./internal/features/users,./internal/features/appointments,./internal/features/slips,./internal/features/students \
    --output ./docs/internal_docs --instanceName internal
RUN swag init -g main.go --parseDependency --parseInternal \
    --dir ./cmd/api,./internal/features/students/integrations,./internal/features/m2mclients \
    --output ./docs/integrations --instanceName integrations
# Implement build cache to speed up successive builds
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=linux go build -o main ./cmd/api

################################################################################
# Production stage
FROM alpine:latest AS prod
RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Manila
WORKDIR /root/
COPY --from=builder /app/main .
COPY --from=builder /app/scripts/migrations ./scripts/migrations/
# COPY --from=builder /app/.env.example ./.env
# Note: You should ideally provide a real .env or use environment variables in compose
EXPOSE 8080
CMD ["./main"]
