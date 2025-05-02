package pkg

import (
	"encoding/binary"
	"log"
	"net"
	"sync"
	"time"
)

const (
	bufferSize = 4096
	workerCount = 10
	messageQueueSize = 1000
	headerSize = 4 // 4 bytes for message length (uint32)
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
				
				// Send response back to client
				response := []byte("Location received!")
				header := make([]byte, 4)
				binary.BigEndian.PutUint32(header, uint32(len(response)))
				
				// Send header
				_, err := msg.Conn.Write(header)
				if err != nil {
					log.Printf("Error sending response header to %s: %v\n", msg.Conn.RemoteAddr(), err)
					continue
				}
				
				// Send response
				_, err = msg.Conn.Write(response)
				if err != nil {
					log.Printf("Error sending response to %s: %v\n", msg.Conn.RemoteAddr(), err)
					continue
				}
			}
		}()
	}

	header := make([]byte, headerSize)
	for {
		// First read the header to get message length
		_, err := conn.Read(header)
		if err != nil {
			log.Printf("Error reading header from %s: %v\n", conn.RemoteAddr(), err)
			break
		}

		// Convert header bytes to message length
		msgLength := binary.BigEndian.Uint32(header)
		
		// Create buffer for the message
		msgData := make([]byte, msgLength)
		
		// Read the full message
		_, err = conn.Read(msgData)
		if err != nil {
			log.Printf("Error reading message from %s: %v\n", conn.RemoteAddr(), err)
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
