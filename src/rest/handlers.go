package rest

import (
	"encoding/json"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/models"
	"github.com/patostickar/go-server-data-viz/src/service"
	"net/http"
	"time"
)

// NewSettingsRequest represents the configuration requested by the client
type NewSettingsRequest struct {
	NumPlotsPerChart int `json:"numPlotsPerChart"`
	NumPoints        int `json:"numPoints"`
	PollInterval     int `json:"pollInterval"`
}

// settingsHandler updates the application configuration
func (s *Server) settingsHandler(w http.ResponseWriter, r *http.Request) {
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
	s.service.SetSettings(service.PlotSettings{
		NumPlots:     settings.NumPlotsPerChart,
		NumPoints:    settings.NumPoints,
		PollInterval: settings.PollInterval,
	})

	// Send response
	w.Header().Set("Content-Type", "s/json")
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"message":  "Configuration updated successfully",
		"settings": settings,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// dataHandler returns the current chart data
func (s *Server) dataHandler(w http.ResponseWriter, _ *http.Request) {
	data, err := s.service.Store.Read(config.ChartsKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := models.ChartDataTimestamp{
		Timestamp: time.Now().UnixMilli(),
		ChartData: data.([]models.ChartData),
	}

	w.Header().Set("Content-Type", "s/json")
	err = json.NewEncoder(w).Encode(res)
	if err != nil {
		s.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	s.log.Debugf("returning %d charts", len(data.([]models.ChartData)))
}

// getSettingsHandler returns the current configuration
func (s *Server) getSettingsHandler(w http.ResponseWriter, _ *http.Request) {
	currentConfig := s.service.GetSettings()

	w.Header().Set("Content-Type", "s/json")
	err := json.NewEncoder(w).Encode(currentConfig)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
