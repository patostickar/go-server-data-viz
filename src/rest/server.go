package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/patostickar/go-server-data-viz/service"
	"net/http"
	"sync"
	"time"
)

func StartHTTPServer(wg *sync.WaitGroup, ctx context.Context, s *service.Service) {
	defer wg.Done()

	// Setup router and routes
	r := mux.NewRouter()
	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		settingsHandler(w, r, s)
	}).Methods("POST")

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		getSettingsHandler(w, r, s)
	}).Methods("GET")

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		dataHandler(w, r, s)
	}).Methods("GET")

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (change this in production)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Timing-Allow-Origin"}),
	)

	// Configure the server
	server := &http.Server{
		Addr: "0.0.0.0:" + s.Config.HttpPort,
		// Good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsHandler(r),
	}

	// Run server in s goroutine so it doesn't block
	go func() {
		s.Logger.Infof("HTTP Server starting on :8080")
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("HTTP server error: %v", err))

		}
	}()

	go func() {
		<-ctx.Done()
		s.Logger.Infof("Shutting down HTTP server")
		if err := server.Shutdown(context.Background()); err != nil {
			s.Logger.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

}
