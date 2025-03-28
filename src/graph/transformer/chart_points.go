package transformer

import (
	gqlmodel "github.com/patostickar/go-server-data-viz/src/graph/model"
	"github.com/patostickar/go-server-data-viz/src/models"
)

func ChartPoints2Gql(storePoints []models.ChartPoint) []*gqlmodel.ChartPoint {
	var gqlPoints []*gqlmodel.ChartPoint
	for _, point := range storePoints {
		gqlPoint := &gqlmodel.ChartPoint{
			Timestamp: point.Timestamp,
			Values:    point.Values,
		}
		gqlPoints = append(gqlPoints, gqlPoint)
	}
	return gqlPoints
}
