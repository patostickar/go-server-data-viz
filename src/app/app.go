package app

import (
	"github.com/patostickar/go-server-data-viz/models"
	"sync"
)

type plotSettings struct {
	NumPlots     int
	NumPoints    int
	PollInterval int
}

type config struct {
	Port string
}

type App struct {
	Config       config
	PlotSettings plotSettings
	LastData     []models.ChartData
	Mutex        sync.RWMutex
}

func New() *App {
	return &App{
		Config: config{
			Port: "8080",
		},
		PlotSettings: plotSettings{
			NumPlots:     1,
			NumPoints:    100,
			PollInterval: 1000,
		},
	}
}
