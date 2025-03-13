package graph

import (
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/service"
	"github.com/sirupsen/logrus"
)

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	config config.Config
	s      *service.Service
	logger *logrus.Logger
}

func NewResolver(cfg config.Config, s *service.Service) *Resolver {
	logger := logrus.New().WithField("service", "graphql")
	logger.Level = logrus.DebugLevel
	return &Resolver{
		config: cfg,
		s:      s,
		logger: logger.Logger,
	}
}
