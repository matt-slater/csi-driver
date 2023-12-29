package storage

import (
	"errors"
)

type MockStorage struct {
	ShouldErr bool
	Path      string
	Volumes   []string
}

func (ms *MockStorage) WriteVolume(id string, vCtx map[string]string) (bool, error) {
	if ms.ShouldErr {
		return false, errors.New("error writing volume")
	}

	return true, nil
}

func (ms *MockStorage) PathForVolume(id string) string {
	return ms.Path
}

func (ms *MockStorage) ListVolumes() ([]string, error) {
	if ms.ShouldErr {
		return nil, errors.New("error getting volumes")
	}

	return ms.Volumes, nil
}

func (ms *MockStorage) RemoveVolume(id string) error {
	if ms.ShouldErr {
		return errors.New("error removing volumes")
	}

	return nil
}
