package rest

import (
	"encoding/json"
	"github.com/patostickar/go-server-data-viz/app"
	"github.com/patostickar/go-server-data-viz/models"
	"net/http"
)

// configHandler updates the application configuration
func configHandler(w http.ResponseWriter, r *http.Request, a *app.App) {
	var newConfig models.ConfigRequest
	if err := json.NewDecoder(r.Body).Decode(&newConfig); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate configuration parameters
	if newConfig.NumPlotsPerChart < 1 || newConfig.NumPlotsPerChart > 100 {
		http.Error(w, "NumPlots must be between 1 and 100", http.StatusBadRequest)
		return
	}
	if newConfig.NumPoints < 10 || newConfig.NumPoints > 1_000_000 {
		http.Error(w, "NumPlots must be between 10 and 1,000,000", http.StatusBadRequest)
		return
	}

	// Update configuration
	a.Mutex.Lock()
	a.Config.NumPlots = newConfig.NumPlotsPerChart
	a.Config.NumPoints = newConfig.NumPoints
	a.Config.PollInterval = newConfig.PollInterval
	a.Mutex.Unlock()

	// Send response
	w.Header().Set("Content-Type", "a/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Configuration updated successfully",
		"config":  newConfig,
	})
}

// dataHandler returns the current chart data
func dataHandler(w http.ResponseWriter, _ *http.Request, a *app.App) {
	a.Mutex.RLock()
	data := a.LastData
	a.Mutex.RUnlock()

	w.Header().Set("Content-Type", "a/json")
	json.NewEncoder(w).Encode(data)
}

// getConfigHandler returns the current configuration
func getConfigHandler(w http.ResponseWriter, _ *http.Request, a *app.App) {
	a.Mutex.RLock()
	currentConfig := models.ConfigRequest{
		NumPlotsPerChart: a.Config.NumPlots,
		NumPoints:        a.Config.NumPoints,
		PollInterval:     a.Config.PollInterval,
	}
	a.Mutex.RUnlock()

	w.Header().Set("Content-Type", "a/json")
	json.NewEncoder(w).Encode(currentConfig)
}
