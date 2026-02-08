# Stage 1: Builder
FROM golang:1.23-alpine AS builder

# Set an alternative Go proxy to avoid timeout issues in some regions
ENV GOPROXY=https://goproxy.io,direct

WORKDIR /app

# No apk add git needed for standard public modules
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build with optimizations
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/airgo .

# Stage 2: Runner
FROM alpine:latest

# Create non-root user
RUN addgroup -g 1000 appgroup && \
    adduser -u 1000 -G appgroup -h /home/appuser -D appuser

WORKDIR /app

# Copy binary and templates
COPY --from=builder /app/airgo .
COPY --from=builder /app/templates ./templates

# Setup uploads directory
RUN mkdir -p uploads && chown -R appuser:appgroup /app

USER appuser
EXPOSE 8081

# Use built-in wget for healthcheck instead of curl to avoid apk add issues
CMD ["./airgo"]
