package server_test

import (
	"csi-driver/internal/pkg/driver"
	"csi-driver/internal/pkg/server"
	"testing"
	"time"

	"go.uber.org/zap/zaptest"
	"golang.org/x/net/nettest"
)

func TestNewExtendedGRPCServer(t *testing.T) {
	t.Parallel()

	listener, err := nettest.NewLocalListener("unix")
	if err != nil {
		t.Fatalf("failed to create test listener: %v", err)
	}

	server := server.NewExtendedGRPCServer(
		listener,
		&driver.IdentityServer{},
		&driver.NodeServer{},
		zaptest.NewLogger(t),
	)

	if server == nil {
		t.Fatal("unexpected nil server")
	}
}

func TestExtendedGRPCSercer_GracefulStop(t *testing.T) {
	t.Parallel()

	listener, err := nettest.NewLocalListener("unix")
	if err != nil {
		t.Fatalf("failed to create test listener: %v", err)
	}

	server := server.NewExtendedGRPCServer(
		listener,
		&driver.IdentityServer{},
		&driver.NodeServer{},
		zaptest.NewLogger(t),
	)

	go func() {
		err := server.Run()
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	}()

	time.Sleep(2 * time.Second)

	server.GracefulStop()
}
