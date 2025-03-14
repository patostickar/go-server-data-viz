package rest

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/service"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	server  *http.Server
	log     *log.Entry
	cfg     config.Config
	service *service.Service
	wg      *sync.WaitGroup
	ctx     context.Context
}

func New(wg *sync.WaitGroup, ctx context.Context, cfg config.Config, s *service.Service) *Server {
	logger := log.WithField("server", "http")

	return &Server{
		cfg:     cfg,
		log:     logger,
		service: s,
		wg:      wg,
		ctx:     ctx,
	}
}

func (s *Server) StartHTTPServer() {
	defer s.wg.Done()

	// Setup router and routes
	r := mux.NewRouter()
	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		s.settingsHandler(w, r)
	}).Methods("POST")

	r.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		s.getSettingsHandler(w, r)
	}).Methods("GET")

	r.HandleFunc("/data", func(w http.ResponseWriter, r *http.Request) {
		s.dataHandler(w, r)
	}).Methods("GET")

	// Configure CORS
	corsHandler := handlers.CORS(
		handlers.AllowedOrigins([]string{"*"}), // Allow all origins (change this in production)
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Authorization", "Timing-Allow-Origin"}),
	)

	// Configure the server
	s.server = &http.Server{
		Addr: "0.0.0.0:" + s.cfg.GetHttpPort(),
		// Good practice to set timeouts to avoid Slowloris attacks
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      corsHandler(r),
	}

	// Run server in s goroutine so it doesn't block
	go func() {
		s.log.Infof("HTTP Server starting on :%s", s.cfg.GetHttpPort())
		if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Errorf("HTTP server error: %v", err))

		}
	}()

	go func() {
		<-s.ctx.Done()
		s.log.Infof("Shutting down HTTP server")
		if err := s.server.Shutdown(context.Background()); err != nil {
			s.log.Errorf("HTTP server shutdown error: %v", err)
		}
	}()

}
