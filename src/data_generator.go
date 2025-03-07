package main

import (
	"fmt"
	"math"
	"time"
)

func generateData() ([]ChartData, int64) {
	config.mutex.RLock()
	defer config.mutex.RUnlock()

	now := time.Now().Unix()
	charts := make([]ChartData, config.NumCharts)

	for chartIndex := 0; chartIndex < config.NumCharts; chartIndex++ {
		points := make([]DataPoint, config.NumPoints)
		frequency := 0.5 + float64(chartIndex)*0.5
		phase := float64(now)/10.0 + float64(chartIndex)*math.Pi/4

		for i := 0; i < config.NumPoints; i++ {
			x := float64(i) * 0.1
			timestamp := now + int64(i)
			points[i] = DataPoint{
				Timestamp: time.Unix(timestamp, 0).Format("15:04:05"),
				Value:     math.Sin(2*math.Pi*frequency*x + phase),
			}
		}

		charts[chartIndex] = ChartData{
			ChartID: fmt.Sprintf("chart%d", chartIndex+1),
			Data:    points,
		}
	}

	config.LastData = charts
	return charts, now
}
