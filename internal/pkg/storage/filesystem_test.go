package storage_test

import (
	"csi-driver/internal/pkg/storage"
	"testing"
	"testing/fstest"

	"go.uber.org/zap/zaptest"
	"k8s.io/mount-utils"
)

func TestNewFilesystem(t *testing.T) {
	t.Parallel()

	fileSystem, err := storage.NewFilesystem(
		zaptest.NewLogger(t),
		"/tmp/test",
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
