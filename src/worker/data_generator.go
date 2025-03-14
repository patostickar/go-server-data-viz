package worker

import (
	"context"
	"github.com/patostickar/go-server-data-viz/src/service"
	log "github.com/sirupsen/logrus"
	"sync"
	"time"
)

// StartDataGenerator initializes and starts the data generation routine
func StartDataGenerator(wg *sync.WaitGroup, ctx context.Context, s *service.Service) {
	logger := log.WithField("worker", "data-generator")

	settings := s.GetSettings()
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				logger.Infof("Data generator stopping due to shutdown signal")
				return
			default:
				time.Sleep(time.Duration(settings.PollInterval) * time.Millisecond)
				timestamp := time.Now().Unix()
				s.GenerateChartsData(settings.NumPlots, settings.NumPoints, timestamp)
				logger.WithFields(log.Fields{
					"timestamp": timestamp,
					"numPlots":  settings.NumPlots,
					"numPoints": settings.NumPoints,
				}).Debug("Generated data")
			}
		}
	}()
}
