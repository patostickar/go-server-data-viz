package service

import (
	"fmt"
	"github.com/patostickar/go-server-data-viz/src/datasource"
	"github.com/patostickar/go-server-data-viz/src/models"
	log "github.com/sirupsen/logrus"
	"math"
	"sync"
	"time"
)

type PlotSettings struct {
	NumPlots     int
	NumPoints    int
	PollInterval int
}

type Service struct {
	Logger       *log.Logger
	plotSettings PlotSettings
	Store        datasource.DataSource
	mu           sync.RWMutex
}

func New(plotSettings PlotSettings, datasource datasource.DataSource) *Service {
	return &Service{
		Logger:       log.New(),
		plotSettings: plotSettings,
		Store:        datasource,
		mu:           sync.RWMutex{},
	}
}

func (s *Service) GetSettings() PlotSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.plotSettings
}

func (s *Service) SetSettings(settings PlotSettings) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.plotSettings = settings
}

func (s *Service) GenerateChartsData(numPlots, numPoints int, timestamp int64) {
	numCharts := 3
	charts := make([]models.ChartData, numCharts)

	waveFunctions := []func(float64, float64, float64) float64{
		s.sineWave,
		s.cosineWave,
		s.rampWave,
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
	s.Store.Create("charts", charts)
}

func (s *Service) sineWave(x, frequency, phase float64) float64 {
	return math.Sin(2*math.Pi*frequency*x + phase)
}

func (s *Service) cosineWave(x, frequency, phase float64) float64 {
	return math.Cos(2*math.Pi*frequency*x + phase)
}

func (s *Service) squareWave(x, frequency, phase float64) float64 {
	return math.Copysign(1, math.Sin(2*math.Pi*frequency*x+phase))
}

func (s *Service) triangleWave(x, frequency, phase float64) float64 {
	return math.Abs(math.Mod(4*x*frequency+phase, 4)-2) / 2
}

func (s *Service) sawtoothWave(x, frequency, phase float64) float64 {
	return math.Mod(2*x*frequency+phase, 2) - 1
}

func (s *Service) rampWave(x, frequency, phase float64) float64 {
	return math.Mod(x*frequency+phase, 1)
}
