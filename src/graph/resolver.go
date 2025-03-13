package graph

import "github.com/patostickar/go-server-data-viz/src/graph/model"

//go:generate go run github.com/99designs/gqlgen generate

type Resolver struct {
	todos []*model.Todo
}
