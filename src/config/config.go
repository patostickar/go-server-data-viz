package config

import (
	"github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

const ChartsKey = "CHARTS"

type Config interface {
	GetHttpPort() string
	GetGraphQlPort() string
	GetGrpcPort() string
}

type config struct {
	HttpPort    string `env:"REST_PORT,default=8080"`
	GraphQlPort string `env:"GRAPHQL_PORT,default=8081"`
	GrpcPort    string `env:"GRPC_PORT,default=8082"`
}

func (c config) GetHttpPort() string {
	return c.HttpPort
}

func (c config) GetGraphQlPort() string {
	return c.GraphQlPort
}

func (c config) GetGrpcPort() string {
	return c.GrpcPort
}

func New() Config {
	var conf config
	_, err := env.UnmarshalFromEnviron(&conf)
	if err != nil {
		log.Fatalf("failed to parse configuration: %v", err)
	}
	return conf
}
