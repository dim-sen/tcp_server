FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum* ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o gps_server cmd/tcp_app/main/main.go

# Use a minimal alpine image for the final stage
FROM alpine:3.19

WORKDIR /app

# Install ca-certificates for HTTPS
RUN apk --no-cache add ca-certificates tzdata

# Copy the binary from builder stage
COPY --from=builder /app/gps_server .

# Expose the TCP port
EXPOSE 8181

# Create a non-root user to run the application
RUN adduser -D -g '' appuser
USER appuser

# Run the application
CMD ["./gps_server"] 