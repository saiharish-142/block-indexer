package grpcserver

import (
	"net"

	"github.com/example/block-indexer/internal/pb"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// New spins up a gRPC server with the supplied services registered.
func New(addr string, logger *zap.Logger, query pb.QueryServiceServer, stream pb.StreamServiceServer) (*grpc.Server, error) {
	s := grpc.NewServer()
	if query != nil {
		pb.RegisterQueryServiceServer(s, query)
	}
	if stream != nil {
		pb.RegisterStreamServiceServer(s, stream)
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}

	go func() {
		logger.Info("grpc server starting", zap.String("addr", addr))
		if err := s.Serve(lis); err != nil {
			logger.Error("grpc server exited", zap.Error(err))
		}
	}()

	return s, nil
}
