package driver

import (
	"context"
	"fmt"
	"os"

	"csi-driver/internal/pkg/storage"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/mount-utils"
)

const (
	roPerms = 0o440
)

// NodeServer implements csi.NodeServer interface.
type NodeServer struct {
	Logger         *zap.Logger
	NodeID         string
	Mounter        mount.Interface
	StorageBackend storage.Storage
}

// NodeStageVolume implements the csi.NodeServer interface.
// Stages volume.
func (ns *NodeServer) NodeStageVolume(
	_ context.Context,
	_ *csi.NodeStageVolumeRequest,
) (*csi.NodeStageVolumeResponse, error) {
	return nil, fmt.Errorf("failed NodeStageVolume: %w",
		status.Error(codes.Unimplemented, "NodeStageVolume not implemented"),
	)
}

// NodeUnstageVolume implements the csi.NodeServer interface.
// Unstages volume.
func (ns *NodeServer) NodeUnstageVolume(
	_ context.Context,
	_ *csi.NodeUnstageVolumeRequest,
) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, fmt.Errorf("failed NodeUnstageVolume: %w",
		status.Error(codes.Unimplemented, "NodeUnstageVolume not implemented"),
	)
}

// NodePublishVolume implements the csi.NodeServer interface.
// Publishes volume.
func (ns *NodeServer) NodePublishVolume(
	_ context.Context,
	req *csi.NodePublishVolumeRequest,
) (*csi.NodePublishVolumeResponse, error) {
	targetPath := req.GetTargetPath()
	vCtx := req.GetVolumeContext()
	volumeID := req.GetVolumeId()

	success := false

	defer func() {
		if !success {
			_ = ns.Mounter.Unmount(targetPath)
			_ = ns.StorageBackend.RemoveVolume(volumeID)
		}
	}()

	_, err := ns.StorageBackend.WriteVolume(volumeID, vCtx)
	if err != nil {
		return nil, fmt.Errorf("unexpected error writing to storage backend: %w", err)
	}

	isMountPoint, err := ns.Mounter.IsMountPoint(targetPath)

	switch {
	case os.IsNotExist(err):
		err := os.MkdirAll(req.GetTargetPath(), roPerms)
		if err != nil {
			return nil, fmt.Errorf("failed to make directories: %w", err)
		}

		isMountPoint = false
	case err != nil:
		return nil, fmt.Errorf("unexpected error checking mount point: %w", err)
	}

	if isMountPoint {
		success = true

		return &csi.NodePublishVolumeResponse{}, nil
	}

	err = ns.Mounter.Mount(ns.StorageBackend.PathForVolume(volumeID), targetPath, "", []string{"bind", "ro"})
	if err != nil {
		return nil, fmt.Errorf("error mounting volume to pod %w", err)
	}

	success = true

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume implements the csi.NodeServer interface.
// Unpublishes volume.
func (ns *NodeServer) NodeUnpublishVolume(
	_ context.Context,
	req *csi.NodeUnpublishVolumeRequest,
) (*csi.NodeUnpublishVolumeResponse, error) {
	// check to see if volume is mounted
	isMounted, err := ns.Mounter.IsMountPoint(req.GetTargetPath())
	if err != nil {
		return nil, fmt.Errorf("failed to determine if target path is mount point: %w", err)
	}

	if isMounted {
		err := ns.Mounter.Unmount(req.GetTargetPath())
		if err != nil {
			return nil, fmt.Errorf("failed to unmount volume: %w", err)
		}
	}

	err = ns.StorageBackend.RemoveVolume(req.GetVolumeId())
	if err != nil {
		return nil, fmt.Errorf("failed to remove directories: %w", err)
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats implements the csi.NodeServer interface.
// Not implemented.
func (ns *NodeServer) NodeGetVolumeStats(
	_ context.Context,
	_ *csi.NodeGetVolumeStatsRequest,
) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, fmt.Errorf("failed NodeGetVolumeStats: %w",
		status.Error(codes.Unimplemented, "NodeGetVolumeStats not implemented"),
	)
}

// NodeExpandVolume implements the csi.NodeServer interface.
// Not implemented.
func (ns *NodeServer) NodeExpandVolume(
	_ context.Context,
	_ *csi.NodeExpandVolumeRequest,
) (*csi.NodeExpandVolumeResponse, error) {
	return nil, fmt.Errorf("failed NodeExpandVolume: %w",
		status.Error(codes.Unimplemented, "NodeExpandVolume not implemented"),
	)
}

// NodeGetCapabilities implements the csi.NodeServer interface.
// Gets node capabilities.
func (ns *NodeServer) NodeGetCapabilities(
	_ context.Context,
	_ *csi.NodeGetCapabilitiesRequest,
) (*csi.NodeGetCapabilitiesResponse, error) {
	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			{
				//nolint:nosnakecase // library code.
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_UNKNOWN,
					},
				},
			},
		},
	}, nil
}

// NodeGetInfo implements the csi.NodeServer interface.
// Returns node name.
func (ns *NodeServer) NodeGetInfo(
	_ context.Context,
	_ *csi.NodeGetInfoRequest,
) (*csi.NodeGetInfoResponse, error) {
	return &csi.NodeGetInfoResponse{
		NodeId: ns.NodeID,
	}, nil
}
