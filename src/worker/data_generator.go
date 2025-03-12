package worker

import (
	"context"
	"fmt"
	"github.com/patostickar/go-server-data-viz/models"
	"github.com/patostickar/go-server-data-viz/service"
	"math"
	"sync"
	"time"
)

func generateChartsData(numPlots, numPoints int, timestamp int64) []models.ChartData {

	numCharts := 3
	charts := make([]models.ChartData, numCharts)

	waveFunctions := []func(float64, float64, float64) float64{
		sineWave,
		cosineWave,
		rampWave,
	}

	for chartIndex := 0; chartIndex < numCharts; chartIndex++ {
		points := make([]models.ChartPoint, numPoints)

		for i := 0; i < numPoints; i++ {
			x := float64(i) * 0.1
			currentTimestamp := timestamp + int64(i)

			values := make(map[string]float64)

			for plotIndex := 0; plotIndex < numPlots; plotIndex++ {
				frequency := 0.5 + float64(chartIndex)*0.5 + float64(plotIndex)*0.2
				phase := float64(timestamp)/10.0 + float64(chartIndex)*math.Pi/4 + float64(plotIndex)*math.Pi/8

				plotID := fmt.Sprintf("plot%d", plotIndex+1)
				values[plotID] = waveFunctions[chartIndex](x, frequency, phase)
			}

			points[i] = models.ChartPoint{
				Timestamp: time.Unix(currentTimestamp, 0).Format("15:04:05"),
				Values:    values,
			}
		}

		charts[chartIndex] = models.ChartData{
			ChartID: fmt.Sprintf("chart%d", chartIndex+1),
			Data:    points,
		}
	}
	return charts
}

// StartDataGenerator initializes and starts the data generation routine
func StartDataGenerator(wg *sync.WaitGroup, ctx context.Context, s *service.Service) {
	settings := s.GetSettings()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				s.Logger.Infof("Data generator stopping due to shutdown signal")
				return
			default:
				time.Sleep(time.Duration(settings.PollInterval) * time.Millisecond)
				timestamp := time.Now().Unix()
				charts := generateChartsData(settings.NumPlots, settings.NumPoints, timestamp)
				s.DataSource.SaveData(charts)
				s.Logger.Infof("Generated data at %settings for %d plot per chart with %d points",
					time.Unix(timestamp, 0).Format("15:04:05"),
					settings.NumPlots,
					settings.NumPoints)
			}
		}
	}()
}

func sineWave(x, frequency, phase float64) float64 {
	return math.Sin(2*math.Pi*frequency*x + phase)
}

func cosineWave(x, frequency, phase float64) float64 {
	return math.Cos(2*math.Pi*frequency*x + phase)
}

func squareWave(x, frequency, phase float64) float64 {
	return math.Copysign(1, math.Sin(2*math.Pi*frequency*x+phase))
}

func triangleWave(x, frequency, phase float64) float64 {
	return math.Abs(math.Mod(4*x*frequency+phase, 4)-2) / 2
}

func sawtoothWave(x, frequency, phase float64) float64 {
	return math.Mod(2*x*frequency+phase, 2) - 1
}

func rampWave(x, frequency, phase float64) float64 {
	return math.Mod(x*frequency+phase, 1)
}
