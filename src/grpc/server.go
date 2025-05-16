package grpc

import (
	"context"
	"fmt"
	pb "github.com/patostickar/go-server-data-viz/models"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/service"
	log "github.com/sirupsen/logrus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"net"
)

type Server struct {
	server  *grpc.Server
	log     *log.Entry
	cfg     config.Config
	service *service.Service
	ctx     context.Context
}

func New(ctx context.Context, cfg config.Config, s *service.Service) *Server {
	logger := log.WithField("server", "grpc")

	return &Server{
		cfg:     cfg,
		log:     logger,
		service: s,
		ctx:     ctx,
	}
}

func (s *Server) StartGrpcServer() error {
	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.cfg.GetGrpcPort()))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	var opts []grpc.ServerOption
	s.server = grpc.NewServer(opts...)

	g := errgroup.Group{}

	g.Go(func() error {
		s.log.Infof("gRPC server starting on :%s", s.cfg.GetGrpcPort())
		pb.RegisterChartServiceServer(s.server, newChartServiceServer(s.ctx, s.cfg, s.service, s.log))
		if err = s.server.Serve(lis); err != nil {
			return fmt.Errorf("gRPC server error: %v", err)
		}
		return nil
	})

	g.Go(func() error {
		<-s.ctx.Done()
		s.log.Infof("gRPC server shutting down")
		s.server.GracefulStop()
		return nil
	})

	return g.Wait()
}
