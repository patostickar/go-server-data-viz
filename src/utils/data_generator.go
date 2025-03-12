package utils

import (
	"context"
	"fmt"
	"github.com/patostickar/go-server-data-viz/app"
	m "github.com/patostickar/go-server-data-viz/models"
	"log"
	"math"
	"sync"
	"time"
)

func generateData(a *app.App, timestamp int64) {
	a.Mutex.RLock()
	defer a.Mutex.RUnlock()

	// Always generate 3 charts
	numCharts := 3
	charts := make([]m.ChartData, numCharts)

	for chartIndex := 0; chartIndex < numCharts; chartIndex++ {
		// Create a slice to hold all points for this chart
		points := make([]m.ChartPoint, a.Config.NumPoints)

		for i := 0; i < a.Config.NumPoints; i++ {
			x := float64(i) * 0.1
			currentTimestamp := timestamp + int64(i)

			// Create a map to store values for each plot at this timestamp
			values := make(map[string]float64)

			for plotIndex := 0; plotIndex < a.Config.NumPlots; plotIndex++ {
				// Calculate unique frequency and phase for each plot
				frequency := 0.5 + float64(chartIndex)*0.5 + float64(plotIndex)*0.2
				phase := float64(timestamp)/10.0 + float64(chartIndex)*math.Pi/4 + float64(plotIndex)*math.Pi/8

				// Calculate the value for this plot
				plotID := fmt.Sprintf("plot%d", plotIndex+1)
				values[plotID] = math.Sin(2*math.Pi*frequency*x + phase)
			}

			// Store this point with its timestamp and values
			points[i] = m.ChartPoint{
				Timestamp: time.Unix(currentTimestamp, 0).Format("15:04:05"),
				Values:    values,
			}
		}

		// Store this chart with its points
		charts[chartIndex] = m.ChartData{
			ChartID: fmt.Sprintf("chart%d", chartIndex+1),
			Data:    points,
		}
	}

	a.LastData = charts
}

// StartDataGenerator initializes and starts the data generation routine
func StartDataGenerator(ctx context.Context, wg *sync.WaitGroup, a *app.App) {
	defer wg.Done()

	for {
		select {
		case <-ctx.Done():
			log.Println("Data generator stopping due to shutdown signal")
			return
		default:
			time.Sleep(time.Duration(a.Config.PollInterval) * time.Millisecond)
			timestamp := time.Now().Unix()
			generateData(a, timestamp)
			log.Printf("Generated data at %s for %d plot per chart with %d points",
				time.Unix(timestamp, 0).Format("15:04:05"),
				a.Config.NumPlots,
				a.Config.NumPoints)
		}
	}
}
