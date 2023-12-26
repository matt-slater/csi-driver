// Package main is the entrypoint for the csi-driver binary.
package main

import (
	"os"
	"os/signal"
	"syscall"

	"csi-driver/internal/pkg/driver"
	"csi-driver/internal/pkg/server"
	"csi-driver/internal/pkg/storage"

	"github.com/caarlos0/env/v10"
	"go.uber.org/zap"
	"k8s.io/mount-utils"
)

const (
	unixDomain = "unix"
	name       = "csi-driver.mattslater.io"
)

type envConfig struct {
	NodeID        string `env:"NODE_ID"`
	CSISocketPath string `env:"CSI_SOCKET_PATH"`
}

var (
	version string
	commit  string
)

func run() int {
	logger := zap.Must(zap.NewProduction(zap.Fields(zap.String("component", "csi-driver"))))
	defer logger.Sync() //nolint:errcheck
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

	storageBackend, err := storage.NewFilesystem(
		logger.With(zap.String("subsystem", "fs storage backend")),
		"/storage-dir",
		os.DirFS("/"),
	)
	if err != nil {
		sugar.Fatal("failed to create storage backend", err)
	}

	grpcServer, err := server.NewExtendedGRPCServer(
		unixDomain,
		envVars.CSISocketPath,
		&driver.IdentityServer{
			Name:    name,
			Version: version,
		},
		&driver.NodeServer{
			NodeID:         envVars.NodeID,
			Logger:         logger,
			Mounter:        mount.New(""),
			StorageBackend: storageBackend,
		},
		logger,
	)
	if err != nil {
		sugar.Fatal("failed to create grpcServer", err)
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
		return 1
	case <-stopChan:
		sugar.Info("caught os signal. shutting down")
	}
	return 0
}

func main() {
	os.Exit(run())
}
