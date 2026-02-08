# Stage 1: Builder (Debian-based for better networking)
FROM golang:1.23 AS builder

# Accept build-time proxy settings
ARG GOPROXY=https://proxy.golang.org,direct
ENV GOPROXY=$GOPROXY
ENV GOSUMDB=off

WORKDIR /app

# Install git (required for some go modules)
RUN apt-get update && apt-get install -y --no-install-recommends git && \
    rm -rf /var/lib/apt/lists/*

# Copy only dependency files first
COPY go.mod go.sum ./
# Attempt to download dependencies
RUN go mod download

# Copy source
COPY . .

# Build optimized binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/airgo .

# Stage 2: Runner (Alpine for production)
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

# Healthcheck will use wget (pre-installed in alpine)
CMD ["./airgo"]