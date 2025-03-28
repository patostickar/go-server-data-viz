package models

// ChartPoint represents the values for a single timestamp on a Plot
type ChartPoint struct {
	Timestamp string    `json:"timestamp"`
	Values    []float64 `json:"values"`
}

// ChartData represents the data for a single chart
type ChartData struct {
	ChartID string       `json:"chartId"`
	Data    []ChartPoint `json:"data"`
}

type ChartDataTimestamp struct {
	Timestamp int64       `json:"timestamp"`
	ChartData []ChartData `json:"chartData"`
}
