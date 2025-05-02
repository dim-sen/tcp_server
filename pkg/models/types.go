package models

import (
	"net"
	"time"
)

// GPSLocation represents a geographical position with timestamp
type GPSLocation struct {
	Latitude  float64
	Longitude float64
	Timestamp time.Time
}

// Message represents a data payload from a client connection
type Message struct {
	Data []byte
	Conn net.Conn
}
