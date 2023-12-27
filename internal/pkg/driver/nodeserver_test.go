package driver_test

import (
	"context"
	"csi-driver/internal/pkg/driver"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func TestNodeServer_NodeStageVolume(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{}

	resp, err := nodeServer.NodeStageVolume(context.Background(), &csi.NodeStageVolumeRequest{})
	if err == nil {
		t.Fatal("unexpected nil error for unimplemented gRPC method")
	}

	if resp != nil {
		t.Fatalf("unexpected non-nil response: %v", resp)
	}
}

func TestNodeServer_NodeUnstageVolume(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{}

	resp, err := nodeServer.NodeUnstageVolume(context.Background(), &csi.NodeUnstageVolumeRequest{})
	if err == nil {
		t.Fatal("unexpected nil error for unimplemented gRPC method")
	}

	if resp != nil {
		t.Fatalf("unexpected non-nil response: %v", resp)
	}
}

func TestNodeServer_NodeGeVolumeStats(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{}

	resp, err := nodeServer.NodeGetVolumeStats(context.Background(), &csi.NodeGetVolumeStatsRequest{})
	if err == nil {
		t.Fatal("unexpected nil error for unimplemented gRPC method")
	}

	if resp != nil {
		t.Fatalf("unexpected non-nil response: %v", resp)
	}
}

func TestNodeServer_NodeExpandVolume(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{}

	resp, err := nodeServer.NodeExpandVolume(context.Background(), &csi.NodeExpandVolumeRequest{})
	if err == nil {
		t.Fatal("unexpected nil error for unimplemented gRPC method")
	}

	if resp != nil {
		t.Fatalf("unexpected non-nil response: %v", resp)
	}
}

func TestNodeServer_NodeGetVolumeStats(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{}

	resp, err := nodeServer.NodeGetVolumeStats(context.Background(), &csi.NodeGetVolumeStatsRequest{})
	if err == nil {
		t.Fatal("unexpected nil error for unimplemented gRPC method")
	}

	if resp != nil {
		t.Fatalf("unexpected non-nil response: %v", resp)
	}
}

func TestNodeServer_NodeGetInfo(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{
		NodeID: "test",
	}

	resp, err := nodeServer.NodeGetInfo(context.Background(), &csi.NodeGetInfoRequest{})
	if err != nil {
		t.Fatalf("unexpected error for gRPC method: %v", err)
	}

	if resp == nil {
		t.Fatal("unexpected nil response")
	}
}

func TestNodeServer_NodeGetCapabilities(t *testing.T) {
	t.Parallel()

	nodeServer := &driver.NodeServer{}

	resp, err := nodeServer.NodeGetCapabilities(context.Background(), &csi.NodeGetCapabilitiesRequest{})
	if err != nil {
		t.Fatalf("unexpected error for gRPC method: %v", err)
	}

	if resp == nil {
		t.Fatal("unexpected nil response")
	}
}
