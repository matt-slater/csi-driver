package driver

import (
	"context"
	"csi-driver/internal/pkg/storage"
	"fmt"
	"os"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/mount-utils"
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
	return nil, status.Error(codes.Unimplemented, "NodeStageVolume not implemented")
}

// NodeUnstageVolume implements the csi.NodeServer interface.
// Unstages volume.
func (ns *NodeServer) NodeUnstageVolume(
	_ context.Context,
	_ *csi.NodeUnstageVolumeRequest,
) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeUnstageVolume not implemented")
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

	volumeLogger := ns.Logger.With(
		zap.String("target path", targetPath),
		zap.String("volume ID", volumeID),
		zap.Any("colume context", vCtx),
		zap.String("pod name", vCtx["csi.storage.k8s.io/pod.name"]),
	)

	err := ns.StorageBackend.WriteVolume(volumeID, targetPath, vCtx, []byte(vCtx["csi-driver.mattslater.io/data"]))
	if err != nil {
		return nil, fmt.Errorf("unexpected error writing to storage backend: %w", err)
	}

	volumeLogger.Info("successfully sotred volume in backend")

	volumeLogger.Info("ensuring volume is mounted to pod")

	isMountPoint, err := ns.Mounter.IsMountPoint(targetPath)
	switch {
	case os.IsNotExist(err):
		if err := os.MkdirAll(req.GetTargetPath(), 0440); err != nil {
			return nil, err
		}
		isMountPoint = false
	case err != nil:
		return nil, fmt.Errorf("unexpected error checking mount point: %w", err)
	}

	if isMountPoint {
		volumeLogger.Info("volume is already mounted to pod, nothing to do")
		success = true
		return &csi.NodePublishVolumeResponse{}, nil
	}

	volumeLogger.Info("bind mounting data directory to the pod's mount namespace")

	err = ns.Mounter.Mount(ns.StorageBackend.PathForVolume(volumeID), targetPath, "", []string{"bind", "ro"})
	if err != nil {
		return nil, fmt.Errorf("error mounting volume to pod %w", err)
	}

	volumeLogger.Info("successfully mounted volume to pod")
	success = true

	return &csi.NodePublishVolumeResponse{}, nil
}

// NodeUnpublishVolume implements the csi.NodeServer interface.
// Unpublishes volume.
func (ns *NodeServer) NodeUnpublishVolume(
	_ context.Context,
	_ *csi.NodeUnpublishVolumeRequest,
) (*csi.NodeUnpublishVolumeResponse, error) {
	return &csi.NodeUnpublishVolumeResponse{}, nil
}

// NodeGetVolumeStats implements the csi.NodeServer interface.
// Not implemented.
func (ns *NodeServer) NodeGetVolumeStats(
	_ context.Context,
	_ *csi.NodeGetVolumeStatsRequest,
) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeGetVolumeStats not implemented")
}

// NodeExpandVolume implements the csi.NodeServer interface.
// Not implemented.
func (ns *NodeServer) NodeExpandVolume(
	_ context.Context,
	_ *csi.NodeExpandVolumeRequest,
) (*csi.NodeExpandVolumeResponse, error) {
	return nil, status.Error(codes.Unimplemented, "NodeExpandVolume not implemented")
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
