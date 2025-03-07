package utils

import (
	"context"
	"fmt"
	"github.com/patostickar/go-server-data-viz/app"
	"github.com/patostickar/go-server-data-viz/models"
	"log"
	"math"
	"sync"
	"time"
)

func generateData(application *app.App) ([]models.ChartData, int64) {
	application.Mutex.RLock()
	defer application.Mutex.RUnlock()

	now := time.Now().Unix()
	charts := make([]models.ChartData, application.Config.NumCharts)

	for chartIndex := 0; chartIndex < application.Config.NumCharts; chartIndex++ {
		points := make([]models.DataPoint, application.Config.NumPoints)
		frequency := 0.5 + float64(chartIndex)*0.5
		phase := float64(now)/10.0 + float64(chartIndex)*math.Pi/4

		for i := 0; i < application.Config.NumPoints; i++ {
			x := float64(i) * 0.1
			timestamp := now + int64(i)
			points[i] = models.DataPoint{
				Timestamp: time.Unix(timestamp, 0).Format("15:04:05"),
				Value:     math.Sin(2*math.Pi*frequency*x + phase),
			}
		}

		charts[chartIndex] = models.ChartData{
			ChartID: fmt.Sprintf("chart%d", chartIndex+1),
			Data:    points,
		}
	}

	application.LastData = charts
	return charts, now
}

// StartDataGenerator initializes and starts the data generation routine
func StartDataGenerator(ctx context.Context, wg *sync.WaitGroup, application *app.App) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("Data generator stopping due to shutdown signal")
			return
		default:
			time.Sleep(time.Duration(application.Config.PollInterval) * time.Millisecond)
			_, now := generateData(application)
			log.Printf("Generated data at %s for %d charts", time.Unix(now, 0).Format("15:04:05"), application.Config.NumCharts)
		}
	}
}
