package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"log"
	"math"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

type DataPoint struct {
	Timestamp string  `json:"timestamp"`
	Value     float64 `json:"value"`
}

type ChartData struct {
	ChartID string      `json:"chartId"`
	Data    []DataPoint `json:"data"`
}

type ConfigRequest struct {
	NumCharts    int `json:"numCharts"`
	NumPoints    int `json:"numPoints"`
	PollInterval int
}

type AppConfig struct {
	NumCharts    int
	NumPoints    int
	PollInterval int
	LastData     []ChartData
	mutex        sync.RWMutex
}

var config = AppConfig{
	NumCharts:    1,
	NumPoints:    100,
	PollInterval: 1000,
}

func generateData() ([]ChartData, int64) {
	config.mutex.RLock()
	defer config.mutex.RUnlock()

	now := time.Now().Unix()
	charts := make([]ChartData, config.NumCharts)

	for chartIndex := 0; chartIndex < config.NumCharts; chartIndex++ {
		points := make([]DataPoint, config.NumPoints)
		frequency := 0.5 + float64(chartIndex)*0.5
		phase := float64(now)/10.0 + float64(chartIndex)*math.Pi/4

		for i := 0; i < config.NumPoints; i++ {
			x := float64(i) * 0.1
			timestamp := now + int64(i)
			points[i] = DataPoint{
				Timestamp: time.Unix(timestamp, 0).Format("15:04:05"),
				Value:     math.Sin(2*math.Pi*frequency*x + phase),
			}
		}

		charts[chartIndex] = ChartData{
			ChartID: fmt.Sprintf("chart%d", chartIndex+1),
			Data:    points,
		}
	}

	config.LastData = charts
	return charts, now
}

func configHandler(w http.ResponseWriter, r *http.Request) {
	var newConfig ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if newConfig.NumCharts < 1 || newConfig.NumCharts > 100 {
		http.Error(w, "NumCharts must be between 1 and 100", http.StatusBadRequest)
		return
	}
	if newConfig.NumPoints < 10 || newConfig.NumPoints > 1_000_000 {
		http.Error(w, "NumPoints must be between 10 and 1.000.000", http.StatusBadRequest)
		return
	}

	config.mutex.Lock()
	config.NumCharts = newConfig.NumCharts
	config.NumPoints = newConfig.NumPoints
	config.PollInterval = newConfig.PollInterval
	config.mutex.Unlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Configuration updated successfully",
		"config":  newConfig,
	})
}

func dataHandler(w http.ResponseWriter, _ *http.Request) {
	config.mutex.RLock()
	data := config.LastData
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func getConfigHandler(w http.ResponseWriter, _ *http.Request) {
	config.mutex.RLock()
	currentConfig := ConfigRequest{
		NumCharts:    config.NumCharts,
		NumPoints:    config.NumPoints,
		PollInterval: config.PollInterval,
	}
	config.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}

func main() {
	var wait time.Duration
	flag.DurationVar(&wait, "graceful-timeout", time.Second*15, "the duration for which the server gracefully wait for existing connections to finish - e.g. 15s or 1m")
	flag.Parse()

	go func() {
		for {
			time.Sleep(time.Duration(config.PollInterval) * time.Millisecond)
			_, now := generateData()
			fmt.Println("Generated data at", time.Unix(now, 0), "for", config.NumCharts, "charts")
		}
	}()

	r := mux.NewRouter()

	r.HandleFunc("/config", configHandler).Methods("POST")
	r.HandleFunc("/config", getConfigHandler).Methods("GET")
	r.HandleFunc("/data", dataHandler).Methods("GET")

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (change this in production)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization"}),
	)

	srv := &http.Server{
		Addr: "0.0.0.0:8080",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsHandler(r), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Printf("Server starting on :8080")
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	log.Println("shutting down")
	os.Exit(0)
}
