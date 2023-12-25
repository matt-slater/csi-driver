// Package driver contains structs and methods to implement the CSI spec.
package driver

import (
	"context"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// IdentityServer implements csi.IdentityServer.
type IdentityServer struct {
	Name    string
	Version string
}

// GetPluginInfo implements csi.IdentityServer.GetPluginInfo.
func (is *IdentityServer) GetPluginInfo(
	_ context.Context,
	_ *csi.GetPluginInfoRequest,
) (*csi.GetPluginInfoResponse, error) {
	if is.Name == "" {
		return nil, status.Error(codes.Unavailable, "driver name not configured")
	}

	if is.Version == "" {
		return nil, status.Error(codes.Unavailable, "driver is missing version")
	}

	return &csi.GetPluginInfoResponse{
		Name:          is.Name,
		VendorVersion: is.Version,
	}, nil
}

// GetPluginCapabilities implements csi.IdentityServer.GetPluginCapabilities.
func (is *IdentityServer) GetPluginCapabilities(
	_ context.Context,
	_ *csi.GetPluginCapabilitiesRequest,
) (*csi.GetPluginCapabilitiesResponse, error) {
	return &csi.GetPluginCapabilitiesResponse{
		Capabilities: []*csi.PluginCapability{
			{
				Type: &csi.PluginCapability_Service_{
					Service: &csi.PluginCapability_Service{
						Type: csi.PluginCapability_Service_CONTROLLER_SERVICE,
					},
				},
			},
		},
	}, nil
}

// Probe implements csi.IdentityServer.Probe.
func (is *IdentityServer) Probe(context.Context, *csi.ProbeRequest) (*csi.ProbeResponse, error) {
	return &csi.ProbeResponse{}, nil
}
