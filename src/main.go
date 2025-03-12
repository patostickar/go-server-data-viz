package main

import (
	"context"
	"github.com/patostickar/go-server-data-viz/datasource"
	"github.com/patostickar/go-server-data-viz/graph"
	"github.com/patostickar/go-server-data-viz/rest"
	"github.com/patostickar/go-server-data-viz/service"
	"github.com/patostickar/go-server-data-viz/worker"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	s := service.New(
		service.Config{HttpPort: "8080", GraphQlPort: "8081"},
		service.PlotSettings{NumPlots: 1, NumPoints: 100, PollInterval: 1000},
		datasource.NewInMemoryStore(),
	)

	var wg sync.WaitGroup

	cancelCtx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	worker.StartDataGenerator(&wg, cancelCtx, s)

	wg.Add(1)
	rest.StartHTTPServer(&wg, cancelCtx, s)

	wg.Add(1)
	graph.StartGqlServer(&wg, cancelCtx, s)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c
	s.Logger.Infof("Shutdown signal received")
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
		s.Logger.Infof("All services shut down properly")
	case <-timeoutCtx.Done():
		s.Logger.Error("Shutdown timed out, forcing exit")
	}

	s.Logger.Infof("Exiting")
	os.Exit(0)
}
