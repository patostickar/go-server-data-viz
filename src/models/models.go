package models

import (
	gqlmodel "github.com/patostickar/go-server-data-viz/src/graph/model"
)

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

type Charts struct {
	Data []ChartData
}

func (d ChartData) toGql() []*gqlmodel.ChartPoint {
	var gqlPoints []*gqlmodel.ChartPoint

	for _, point := range d.Data {
		gqlPoint := &gqlmodel.ChartPoint{
			Timestamp: point.Timestamp,
			Values:    point.Values,
		}
		gqlPoints = append(gqlPoints, gqlPoint)
	}
	return gqlPoints
}

func (c Charts) ToGql() []*gqlmodel.ChartData {
	var gqlCharts []*gqlmodel.ChartData

	for _, data := range c.Data {
		gqlChart := &gqlmodel.ChartData{
			ChartID: data.ChartID,
			Data:    data.toGql(),
		}
		gqlCharts = append(gqlCharts, gqlChart)
	}

	return gqlCharts
}
