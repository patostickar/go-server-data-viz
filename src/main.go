package main

import (
	"context"
	"github.com/patostickar/go-server-data-viz/app"
	"github.com/patostickar/go-server-data-viz/graph"
	"github.com/patostickar/go-server-data-viz/rest"
	"github.com/patostickar/go-server-data-viz/utils"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	a := app.New()

	var wg sync.WaitGroup

	cancelCtx, cancel := context.WithCancel(context.Background())
	wg.Add(1)
	go utils.StartDataGenerator(cancelCtx, &wg, a)

	wg.Add(1)
	rest.StartHTTPServer(&wg, cancelCtx, a)

	wg.Add(1)
	graph.StartGqlServer(&wg, cancelCtx)

	// Setup graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal
	<-c
	log.Println("Shutdown signal received")
	cancel()

	// Create a deadline to wait for
	timeoutCtx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Wait for all goroutines to finish (with a timeout)
	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		log.Println("All services shut down properly")
	case <-timeoutCtx.Done():
		log.Println("Shutdown timed out, forcing exit")
	}

	log.Println("Exiting")
	os.Exit(0)
}
