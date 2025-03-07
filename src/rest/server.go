package rest

import (
	"errors"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/patostickar/go-server-data-viz/app"
	"log"
	"net/http"
	"sync"
	"time"
)

func StartHTTPServer(wg *sync.WaitGroup, application *app.App) *http.Server {
	defer wg.Done()

	// Setup router and routes
	r := mux.NewRouter()
	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		configHandler(w, r, application)
	}).Methods("POST")

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		getConfigHandler(w, r, application)
	}).Methods("GET")

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		dataHandler(w, r, application)
	}).Methods("GET")

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (change this in production)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Timing-Allow-Origin"}),
	)

	// Configure the server
	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsHandler(r),
	}

	// Run server in a goroutine so it doesn't block
	go func() {
		log.Printf("HTTP Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	return srv
}
