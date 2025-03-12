package main

import (
	"context"
	"github.com/patostickar/go-server-data-viz/app"
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
	httpServer := rest.StartHTTPServer(&wg, a)

	//TODO: add another wg with the gql server in awesomeProject folder

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

	// Shutdown HTTP server if it exists
	if httpServer != nil {
		log.Println("Shutting down HTTP server")
		if err := httpServer.Shutdown(timeoutCtx); err != nil {
			log.Printf("HTTP server shutdown error: %v", err)
		}
	}

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

	log.Println("Server exiting")
	os.Exit(0)
}
