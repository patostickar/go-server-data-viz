package app

import (
	"github.com/patostickar/go-server-data-viz/models"
	"sync"
)

type config struct {
	NumPlots     int
	NumPoints    int
	PollInterval int
}

type App struct {
	Config   config
	LastData []models.ChartData
	Mutex    sync.RWMutex
}

func New() *App {
	return &App{
		Config: config{
			NumPlots:     1,
			NumPoints:    100,
			PollInterval: 1000,
		},
	}
}
