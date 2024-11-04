# Build stage
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git make

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN go build -o main main.go &&\
    go build -buildmode=plugin -o output/01_plugin.so plugins/01_plugin.go &&\
    go build -buildmode=plugin -o output/02_plugin.so plugins/02_plugin.go

# Final stage
FROM alpine:latest

WORKDIR /app

# Install required runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/main /app/main
COPY --from=builder /app/guardhouse /app/guardhouse

# Create output directory for plugins
RUN mkdir -p /app/output

# Set executable permissions
RUN chmod +x /app/main /app/noah

# Expose port (assuming default web server port)
EXPOSE 8080

# Run the main binary
CMD ["/app/main"]