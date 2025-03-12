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

func generateData(application *app.App, timestamp int64) {
	application.Mutex.RLock()
	defer application.Mutex.RUnlock()

	charts := make([]models.ChartData, application.Config.NumCharts)

	for chartIndex := 0; chartIndex < application.Config.NumCharts; chartIndex++ {
		points := make([]models.DataPoint, application.Config.NumPoints)
		frequency := 0.5 + float64(chartIndex)*0.5
		phase := float64(timestamp)/10.0 + float64(chartIndex)*math.Pi/4

		for i := 0; i < application.Config.NumPoints; i++ {
			x := float64(i) * 0.1
			timestamp = timestamp + int64(i)
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
			timestamp := time.Now().Unix()
			generateData(application, timestamp)
			log.Printf("Generated data at %s for %d charts", time.Unix(timestamp, 0).Format("15:04:05"), application.Config.NumCharts)
		}
	}
}
