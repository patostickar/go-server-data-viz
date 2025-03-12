package rest

import (
	"encoding/json"
	"github.com/patostickar/go-server-data-viz/app"
	"github.com/patostickar/go-server-data-viz/models"
	"net/http"
)

// configHandler updates the application configuration
func configHandler(w http.ResponseWriter, r *http.Request, application *app.App) {
	var newConfig models.ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate configuration parameters
	if newConfig.NumCharts < 1 || newConfig.NumCharts > 100 {
		http.Error(w, "NumCharts must be between 1 and 100", http.StatusBadRequest)
		return
	}
	if newConfig.NumPoints < 10 || newConfig.NumPoints > 1_000_000 {
		http.Error(w, "NumPoints must be between 10 and 1,000,000", http.StatusBadRequest)
		return
	}

	// Update configuration
	application.Mutex.Lock()
	application.Config.NumCharts = newConfig.NumCharts
	application.Config.NumPoints = newConfig.NumPoints
	application.Config.PollInterval = newConfig.PollInterval
	application.Mutex.Unlock()

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Configuration updated successfully",
		"config":  newConfig,
	})
}

// dataHandler returns the current chart data
func dataHandler(w http.ResponseWriter, _ *http.Request, application *app.App) {
	application.Mutex.RLock()
	data := application.LastData
	application.Mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// getConfigHandler returns the current configuration
func getConfigHandler(w http.ResponseWriter, _ *http.Request, application *app.App) {
	application.Mutex.RLock()
	currentConfig := models.ConfigRequest{
		NumCharts:    application.Config.NumCharts,
		NumPoints:    application.Config.NumPoints,
		PollInterval: application.Config.PollInterval,
	}
	application.Mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(currentConfig)
}
