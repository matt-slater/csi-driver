package driver_test

import (
	"context"
	"csi-driver/internal/pkg/driver"
	"csi-driver/internal/pkg/storage"
	"os"
	"reflect"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest"
	"k8s.io/mount-utils"
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

func TestNodeServer_NodePublishVolume(t *testing.T) {
	t.Parallel()

	type fields struct {
		Logger         *zap.Logger
		NodeID         string
		Mounter        mount.Interface
		StorageBackend storage.Storage
	}

	type args struct {
		req *csi.NodePublishVolumeRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodePublishVolumeResponse
		wantErr bool
	}{
		{
			name: "success publish volume",
			fields: fields{
				Logger: zaptest.NewLogger(t),
				NodeID: "test-node",
				Mounter: mount.NewFakeMounter([]mount.MountPoint{
					{
						Path: "/var/lib/kubelet/pods/test-pod/volumes/test-volume/mount",
					},
				}),
				StorageBackend: &storage.MockStorage{},
			},
			args: args{
				req: &csi.NodePublishVolumeRequest{
					VolumeId:       "x1b3n4",
					PublishContext: map[string]string{"a": "b", "c": "d"},
					TargetPath:     "/data",
				},
			},
			want:    &csi.NodePublishVolumeResponse{},
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			tmpDir, err := os.MkdirTemp("", testCase.name+"*")
			if err != nil {
				t.Fatalf("failed to create a temp dir: %v", err)
			}

			defer func() {
				err := os.RemoveAll(tmpDir)
				if err != nil {
					t.Fatalf("failed to remove temp dir: %v", err)
				}
			}()

			testCase.args.req.TargetPath = tmpDir + testCase.args.req.GetTargetPath()

			nodeServer := &driver.NodeServer{
				Logger:         testCase.fields.Logger,
				NodeID:         testCase.fields.NodeID,
				Mounter:        testCase.fields.Mounter,
				StorageBackend: testCase.fields.StorageBackend,
			}
			got, err := nodeServer.NodePublishVolume(context.Background(), testCase.args.req)
			if (err != nil) != testCase.wantErr {
				t.Errorf("NodeServer.NodePublishVolume() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("NodeServer.NodePublishVolume() = %v, want %v", got, testCase.want)
			}

			t.Logf("mounts: %+v", testCase.fields.Mounter)
		})
	}
}

func TestNodeServer_NodeUnpublishVolume(t *testing.T) {
	t.Parallel()

	type fields struct {
		Logger         *zap.Logger
		NodeID         string
		Mounter        mount.Interface
		StorageBackend storage.Storage
	}

	type args struct {
		req *csi.NodeUnpublishVolumeRequest
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *csi.NodeUnpublishVolumeResponse
		wantErr bool
	}{
		{
			name: "unmount success, mount doesnt exist",
			fields: fields{
				Logger:         zaptest.NewLogger(t),
				NodeID:         "test-node",
				Mounter:        &mount.FakeMounter{},
				StorageBackend: &storage.MockStorage{},
			},
			args: args{
				req: &csi.NodeUnpublishVolumeRequest{},
			},
			want:    &csi.NodeUnpublishVolumeResponse{},
			wantErr: false,
		},
	}

	for _, testCase := range tests {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			mtPt, err := os.MkdirTemp("", testCase.name+"mount*")
			if err != nil {
				t.Fatalf("failed to create a temp dir: %v", err)
			}

			defer func() {
				err := os.RemoveAll(mtPt)
				if err != nil {
					t.Fatalf("failed to remove temp dir: %v", err)
				}
			}()

			dataPt, err := os.MkdirTemp("", testCase.name+"data*")
			if err != nil {
				t.Fatalf("failed to create a temp dir: %v", err)
			}

			defer func() {
				err := os.RemoveAll(dataPt)
				if err != nil {
					t.Fatalf("failed to remove temp dir: %v", err)
				}
			}()

			err = testCase.fields.Mounter.Mount(dataPt, mtPt, "tmpfs", []string{})
			if err != nil {
				t.Fatalf("failed to mount tmpdir: %v", err)
			}

			testCase.args.req.TargetPath = mtPt

			nodeServer := &driver.NodeServer{
				Logger:         testCase.fields.Logger,
				NodeID:         testCase.fields.NodeID,
				Mounter:        testCase.fields.Mounter,
				StorageBackend: testCase.fields.StorageBackend,
			}

			got, err := nodeServer.NodeUnpublishVolume(context.Background(), testCase.args.req)
			if (err != nil) != testCase.wantErr {
				t.Errorf("NodeServer.NodeUnpublishVolume() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if !reflect.DeepEqual(got, testCase.want) {
				t.Errorf("NodeServer.NodeUnpublishVolume() = %v, want %v", got, testCase.want)
			}
		})
	}
}
