package server

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/kubernetes-csi/csi-lib-utils/protosanitizer"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type ExtendedGRPCServer struct {
	server   *grpc.Server
	logger   *zap.Logger
	listener net.Listener
}

func NewExtendedGRPCServer(protocol, endpoint string, is csi.IdentityServer, ns csi.NodeServer, logger *zap.Logger) (*ExtendedGRPCServer, error) {
	err := os.Remove(endpoint)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to remove unix socket file: %w", err)
	}

	listener, err := net.Listen(protocol, endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	opts := []grpc.ServerOption{
		grpc.UnaryInterceptor(loggingInterceptor(logger)),
	}

	server := grpc.NewServer(opts...)

	if is != nil {
		csi.RegisterIdentityServer(server, is)
	}

	if ns != nil {
		csi.RegisterNodeServer(server, ns)
	}

	return &ExtendedGRPCServer{
		server:   server,
		logger:   logger,
		listener: listener,
	}, nil

}

func (gs *ExtendedGRPCServer) Run() error {
	return gs.server.Serve(gs.listener)
}

func (gs *ExtendedGRPCServer) GracefulStop() {
	gs.server.GracefulStop()
}

func (gs *ExtendedGRPCServer) ForceStop() {
	gs.server.Stop()
}

func loggingInterceptor(logger *zap.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger.Info("request recieved",
			zap.String("rpc_method", info.FullMethod),
			zap.Any("request", protosanitizer.StripSecrets(req)),
		)
		resp, err := handler(ctx, req)
		if err != nil {
			logger.Error("failed processing request", zap.Error(err))
		} else {
			logger.Info("request completed",
				zap.Any("response", protosanitizer.StripSecrets(resp)),
			)
		}
		return resp, err
	}
}