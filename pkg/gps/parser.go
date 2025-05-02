package gps

import (
	"encoding/binary"
	"math"
	"tcp_server/pkg/config"
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

	// Round to configured precision
	latitude = roundToDecimalPlaces(latitude, config.LocationPrecision)
	longitude = roundToDecimalPlaces(longitude, config.LocationPrecision)

	return models.GPSLocation{
		Latitude:  latitude,
		Longitude: longitude,
		Timestamp: time.Unix(int64(timestamp), 0),
	}
}

// ValidateGPSData checks if the data is valid GPS data
func ValidateGPSData(data []byte) bool {
	// Skip validation if disabled
	if !config.GPSValidationEnabled {
		return true
	}

	// Minimum length check (24 bytes)
	if len(data) < 24 {
		return false
	}

	// Parse the data
	lat := binary.BigEndian.Uint64(data[0:8])
	lon := binary.BigEndian.Uint64(data[8:16])
	timestamp := binary.BigEndian.Uint64(data[16:24])

	// Convert from microdegrees to degrees
	latitude := float64(int64(lat)) / 1e6
	longitude := float64(int64(lon)) / 1e6

	// Validate coordinates
	if latitude < config.MinLatitude || latitude > config.MaxLatitude {
		return false
	}

	if longitude < config.MinLongitude || longitude > config.MaxLongitude {
		return false
	}

	// Validate timestamp
	if config.MaxTimestampDrift > 0 {
		timestampTime := time.Unix(int64(timestamp), 0)
		now := time.Now()
		diff := now.Sub(timestampTime)
		if diff < 0 {
			diff = -diff // Handle future timestamps
		}

		if diff > config.MaxTimestampDrift {
			return false
		}
	}

	return true
}

// roundToDecimalPlaces rounds a float to the specified number of decimal places
func roundToDecimalPlaces(value float64, decimalPlaces int) float64 {
	shift := math.Pow(10, float64(decimalPlaces))
	return math.Round(value*shift) / shift
}
