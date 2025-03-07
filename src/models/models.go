package models

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
