package grpc

import (
	"context"
	pb "github.com/patostickar/go-server-data-viz/models"
	"github.com/patostickar/go-server-data-viz/src/config"
	"github.com/patostickar/go-server-data-viz/src/models"
	"github.com/patostickar/go-server-data-viz/src/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/emptypb"
)

type chartServiceServer struct {
	pb.ChartServiceServer
	log     *log.Entry
	cfg     config.Config
	service *service.Service
	ctx     context.Context
}

func newChartServiceServer(ctx context.Context, cfg config.Config, s *service.Service, logger *log.Entry) pb.ChartServiceServer {
	return &chartServiceServer{
		cfg:     cfg,
		log:     logger,
		service: s,
		ctx:     ctx,
	}
}

func (s *chartServiceServer) GetChartData(context.Context, *emptypb.Empty) (*pb.ChartDataList, error) {
	data, err := s.service.Store.Read(config.ChartsKey)
	if err != nil {
		return nil, err
	}

	s.log.Debugf("returning %d charts", len(data.([]models.ChartData)))
	return &pb.ChartDataList{
		Items: data.(models.Charts).ToProto(),
	}, nil
}
