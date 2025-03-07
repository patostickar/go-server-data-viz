package main

import (
	"encoding/json"
	"net/http"
)

// configHandler updates the application configuration
func configHandler(w http.ResponseWriter, r *http.Request) {
	var newConfig ConfigRequest
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
	config.mutex.Lock()
	config.NumCharts = newConfig.NumCharts
	config.NumPoints = newConfig.NumPoints
	config.PollInterval = newConfig.PollInterval
	config.mutex.Unlock()

	// Send response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Configuration updated successfully",
		"config":  newConfig,
	})
}

// dataHandler returns the current chart data
func dataHandler(w http.ResponseWriter, _ *http.Request) {
	config.mutex.RLock()
	data := config.LastData
	config.mutex.RUnlock()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

// getConfigHandler returns the current configuration
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
