package storage

import (
	"errors"
)

type MockStorage struct {
	ShouldErr bool
	Path      string
	Volumes   []string
}

var errMock = errors.New("mock error")

func (ms *MockStorage) WriteVolume(_ string, _ map[string]string) (bool, error) {
	if ms.ShouldErr {
		return false, errMock
	}

	return true, nil
}

func (ms *MockStorage) PathForVolume(_ string) string {
	return ms.Path
}

func (ms *MockStorage) ListVolumes() ([]string, error) {
	if ms.ShouldErr {
		return nil, errMock
	}

	return ms.Volumes, nil
}

func (ms *MockStorage) RemoveVolume(_ string) error {
	if ms.ShouldErr {
		return errMock
	}

	return nil
}
