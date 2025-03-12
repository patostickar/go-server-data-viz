package graph

import (
	"context"
	"errors"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
	"log"
	"net/http"
	"os"
	"sync"
)

const defaultPort = "8081" // Use a different port than the HTTP server

func StartGqlServer(wg *sync.WaitGroup, ctx context.Context) {
	defer wg.Done()

	port := os.Getenv("GRAPHQL_PORT")
	if port == "" {
		port = defaultPort
	}

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
		Addr:    ":" + port,
		Handler: nil,
	}

	go func() {
		log.Printf("GraphQL Server starting on :%s", port)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("GraphQL server error: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		log.Println("Shutting down GraphQL server")
		if err := server.Shutdown(context.Background()); err != nil {
			log.Printf("GraphQL server shutdown error: %v", err)
		}
	}()
}
