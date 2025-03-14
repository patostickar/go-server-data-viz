package main

import (
	"context"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/datasource"
	"github.com/patostickar/go-server-data-viz/src/graph"
	"github.com/patostickar/go-server-data-viz/src/rest"
	"github.com/patostickar/go-server-data-viz/src/service"
	"github.com/patostickar/go-server-data-viz/src/worker"
	log "github.com/sirupsen/logrus"

	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	logger := log.WithField("service", "go-server-data-viz")

	cfg := config.New()

	s := service.New(
		service.PlotSettings{NumPlots: 1, NumPoints: 100, PollInterval: 1000},
		datasource.NewInMemoryDB(),
	)

	var wg sync.WaitGroup

	cancelCtx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	worker.StartDataGenerator(&wg, cancelCtx, s)

	wg.Add(1)
	rest.New(&wg, cancelCtx, cfg, s).StartHTTPServer()

	wg.Add(1)
	graph.New(&wg, cancelCtx, cfg, s).StartGqlServer()

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c
	logger.Infof("Shutdown signal received")
	cancel()

	// Create s deadline to wait for
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Wait for all goroutines to finish (with s timeout)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		logger.Infof("All services shut down properly")
	case <-timeoutCtx.Done():
		logger.Error("Shutdown timed out, forcing exit")
	}

	logger.Infof("Exiting")
	os.Exit(0)
}
