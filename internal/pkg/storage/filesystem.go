package storage

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"go.uber.org/zap"
	"k8s.io/mount-utils"
)

type Filesystem struct {
	logger  *zap.Logger
	storage fs.FS
	mounter mount.Interface
	baseDir string
}

const (
	rwePerms = 0o700
)

func NewFilesystem(
	logger *zap.Logger,
	baseDir string,
	rootFS fs.FS,
	mounter mount.Interface,
) (*Filesystem, error) {
	filesystem := &Filesystem{
		logger:  logger,
		storage: rootFS,
		mounter: mounter,
		baseDir: baseDir,
	}

	isMount, err := filesystem.mounter.IsMountPoint(filesystem.baseDir)
	if err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("unexpected error checking mount point: %w", err)
		}

		err := os.MkdirAll(filesystem.baseDir, rwePerms)
		if err != nil {
			return nil, fmt.Errorf("failed to create tempfs directories: %w", err)
		}
	}

	if !isMount {
		err := filesystem.mounter.Mount("tempfs", filesystem.baseDir, "tmpfs", []string{})
		if err != nil {
			return nil, fmt.Errorf("failed to mount tmpfs: %w", err)
		}

		logger.Info("mounted new tmpfs",
			zap.String("path", filesystem.baseDir),
		)
	}

	return filesystem, nil
}

func (f *Filesystem) WriteVolume(id string, vCtx map[string]string) (bool, error) {
	datapath := filepath.Join(f.baseDir, id, "data")
	// check if volume already exists
	err := os.MkdirAll(datapath, rwePerms)
	if err != nil {
		if !errors.Is(err, os.ErrExist) {
			return false, fmt.Errorf("unexpected error creating data dir: %w", err)
		}
		// if yes, return false
		return false, nil
	}

	// if no, create volume and files in it, return true.
	f.logger.Info("creating file",
		zap.String("path", filepath.Join(datapath, vCtx["csi-driver.mattslater.io/filename"])),
	)

	file, err := os.Create(filepath.Join(datapath, vCtx["csi-driver.mattslater.io/filename"]))
	if err != nil {
		return false, fmt.Errorf("failed to create file: %w", err)
	}

	defer file.Close()

	_, err = file.WriteString(vCtx["csi-driver.mattslater.io/data"])
	if err != nil {
		return false, fmt.Errorf("failed to write data to file: %w", err)
	}

	f.logger.Info("wrote volume datapath and file successfully")

	return true, nil
}

func (f *Filesystem) PathForVolume(id string) string {
	return filepath.Join(f.baseDir, id, "data")
}

func (f *Filesystem) ListVolumes() ([]string, error) {
	dirs, err := fs.ReadDir(f.storage, f.baseDir)
	if err != nil {
		return nil, fmt.Errorf("failed to list volumes: %w", err)
	}

	var volumes []string

	for _, dir := range dirs {
		volumes = append(volumes, dir.Name())
	}

	return volumes, nil
}

func (f *Filesystem) RemoveVolume(id string) error {
	err := os.RemoveAll(filepath.Join(f.baseDir, id))
	if err != nil {
		return fmt.Errorf("failed to remove volume: %w", err)
	}

	return nil
}
