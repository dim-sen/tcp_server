package pkg

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"
)

const (
	bufferSize       = 4096
	workerCount      = 10
	messageQueueSize = 1000
	headerSize       = 4 // 4 bytes for message length (uint32)
	maxRetries       = 3
	retryDelay       = 2 * time.Second
)

type GPSLocation struct {
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

type Message struct {
	Data []byte
	Conn net.Conn
}

func HandleClient(conn net.Conn) {
	defer conn.Close()

	// Set read and write deadlines
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	conn.SetWriteDeadline(time.Now().Add(30 * time.Second))

	// Create a buffered channel for messages
	messageChan := make(chan Message, messageQueueSize)

	// Create a worker pool
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for msg := range messageChan {
				// Process the message
				location := parseGPSData(msg.Data)
				log.Printf("Received GPS location from %s: Lat=%.6f, Lon=%.6f, Time=%s\n",
					msg.Conn.RemoteAddr(),
					location.Latitude,
					location.Longitude,
					location.Timestamp.Format(time.RFC3339))

				// Send response back to client with retry mechanism
				response := []byte("Location received!")
				header := make([]byte, 4)
				binary.BigEndian.PutUint32(header, uint32(len(response)))

				// Retry sending response if failed
				for retry := 0; retry < maxRetries; retry++ {
					// Send header
					_, err := msg.Conn.Write(header)
					if err != nil {
						log.Printf("Error sending response header to %s (attempt %d/%d): %v\n",
							msg.Conn.RemoteAddr(), retry+1, maxRetries, err)
						if retry < maxRetries-1 {
							time.Sleep(retryDelay)
							continue
						}
						break
					}

					// Send response
					_, err = msg.Conn.Write(response)
					if err != nil {
						log.Printf("Error sending response to %s (attempt %d/%d): %v\n",
							msg.Conn.RemoteAddr(), retry+1, maxRetries, err)
						if retry < maxRetries-1 {
							time.Sleep(retryDelay)
							continue
						}
						break
					}

					// Success, reset read deadline
					msg.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
					break
				}
			}
		}()
	}

	header := make([]byte, headerSize)
	for {
		// Reset read deadline
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		// First read the header to get message length
		_, err := conn.Read(header)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Read timeout from %s, closing connection\n", conn.RemoteAddr())
			} else {
				log.Printf("Error reading header from %s: %v\n", conn.RemoteAddr(), err)
			}
			break
		}

		// Convert header bytes to message length
		msgLength := binary.BigEndian.Uint32(header)

		// Validate message length
		if msgLength > bufferSize {
			log.Printf("Invalid message length %d from %s, closing connection\n", msgLength, conn.RemoteAddr())
			break
		}

		// Create buffer for the message
		msgData := make([]byte, msgLength)

		// Read the full message
		_, err = conn.Read(msgData)
		if err != nil {
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				log.Printf("Read timeout from %s, closing connection\n", conn.RemoteAddr())
			} else {
				log.Printf("Error reading message from %s: %v\n", conn.RemoteAddr(), err)
			}
			break
		}

		// Try to send to message channel
		select {
		case messageChan <- Message{Data: msgData, Conn: conn}:
			// Message sent successfully
		default:
			// Channel is full, handle backpressure
			log.Printf("Message queue full for %s, dropping message\n", conn.RemoteAddr())
		}
	}

	// Close the message channel and wait for workers to finish
	close(messageChan)
	wg.Wait()
}

func parseGPSData(data []byte) GPSLocation {
	// Assuming data format: [latitude(8 bytes)][longitude(8 bytes)][timestamp(8 bytes)]
	// All values are in big-endian format
	lat := binary.BigEndian.Uint64(data[0:8])
	lon := binary.BigEndian.Uint64(data[8:16])
	timestamp := binary.BigEndian.Uint64(data[16:24])

	return GPSLocation{
		Latitude:  float64(lat) / 1e6, // Convert from microdegrees to degrees
		Longitude: float64(lon) / 1e6, // Convert from microdegrees to degrees
		Timestamp: time.Unix(int64(timestamp), 0),
	}
}
