package worker

import (
	"context"
	"github.com/patostickar/go-server-data-viz/src/service"
	"sync"
	"time"
)

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
				s.GenerateChartsData(settings.NumPlots, settings.NumPoints, timestamp)
				s.Logger.Infof("Generated data at %s for %d plot per chart with %d points",
					time.Unix(timestamp, 0).Format("15:04:05"),
					settings.NumPlots,
					settings.NumPoints)
			}
		}
	}()
}
