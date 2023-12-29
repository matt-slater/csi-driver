package storage_test

import (
	"csi-driver/internal/pkg/storage"
	"io/fs"
	"os"
	"strings"
	"testing"
	"testing/fstest"

	"go.uber.org/zap/zaptest"
	"k8s.io/mount-utils"
)

func TestNewFilesystem(t *testing.T) {
	t.Parallel()

	fileSystem, err := storage.NewFilesystem(
		zaptest.NewLogger(t),
		"/",
		fstest.MapFS{},
		mount.NewFakeMounter([]mount.MountPoint{}),
	)
	if err != nil {
		t.Fatalf("unexpected err creating filesystem backend: %v", err)
	}

	if fileSystem == nil {
		t.Fatal("unexpected nil filesystem backend")
	}
}

func TestFilesystem_PathForVolume(t *testing.T) {
	t.Parallel()

	fileSystem, err := storage.NewFilesystem(
		zaptest.NewLogger(t),
		"/",
		fstest.MapFS{},
		mount.NewFakeMounter([]mount.MountPoint{}),
	)
	if err != nil {
		t.Fatalf("unexpected err creating filesystem backend: %v", err)
	}

	path := fileSystem.PathForVolume("/dev/test")

	if path != "/dev/test/data" {
		t.Fatalf("unexpected data path: %s", path)
	}
}

func TestFilesystem_WriteVolume(t *testing.T) {
	t.Parallel()

	type fields struct {
		storage fs.FS
		mounter mount.Interface
	}

	type args struct {
		id   string
		vCtx map[string]string
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    bool
		wantErr bool
	}{
		{
			name: "success write volume",
			fields: fields{
				storage: fstest.MapFS{},
				mounter: mount.NewFakeMounter([]mount.MountPoint{}),
			},
			args: args{
				id: "test-id",
				vCtx: map[string]string{
					"csi-driver.mattslater.io/filename": "yolo.txt",
					"csi-driver.mattslater.io/data":     "you only live once",
				},
			},
			want: true,
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

			f, err := storage.NewFilesystem(zaptest.NewLogger(t), tmpDir, testCase.fields.storage, testCase.fields.mounter)
			if err != nil {
				t.Fatalf("failed to create filesystem: %v", err)
			}

			got, err := f.WriteVolume(testCase.args.id, testCase.args.vCtx)
			if (err != nil) != testCase.wantErr {
				t.Errorf("Filesystem.WriteVolume() error = %v, wantErr %v", err, testCase.wantErr)

				return
			}

			if got != testCase.want {
				t.Errorf("Filesystem.WriteVolume() = %v, want %v", got, testCase.want)
			}
		})
	}
}

func TestFilesystem_ListVolumes(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "list-volumes*")
	if err != nil {
		t.Fatalf("failed to create a temp dir: %v", err)
	}

	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}()

	err = os.Mkdir(tmpDir+"/matt", 0644)
	if err != nil {
		t.Fatalf("failed to create dir in temp dir: %v", err)
	}

	trimmedDir := strings.TrimPrefix(tmpDir, "/")

	file, err := storage.NewFilesystem(
		zaptest.NewLogger(t),
		trimmedDir,
		os.DirFS("/"),
		mount.NewFakeMounter([]mount.MountPoint{}),
	)
	if err != nil {
		t.Fatalf("failed to create filesystem: %v", err)
	}

	vols, err := file.ListVolumes()
	if err != nil {
		t.Fatalf("failed to list volumes: %v", err)
	}

	if vols == nil {
		t.Fatalf("unexpected nil volumes")
	}

	t.Logf("volumes: %s", vols)
}

func TestFilesystem_RemoveVolume(t *testing.T) {
	t.Parallel()

	tmpDir, err := os.MkdirTemp("", "list-volumes*")
	if err != nil {
		t.Fatalf("failed to create a temp dir: %v", err)
	}

	defer func() {
		err := os.RemoveAll(tmpDir)
		if err != nil {
			t.Fatalf("failed to remove temp dir: %v", err)
		}
	}()

	err = os.Mkdir(tmpDir+"/matt", 0644)
	if err != nil {
		t.Fatalf("failed to create dir in temp dir: %v", err)
	}

	trimmedDir := strings.TrimPrefix(tmpDir, "/")

	file, err := storage.NewFilesystem(
		zaptest.NewLogger(t),
		trimmedDir,
		os.DirFS("/"),
		mount.NewFakeMounter([]mount.MountPoint{}),
	)
	if err != nil {
		t.Fatalf("failed to create filesystem: %v", err)
	}

	err = file.RemoveVolume("matt")
	if err != nil {
		t.Fatalf("failed to remove volume: %v", err)
	}
}
