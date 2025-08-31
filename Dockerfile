# Multi-stage build for GoKit CLI
FROM golang:1.23-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the binary
RUN cd cli && make build

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata

# Create non-root user
RUN addgroup -g 1001 -S gokit && \
    adduser -u 1001 -S gokit -G gokit

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/cli/gokit /usr/local/bin/gokit

# Change ownership to non-root user
RUN chown gokit:gokit /usr/local/bin/gokit

# Switch to non-root user
USER gokit

# Set entrypoint
ENTRYPOINT ["gokit"]

# Default command
CMD ["--help"]
