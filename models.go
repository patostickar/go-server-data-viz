package main

import (
	"sync"
)

// DataPoint represents a single data point on a chart
type DataPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

// ChartData represents the data for a single chart
type ChartData struct {
	ChartID string      `json:"chartId"`
	Data    []DataPoint `json:"data"`
}

// ConfigRequest represents the configuration requested by the client
type ConfigRequest struct {
	NumCharts    int `json:"numCharts"`
	NumPoints    int `json:"numPoints"`
	PollInterval int `json:"pollInterval"`
}

// AppConfig holds the application configuration and state
type AppConfig struct {
	NumCharts    int
	NumPoints    int
	PollInterval int
	LastData     []ChartData
	mutex        sync.RWMutex
}

// Global configuration with defaults
var config = AppConfig{
	NumCharts:    1,
	NumPoints:    100,
	PollInterval: 1000,
}
