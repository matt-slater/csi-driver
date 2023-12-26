package storage

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"k8s.io/mount-utils"
)

type Filesystem struct {
	logger        *zap.Logger
	baseDirectory string
	storage       fs.FS
	mounter       mount.Interface
	tempfsPath    string
}

func NewFilesystem(logger *zap.Logger, baseDirectory string, rootFS fs.FS) (*Filesystem, error) {
	tempfsPath := filepath.Join(baseDirectory, "inmemfs")

	filesystem := &Filesystem{
		logger:        logger,
		baseDirectory: baseDirectory,
		storage:       rootFS,
		mounter:       mount.New(""),
		tempfsPath:    tempfsPath,
	}

	filesystem.tempfsPath = filesystem.baseDirectory

	isMount, err := filesystem.mounter.IsMountPoint(filesystem.tempfsPath)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("unexpected error checking mount point: %w", err)
		}

		err := os.MkdirAll(filesystem.tempfsPath, 0700)
		if err != nil {
			return nil, fmt.Errorf("failed to create tempfs directories: %w", err)
		}
	}

	if !isMount {
		err := filesystem.mounter.Mount("tempfs", filesystem.tempfsPath, "tmpfs", []string{})
		if err != nil {
			return nil, fmt.Errorf("failed to mount tmpfs: %w", err)
		}

		logger.Info("mounted new tmpfs",
			zap.String("path", filesystem.tempfsPath),
		)
	}

	return filesystem, nil
}

func (fs *Filesystem) WriteVolume(id string, targetPath string, vCtx map[string]string, data []byte) error {
	panic("not implemented") // TODO: Implement
}

func (fs *Filesystem) PathForVolume(id string) string {
	panic("not implemented") // TODO: Implement
}

func (fs *Filesystem) ListVolumes() []string {
	panic("not implemented") // TODO: Implement
}

func (fs *Filesystem) RemoveVolume(id string) error {
	panic("not implemented") // TODO: Implement
}
