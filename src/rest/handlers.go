package rest

import (
	"encoding/json"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/models"
	"net/http"
)

// NewSettingsRequest represents the configuration requested by the client
type NewSettingsRequest struct {
	NumPlotsPerChart int `json:"numPlotsPerChart"`
	NumPoints        int `json:"numPoints"`
	PollInterval     int `json:"pollInterval"`
}

// dataHandler returns the current chart data
func (s *Server) dataHandler(w http.ResponseWriter, _ *http.Request) {
	data, err := s.service.Store.Read(config.ChartsKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	res := data.(models.Charts)

	w.Header().Set("Content-Type", "s/json")
	err = json.NewEncoder(w).Encode(res.Data)
	if err != nil {
		s.log.Error(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	s.log.Debugf("returning %d charts", len(res.Data))
}
