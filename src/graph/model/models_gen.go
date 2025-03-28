// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type ChartData struct {
	ChartID string        `json:"chartId"`
	Data    []*ChartPoint `json:"data"`
}

type ChartDataTimestamp struct {
	Timestamp int32        `json:"timestamp"`
	ChartData []*ChartData `json:"chartData"`
}

type ChartPoint struct {
	Timestamp string          `json:"timestamp"`
	Values    []*KeyValuePair `json:"values"`
}

type KeyValuePair struct {
	Key   string  `json:"key"`
	Value float64 `json:"value"`
}

type Query struct {
}
