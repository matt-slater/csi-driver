// Package main is the entrypoint for the csi-driver binary.
package main

import (
	"fmt"
	"net"
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

func run() error {
	logger := zap.Must(zap.NewProduction(zap.Fields(zap.String("component", "csi-driver"))))
	defer logger.Sync() //nolint:errcheck
	sugar := logger.Sugar()

	sugar.Infow("starting up...",
		"commit", commit,
		"version", version,
	)

	envVars := &envConfig{} //nolint:exhaustivestruct,exhaustruct

	err := env.Parse(envVars)
	if err != nil {
		sugar.Fatal("failed to parse env vars", err)
	}

	storageBackend, err := storage.NewFilesystem(
		logger.With(zap.String("subsystem", "fs storage backend")),
		"/storage-dir",
		os.DirFS("/"),
		mount.New(""),
	)
	if err != nil {
		sugar.Fatal("failed to create storage backend", err)
	}

	err = os.Remove(envVars.CSISocketPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove unix socket file: %w", err)
	}

	listener, err := net.Listen(unixDomain, envVars.CSISocketPath)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	grpcServer := server.NewExtendedGRPCServer(
		listener,
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
		sugar.Errorw("caught error", "error", err)

		return err
	case <-stopChan:
		sugar.Info("caught os signal. shutting down")
	}

	return nil
}

func main() {
	err := run()
	if err != nil {
		os.Exit(1)
	}
}
