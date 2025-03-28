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
	"golang.org/x/sync/errgroup"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	log.SetLevel(log.DebugLevel)
	logger := log.NewEntry(log.StandardLogger())

	mainCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.New()

	s := service.New(
		service.PlotSettings{NumPlots: 1, NumPoints: 100, PollInterval: 1000},
		datasource.NewInMemoryDB(),
	)

	g, gCtx := errgroup.WithContext(mainCtx)

	// Start data generator worker
	g.Go(func() error {
		return worker.StartDataGenerator(gCtx, s)
	})

	// Start HTTP server
	g.Go(func() error {
		return rest.New(gCtx, cfg, s).StartHTTPServer()
	})

	// Start GraphQL server
	g.Go(func() error {
		return graph.New(gCtx, cfg, s).StartGqlServer()
	})

	shutDownErrorChan := make(chan error, 1)
	go func() {
		err := g.Wait()
		if err != nil {
			shutDownErrorChan <- err
		}
		close(shutDownErrorChan)
	}()

	forceShutdownTimeout := time.Second * 1

	select {
	case <-mainCtx.Done():
		logger.Info("Termination signal received. Initiating shutdown...")

		timeoutCtx, cancel := context.WithTimeout(context.Background(), forceShutdownTimeout)
		defer cancel()

		select {
		case err := <-shutDownErrorChan:
			if err != nil {
				logger.Errorf("Application error: %v", err)
				os.Exit(1)
			}
		case <-timeoutCtx.Done():
			logger.Errorf("Application forced shutdown after %s", forceShutdownTimeout)
			os.Exit(1)
		}

	case err := <-shutDownErrorChan:
		if err != nil {
			logger.Errorf("Application error: %v", err)
			os.Exit(1)
		}
	}

	logger.Info("Application shutdown complete")
}
