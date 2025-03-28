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
	"github.com/gorilla/handlers"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/service"
	log "github.com/sirupsen/logrus"
	"github.com/vektah/gqlparser/v2/ast"
	"golang.org/x/sync/errgroup"
	"net/http"
)

type Server struct {
	server  *http.Server
	log     *log.Entry
	cfg     config.Config
	service *service.Service
	ctx     context.Context
}

func New(ctx context.Context, cfg config.Config, s *service.Service) *Server {
	logger := log.WithField("server", "graphql")

	return &Server{
		cfg:     cfg,
		log:     logger,
		service: s,
		ctx:     ctx,
	}
}
func (s *Server) StartGqlServer() error {

	srv := handler.New(NewExecutableSchema(Config{
		Resolvers: NewResolver(&s.cfg, s.log, s.service),
	}))

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

	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (change this in production)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Timing-Allow-Origin"}),
	)

	server := &http.Server{
		Addr:    ":" + s.cfg.GetGraphQlPort(),
		Handler: corsHandler(srv),
	}

	g := errgroup.Group{}

	g.Go(func() error {
		s.log.Infof("GraphQL Server starting on :%s", s.cfg.GetGraphQlPort())
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("GraphQL server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		<-s.ctx.Done()
		s.log.Infof("Shutting down GraphQL server")
		if err := server.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("GraphQL server shutdown error: %v", err)
		}
		return nil
	})

	return g.Wait()
}
