package config

import (
	"log"
	"os"
	"strconv"
	"time"
)

// Configuration constants with their default values
var (
	// BufferSize is the maximum size of the read buffer
	BufferSize = 4096

	// WorkerCount is the number of worker goroutines in the pool
	WorkerCount = 10

	// MessageQueueSize is the size of the buffered message channel
	MessageQueueSize = 1000

	// HeaderSize is the size of the message header in bytes
	HeaderSize = 4 // 4 bytes for message length (uint32)

	// MaxRetries is the maximum number of retry attempts
	MaxRetries = 3

	// RetryDelay is the delay between retry attempts
	RetryDelay = 2 * time.Second

	// ReadTimeout is the timeout for read operations
	ReadTimeout = 60 * time.Second

	// WriteTimeout is the timeout for write operations
	WriteTimeout = 10 * time.Second
)

// init loads configuration from environment variables
func init() {
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

	log.Println("Configuration loaded:")
	log.Printf("  BufferSize: %d bytes", BufferSize)
	log.Printf("  WorkerCount: %d workers", WorkerCount)
	log.Printf("  MessageQueueSize: %d messages", MessageQueueSize)
	log.Printf("  HeaderSize: %d bytes", HeaderSize)
	log.Printf("  MaxRetries: %d attempts", MaxRetries)
	log.Printf("  RetryDelay: %s", RetryDelay)
	log.Printf("  ReadTimeout: %s", ReadTimeout)
	log.Printf("  WriteTimeout: %s", WriteTimeout)
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
