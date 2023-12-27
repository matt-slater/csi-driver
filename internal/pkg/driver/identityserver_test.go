package driver_test

import (
	"context"
	"csi-driver/internal/pkg/driver"
	"reflect"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func TestIdentityServer_GetPluginInfo(t *testing.T) {
	t.Parallel()

	type fields struct {
		Name    string
		Version string
	}

	tests := []struct {
		name    string
		fields  fields
		want    *csi.GetPluginInfoResponse
		wantErr bool
	}{
		{
			name: "correct response",
			fields: fields{
				Name:    "test-driver.csi.io",
				Version: "v1.0.0",
			},
			want: &csi.GetPluginInfoResponse{
				Name:          "test-driver.csi.io",
				VendorVersion: "v1.0.0",
			},
		},
		{
			name: "empty name",
			fields: fields{
				Name:    "",
				Version: "v1.0.0",
			},
			wantErr: true,
		},
		{
			name: "empty version",
			fields: fields{
				Name:    "test-driver.csi.io",
				Version: "",
			},
			wantErr: true,
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			is := &driver.IdentityServer{
				Name:    testCase.fields.Name,
				Version: testCase.fields.Version,
			}
			got, err := is.GetPluginInfo(context.Background(), &csi.GetPluginInfoRequest{})
			if (err != nil) != testCase.wantErr {
				t.Errorf("IdentityServer.GetPluginInfo() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}
			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("IdentityServer.GetPluginInfo() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestIdentityServer_GetPluginCapabilities(t *testing.T) {
	t.Parallel()

	identityServer := driver.IdentityServer{}

	resp, err := identityServer.GetPluginCapabilities(context.Background(), &csi.GetPluginCapabilitiesRequest{})
	if err != nil {
		t.Fatalf("unexpected error getting capabilities: %v", err)
	}

	if resp == nil {
		t.Fatal("unexpected nil response")
	}
}

func TestIdentityServer_Probe(t *testing.T) {
	t.Parallel()

	identityServer := driver.IdentityServer{}

	resp, err := identityServer.Probe(context.Background(), &csi.ProbeRequest{})
	if err != nil {
		t.Fatalf("unexpected error during probe: %v", err)
	}

	if resp == nil {
		t.Fatal("unexpected nil response")
	}
}
