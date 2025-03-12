package service

import (
	"github.com/patostickar/go-server-data-viz/datasource"
	log "github.com/sirupsen/logrus"
	"sync"
)

type Config struct {
	HttpPort    string
	GraphQlPort string
}

type PlotSettings struct {
	NumPlots     int
	NumPoints    int
	PollInterval int
}

type Service struct {
	Config       Config
	Logger       *log.Logger
	plotSettings PlotSettings
	DataSource   datasource.DataSource
	mu           sync.RWMutex
}

func New(config Config, plotSettings PlotSettings, datasource datasource.DataSource) *Service {
	return &Service{
		Config:       config,
		Logger:       log.New(),
		plotSettings: plotSettings,
		DataSource:   datasource,
		mu:           sync.RWMutex{},
	}
}

func (a *Service) GetSettings() PlotSettings {
	a.mu.RLock()
	defer a.mu.RUnlock()
	return a.plotSettings
}

func (a *Service) SetSettings(settings PlotSettings) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.plotSettings = settings
}
