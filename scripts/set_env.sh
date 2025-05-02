#!/bin/bash
# TCP GPS Server - Environment Configuration Script for Unix/Linux/macOS
# Set default environment variables for development

# Server settings
export TCP_SERVER_HOST="localhost"
export TCP_SERVER_PORT="8181"
export TCP_MAX_CONNECTIONS="1000"
export TCP_ENABLE_TLS="false"
export TCP_CERT_FILE="cert.pem"
export TCP_KEY_FILE="key.pem"

# Performance settings
export TCP_BUFFER_SIZE="4096"
export TCP_WORKER_COUNT="10"
export TCP_MESSAGE_QUEUE_SIZE="1000"
export TCP_HEADER_SIZE="4"

# Timeout settings
export TCP_MAX_RETRIES="3"
export TCP_RETRY_DELAY="2s"
export TCP_READ_TIMEOUT="60s"
export TCP_WRITE_TIMEOUT="10s"
export TCP_SHUTDOWN_TIMEOUT="30s"
export TCP_IDLE_TIMEOUT="300s"

# GPS data validation
export TCP_GPS_VALIDATION_ENABLED="true"
export TCP_MIN_LATITUDE="-90.0"
export TCP_MAX_LATITUDE="90.0"
export TCP_MIN_LONGITUDE="-180.0"
export TCP_MAX_LONGITUDE="180.0"
export TCP_LOCATION_PRECISION="6"
export TCP_MAX_TIMESTAMP_DRIFT="24h"

# Security settings
export TCP_RATE_LIMIT_ENABLED="false"
export TCP_RATE_LIMIT_PER_MINUTE="300"
export TCP_IP_WHITELIST_ENABLED="false"
export TCP_ALLOWED_IPS="127.0.0.1"

# Logging settings
export TCP_LOG_LEVEL="info"
export TCP_LOG_TIMESTAMPS="true"
export TCP_LOG_CLIENT_IP="true"
export TCP_LOG_GPS_COORDINATES="true"
export TCP_LOG_CONNECTION_EVENTS="true"
export TCP_LOG_MESSAGE_STATISTICS="true"
export TCP_STATISTICS_INTERVAL="60s"

echo "Environment variables set for TCP GPS Server"
echo "Run the server with: go run cmd/tcp_app/main/main.go" 