# Build stage
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o productApp ./cmd/api

# Final stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates && \
    adduser -D appuser

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/productApp .

# Use non-root user
USER appuser

# Expose port
EXPOSE 8083

# Start the application
CMD ["./productApp"]
