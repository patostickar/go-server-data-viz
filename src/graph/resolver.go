package graph

import (
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/service"
	log "github.com/sirupsen/logrus"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	config *config.Config
	s      *service.Service
	logger *log.Entry
}

func NewResolver(cfg *config.Config, logger *log.Entry, s *service.Service) *Resolver {
	return &Resolver{
		config: cfg,
		s:      s,
		logger: logger,
	}
}
