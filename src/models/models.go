package models

import (
	pb "github.com/patostickar/go-server-data-viz/models"
	gqlmodel "github.com/patostickar/go-server-data-viz/src/graph/model"
	"time"
)

// ChartPoint represents the values for a single timestamp on a Plot
type ChartPoint struct {
	Timestamp int64     `json:"timestamp"`
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
			Timestamp: time.Unix(point.Timestamp, 0).Format("15:04:05"),
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

func (c Charts) ToProto() []*pb.ChartData {
	grpcCharts := make([]*pb.ChartData, len(c.Data))

	for i, chart := range c.Data {
		grpcPoints := make([]*pb.ChartPoint, len(chart.Data))

		for j, point := range chart.Data {
			grpcPoints[j] = &pb.ChartPoint{
				TimestampUnixSeconds: point.Timestamp,
				Values:               point.Values,
			}
		}

		grpcCharts[i] = &pb.ChartData{
			ChartId: chart.ChartID,
			Points:  grpcPoints,
		}
	}

	return grpcCharts
}
