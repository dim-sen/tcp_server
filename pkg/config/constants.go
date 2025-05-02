package config

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Configuration constants with their default values
var (
	// Server settings
	ServerHost     = "localhost"
	ServerPort     = 8181
	MaxConnections = 1000
	EnableTLS      = false
	CertFile       = "cert.pem"
	KeyFile        = "key.pem"

	// Buffer and message handling
	BufferSize       = 4096
	WorkerCount      = 10
	MessageQueueSize = 1000
	HeaderSize       = 4 // 4 bytes for message length (uint32)

	// Retry and timeout settings
	MaxRetries      = 3
	RetryDelay      = 2 * time.Second
	ReadTimeout     = 60 * time.Second
	WriteTimeout    = 10 * time.Second
	ShutdownTimeout = 30 * time.Second
	IdleTimeout     = 300 * time.Second

	// GPS data settings
	GPSValidationEnabled = true
	MinLatitude          = -90.0
	MaxLatitude          = 90.0
	MinLongitude         = -180.0
	MaxLongitude         = 180.0
	LocationPrecision    = 6 // Decimal places for coordinates
	MaxTimestampDrift    = 24 * time.Hour

	// Security settings
	RateLimitEnabled   = false
	RateLimitPerMinute = 300
	IPWhitelistEnabled = false
	AllowedIPs         = []string{"127.0.0.1"}

	// Logging settings
	LogLevel             = "info"
	LogTimestamps        = true
	LogClientIP          = true
	LogGPSCoordinates    = true
	LogConnectionEvents  = true
	LogMessageStatistics = true
	StatisticsInterval   = 60 * time.Second
)

