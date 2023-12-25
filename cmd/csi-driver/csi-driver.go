package main

import (
	"csi-driver/internal/pkg/driver"
	"csi-driver/internal/pkg/server"
	"os"
	"os/signal"
	"syscall"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
)

const (
	unixDomain = "unix"
)

type envConfig struct {
	NodeID        string `env:"NODE_ID"`
	CSISocketPath string `env:"CSI_SOCKET_PATH"`
}

var version string
var commit string

func main() {
	logger := zap.Must(zap.NewProduction(zap.Fields(zap.String("component", "csi-driver"))))
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infow("starting up...",
		"commit", commit,
		"version", version,
	)

	envVars := &envConfig{}

	err := env.Parse(envVars)
	if err != nil {
		sugar.Fatal("failed to parse env vars", err)
	}

	grpcServer, err := server.NewExtendedGRPCServer(
		unixDomain,
		envVars.CSISocketPath,
		&driver.IdentityServer{},
		&driver.NodeServer{NodeID: envVars.NodeID},
		logger,
	)
	if err != nil {
		sugar.Fatal("fsiled to create grpcServer", err)
	}

	errChan := make(chan error)
	stopChan := make(chan os.Signal, 1)

	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)

	go func() {
		err = grpcServer.Run()
		if err != nil {
			errChan <- err
		}
	}()

	defer func() {
		grpcServer.GracefulStop()
		sugar.Info("shutdown gRPC server gracefully")
	}()

	select {
	case err := <-errChan:
		sugar.Errorw("caught error", err)
	case <-stopChan:
		sugar.Info("caught os signal. shutting down")
	}
}
