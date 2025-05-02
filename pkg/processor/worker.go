package processor

import (
	"encoding/binary"
	"log"
	"sync"
	"tcp_server/pkg/config"
	"tcp_server/pkg/gps"
	"tcp_server/pkg/models"
	"time"
)

// Worker represents a message processing worker
type Worker struct {
	messageChan chan models.Message
	wg          *sync.WaitGroup
}

// NewWorkerPool creates a new pool of workers
func NewWorkerPool(messageQueueSize int) (*Worker, *sync.WaitGroup) {
	var wg sync.WaitGroup
	messageChan := make(chan models.Message, messageQueueSize)

	worker := &Worker{
		messageChan: messageChan,
		wg:          &wg,
	}

	// Start the worker goroutines
	for i := 0; i < config.WorkerCount; i++ {
		wg.Add(1)
		go worker.processMessages()
	}

	return worker, &wg
}

// Submit puts a message in the processing queue
func (w *Worker) Submit(msg models.Message) bool {
	select {
	case w.messageChan <- msg:
		return true
	default:
		// Channel is full, handle backpressure
		log.Printf("Message queue full for %s, dropping message\n", msg.Conn.RemoteAddr())
		return false
	}
}

// Close stops the worker pool
func (w *Worker) Close() {
	close(w.messageChan)
}

// processMessages handles incoming messages
func (w *Worker) processMessages() {
	defer w.wg.Done()

	for msg := range w.messageChan {
		// Validate data format
		if !gps.ValidateGPSData(msg.Data) {
			log.Printf("Invalid GPS data format from %s, ignoring\n", msg.Conn.RemoteAddr())
			continue
		}

		// Process the message
		location := gps.ParseData(msg.Data)
		log.Printf("Received GPS location from %s: Lat=%.6f, Lon=%.6f, Time=%s\n",
			msg.Conn.RemoteAddr(),
			location.Latitude,
			location.Longitude,
			location.Timestamp.Format(time.RFC3339))

		sendResponse(msg)
	}
}

// sendResponse sends an acknowledgment back to the client
func sendResponse(msg models.Message) {
	// Send response back to client with retry mechanism
	response := []byte("Location received!")
	header := make([]byte, 4)
	binary.BigEndian.PutUint32(header, uint32(len(response)))

	// Set write deadline for response
	msg.Conn.SetWriteDeadline(time.Now().Add(config.WriteTimeout))

	// Retry sending response if failed
	for retry := 0; retry < config.MaxRetries; retry++ {
		// Send header
		_, err := msg.Conn.Write(header)
		if err != nil {
			log.Printf("Error sending response header to %s (attempt %d/%d): %v\n",
				msg.Conn.RemoteAddr(), retry+1, config.MaxRetries, err)
			if retry < config.MaxRetries-1 {
				time.Sleep(config.RetryDelay)
				continue
			}
			break
		}

		// Send response
		_, err = msg.Conn.Write(response)
		if err != nil {
			log.Printf("Error sending response to %s (attempt %d/%d): %v\n",
				msg.Conn.RemoteAddr(), retry+1, config.MaxRetries, err)
			if retry < config.MaxRetries-1 {
				time.Sleep(config.RetryDelay)
				continue
			}
			break
		}

		// Success, reset read deadline
		msg.Conn.SetReadDeadline(time.Now().Add(config.ReadTimeout))
		break
	}
}
