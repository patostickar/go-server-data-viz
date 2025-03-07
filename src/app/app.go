package app

import (
	"github.com/patostickar/go-server-data-viz/models"
	"sync"
)

type config struct {
	NumCharts    int
	NumPoints    int
	PollInterval int
}

type App struct {
	Config   config
	LastData []models.ChartData
	Mutex    sync.RWMutex
}

func NewApp() *App {
	return &App{
		Config: config{
			NumCharts:    1,
			NumPoints:    100,
			PollInterval: 1000,
		},
	}
}
