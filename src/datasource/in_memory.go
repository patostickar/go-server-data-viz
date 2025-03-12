package datasource

import (
	"github.com/patostickar/go-server-data-viz/models"
	"sync"
)

type DataSource interface {
	GetData() []models.ChartData
	SaveData(data []models.ChartData)
}

type InMemoryStore struct {
	data  []models.ChartData
	mutex sync.RWMutex
}

func NewInMemoryStore() DataSource {
	return &InMemoryStore{
		data: []models.ChartData{},
	}
}

func (s *InMemoryStore) GetData() []models.ChartData {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.data
}

func (s *InMemoryStore) SaveData(newData []models.ChartData) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.data = newData
}
