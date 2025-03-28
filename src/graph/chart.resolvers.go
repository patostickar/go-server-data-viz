package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.
// Code generated by github.com/99designs/gqlgen version v0.17.70

import (
	"context"
	"time"

	"github.com/patostickar/go-server-data-viz/src/config"
	gqlmodel "github.com/patostickar/go-server-data-viz/src/graph/model"
	"github.com/patostickar/go-server-data-viz/src/models"
)

// GetCharts is the resolver for the getCharts field.
func (r *queryResolver) GetCharts(ctx context.Context) (*gqlmodel.ChartDataTimestamp, error) {
	charts, err := r.s.Store.Read(config.ChartsKey)
	if err != nil {
		return nil, err
	}

	var chartData []*gqlmodel.ChartData
	for _, chart := range charts.([]models.ChartData) {
		gqlChart := &gqlmodel.ChartData{
			ChartID: chart.ChartID,
			Data:    convertChartPoints(chart.Data),
		}
		chartData = append(chartData, gqlChart)
	}
	r.logger.Debugf("returning %d charts", len(chartData))
	res := gqlmodel.ChartDataTimestamp{
		Timestamp: int32(time.Now().UnixMilli()),
		ChartData: chartData,
	}
	return &res, nil
}

func convertChartPoints(storePoints []models.ChartPoint) []*gqlmodel.ChartPoint {
	var gqlPoints []*gqlmodel.ChartPoint
	for _, point := range storePoints {
		gqlPoint := &gqlmodel.ChartPoint{
			Timestamp: point.Timestamp,
			Values:    convertValues(point.Values),
		}
		gqlPoints = append(gqlPoints, gqlPoint)
	}
	return gqlPoints
}
func convertValues(values map[string]float64) []*gqlmodel.KeyValuePair {
	var keyValuePairs []*gqlmodel.KeyValuePair
	for key, value := range values {
		keyValuePairs = append(keyValuePairs, &gqlmodel.KeyValuePair{
			Key:   key,
			Value: value,
		})
	}
	return keyValuePairs
}

// Query returns QueryResolver implementation.
func (r *Resolver) Query() QueryResolver { return &queryResolver{r} }

type queryResolver struct{ *Resolver }
