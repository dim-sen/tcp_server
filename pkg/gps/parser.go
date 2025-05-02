package gps

import (
	"encoding/binary"
	"tcp_server/pkg/models"
	"time"
)

// ParseData extracts GPS location from binary data
// Data format: [latitude(8 bytes)][longitude(8 bytes)][timestamp(8 bytes)]
// All values are in big-endian format
func ParseData(data []byte) models.GPSLocation {
	lat := binary.BigEndian.Uint64(data[0:8])
	lon := binary.BigEndian.Uint64(data[8:16])
	timestamp := binary.BigEndian.Uint64(data[16:24])

	// Convert from microdegrees to degrees and handle potential overflow
	latitude := float64(int64(lat)) / 1e6
	longitude := float64(int64(lon)) / 1e6

	return models.GPSLocation{
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Unix(int64(timestamp), 0),
	}
}

// ValidateGPSData checks if the data is valid GPS data
func ValidateGPSData(data []byte) bool {
	// Minimum length check (24 bytes)
	if len(data) < 24 {
		return false
	}

	// Could add more validation logic here
	// - Check if coordinates are in valid range
	// - Check if timestamp is reasonable

	return true
}
