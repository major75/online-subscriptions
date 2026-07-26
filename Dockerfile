# Build stage
FROM golang:1.26-alpine AS builder

# Install build dependencies
RUN apk add --no-cache git make

# Set working directory
WORKDIR /app

# Copy source code
COPY . .

# Build the application
RUN GOOS=linux make swagger && make install && make build

# Final stage
FROM alpine:latest

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates make

WORKDIR /root/

# Copy binary from builder
COPY --from=builder /app/bin/subscriptions .

# Run the application
CMD ["./subscriptions"]