// init loads configuration from environment variables
func init() {
	// Server settings
	ServerHost = getEnvString("TCP_SERVER_HOST", ServerHost)
	ServerPort = getEnvInt("TCP_SERVER_PORT", ServerPort)
	MaxConnections = getEnvInt("TCP_MAX_CONNECTIONS", MaxConnections)
	EnableTLS = getEnvBool("TCP_ENABLE_TLS", EnableTLS)
	CertFile = getEnvString("TCP_CERT_FILE", CertFile)
	KeyFile = getEnvString("TCP_KEY_FILE", KeyFile)

	// Buffer and message handling
	BufferSize = getEnvInt("TCP_BUFFER_SIZE", BufferSize)
	WorkerCount = getEnvInt("TCP_WORKER_COUNT", WorkerCount)
	MessageQueueSize = getEnvInt("TCP_MESSAGE_QUEUE_SIZE", MessageQueueSize)
	HeaderSize = getEnvInt("TCP_HEADER_SIZE", HeaderSize)

	// Retry and timeout settings
	MaxRetries = getEnvInt("TCP_MAX_RETRIES", MaxRetries)
	RetryDelay = getEnvDuration("TCP_RETRY_DELAY", RetryDelay)
	ReadTimeout = getEnvDuration("TCP_READ_TIMEOUT", ReadTimeout)
	WriteTimeout = getEnvDuration("TCP_WRITE_TIMEOUT", WriteTimeout)
	ShutdownTimeout = getEnvDuration("TCP_SHUTDOWN_TIMEOUT", ShutdownTimeout)
	IdleTimeout = getEnvDuration("TCP_IDLE_TIMEOUT", IdleTimeout)

	// GPS data settings
	GPSValidationEnabled = getEnvBool("TCP_GPS_VALIDATION_ENABLED", GPSValidationEnabled)
	MinLatitude = getEnvFloat("TCP_MIN_LATITUDE", MinLatitude)
	MaxLatitude = getEnvFloat("TCP_MAX_LATITUDE", MaxLatitude)
	MinLongitude = getEnvFloat("TCP_MIN_LONGITUDE", MinLongitude)
	MaxLongitude = getEnvFloat("TCP_MAX_LONGITUDE", MaxLongitude)
	LocationPrecision = getEnvInt("TCP_LOCATION_PRECISION", LocationPrecision)
	MaxTimestampDrift = getEnvDuration("TCP_MAX_TIMESTAMP_DRIFT", MaxTimestampDrift)

	// Security settings
	RateLimitEnabled = getEnvBool("TCP_RATE_LIMIT_ENABLED", RateLimitEnabled)
	RateLimitPerMinute = getEnvInt("TCP_RATE_LIMIT_PER_MINUTE", RateLimitPerMinute)
	IPWhitelistEnabled = getEnvBool("TCP_IP_WHITELIST_ENABLED", IPWhitelistEnabled)
	if envIPs := getEnvString("TCP_ALLOWED_IPS", ""); envIPs != "" {
		AllowedIPs = strings.Split(envIPs, ",")
	}

	// Logging settings
	LogLevel = getEnvString("TCP_LOG_LEVEL", LogLevel)
	LogTimestamps = getEnvBool("TCP_LOG_TIMESTAMPS", LogTimestamps)
	LogClientIP = getEnvBool("TCP_LOG_CLIENT_IP", LogClientIP)
	LogGPSCoordinates = getEnvBool("TCP_LOG_GPS_COORDINATES", LogGPSCoordinates)
	LogConnectionEvents = getEnvBool("TCP_LOG_CONNECTION_EVENTS", LogConnectionEvents)
	LogMessageStatistics = getEnvBool("TCP_LOG_MESSAGE_STATISTICS", LogMessageStatistics)
	StatisticsInterval = getEnvDuration("TCP_STATISTICS_INTERVAL", StatisticsInterval)

	log.Println("Configuration loaded:")
	// Server settings
	log.Printf("  ServerHost: %s", ServerHost)
	log.Printf("  ServerPort: %d", ServerPort)
	log.Printf("  MaxConnections: %d", MaxConnections)
	log.Printf("  EnableTLS: %t", EnableTLS)

	// Buffer and message handling
	log.Printf("  BufferSize: %d bytes", BufferSize)
	log.Printf("  WorkerCount: %d workers", WorkerCount)
	log.Printf("  MessageQueueSize: %d messages", MessageQueueSize)
	log.Printf("  HeaderSize: %d bytes", HeaderSize)

	// Retry and timeout settings
	log.Printf("  MaxRetries: %d attempts", MaxRetries)
	log.Printf("  RetryDelay: %s", RetryDelay)
	log.Printf("  ReadTimeout: %s", ReadTimeout)
	log.Printf("  WriteTimeout: %s", WriteTimeout)
	log.Printf("  ShutdownTimeout: %s", ShutdownTimeout)
	log.Printf("  IdleTimeout: %s", IdleTimeout)

	// GPS data settings
	log.Printf("  GPSValidationEnabled: %t", GPSValidationEnabled)
	log.Printf("  Coordinate Range: [%f,%f] to [%f,%f]", MinLatitude, MinLongitude, MaxLatitude, MaxLongitude)
	log.Printf("  LocationPrecision: %d decimal places", LocationPrecision)

	// Security settings
	log.Printf("  RateLimitEnabled: %t", RateLimitEnabled)
	if RateLimitEnabled {
		log.Printf("  RateLimitPerMinute: %d", RateLimitPerMinute)
	}
	log.Printf("  IPWhitelistEnabled: %t", IPWhitelistEnabled)
	if IPWhitelistEnabled {
		log.Printf("  AllowedIPs: %s", strings.Join(AllowedIPs, ", "))
	}

	// Logging settings
	log.Printf("  LogLevel: %s", LogLevel)
	log.Printf("  LogMessageStatistics: %t (interval: %s)", LogMessageStatistics, StatisticsInterval)
}

// Helper functions for environment variables

// getEnvString gets a string from environment variable or returns the default value
func getEnvString(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

// getEnvInt gets an integer from environment variable or returns the default value
func getEnvInt(key string, defaultVal int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
		log.Printf("Warning: Invalid value for %s, using default: %d", key, defaultVal)
	}
	return defaultVal
}

// getEnvFloat gets a float from environment variable or returns the default value
func getEnvFloat(key string, defaultVal float64) float64 {
	if value, exists := os.LookupEnv(key); exists {
		if floatVal, err := strconv.ParseFloat(value, 64); err == nil {
			return floatVal
		}
		log.Printf("Warning: Invalid value for %s, using default: %f", key, defaultVal)
	}
	return defaultVal
}

// getEnvBool gets a boolean from environment variable or returns the default value
func getEnvBool(key string, defaultVal bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
		log.Printf("Warning: Invalid value for %s, using default: %t", key, defaultVal)
	}
	return defaultVal
}

// getEnvDuration gets a duration from environment variable or returns the default value
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	if value, exists := os.LookupEnv(key); exists {
		if duration, err := time.ParseDuration(value); err == nil {
			return duration
		}
		log.Printf("Warning: Invalid value for %s, using default: %s", key, defaultVal)
	}
	return defaultVal
}
