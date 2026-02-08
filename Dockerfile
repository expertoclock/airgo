# Stage 1: Builder
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install git for fetching dependencies
RUN apk add --no-cache git

# CRITICAL FIX: Copy source code BEFORE running 'go mod tidy'
# Since we don't have a populated go.mod or go.sum locally, we must provide
# the source code (main.go) so 'go mod tidy' can scan it and detect dependencies.
COPY . .

# Download dependencies and generate go.sum
RUN go mod tidy

# Build the binary
# -ldflags="-w -s" strips debug information, reducing binary size
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/airgo .

# Stage 2: Runner
FROM alpine:latest

# Security: Create a non-root user and group with strict UID 1000
# This ensures file ownership matches the host user (Talha) for bind mounts
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -h /home/appuser -D appuser

WORKDIR /app

# Install certificates for HTTPS
RUN apk --no-cache add ca-certificates curl

# Copy binary and assets
COPY --from=builder /app/airgo .
COPY --from=builder /app/templates ./templates

# Create uploads directory and assign permissions to the non-root user
RUN mkdir -p uploads && \
    chown -R appuser:appgroup /app

# Switch to non-root user (Security Best Practice)
USER appuser

EXPOSE 8080

CMD ["./airgo"]