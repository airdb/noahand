# Build stage
FROM golang:1.26-alpine AS builder

WORKDIR /opt/

# Install build dependencies
RUN apk add --no-cache git make build-base binutils binutils-gold

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN mkdir -p output/plugins && \
    CGO_ENABLED=0 go build -o output/main ./main.go && \
    CGO_ENABLED=1 go build -buildmode=plugin -o output/plugins/01_plugin.so ./plugins/01_plugin && \
    CGO_ENABLED=1 go build -buildmode=plugin -o output/plugins/02_plugin.so ./plugins/02_plugin

# Final stage
#FROM alpine:latest
#
#WORKDIR /opt
#
## Install required runtime dependencies
#RUN apk add --no-cache ca-certificates tzdata
#
## Copy binary from builder
#COPY --from=builder /opt/main /opt/main
#COPY --from=builder /opt/noahand /opt/noahand
#
## Create output directory for plugins
#RUN mkdir -p /opt/output
#
## Set executable permissions
#RUN chmod +x /opt/main /opt/noah
#
## Expose port (assuming default web server port)
#EXPOSE 8080
#
## Run the main binary
#CMD ["/app/main"]
