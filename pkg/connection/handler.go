package connection

import (
	"encoding/binary"
	"log"
	"net"
	"tcp_server/pkg/config"
	"tcp_server/pkg/models"
	"tcp_server/pkg/processor"
	"time"
)

// HandleClient manages a client connection
func HandleClient(conn net.Conn) {
	defer conn.Close()

	// Set initial read deadline
	conn.SetReadDeadline(time.Now().Add(config.ReadTimeout))

	// Create the worker pool
	worker, wg := processor.NewWorkerPool(config.MessageQueueSize)
	defer func() {
		worker.Close()
		wg.Wait()
	}()

	header := make([]byte, config.HeaderSize)
	for {
		// Reset read deadline
		conn.SetReadDeadline(time.Now().Add(config.ReadTimeout))

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
		if msgLength > uint32(config.BufferSize) {
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

		// Submit message to worker pool
		worker.Submit(models.Message{Data: msgData, Conn: conn})
	}
}
