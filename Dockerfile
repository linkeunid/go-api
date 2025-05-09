# Build stage
FROM golang:1.24.2-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git ca-certificates tzdata

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o api ./cmd/api

# Create health check script
RUN echo '#!/bin/sh' > /app/healthcheck.sh && \
    echo 'wget -q --spider http://localhost:8080/health || exit 1' >> /app/healthcheck.sh && \
    chmod +x /app/healthcheck.sh

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates tzdata wget

# Set timezone
ENV TZ=UTC

# Create non-root user
RUN addgroup -g 1000 app && \
    adduser -u 1000 -G app -s /bin/sh -D app

WORKDIR /app

# Copy the binary and health check script from builder
COPY --from=builder /app/api /app/
COPY --from=builder /app/healthcheck.sh /

# Copy configuration files
COPY --from=builder /app/config /app/config

# Set user
USER app

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/api"]

# Health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD /healthcheck.sh 