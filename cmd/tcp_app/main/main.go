package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"tcp_server/pkg"
	"tcp_server/pkg/config"
	"time"
)

func main() {
	// Configure logging
	log.SetFlags(0)
	if config.LogTimestamps {
		log.SetFlags(log.LstdFlags | log.Lmicroseconds)
	}

	log.Println("Starting TCP GPS server...")
	serverAddress := fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort)

	var listener net.Listener
	var err error

	if config.EnableTLS {
		// Configure TLS
		cert, err := tls.LoadX509KeyPair(config.CertFile, config.KeyFile)
		if err != nil {
			log.Fatalf("Error loading certificate: %v", err)
		}

		tlsConfig := &tls.Config{
			Certificates: []tls.Certificate{cert},
			MinVersion:   tls.VersionTLS12,
		}

		listener, err = tls.Listen("tcp", serverAddress, tlsConfig)
		if err != nil {
			log.Fatalf("Error creating TLS listener: %v", err)
		}
		log.Printf("Server started with TLS on %s", serverAddress)
	} else {
		// Non-TLS listener
		listener, err = net.Listen("tcp", serverAddress)
		if err != nil {
			log.Fatalf("Error creating listener: %v", err)
		}
		log.Printf("Server started on %s", serverAddress)
	}

	defer listener.Close()

	// Setup connection counter and active connections tracking
	var connectionCounter int
	var activeConnections sync.WaitGroup
	var activeConnectionsMutex sync.Mutex
	activeConnectionsMap := make(map[string]net.Conn)

	// Setup signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Channel to receive OS signals
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start goroutine for handling signals
	go func() {
		sig := <-signalChan
		log.Printf("Received signal: %v. Initiating graceful shutdown...", sig)
		cancel() // Trigger shutdown
	}()

	// Start statistics reporting if enabled
	if config.LogMessageStatistics {
		go reportStatistics(ctx, &activeConnectionsMutex, activeConnectionsMap)
	}

	// Accept connections in a separate goroutine
	acceptError := make(chan error, 1)
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				acceptError <- err
				return
			}

			// Check if we're shutting down
			select {
			case <-ctx.Done():
				conn.Close()
				continue
			default:
				// Continue processing
			}

			// Track the connection
			connectionCounter++
			clientIP := conn.RemoteAddr().String()

			// Apply connection limits
			activeConnectionsMutex.Lock()
			currentCount := len(activeConnectionsMap)
			if currentCount >= config.MaxConnections {
				activeConnectionsMutex.Unlock()
				log.Printf("Connection limit reached (%d), rejecting connection from %s",
					config.MaxConnections, clientIP)
				conn.Close()
				continue
			}

			// IP whitelist check
			if config.IPWhitelistEnabled {
				allowed := false
				clientHost, _, _ := net.SplitHostPort(clientIP)
				for _, ip := range config.AllowedIPs {
					if ip == clientHost {
						allowed = true
						break
					}
				}

				if !allowed {
					activeConnectionsMutex.Unlock()
					log.Printf("Connection from %s rejected (not in whitelist)", clientIP)
					conn.Close()
					continue
				}
			}

			// Update active connections map
			activeConnectionsMap[clientIP] = conn
			activeConnectionsMutex.Unlock()

			// Increment wait group counter
			activeConnections.Add(1)

			if config.LogConnectionEvents {
				log.Printf("Client #%d connected from %s", connectionCounter, clientIP)
			}

			// Handle client in a goroutine
			go func(clientNum int, conn net.Conn) {
				defer func() {
					// Clean up connection tracking
					activeConnectionsMutex.Lock()
					delete(activeConnectionsMap, clientIP)
					activeConnectionsMutex.Unlock()

					conn.Close()
					activeConnections.Done()

					if config.LogConnectionEvents {
						log.Printf("Client #%d from %s disconnected", clientNum, clientIP)
					}
				}()

				pkg.HandleClient(conn)
			}(connectionCounter, conn)
		}
	}()

	// Wait for a shutdown signal or accept error
	select {
	case <-ctx.Done():
		// Graceful shutdown
		log.Println("Shutting down server...")
		listener.Close()

		// Create a context with timeout for shutdown
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), config.ShutdownTimeout)
		defer shutdownCancel()

		// Wait for active connections to finish or timeout
		done := make(chan struct{})
		go func() {
			activeConnections.Wait()
			close(done)
		}()

		select {
		case <-done:
			log.Println("All connections closed gracefully")
		case <-shutdownCtx.Done():
			log.Println("Shutdown timeout reached, forcing exit")
		}
	case err := <-acceptError:
		log.Fatalf("Error accepting connection: %v", err)
	}
}

// reportStatistics periodically logs connection and usage statistics
func reportStatistics(ctx context.Context, mutex *sync.Mutex, connections map[string]net.Conn) {
	ticker := time.NewTicker(config.StatisticsInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			mutex.Lock()
			count := len(connections)
			mutex.Unlock()

			log.Printf("STATISTICS: Active connections: %d", count)

			// Could add more statistics here
			// - Memory usage
			// - CPU usage
			// - Message throughput

		case <-ctx.Done():
			return
		}
	}
}
