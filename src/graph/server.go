package graph

import (
	"context"
	"errors"
	"fmt"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/service"
	"github.com/vektah/gqlparser/v2/ast"
	"net/http"
	"sync"
)

func StartGqlServer(wg *sync.WaitGroup, ctx context.Context, cfg config.Config, s *service.Service) {
	defer wg.Done()

	srv := handler.New(NewExecutableSchema(Config{Resolvers: &Resolver{}}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	server := &http.Server{
		Addr:    ":" + cfg.GetGraphQlPort(),
		Handler: nil,
	}

	go func() {
		s.Logger.Infof("GraphQL Server starting on :%s", cfg.GetGraphQlPort())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("GraphQL server error: %v", err))
		}
	}()

	go func() {
		<-ctx.Done()
		s.Logger.Infof("Shutting down GraphQL server")
		if err := server.Shutdown(context.Background()); err != nil {
			s.Logger.Errorf("GraphQL server shutdown error: %v", err)
		}
	}()
}
