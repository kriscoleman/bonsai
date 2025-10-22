# Dockerfile for Bonsai GitHub Action
# Multi-stage build for minimal image size

FROM golang:1.24-alpine AS builder

WORKDIR /build

# Install git (required for bonsai to work)
RUN apk add --no-cache git

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o bonsai ./cmd/bonsai

# Final minimal image
FROM alpine:latest

# Install git (required for bonsai to execute git commands)
RUN apk add --no-cache git bash

# Copy the binary from builder
COPY --from=builder /build/bonsai /usr/local/bin/bonsai

# Copy the entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# Set the entrypoint
ENTRYPOINT ["/entrypoint.sh"]
