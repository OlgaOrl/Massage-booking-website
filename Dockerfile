# Multi-stage build for Go application
FROM golang:1.23-alpine AS builder

# Set working directory
WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev sqlite-dev

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY backend/ ./backend/

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./backend/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add ca-certificates sqlite

# Set working directory
WORKDIR /root/

# Copy the binary from builder stage
COPY --from=builder /app/main .

# Copy static files
COPY --from=builder /app/backend/static ./backend/static

# Create directory for database
RUN mkdir -p /data

# Expose port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=release

# Run the application
CMD ["./main"]
