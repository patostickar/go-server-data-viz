package rest

import (
	"encoding/json"
	"github.com/patostickar/go-server-data-viz/src/service"
	"net/http"
)

// NewSettingsRequest represents the configuration requested by the client
type NewSettingsRequest struct {
	NumPlotsPerChart int `json:"numPlotsPerChart"`
	NumPoints        int `json:"numPoints"`
	PollInterval     int `json:"pollInterval"`
}

// settingsHandler updates the application configuration
func settingsHandler(w http.ResponseWriter, r *http.Request, a *service.Service) {
	var settings NewSettingsRequest
	if err := json.NewDecoder(r.Body).Decode(&settings); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate configuration parameters
	if settings.NumPlotsPerChart < 1 || settings.NumPlotsPerChart > 100 {
		http.Error(w, "NumPlots must be between 1 and 100", http.StatusBadRequest)
		return
	}
	if settings.NumPoints < 10 || settings.NumPoints > 1_000_000 {
		http.Error(w, "NumPlots must be between 10 and 1,000,000", http.StatusBadRequest)
		return
	}

	// Update configuration
	a.SetSettings(service.PlotSettings{
		NumPlots:     settings.NumPlotsPerChart,
		NumPoints:    settings.NumPoints,
		PollInterval: settings.PollInterval,
	})

	// Send response
	w.Header().Set("Content-Type", "a/json")
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Configuration updated successfully",
		"settings": settings,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// dataHandler returns the current chart data
func dataHandler(w http.ResponseWriter, _ *http.Request, a *service.Service) {
	data, err := a.Store.Read("charts")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "a/json")
	err = json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// getSettingsHandler returns the current configuration
func getSettingsHandler(w http.ResponseWriter, _ *http.Request, a *service.Service) {
	currentConfig := a.GetSettings()

	w.Header().Set("Content-Type", "a/json")
	err := json.NewEncoder(w).Encode(currentConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
