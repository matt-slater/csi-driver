package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/mount-utils"
)

// NodeServer implements csi.NodeServer interface.
type NodeServer struct {
	NodeID  string
	Mounter mount.Interface
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
	_ *csi.NodePublishVolumeRequest,
) (*csi.NodePublishVolumeResponse, error) {
	panic("not implemented") // TODO: Implement
}

// NodeUnpublishVolume implements the csi.NodeServer interface.
// Unpublishes volume.
func (ns *NodeServer) NodeUnpublishVolume(
	_ context.Context,
	_ *csi.NodeUnpublishVolumeRequest,
) (*csi.NodeUnpublishVolumeResponse, error) {
	panic("not implemented") // TODO: Implement
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
